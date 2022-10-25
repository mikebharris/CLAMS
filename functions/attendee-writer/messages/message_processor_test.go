package messages

import (
	"attendee-writer/attendee"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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

func Test_processMessage_ShouldStoreMessage(t *testing.T) {
	// Given
	ctx := context.Background()

	mockClock := MockClock{}
	now := time.Now()
	mockClock.On("Now").Return(now)

	mockAttendeesStore := MockAttendeesStore{}
	attendee := attendee.Attendee{
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
	mockAttendeesStore.On("Store", ctx, attendee).Return(nil)

	mp := MessageProcessor{AttendeesStore: &mockAttendeesStore, Clock: mockClock}

	message := Message{
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

	// When
	err := mp.ProcessMessage(ctx, events.SQSMessage{Body: string(body)})

	// Then
	assert.Nil(t, err)
	mockAttendeesStore.AssertCalled(t, "Store", ctx, attendee)
}

func Test_processMessage_ShouldReturnErrorIfUnableToStoreMessage(t *testing.T) {
	// Given
	ctx := context.Background()

	mockClock := MockClock{}
	now := time.Now()
	mockClock.On("Now").Return(now)

	mockAttendeesStore := MockAttendeesStore{}
	mockAttendeesStore.On("Store", ctx, mock.Anything).Return(fmt.Errorf("some storage error"))

	mp := MessageProcessor{AttendeesStore: &mockAttendeesStore, Clock: mockClock}

	message := Message{
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

	// When
	err := mp.ProcessMessage(ctx, events.SQSMessage{Body: string(body)})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("storing attendee in datastore: some storage error"), err)
	mockAttendeesStore.AssertCalled(t, "Store", ctx, mock.Anything)
}

func Test_processMessage_ShouldReturnErrorIfUnableToParseMessage(t *testing.T) {
	// Given
	ctx := context.Background()

	mockAttendeesStore := MockAttendeesStore{}
	mockClock := MockClock{}

	mp := MessageProcessor{AttendeesStore: &mockAttendeesStore, Clock: &mockClock}

	// When
	err := mp.ProcessMessage(ctx, events.SQSMessage{Body: ""})

	// Then
	assert.NotNil(t, err)
	assert.Regexp(t, "^reading message.*", err)
	mockAttendeesStore.AssertNotCalled(t, "Store")
	mockClock.AssertNotCalled(t, "Now")
}
