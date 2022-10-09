package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"sync"
)

type Message struct {
	AuthCode     string
	Name         string
	Email        string
	AmountToPay  int
	AmountPaid   int
	DatePaid     string
	Telephone    string
	ArrivalDay   string
	StayingLate  string
	NumberOfKids int
	Diet         string
}

type IMessageProcessor interface {
	processMessage(ctx context.Context, message events.SQSMessage) error
}

type Handler struct {
	messageProcessor IMessageProcessor
}

func (h Handler) handleRequest(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
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
			if err := h.messageProcessor.processMessage(ctx, sqsMessage); err != nil {
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
