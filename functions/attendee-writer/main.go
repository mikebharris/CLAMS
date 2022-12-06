package main

import (
	"attendee-writer/dynds"
	"attendee-writer/messages"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main() {
	ds := dynds.DynamoDatastore{
		Table:    os.Getenv("ATTENDEES_TABLE_NAME"),
		Endpoint: os.Getenv("DYNAMO_ENDPOINT_OVERRIDE"),
	}
	ds.Init()

	lambda.Start(Handler{
		MessageProcessor: messages.MessageProcessor{
			AttendeesStore: &ds,
		},
	}.HandleRequest)
}
