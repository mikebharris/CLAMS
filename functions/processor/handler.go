package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/net/context"
)

type IDatastore interface {
	Store(thing interface{}) error
}

type handler struct {
	dbConx    *sql.DB
	datastore IDatastore
}

type message struct {
	WorkshopSignupId int
}

type WorkshopSignupRecord struct {
	WorkshopSignupId int
	WorkshopTitle    string
	Role             string
	Name             string
}

func (h handler) handleRequest(_ context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	if len(sqsEvent.Records) == 0 {
		return events.SQSEventResponse{}, errors.New("sqs event contained no records")
	}
	repository := repository{dbConx: h.dbConx}
	var batchItemFailures []events.SQSBatchItemFailure
	for _, record := range sqsEvent.Records {
		msg := message{}
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		} else {
			signupRecord, _ := repository.getSignupRecord(msg.WorkshopSignupId)
			if signupRecord == (WorkshopSignupRecord{}) {
				batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
				continue
			}
			if err := h.datastore.Store(signupRecord); err != nil {
				batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			}
		}
	}

	return events.SQSEventResponse{BatchItemFailures: batchItemFailures}, nil
}
