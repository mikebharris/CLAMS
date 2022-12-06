package main

import (
	"attendee-writer/messages"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
)

type IMessageProcessor interface {
	ProcessMessage(message messages.Message) error
}

type Handler struct {
	MessageProcessor IMessageProcessor
}

func (h Handler) HandleRequest(_ context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	if len(sqsEvent.Records) == 0 {
		return events.SQSEventResponse{}, errors.New("sqs event contained no records")
	}

	var batchItemFailures []events.SQSBatchItemFailure
	for _, record := range sqsEvent.Records {
		message := messages.Message{}
		if err := json.Unmarshal([]byte(record.Body), &message); err != nil {
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		} else {
			message.MessageId = record.MessageId
			if err := h.MessageProcessor.ProcessMessage(message); err != nil {
				batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			}
		}
	}

	return events.SQSEventResponse{BatchItemFailures: batchItemFailures}, nil
}
