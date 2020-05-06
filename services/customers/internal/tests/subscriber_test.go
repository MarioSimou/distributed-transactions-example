package tests

import (
	"context"
	m "customers/internal/mocks"
	r "customers/internal/rabbitmq"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)


type SubscriberSuite struct {
	suite.Suite
}

func (ss *SubscriberSuite) TestGetQueue(){
	var sub = r.NewSubscriber("queue", nil, nil)

	assert.Equal(ss.T(),sub.GetQueue(), "queue")
}

func (ss *SubscriberSuite) TestGetContext(){
	var ctx = context.Background()
	var sub = r.NewSubscriber("", nil, ctx)

	assert.Equal(ss.T(),sub.GetContext(),ctx)
}


func (ss *SubscriberSuite) TestGetHandlerFunc(){
	var msg = r.Message{}
	var dummyHandler = func(m r.Message,ctx context.Context) error {
		return fmt.Errorf("Some")
	}
	var sub = r.NewSubscriber("", dummyHandler, nil)
	var ctx = context.Background()

	assert.EqualError(ss.T(),sub.GetHandlerFunc()(msg,ctx), dummyHandler(msg,ctx).Error())
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
		var handler = func(msg r.Message, ctx context.Context) error { 
			return nil
		}

		var subscribers = r.NewSubscribers(
			r.NewSubscriber("first",handler,nil),
		)
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