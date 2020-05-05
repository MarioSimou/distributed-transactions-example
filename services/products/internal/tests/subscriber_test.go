package tests

import (
	"fmt"
	r "products/internal/rabbitmq"
	"testing"

	m "products/internal/mocks"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)


type SubscriberSuite struct {
	suite.Suite
}

func (ss *SubscriberSuite) TestGetQueue(){
	var sub = r.Subscriber{
		QueueName: "queue",
	} 

	assert.Equal(ss.T(),sub.GetQueue(), "queue")
}


func (ss *SubscriberSuite) TestGetHandlerFunc(){
	var msg = r.Message{}
	var dummyHandler = func(m r.Message) error {
		return fmt.Errorf("Some")
	}
	var sub = r.Subscriber{
		HandlerFunc: dummyHandler,
	} 

	assert.EqualError(ss.T(),sub.GetHandlerFunc()(msg), dummyHandler(msg).Error())
}

func (ss *SubscriberSuite) TestNewSubscription(){
	var table = []struct{
		setAssertions func(conn *m.ConnectionInterface) []*m.ChannelInterface
		e error
	}{
		{
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var e = fmt.Errorf("Channel Error")
				conn.On("Channel", "first").Return(nil, e)

				return nil
			},
			e: fmt.Errorf("Channel Error"),
		},
		{
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface{
				var e = fmt.Errorf("Consume Error")
				var channel = &m.ChannelInterface{}
				conn.On("Channel", "first").Return(channel,nil)
				channel.On("Consume", "first", "", true, false, false, false, amqpTable).Return(nil,e)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Consume Error"),
		},
		{
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface{
				var channel = &m.ChannelInterface{}
				var deliveryChan = make(<-chan amqp.Delivery)

				conn.On("Channel", "first").Return(channel,nil)
				channel.On("Consume", "first", "", true, false, false, false, amqpTable).Return(deliveryChan,nil)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Consume Error"),
		},
	}

	var t = ss.T()
	var assert = assert.New(t)
	for _, row := range table {
		var conn = &m.ConnectionInterface{}
		var subscribers = r.Subscribers{
			r.Subscriber{
				QueueName: "first",
				HandlerFunc: func(msg r.Message) error { 
					return nil
				},
			},
		}
		// set expectations
		var channels = row.setAssertions(conn)

		// assertions
		var subChan, e = r.NewSubscription(subscribers, conn)
		if e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.NotNil(subChan)	
		}
		

		conn.AssertExpectations(t)
		for _, channel := range channels {
			channel.AssertExpectations(t)
		}
	}
}

func TestSubscriberSuite(t *testing.T){
	suite.Run(t, new(SubscriberSuite))
}