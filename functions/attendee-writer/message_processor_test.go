package main

import (
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

func (s *MockAttendeesStore) Store(ctx context.Context, attendee Attendee) error {
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
	attendee := Attendee{
		AuthCode:     "123456",
		Name:         "Frank Ostrowski",
		Email:        "frank.o@gfa.de",
		Telephone:    "123456789",
		NumberOfKids: 1,
		Diet:         "I eat BASIC code for lunch",
		Financials: Financials{
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

	mp := MessageProcessor{attendeesStore: &mockAttendeesStore, clock: mockClock}

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
	err := mp.processMessage(ctx, events.SQSMessage{Body: string(body)})

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

	mp := MessageProcessor{attendeesStore: &mockAttendeesStore, clock: mockClock}

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
	err := mp.processMessage(ctx, events.SQSMessage{Body: string(body)})

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

	mp := MessageProcessor{attendeesStore: &mockAttendeesStore, clock: &mockClock}

	// When
	err := mp.processMessage(ctx, events.SQSMessage{Body: ""})

	// Then
	assert.NotNil(t, err)
	assert.Regexp(t, "^reading message.*", err)
	mockAttendeesStore.AssertNotCalled(t, "Store")
	mockClock.AssertNotCalled(t, "Now")
}

func Test_handler_computeNights(t *testing.T) {
	type args struct {
		arrival     string
		stayingLate string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Should report 4 days if arriving Wednesday and not staying late",
			args: args{
				arrival:     "Wednesday AM",
				stayingLate: "No",
			},
			want: 4,
		},
		{
			name: "Should report 3 days if arriving Thursday and not staying late",
			args: args{
				arrival:     "Thursday PM",
				stayingLate: "No",
			},
			want: 3,
		},
		{
			name: "Should report 2 days if arriving Friday and not staying late",
			args: args{
				arrival:     "Friday AM",
				stayingLate: "No",
			},
			want: 2,
		},
		{
			name: "Should report 1 day if arriving Saturday and not staying late",
			args: args{
				arrival:     "Saturday",
				stayingLate: "No",
			},
			want: 1,
		},
		{
			name: "Should report 3 days if arriving Friday and staying late",
			args: args{
				arrival:     "Friday",
				stayingLate: "Yes",
			},
			want: 3,
		},
		{
			name: "Should default to five days if unknown number of nights",
			args: args{
				arrival:     "Miércoles",
				stayingLate: "No",
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &MessageProcessor{}
			if got := mp.computeNights(tt.args.arrival, tt.args.stayingLate); got != tt.want {
				t.Errorf("computeNights() = %v, want %v", got, tt.want)
			}
		})
	}
}
