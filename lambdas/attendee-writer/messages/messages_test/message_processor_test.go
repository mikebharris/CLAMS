package messages_test

import (
	"clams/attendee"
	"clams/attendee-writer/messages"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_processMessage_ShouldStoreMessage(t *testing.T) {
	// Given
	boxForSpyToPutStoredAttendeesIn := &[]attendee.Attendee{}

	mp := messages.MessageProcessor{AttendeesStore: spyingAttendeesStore{boxForSpyToPutStoredAttendeesIn}}

	// When
	_ = mp.ProcessMessage(aMessage())

	// Then
	assertValidAttendeeReturned(t, anAttendee(), *boxForSpyToPutStoredAttendeesIn)
}

func assertValidAttendeeReturned(t *testing.T, a attendee.Attendee, b []attendee.Attendee) {
	assert.Equal(t, a.ArrivalDay, b[0].ArrivalDay)
	assert.Equal(t, a.StayingLate, b[0].StayingLate)
	assert.Equal(t, a.Name, b[0].Name)
	assert.Equal(t, a.NumberOfNights, b[0].NumberOfNights)
	assert.Equal(t, a.NumberOfKids, b[0].NumberOfKids)
	assert.Equal(t, a.AuthCode, b[0].AuthCode)
	assert.Equal(t, a.Email, b[0].Email)
	assert.Equal(t, a.Diet, b[0].Diet)
	assert.Equal(t, a.Telephone, b[0].Telephone)
	assert.Equal(t, a.Financials, b[0].Financials)
	assert.NotNil(t, b[0].CreatedTime)
}

func Test_processMessage_ShouldReturnErrorIfUnableToStoreMessage(t *testing.T) {
	// Given
	mp := messages.MessageProcessor{AttendeesStore: failingAttendeeStore{}}

	// When
	err := mp.ProcessMessage(messages.Message{})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("storing attendee in datastore: some storage error"), err)
}

func aMessage() messages.Message {
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
	return message
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
