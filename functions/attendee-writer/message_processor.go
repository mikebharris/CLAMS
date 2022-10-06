package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"strings"
	"time"
)

type IAttendeesStore interface {
	Store(ctx context.Context, attendee Attendee) error
}

type IClock interface {
	Now() time.Time
}

type MessageProcessor struct {
	attendeesStore IAttendeesStore
	clock          IClock
}

func (mp MessageProcessor) processMessage(ctx context.Context, message events.SQSMessage) error {
	msg, err := mp.jsonToMessageObject(message)
	if err != nil {
		return fmt.Errorf("reading message %v: %v", message, err)
	}

	attendee := Attendee{
		AuthCode:     msg.AuthCode,
		Name:         msg.Name,
		Email:        msg.Email,
		Telephone:    msg.Telephone,
		NumberOfKids: msg.NumberOfKids,
		Diet:         msg.Diet,
		Financials: Financials{
			AmountToPay: msg.AmountToPay,
			AmountPaid:  msg.AmountPaid,
			AmountDue:   msg.AmountToPay - msg.AmountPaid,
			DatePaid:    msg.DatePaid,
		},
		ArrivalDay:     msg.ArrivalDay,
		NumberOfNights: mp.computeNights(msg.ArrivalDay, msg.StayingLate),
		StayingLate:    msg.StayingLate,
		CreatedTime:    mp.clock.Now(),
	}

	if err := mp.attendeesStore.Store(ctx, attendee); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}

func (mp MessageProcessor) jsonToMessageObject(message events.SQSMessage) (*Message, error) {
	r := Message{}
	if err := json.Unmarshal([]byte(message.Body), &r); err != nil {
		return nil, fmt.Errorf("unmarshalling message body %s: %v", message.Body, err)
	}
	return &r, nil
}

func (mp MessageProcessor) computeNights(arrival string, stayingLate string) int {
	var nights int

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
