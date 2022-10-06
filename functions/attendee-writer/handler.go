package main

import (
	"attendee-writer/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Message struct {
	AuthCode     string
	Name         string
	Email        string
	AmountToPay  int
	AmountPaid   int
	DatePaid     string
	Telephone    string
	ArrivalDay   string
	StayingLate  string
	NumberOfKids int
	Diet         string
}

type handler struct {
	attendees storage.Attendees
}

func (h *handler) handleRequest(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	if len(sqsEvent.Records) == 0 {
		return events.SQSEventResponse{}, errors.New("sqs event contained no records")
	}

	var batchItemFailures []events.SQSBatchItemFailure
	for _, message := range sqsEvent.Records {
		log.Printf("processing a message with id %s for event source %s\n", message.MessageId, message.EventSource)
		if err := h.handleMessage(ctx, message); err != nil {
			log.Printf("handling error: %v", err)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: message.MessageId})
		}
	}

	return events.SQSEventResponse{BatchItemFailures: batchItemFailures}, nil
}

func (h *handler) handleMessage(ctx context.Context, message events.SQSMessage) error {
	m, err := h.jsonToMessageObject(message)
	if err != nil {
		return fmt.Errorf("reading m message %v: %v", message, err)
	}

	attendee := storage.Attendee{
		AuthCode:     m.AuthCode,
		Name:         m.Name,
		Email:        m.Email,
		Telephone:    m.Telephone,
		NumberOfKids: m.NumberOfKids,
		Diet:         m.Diet,
		Financials: storage.Financials{
			AmountToPay: m.AmountToPay,
			AmountPaid:  m.AmountPaid,
			AmountDue:   m.AmountToPay - m.AmountPaid,
			DatePaid:    m.DatePaid,
		},
		ArrivalDay:     m.ArrivalDay,
		NumberOfNights: h.computeNights(m.ArrivalDay, m.StayingLate),
		StayingLate:    m.StayingLate,
		CreatedTime:    time.Now(),
	}

	if err := h.storeAttendee(ctx, attendee); err != nil {
		return fmt.Errorf("storing result: %v", err)
	}

	return nil
}

func (h *handler) storeAttendee(ctx context.Context, attendee storage.Attendee) error {
	if err := h.attendees.Store(ctx, attendee); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}
	return nil
}

func (h *handler) jsonToMessageObject(message events.SQSMessage) (*Message, error) {
	r := Message{}
	if err := json.Unmarshal([]byte(message.Body), &r); err != nil {
		return nil, fmt.Errorf("unmarshalling message body %s: %v", message.Body, err)
	}
	return &r, nil
}

func (h *handler) computeNights(arrival string, stayingLate string) int {
	var nights = 1

	if strings.Contains(arrival, "Wednesday") {
		nights = 4
	} else if strings.Contains(arrival, "Thursday") {
		nights = 3
	} else if strings.Contains(arrival, "Friday") {
		nights = 2
	} else if strings.Contains(arrival, "Saturday") {
		nights = 1
	} else {
		nights = 5
	}

	if stayingLate == "Yes" {
		nights += 1
	}

	return nights
}
