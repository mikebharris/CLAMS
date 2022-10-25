package main

import (
	"attendees-api/attendee"
	"attendees-api/handler"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	awsConfig, err := newAwsConfig("us-east-1")
	if err != nil {
		panic(err)
	}

	lambdaHandler := handler.Handler{
		AttendeesStore: &attendee.AttendeesStore{
			Db:    dynamodb.NewFromConfig(*awsConfig),
			Table: os.Getenv("ATTENDEES_TABLE_NAME"),
		},
	}
	lambda.Start(lambdaHandler.HandleRequest)
}
