package tests

import (
	"fmt"
	m "products/internal/mocks"
	r "products/internal/rabbitmq"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ChannelSuite struct {
	suite.Suite
	cid *m.ChannelInterfaceDefault
	ci *m.ChannelInterface
}
func (cs *ChannelSuite) TestGetQueueName(){
	var ch = r.Channel{QueueName: "queue"}
	var queueName = ch.GetQueueName()
	assert.Equal(cs.T(), "queue", queueName)
}

func (cs *ChannelSuite) TestClose(){
	var table = []struct{
		setExpectations func(*m.ChannelInterfaceDefault)
		e error		
	}{
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var e = fmt.Errorf("Internal Error")
				cid.On("Close").Return(e)
			},
			e: fmt.Errorf("Internal Error"),
		},
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				cid.On("Close").Return(nil)
			},
			e: nil,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var cid = &m.ChannelInterfaceDefault{}
		var ch = r.Channel{C: cid}
		row.setExpectations(cid)


		if e := ch.Close(); e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.Nil(e)
		}

		cid.AssertExpectations(t)
	}
}

func (cs *ChannelSuite) TestQueueDeclare(){
	var table = []struct{
		setExpectations func(*m.ChannelInterfaceDefault)
		e error		
	}{
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var e = fmt.Errorf("Internal Error")
				var queue = amqp.Queue{}
				cid.On("QueueDeclare", "queue", mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.AnythingOfType("bool"),  mock.AnythingOfType("amqp.Table")).Return(queue,e)
			},
			e: fmt.Errorf("Internal Error"),
		},
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var queue = amqp.Queue{Name: "queue"}
				cid.On("QueueDeclare", "queue", mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.AnythingOfType("bool"), mock.AnythingOfType("bool"),  mock.AnythingOfType("amqp.Table")).Return(queue,nil)
			},
			e: nil,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var queueName = "queue"
		var cid = &m.ChannelInterfaceDefault{}
		var ch = r.Channel{C: cid}

		row.setExpectations(cid)

		if e := ch.QueueDeclare(queueName, false, false, false,false, nil); e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.Nil(e)
		}

		cid.AssertExpectations(t)
	}
}

func (cs *ChannelSuite) TestExchangeDeclare(){
	var table = []struct{
		setExpectations func(*m.ChannelInterfaceDefault)
		e error		
	}{
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var e = fmt.Errorf("Internal Error")
				cid.On("ExchangeDeclare", "queue", "fanout",false,false,false,false,mock.AnythingOfType("amqp.Table")).Return(e)
			},
			e: fmt.Errorf("Internal Error"),
		},
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				cid.On("ExchangeDeclare", "queue", "fanout",false,false,false,false,mock.AnythingOfType("amqp.Table")).Return(nil)
			},
			e: nil,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var exchangeName = "queue"
		var kind = "fanout"
		var cid = &m.ChannelInterfaceDefault{}
		var ch = r.Channel{C: cid}

		row.setExpectations(cid)

		if e := ch.ExchangeDeclare(exchangeName, kind, false, false, false,false, nil); e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.Nil(e)
		}

		cid.AssertExpectations(t)
	}
}

func (cs *ChannelSuite) TestQueueBind(){
	var table = []struct{
		setExpectations func(*m.ChannelInterfaceDefault)
		e error		
	}{
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var e = fmt.Errorf("Internal Error")
				cid.On("QueueBind", "queue","","queue", false, mock.AnythingOfType("amqp.Table")).Return(e)
			},
			e: fmt.Errorf("Internal Error"),
		},
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				cid.On("QueueBind", "queue","","queue", false, mock.AnythingOfType("amqp.Table")).Return(nil)
			},
			e: nil,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var queueName = "queue"
		var key= ""
		var exchangeName = "queue"
		var cid = &m.ChannelInterfaceDefault{}
		var ch = r.Channel{C: cid}

		row.setExpectations(cid)

		if e := ch.QueueBind(queueName, key,exchangeName, false,nil); e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.Nil(e)
		}

		cid.AssertExpectations(t)
	}
}

func (cs *ChannelSuite) TestPublish(){
	var table = []struct{
		setExpectations func(*m.ChannelInterfaceDefault)
		e error		
	}{
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var e = fmt.Errorf("Internal Error")
				cid.On("Publish", "queue","",false, false, mock.AnythingOfType("amqp.Publishing")).Return(e)
			},
			e: fmt.Errorf("Internal Error"),
		},
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				cid.On("Publish", "queue","",false, false, mock.AnythingOfType("amqp.Publishing")).Return(nil)
			},
			e: nil,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var key= ""
		var exchangeName = "queue"
		var msg = amqp.Publishing{}
		var cid = &m.ChannelInterfaceDefault{}
		var ch = r.Channel{C: cid}

		row.setExpectations(cid)

		if e := ch.Publish(exchangeName, key,false,false, msg); e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.Nil(e)
		}

		cid.AssertExpectations(t)
	}
}

func (cs *ChannelSuite) TestConsume(){
	var sampleDeliveryChan = make(<- chan amqp.Delivery)
	var table = []struct{
		setExpectations func(*m.ChannelInterfaceDefault)
		e error		
	}{
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				var e = fmt.Errorf("Internal Error")
				cid.On("Consume", "queue","queue",false, false,false, false, mock.AnythingOfType("amqp.Table")).Return(nil, e)
			},
			e: fmt.Errorf("Internal Error"),
		},
		{
			setExpectations: func(cid *m.ChannelInterfaceDefault){
				cid.On("Consume", "queue","queue",false, false,false, false, mock.AnythingOfType("amqp.Table")).Return(sampleDeliveryChan,nil)
			},
			e: nil,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var consumer = "queue"
		var queueName = "queue"
		var cid = &m.ChannelInterfaceDefault{}
		var ch = r.Channel{C: cid}

		row.setExpectations(cid)

		if resultDeliveryChan, e := ch.Consume(queueName, consumer,false,false,false,false, nil); e != nil {
			assert.EqualError(e, row.e.Error())
		} else {
			assert.Equal(resultDeliveryChan,sampleDeliveryChan)
		}

		cid.AssertExpectations(t)
	}
}

func TestChannelSuite(t *testing.T){
	suite.Run(t, new(ChannelSuite))
}