package main

import (
	"clams/attendee"
	"clams/attendees-api/handler"
	"clams/awscfg"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	awsConfig := awscfg.GetAwsConfig(dynamodb.ServiceID, os.Getenv("DYNAMO_ENDPOINT_OVERRIDE"))

	lambdaHandler := handler.Handler{
		AttendeesStore: &attendee.AttendeesStore{
			Db:    dynamodb.NewFromConfig(*awsConfig),
			Table: os.Getenv("ATTENDEES_TABLE_NAME"),
		},
	}
	lambda.Start(lambdaHandler.HandleRequest)
}
