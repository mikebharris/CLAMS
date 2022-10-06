package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"log"
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

type handler struct {
	attendees        AttendeesStore
	messageProcessor IMessageProcessor
}

func (h *handler) handleRequest(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	if len(sqsEvent.Records) == 0 {
		return events.SQSEventResponse{}, errors.New("sqs event contained no records")
	}

	var batchItemFailures []events.SQSBatchItemFailure
	for _, message := range sqsEvent.Records {
		log.Printf("processing a message with id %s for event source %s\n", message.MessageId, message.EventSource)
		if err := h.messageProcessor.processMessage(ctx, message); err != nil {
			log.Printf("handling error: %v", err)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: message.MessageId})
		}
	}

	return events.SQSEventResponse{BatchItemFailures: batchItemFailures}, nil
}
