package main

import (
	"github.com/aws/aws-lambda-go/lambda"
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

	lambdaHandler := Handler{
		messageProcessor: MessageProcessor{
			attendeesStore: AttendeesStore{
				Db:    dynamodb.NewFromConfig(*awsConfig),
				Table: os.Getenv("ATTENDEES_TABLE_NAME"),
			},
			clock: Clock{},
		},
	}

	lambda.Start(lambdaHandler.handleRequest)
}
