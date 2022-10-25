package main

import (
	"attendee-writer/attendee"
	"attendee-writer/handler"
	"attendee-writer/messages"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
	"time"
)

const (
	awsRegion = "us-east-1"
)

type Clock struct{}

func (Clock) Now() time.Time { return time.Now() }

func main() {
	awsConfig, err := newAwsConfig(awsRegion)
	if err != nil {
		panic(err)
	}

	lambda.Start(newDefaultHandler(awsConfig).HandleRequest)
}

func newDefaultHandler(awsConfig *aws.Config) handler.Handler {
	return handler.Handler{
		MessageProcessor: messages.MessageProcessor{
			AttendeesStore: attendee.AttendeesStore{
				Db:    dynamodb.NewFromConfig(*awsConfig),
				Table: os.Getenv("ATTENDEES_TABLE_NAME"),
			},
			Clock: Clock{},
		},
	}
}
