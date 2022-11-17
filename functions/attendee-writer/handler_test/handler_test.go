package handler_test

import (
	"attendee-writer/handler"
	"attendee-writer/messages"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/CLAMS/attendee"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SpyingAttendeesStore struct {
	attendees *[]attendee.Attendee
}

func (s SpyingAttendeesStore) Store(attendee attendee.Attendee) error {
	if attendee == anotherAttendee() {
		return errors.New("some error")
	}
	*s.attendees = append(*s.attendees, attendee)
	return nil
}

func Test_ShouldProcessMessagesPuttingFailuresOnInBatchItemFailures(t *testing.T) {
	// Given
	ctx := context.Background()

	var attendees []attendee.Attendee
	h := handler.Handler{
		MessageProcessor: messages.MessageProcessor{
			AttendeesStore: &SpyingAttendeesStore{&attendees},
		},
	}

	// When
	request, err := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{aMessage(), anotherMessage()}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, anAttendee(), attendees[0])
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{{ItemIdentifier: anotherMessage().MessageId}}}, request)
}

func Test_handleRequest_ShouldReturnErrorIfThereSqsEventContainsNoMessages(t *testing.T) {
	// Given
	h := handler.Handler{
		MessageProcessor: messages.MessageProcessor{},
	}

	// When
	_, err := h.HandleRequest(context.Background(), events.SQSEvent{Records: []events.SQSMessage{}})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("sqs event contained no records"), err)
}

func anAttendee() attendee.Attendee {
	return attendee.Attendee{
		AuthCode:     "123456",
		Name:         "Frank Ostrowski",
		Email:        "frank.o@gfa.de",
		Telephone:    "123456789",
		NumberOfKids: 1,
		Diet:         "I eat BASIC code for lunch",
		Financials: attendee.Financials{
			DatePaid:    "29/05/2022",
			AmountPaid:  75,
			AmountToPay: 75,
		},
		ArrivalDay:     "Wednesday",
		NumberOfNights: 4,
		StayingLate:    "No",
	}
}

func anotherAttendee() attendee.Attendee {
	return attendee.Attendee{
		AuthCode:     "0101010",
		Name:         "Grace Hopper",
		Email:        "grace@nasa.gov",
		Telephone:    "123456789",
		NumberOfKids: 1,
		Diet:         "I eat COBOL code for lunch",
		Financials: attendee.Financials{
			DatePaid:    "29/05/2022",
			AmountPaid:  75,
			AmountToPay: 75,
		},
		ArrivalDay:     "Wednesday",
		NumberOfNights: 4,
		StayingLate:    "No",
	}
}

func aMessage() events.SQSMessage {
	message := messages.Message{
		Name:         "Frank Ostrowski",
		Email:        "frank.o@gfa.de",
		AuthCode:     "123456",
		AmountToPay:  75,
		AmountPaid:   75,
		DatePaid:     "29/05/2022",
		Telephone:    "123456789",
		ArrivalDay:   "Wednesday",
		Diet:         "I eat BASIC code for lunch",
		StayingLate:  "No",
		NumberOfKids: 1,
	}
	body, _ := json.Marshal(message)
	return events.SQSMessage{MessageId: "abcdef", Body: string(body)}
}

func anotherMessage() events.SQSMessage {
	message := messages.Message{
		Name:         "Grace Hopper",
		Email:        "grace@nasa.gov",
		AuthCode:     "0101010",
		AmountToPay:  75,
		AmountPaid:   75,
		DatePaid:     "29/05/2022",
		Telephone:    "123456789",
		ArrivalDay:   "Wednesday",
		Diet:         "I eat COBOL code for lunch",
		StayingLate:  "No",
		NumberOfKids: 1,
	}
	body, _ := json.Marshal(message)
	return events.SQSMessage{MessageId: "123456", Body: string(body)}
}
