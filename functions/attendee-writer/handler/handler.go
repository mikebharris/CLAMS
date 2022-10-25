package handler

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"sync"
)

type IMessageProcessor interface {
	ProcessMessage(ctx context.Context, message events.SQSMessage) error
}

type Handler struct {
	MessageProcessor IMessageProcessor
}

func (h Handler) HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	if len(sqsEvent.Records) == 0 {
		return events.SQSEventResponse{}, errors.New("sqs event contained no records")
	}

	var wg sync.WaitGroup
	failedEventChan := make(chan events.SQSMessage)

	for _, record := range sqsEvent.Records {
		wg.Add(1)
		sqsMessage := record
		go func() {
			defer wg.Done()
			if err := h.MessageProcessor.ProcessMessage(ctx, sqsMessage); err != nil {
				failedEventChan <- sqsMessage
			}
		}()
	}

	go func() {
		wg.Wait()
		close(failedEventChan)
	}()

	var batchItemFailures []events.SQSBatchItemFailure
	for i := range failedEventChan {
		batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: i.MessageId})
	}

	return events.SQSEventResponse{BatchItemFailures: batchItemFailures}, nil
}
