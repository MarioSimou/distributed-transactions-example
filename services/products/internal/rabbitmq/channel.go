package rabbitmq

import (
	"github.com/streadway/amqp"
)

type ChannelInterfaceDefault interface {
	Close() error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type ChannelInterface interface {
	Close() error
	QueueDeclare(string, bool, bool, bool, bool, amqp.Table) error
	ExchangeDeclare(string, string, bool, bool, bool, bool, amqp.Table) error 
	QueueBind(string, string, string, bool,amqp.Table) error
	GetQueueName() string
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(string, string, bool, bool, bool, bool, amqp.Table) (<-chan amqp.Delivery, error)
}

type Channel struct {
	QueueName string
	C ChannelInterfaceDefault
}

func (ch *Channel) GetQueueName() string {
	return ch.QueueName
}

func (ch *Channel) Close() error {
	return ch.C.Close()
}

func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) error {
	if _, e := ch.C.QueueDeclare(name, durable,autoDelete, exclusive,noWait, args); e != nil {
		return e
	}
	return nil
}

func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	if e := ch.C.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args); e != nil {
		return e
	}
	return nil
}

func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	if e := ch.C.QueueBind(name, key, exchange, noWait, args); e != nil {
		return e
	}
	return nil
}

func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return ch.C.Publish(exchange, key, mandatory, immediate, msg)
}

func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	var deliveryChan <-chan amqp.Delivery
	var e error
	if deliveryChan, e = ch.C.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args); e != nil {
		return nil, e
	}
	return deliveryChan, nil
}