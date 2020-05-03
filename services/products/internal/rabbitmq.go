package internal

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

var (
	FANOUT = "fanout"
)

type Publishing interface {
	GetChannels() []*Channel
	GetConnection() *amqp.Connection
	Pub(queueName string, body interface{}) error
}

type Subscribing interface {
	GetChannels() []*Channel
	GetConnection() *amqp.Connection
}

type Channel struct {
	QueueName string
	Q amqp.Queue
	C *amqp.Channel
}

func (ch *Channel) Close() error {
	return ch.C.Close()
}

type Publisher struct {
	ConnString string
	Connection *amqp.Connection
	Channels []*Channel
}

func (pub *Publisher) GetChannels() []*Channel {
	return pub.Channels
}
func (pub *Publisher) GetConnection() *amqp.Connection {
	return pub.Connection	
}

func (pub *Publisher) Pub(queueName string, body interface{}) error {
	var channel *Channel
	var e error
	var bf []byte
	if channel = pub.getChannel(queueName); channel == nil {
		return fmt.Errorf("No Channel found")
	}
	if bf, e = json.Marshal(body); e != nil {
		return e
	}

	var message = amqp.Publishing{
		ContentType: "application/json",
		Body: bf,
	}
	if e = channel.C.Publish(queueName, "", false, false, message); e != nil {
		return e
	}
	return nil
}

func (pub *Publisher) getChannel(queueName string) *Channel {
	for _, channel := range pub.Channels {
		if channel.QueueName == queueName {
			return channel
		}
	}
	return nil
}

func NewPublisher(connString string, queuesNames []string) (Publishing, error) {
	var conn *amqp.Connection
	var e error
	
	if conn, e = amqp.Dial(connString); e != nil {
		return nil,e
	}

	var channels []*Channel
	for _, queueName := range queuesNames {
		var channel *Channel
		if channel, e = getChannelQueue(queueName, conn);  e != nil {
			if e := closeChannels(channels); e != nil {
				return nil, e 
			}
			return nil, e
		}

		channels = append(channels, channel)
	}

	return &Publisher{
		ConnString: connString,
		Connection: conn,
		Channels: channels,
	}, nil
}

type Message amqp.Delivery
type HandlerFunc func(Message) error

type Subscription interface {
	GetQueue() string
	GetHandlerFunc() HandlerFunc
}

type SubResponse struct {
	Ok bool
	Err error
}

type Subscriber struct {
	QueueName string
	HandlerFunc HandlerFunc
}
type Subscribers []Subscriber

func (sub Subscriber) GetQueue() string {
	return sub.QueueName
}

func (sub Subscriber) GetHandlerFunc() HandlerFunc {
	return sub.HandlerFunc
}

func NewSubscription(connString string, subscribers Subscribers)  (chan SubResponse,error) {
	var conn *amqp.Connection
	var e error
	var subResChan = make(chan SubResponse, 1000)

	if conn, e = amqp.Dial(connString); e != nil {
		return nil, e
	}
	for _, subscriber := range subscribers {
		if e = initSubscription(subscriber, subResChan, conn); e != nil {
			return nil, e
		}
	}
	
	return subResChan, nil
}

func initSubscription(subscription Subscription, subResChan chan SubResponse, conn *amqp.Connection) error  {
	var queueName = subscription.GetQueue()
	var channel *amqp.Channel
	var deliveryChan <-chan amqp.Delivery
	var e error
	var nWorkers = 10
	var localDeliveryChan = make(chan Message, nWorkers)
	var handlerFunc = subscription.GetHandlerFunc()

	worker := func(localDeliveryChan chan Message){
		var message = <- localDeliveryChan
		if e := handlerFunc(message); e != nil {
			subResChan <- SubResponse{
				Ok: false,
				Err: e,
			}
		} else {
			subResChan <- SubResponse{Ok: true}
		}
	}
	manager := func(deliveryChan <-chan amqp.Delivery){
		for delivery := range deliveryChan {
			localDeliveryChan <- Message(delivery)	
		}
	}

	if channel, e = conn.Channel(); e != nil {
		return e
	}
	if deliveryChan, e = channel.Consume(queueName, "", true, false, false, false, nil); e != nil {
		return e
	}

	// create n workers that will consume the queue's messages
	for i:=0; i < nWorkers; i++ {
		go worker(localDeliveryChan)
	}
	go manager(deliveryChan)

	return nil
}

func closeChannels(channels []*Channel) error {
	for _, channel := range channels {
		if e := channel.Close(); e != nil {
			return e
		}
	}
	return nil
}

func getChannelQueue(queueName string, conn *amqp.Connection) (*Channel, error) {
	var channel *amqp.Channel
	var queue amqp.Queue
	var e error

	if channel, e = conn.Channel(); e != nil {
		return nil, e
	}
	if queue, e = channel.QueueDeclare(queueName, true, false, false, false, nil); e != nil {
		return nil, e
	}
	if e = channel.ExchangeDeclare(queueName, FANOUT, true, false, false, false, nil); e != nil {
		return nil, e
	}
	if e = channel.QueueBind(queueName, "", queueName, false, nil); e != nil {
		return nil, e
	}

	return &Channel{
		QueueName: queueName,
		Q: queue,
		C: channel,
	}, nil
}