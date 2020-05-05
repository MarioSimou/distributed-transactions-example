package tests

import (
	"fmt"
	"math"
	m "products/internal/mocks"
	r "products/internal/rabbitmq"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var amqpTable = mock.AnythingOfType("amqp.Table")


type PublisherSuite struct {
	suite.Suite
	ci *m.ConnectionInterface
	chi []m.ChannelInterface
}
func (ps *PublisherSuite) TestGetChannels(){
	var channels = []r.ChannelInterface{
		&m.ChannelInterface{},
		&m.ChannelInterface{},
	}
	var pub = r.Publisher{Channels: channels}

	assert.Equal(ps.T(), len(channels), len(pub.GetChannels()))
}

func (ps *PublisherSuite) TestGetConnection(){
	var conn = &m.ConnectionInterface{}
	var table = []struct{
		expectedRes r.ConnectionInterface
	}{
		{
			expectedRes: conn,
		},
	}
	var t = ps.T()
	var assert = assert.New(t)
	for _, row := range table {
		var pub = r.Publisher{Conn: conn}

		assert.Equal(pub.GetConnection(), row.expectedRes)
	}
}

func (ps *PublisherSuite) TestGetChannel(){
	var firstChan = &r.Channel{QueueName: "first"}
	var secondChan = &r.Channel{QueueName: "second"}
	var table = []struct{
		queueName string
		expectedRes r.ChannelInterface
	}{
		{
			queueName: "second",
			expectedRes: secondChan,
		},
		{
			queueName: "third",
			expectedRes: nil,
		},
	}
	var t = ps.T()
	var assert = assert.New(t)

	for _, row := range table {
		var pub = &r.Publisher{
			Channels: []r.ChannelInterface{firstChan, secondChan},
		}

		var channel = pub.GetChannel(row.queueName)
		assert.Equal(channel, row.expectedRes)
	}
}

func (ps *PublisherSuite) TestPub(){
	var table = []struct{
		queueName string
		body interface{}
		e error 
		setExpectations func(ci *m.ChannelInterface, conni *m.ConnectionInterface)
	}{
		{
			queueName: "first",
			body: nil,
			e: r.CHANNEL_NOT_FOUND,
			setExpectations: func (channel *m.ChannelInterface, conn *m.ConnectionInterface){
				channel.On("GetQueueName").Return("")
			},
		},
		{
			queueName: "first",
			body: math.Inf(1),
			e: fmt.Errorf("json: unsupported value: +Inf"),
			setExpectations: func (channel *m.ChannelInterface, conn *m.ConnectionInterface){
				channel.On("GetQueueName").Return("first")
			},
		},
		{
			queueName: "first",
			body: []byte(`{"message": "some text"}`),
			e: fmt.Errorf("Queue error"),
			setExpectations: func (channel *m.ChannelInterface, conn *m.ConnectionInterface){
				var e = fmt.Errorf("Queue error")
				channel.On("GetQueueName").Return("first")
				channel.On("Publish", "first","", false, false, mock.AnythingOfType("amqp.Publishing")).Return(e)
			},
		},
		{
			queueName: "first",
			body: []byte(`{"message": "some text"}`),
			e: nil,
			setExpectations: func (channel *m.ChannelInterface, conn *m.ConnectionInterface){
				channel.On("GetQueueName").Return("first")
				channel.On("Publish", "first","", false, false, mock.AnythingOfType("amqp.Publishing")).Return(nil)
			},
		},
	}

	var t = ps.T()
	var assert = assert.New(t)

	for _, row := range table {
		var conn = &m.ConnectionInterface{}
		var channel = &m.ChannelInterface{}
		var channels = []r.ChannelInterface{channel}
		var pub = r.Publisher{Conn: conn, Channels: channels}
		
		row.setExpectations(channel,conn)

		if e := pub.Pub(row.queueName, row.body); e != nil {
			assert.EqualError(e,row.e.Error())
		}

		conn.AssertExpectations(t)
		channel.AssertExpectations(t)
	}
}

func (ps *PublisherSuite) TestNewPublisher(){
	var table = []struct{
		queueNames []string
		setAssertions func(*m.ConnectionInterface) []*m.ChannelInterface
		e error
	}{
		{
			queueNames: []string{"first"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var e = fmt.Errorf("Channel error")
				conn.On("Channel", "first").Return(nil, e)
				return nil
			},
			e: fmt.Errorf("Channel error"),
		},
		{
			queueNames: []string{"first"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var channel = &m.ChannelInterface{}
				var e = fmt.Errorf("Queue Declare error")
				conn.On("Channel", "first").Return(channel, nil)
				channel.On("QueueDeclare", "first", true,false,false,false, amqpTable).Return(e)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Queue Declare error"),
		},
		{
			queueNames: []string{"first"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var channel = &m.ChannelInterface{}
				var e = fmt.Errorf("Exchange Declare Error")

				conn.On("Channel", "first").Return(channel, nil)
				channel.On("QueueDeclare", "first", true,false,false,false, amqpTable).Return(nil)
				channel.On("ExchangeDeclare", "first", r.FANOUT, true, false ,false, false, amqpTable).Return(e)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Exchange Declare Error"),
		},
		{
			queueNames: []string{"first"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var channel = &m.ChannelInterface{}
				var e = fmt.Errorf("Queue Bind Error")

				conn.On("Channel", "first").Return(channel, nil)
				channel.On("QueueDeclare", "first", true,false,false,false, amqpTable).Return(nil)
				channel.On("ExchangeDeclare", "first", r.FANOUT, true, false ,false, false, amqpTable).Return(nil)
				channel.On("QueueBind", "first", "", "first", false, amqpTable).Return(e)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Queue Bind Error"),
		},
		{
			queueNames: []string{"first"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var channel = &m.ChannelInterface{}
				
				conn.On("Channel", "first").Return(channel, nil)
				channel.On("QueueDeclare", "first", true,false,false,false, amqpTable).Return(nil)
				channel.On("ExchangeDeclare", "first", r.FANOUT, true, false ,false, false, amqpTable).Return(nil)
				channel.On("QueueBind", "first", "", "first", false, amqpTable).Return(nil)

				return []*m.ChannelInterface{channel}
			},
			e: nil,
		},
		{
			queueNames: []string{"first", "second"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var channel = &m.ChannelInterface{}
				var e = fmt.Errorf("Channel Error")

				conn.On("Channel", "first").Return(channel, nil)
				conn.On("Channel", "second").Return(nil, e)

				channel.On("QueueDeclare", "first", true,false,false,false, amqpTable).Return(nil)
				channel.On("ExchangeDeclare", "first", r.FANOUT, true, false ,false, false, amqpTable).Return(nil)
				channel.On("QueueBind", "first", "", "first", false, amqpTable).Return(nil)
				channel.On("Close").Return(nil)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Channel Error"),
		},
		{
			queueNames: []string{"first", "second"},
			setAssertions: func(conn *m.ConnectionInterface) []*m.ChannelInterface {
				var channel = &m.ChannelInterface{}
				var e = fmt.Errorf("Closing Channel Error")

				conn.On("Channel", "first").Return(channel, nil)
				conn.On("Channel", "second").Return(nil, fmt.Errorf("Some random error"))

				channel.On("QueueDeclare", "first", true,false,false,false, amqpTable).Return(nil)
				channel.On("ExchangeDeclare", "first", r.FANOUT, true, false ,false, false, amqpTable).Return(nil)
				channel.On("QueueBind", "first", "", "first", false, amqpTable).Return(nil)
				channel.On("Close").Return(e)

				return []*m.ChannelInterface{channel}
			},
			e: fmt.Errorf("Closing Channel Error"),
		},
	}

	var t = ps.T()
	var assert = assert.New(t)
	for _, row := range table {
		var conn = &m.ConnectionInterface{}
		var channels = row.setAssertions(conn)

		if pub, e := r.NewPublisher(row.queueNames,conn); e != nil {
			assert.EqualError(e,row.e.Error())
		} else {
			assert.EqualValues(pub.GetConnection(), conn)
			assert.Len(pub.GetChannels(), len(row.queueNames))
		}

		conn.AssertExpectations(t)

		for _, channel := range channels {
			channel.AssertExpectations(t)
		}
	}
}


func TestPublisherSuite(t *testing.T){
	suite.Run(t, new(PublisherSuite))
}