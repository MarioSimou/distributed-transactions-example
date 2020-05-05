package rabbitmq

import "github.com/streadway/amqp"

var (
	N_WORKERS = 10
)

type Subscription interface {
	GetQueue() string
	GetHandlerFunc() HandlerFunc
}

type Message amqp.Delivery
type HandlerFunc func(Message) error

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

func NewSubscription(subscribers Subscribers, conn ConnectionInterface)  (chan SubResponse,error) {
	var subResChan = make(chan SubResponse, 1000)

	for _, subscriber := range subscribers {
		if e := initSubscription(subscriber, subResChan, conn); e != nil {
			return nil, e
		}
	}
	
	return subResChan, nil
}

func subWorker(localDeliveryChan chan Message, resChan chan SubResponse, handlerFunc HandlerFunc){
	var message = <- localDeliveryChan
	if e := handlerFunc(message); e != nil {
		resChan <- SubResponse{
			Ok: false,
			Err: e,
		}
	} else {
		resChan <- SubResponse{Ok: true}
	}
}

func subManager(deliveryChan <-chan amqp.Delivery, localDeliveryChan chan Message){
	for delivery := range deliveryChan {
		localDeliveryChan <- Message(delivery)	
	}
}

func initSubscription(subscription Subscription, subResChan chan SubResponse, conn ConnectionInterface) error  {
	var queueName = subscription.GetQueue()
	var channel ChannelInterface
	var deliveryChan <-chan amqp.Delivery
	var e error
	var localDeliveryChan = make(chan Message, N_WORKERS)
	var handlerFunc = subscription.GetHandlerFunc()

	if channel, e = conn.Channel(queueName); e != nil {
		return e
	}
	if deliveryChan, e = channel.Consume(queueName, "", true, false, false, false, nil); e != nil {
		return e
	}

	// create n workers that will consume the queue's messages
	for i:=0; i < N_WORKERS; i++ {
		go subWorker(localDeliveryChan, subResChan, handlerFunc)
	}
	go subManager(deliveryChan, localDeliveryChan)

	return nil
}