package handler

import (
	"attendee-writer/messages"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/CLAMS/attendee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockAttendeesStore struct {
	mock.Mock
}

func (s *MockAttendeesStore) Store(ctx context.Context, attendee attendee.Attendee) error {
	args := s.Called(ctx, attendee)
	return args.Error(0)
}

type MockClock struct {
	mock.Mock
}

func (m MockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func Test_ShouldProcessMessagesPuttingFailuresOnInBatchItemFailures(t *testing.T) {
	// Given
	ctx := context.Background()
	mockAttendeesStore := MockAttendeesStore{}
	clock := MockClock{}

	h := Handler{
		MessageProcessor: messages.MessageProcessor{
			AttendeesStore: &mockAttendeesStore,
			Clock:          &clock,
		},
	}

	// mocking the clock - this is only so we can test the mock calls to Store below!
	now := time.Now()
	clock.On("Now").Return(now)

	// mocking the adaptor
	anAttendee := anAttendee(now)
	anotherAttendee := anotherAttendee(now)
	mockAttendeesStore.On("Store", ctx, anAttendee).Return(nil)
	mockAttendeesStore.On("Store", ctx, anotherAttendee).Return(fmt.Errorf("some error"))

	// When
	request, err := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{aMessage(), anotherMessage()}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{{ItemIdentifier: anotherMessage().MessageId}}}, request)

	mockAttendeesStore.AssertCalled(t, "Store", ctx, anAttendee)
	mockAttendeesStore.AssertCalled(t, "Store", ctx, anotherAttendee)
}

func Test_handleRequest_ShouldReturnErrorIfThereSqsEventContainsNoMessages(t *testing.T) {
	// Given
	h := Handler{
		MessageProcessor: messages.MessageProcessor{},
	}

	// When
	_, err := h.HandleRequest(context.Background(), events.SQSEvent{Records: []events.SQSMessage{}})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("sqs event contained no records"), err)
}

func anAttendee(now time.Time) attendee.Attendee {
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
		CreatedTime:    now,
	}
}

func anotherAttendee(now time.Time) attendee.Attendee {
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
		CreatedTime:    now,
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
