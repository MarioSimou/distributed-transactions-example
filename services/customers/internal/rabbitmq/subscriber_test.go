package rabbitmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)
type SubscriberSuite struct {
	suite.Suite
}

func (ss *SubscriberSuite) TestSubWorker(){	
	var table = []struct{
		handlerFunc HandlerFunc
		expectedRes SubscriptionResponse
	}{
		{
			handlerFunc: func(Message,context.Context) error {
				return nil
			},
			expectedRes: SubscriptionResponse{Ok: true, WorkerID: "worker1",RequestID: "id"},
		},
		{
			handlerFunc: func(Message, context.Context) error {
				return fmt.Errorf("Some Error")
			},
			expectedRes: SubscriptionResponse{Ok: false, Err: fmt.Errorf("Some Error"), WorkerID: "worker1",RequestID: "id"},
		},
	}

	var t = ss.T()
	var assert = assert.New(t)
	for _, row := range table {
		var localDeliveryChan = make(chan Message)
		var resChan = make(chan SubscriptionResponse)
		var ctx = context.Background()

		go subWorker(localDeliveryChan, resChan, row.handlerFunc, ctx, "worker1")
		localDeliveryChan <- Message{CorrelationId: "id",ContentType: "application/json"}

		var res = <- resChan

		assert.EqualValues(res, row.expectedRes)
	}
}

func TestSubscriberSuite(t *testing.T){
	suite.Run(t, new(SubscriberSuite))
}