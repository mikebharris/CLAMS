package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"registrar/storage"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type IncomingRequestPayload struct {
	Name        string
	Email       string
	Code        string
	ToPay       uint
	Paid        uint
	PaidDate    string
	Phone       string
	Arrival     string
	Diet        string
	StayingLate string
	Kids        uint
}

type handler struct {
	attendees storage.IAttendees
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
	payload, err := h.unmarshalToIncomingRequest(message)
	if err != nil {
		return fmt.Errorf("reading payload message %v: %v", message, err)
	}

	attendee := storage.Attendee{
		Code:  payload.Code,
		Name:  payload.Name,
		Email: payload.Email,
		Phone: payload.Phone,
		Kids:  payload.Kids,
		Diet:  payload.Diet,
		Financials: storage.Financials{
			ToPay:    payload.ToPay,
			Paid:     payload.Paid,
			Due:      int(payload.ToPay - payload.Paid),
			PaidDate: payload.PaidDate,
		},
		Arrival:     payload.Arrival,
		Nights:      h.computeNights(payload.Arrival, payload.StayingLate),
		StayingLate: payload.StayingLate,
		CreatedTime: time.Now(),
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

func (h *handler) unmarshalToIncomingRequest(message events.SQSMessage) (*IncomingRequestPayload, error) {
	r := IncomingRequestPayload{}
	if err := json.Unmarshal([]byte(message.Body), &r); err != nil {
		return nil, fmt.Errorf("unmarshalling payload body %s: %v", message.Body, err)
	}
	return &r, nil
}

func (h *handler) computeNights(arrival string, stayingLate string) uint {
	var nights uint = 1

	switch arrival {
	case "Wednesday":
		nights = 4
	case "Thursday":
		nights = 3
	case "Friday":
		nights = 2
	}

	if stayingLate == "Yes" {
		nights += 1
	}

	return nights
}
