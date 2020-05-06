package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

var (
	N_WORKERS = 10
)

type Subscriber interface {
	GetQueue() string
	GetHandlerFunc() HandlerFunc
	GetContext() context.Context
}

type Subscribers []Subscriber

type Message amqp.Delivery
type HandlerFunc func(Message, context.Context) error

type SubscriptionResponse struct {
	Ok bool
	Err error
	RequestID string
	WorkerID string
}

type subscriber struct {
	queueName string
	handlerFunc HandlerFunc
	ctx context.Context
}

func (sub subscriber) GetQueue() string {
	return sub.queueName
}

func (sub subscriber) GetHandlerFunc() HandlerFunc {
	return sub.handlerFunc
}

func (sub subscriber) GetContext() context.Context {
	return sub.ctx
}

func NewSubscriber(queueName string, handlerFunc HandlerFunc, ctx context.Context) Subscriber {

	return subscriber{
		queueName: queueName,
		handlerFunc: handlerFunc,
		ctx: ctx,
	}
}

func NewSubscribers(subscribers... Subscriber) Subscribers{
	var s = Subscribers{}
	for _, subscriber := range subscribers {
		s = append(s, subscriber)
	}
	return s
}


func NewSubscription(subscribers Subscribers, conn ConnectionInterface)  (chan SubscriptionResponse,error) {
	var subResChan = make(chan SubscriptionResponse, 1000)

	for _, subscriber := range subscribers {
		if e := initSubscription(subscriber, subResChan, conn); e != nil {
			return nil, e
		}
	}
	
	return subResChan, nil
}

func subWorker(localDeliveryChan chan Message, resChan chan SubscriptionResponse, handlerFunc HandlerFunc, ctx context.Context, workerID string){
	var message = <- localDeliveryChan
	if e := handlerFunc(message, ctx); e != nil {
		resChan <- SubscriptionResponse{
			Ok: false,
			Err: e,
			RequestID: message.CorrelationId,
			WorkerID: workerID,
		}
	} else {
		resChan <- SubscriptionResponse{
			Ok: true,
			RequestID: message.CorrelationId,
			WorkerID: workerID,
		}
	}
}

func subManager(deliveryChan <-chan amqp.Delivery, localDeliveryChan chan Message){
	for delivery := range deliveryChan {
		localDeliveryChan <- Message(delivery)	
	}
}

func initSubscription(subscriber Subscriber, subResChan chan SubscriptionResponse, conn ConnectionInterface) error  {
	var queueName = subscriber.GetQueue()
	var channel ChannelInterface
	var deliveryChan <-chan amqp.Delivery
	var e error
	var localDeliveryChan = make(chan Message, N_WORKERS)
	var handlerFunc = subscriber.GetHandlerFunc()
	var context = subscriber.GetContext()

	if channel, e = conn.Channel(queueName); e != nil {
		return e
	}
	if deliveryChan, e = channel.Consume(queueName, "", true, false, false, false, nil); e != nil {
		return e
	}

	// create n workers that will consume the queue's messages
	for i:=0; i < N_WORKERS; i++ {
		var workerID = fmt.Sprintf("%s - %d", queueName, i)
		go subWorker(localDeliveryChan, subResChan, handlerFunc, context, workerID)
	}
	go subManager(deliveryChan, localDeliveryChan)

	return nil
}