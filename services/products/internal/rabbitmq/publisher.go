package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

var (
	FANOUT = "fanout"
)

type PublisherInterface interface {
	GetChannels() []ChannelInterface
	GetChannel(string) ChannelInterface
	GetConnection() ConnectionInterface
	Pub(queueName string, body interface{}) error
}

type Publisher struct {
	Conn ConnectionInterface
	Channels []ChannelInterface
}

func (pub *Publisher) GetChannels() []ChannelInterface {
	return pub.Channels
}

func (pub *Publisher) GetChannel(queueName string) ChannelInterface {
	for _, channel := range pub.Channels {
		if channel.GetQueueName() == queueName {
			return channel
		}
	}
	return nil
}

func (pub *Publisher) GetConnection() ConnectionInterface {
	return pub.Conn	
}

func (pub *Publisher) Pub(queueName string, body interface{}) error {
	var channel ChannelInterface
	var e error
	var bf []byte
	if channel = pub.GetChannel(queueName); channel == nil {
		return CHANNEL_NOT_FOUND
	}
	if bf, e = json.Marshal(body); e != nil {
		return e
	}

	var message = amqp.Publishing{
		ContentType: "application/json",
		Body: bf,
	}
	if e = channel.Publish(queueName, "", false, false, message); e != nil {
		return e
	}
	return nil
}

func NewPublisher(queuesNames []string, conn ConnectionInterface) (PublisherInterface, error) {
	var channels []ChannelInterface
	
	for _, queueName := range queuesNames {
		var channel ChannelInterface
		var e error
		if channel, e = getChannelQueue(queueName, conn);  e != nil {
			if e := closeChannels(channels); e != nil {
				return nil, e 
			}
			return nil, e
		}

		channels = append(channels, channel)
	}

	return &Publisher{
		Conn: conn,
		Channels: channels,
	}, nil
}

func closeChannels(channels []ChannelInterface) error {
	for _, channel := range channels {
		if e := channel.Close(); e != nil {
			return e
		}
	}
	return nil
}

func getChannelQueue(queueName string, conn ConnectionInterface) (ChannelInterface, error) {
	var channel ChannelInterface
	var e error

	if channel, e = conn.Channel(queueName); e != nil {
		return nil, e
	}
	if e := channel.QueueDeclare(queueName, true, false, false, false, nil); e != nil {
		return nil, e
	}
	if e := channel.ExchangeDeclare(queueName, FANOUT, true, false, false, false, nil); e != nil {
		return nil, e
	}
	if e := channel.QueueBind(queueName, "", queueName, false, nil); e != nil {
		return nil, e
	}

	return channel, nil
}