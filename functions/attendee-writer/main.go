package main

import (
	"attendee-writer/attendee"
	"attendee-writer/messages"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
	"sync"
)

const (
	awsRegion = "us-east-1"
)

func main() {
	awsConfig, err := newAwsConfig(awsRegion)
	if err != nil {
		panic(err)
	}

	lambda.Start(newDefaultHandler(awsConfig).HandleRequest)
}

type IMessageProcessor interface {
	ProcessMessage(msg events.SQSMessage) error
}

type Handler struct {
	MessageProcessor IMessageProcessor
}

func (h Handler) HandleRequest(_ context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
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
			if err := h.MessageProcessor.ProcessMessage(sqsMessage); err != nil {
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

func newDefaultHandler(awsConfig *aws.Config) Handler {
	return Handler{
		MessageProcessor: messages.MessageProcessor{
			AttendeesStore: &attendee.AttendeesStore{
				Db:    dynamodb.NewFromConfig(*awsConfig),
				Table: os.Getenv("ATTENDEES_TABLE_NAME"),
			},
		},
	}
}
