package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type AttendeeFactory struct {
	clock IClock
}

func (af AttendeeFactory) NewFromMessage(message events.SQSMessage) (Attendee, error) {
	msg, err := af.jsonToMessageObject(message)
	if err != nil {
		return Attendee{}, fmt.Errorf("reading message %v: %v", message, err)
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
		NumberOfNights: af.computeNights(msg.ArrivalDay, msg.StayingLate),
		StayingLate:    msg.StayingLate,
		CreatedTime:    af.clock.Now(),
	}
	return attendee, nil
}

func (af AttendeeFactory) jsonToMessageObject(message events.SQSMessage) (*Message, error) {
	r := Message{}
	if err := json.Unmarshal([]byte(message.Body), &r); err != nil {
		return nil, fmt.Errorf("unmarshalling message body %s: %v", message.Body, err)
	}
	return &r, nil
}

func (af AttendeeFactory) computeNights(arrival string, stayingLate string) int {
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
