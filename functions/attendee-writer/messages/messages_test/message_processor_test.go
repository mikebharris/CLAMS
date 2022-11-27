package messages_test

import (
	"attendee-writer/attendee"
	"attendee-writer/messages"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_processMessage_ShouldStoreMessage(t *testing.T) {
	// Given
	boxForSpyToPutStoredAttendeesIn := &[]attendee.Attendee{}

	mp := messages.MessageProcessor{AttendeesStore: spyingAttendeesStore{boxForSpyToPutStoredAttendeesIn}}

	// When
	mp.ProcessMessage(aMessage())

	// Then
	assert.Equal(t, []attendee.Attendee{anAttendee()}, *boxForSpyToPutStoredAttendeesIn)
}

func Test_processMessage_ShouldReturnErrorIfUnableToStoreMessage(t *testing.T) {
	// Given
	mp := messages.MessageProcessor{AttendeesStore: failingAttendeeStore{}}

	// When
	body, _ := json.Marshal(messages.Message{})
	err := mp.ProcessMessage(events.SQSMessage{Body: string(body)})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("storing attendee in datastore: some storage error"), err)
}

func Test_processMessage_ShouldReturnErrorIfUnableToParseMessage(t *testing.T) {
	// Given
	mp := messages.MessageProcessor{AttendeesStore: &attendee.AttendeesStore{}}

	// When
	err := mp.ProcessMessage(events.SQSMessage{Body: ""})

	// Then
	assert.NotNil(t, err)
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
