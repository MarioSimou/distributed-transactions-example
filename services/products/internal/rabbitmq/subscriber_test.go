package rabbitmq

import (
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
		expectedRes SubResponse
	}{
		{
			handlerFunc: func(Message) error {
				return nil
			},
			expectedRes: SubResponse{Ok: true},
		},
		{
			handlerFunc: func(Message) error {
				return fmt.Errorf("Some Error")
			},
			expectedRes: SubResponse{Ok: false, Err: fmt.Errorf("Some Error")},
		},
	}

	var t = ss.T()
	var assert = assert.New(t)
	for _, row := range table {
		var localDeliveryChan = make(chan Message)
		var resChan = make(chan SubResponse)

		go subWorker(localDeliveryChan, resChan, row.handlerFunc)
		localDeliveryChan <- Message{ContentType: "application/json"}

		var res = <- resChan

		assert.EqualValues(res, row.expectedRes)
	}
}

func TestSubscriberSuite(t *testing.T){
	suite.Run(t, new(SubscriberSuite))
}