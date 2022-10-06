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

type clock struct{}

func (clock) Now() time.Time { return time.Now() }

func main() {
	awsConfig, err := newAwsConfig()
	if err != nil {
		panic(err)
	}

	lambdaHandler := handler{
		messageProcessor: messageProcessor{
			attendeesStore: attendeesStore{
				Db:    dynamodb.NewFromConfig(*awsConfig),
				Table: os.Getenv("ATTENDEES_TABLE_NAME"),
			},
			clock: clock{},
		},
	}

	lambda.Start(lambdaHandler.handleRequest)
}
