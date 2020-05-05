package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	QUEUE_NOT_FOUND = fmt.Errorf("Queue not found")
	CHANNEL_NOT_FOUND = fmt.Errorf("Channel not found")
)

type ConnectionInterface interface {
	Close() error
	Channel(string) (ChannelInterface, error)
}

type ConnectionStruct struct {
	Conn *amqp.Connection
}

func (c *ConnectionStruct) Close() error {
	return c.Conn.Close()
} 

func (c *ConnectionStruct) Channel(queueName string) (ChannelInterface, error){
	var channel *amqp.Channel
	var e error
	
	if queueName == "" {
		return nil, QUEUE_NOT_FOUND
	}
	if channel, e = c.Conn.Channel(); e != nil {
		return nil, e
	}

	return &Channel{
		QueueName: queueName,
		C: channel,
	}, nil
}

func (c *ConnectionStruct) Start(connString string) error{
	var conn *amqp.Connection
	var e error

	if c.Conn != nil {
		return nil
	}
	if conn, e = amqp.Dial(connString); e != nil {
		return e
	}
	
	c.Conn = conn
	return nil
} 
