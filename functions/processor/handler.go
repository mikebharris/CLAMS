package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/net/context"
)

type handler struct {
	db IRepository
}

type message struct {
	WorkshopSignupId int
	WorkshopId       int
	PeopleId         int
	RoleId           int
}

type Person struct {
	Name  string
	Email string
	Diet  string
}

type WorkshopSignupRecord struct {
	WorkshopId       int
	WorkshopTitle    string
	FacilitatorName  string
	FacilitatorEmail string
	People           []Person
}

func (h handler) handleRequest(_ context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	if len(sqsEvent.Records) == 0 {
		return events.SQSEventResponse{}, errors.New("sqs event contained no records")
	}

	var batchItemFailures []events.SQSBatchItemFailure
	for _, record := range sqsEvent.Records {
		msg := message{}
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		} else {
			fmt.Println("I have a message: ", msg)
			signupRecord, _ := h.db.getSignupRecord(msg.WorkshopSignupId)
			fmt.Println("Someone signed up for ", signupRecord.WorkshopTitle)
		}
	}

	return events.SQSEventResponse{BatchItemFailures: batchItemFailures}, nil
}

type IRepository interface {
	getSignupRecord(signupId int) (WorkshopSignupRecord, error)
}
