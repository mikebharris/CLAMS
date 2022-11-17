package messages_test

import (
	"attendee-writer/messages"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/CLAMS/attendee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_processMessage_ShouldStoreMessage(t *testing.T) {
	// Given
	mockAttendeesStore := MockAttendeesStore{}
	mockAttendeesStore.On("Store", mock.Anything).Return(nil)
	mp := messages.MessageProcessor{AttendeesStore: &mockAttendeesStore}

	body, _ := json.Marshal(messages.Message{Name: "Frank Ostrowski"})

	// When
	err := mp.ProcessMessage(events.SQSMessage{Body: string(body)})

	// Then
	assert.Nil(t, err)
	mockAttendeesStore.AssertCalled(t, "Store",
		attendee.Attendee{
			Name:           "Frank Ostrowski",
			NumberOfNights: 5,
		},
	)
}

func Test_processMessage_ShouldReturnErrorIfUnableToStoreMessage(t *testing.T) {
	// Given
	mockAttendeesStore := MockAttendeesStore{}
	mockAttendeesStore.On("Store", mock.Anything).Return(fmt.Errorf("some storage error"))

	mp := messages.MessageProcessor{AttendeesStore: &mockAttendeesStore}

	// When
	body, _ := json.Marshal(messages.Message{})
	err := mp.ProcessMessage(events.SQSMessage{Body: string(body)})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("storing attendee in datastore: some storage error"), err)
	mockAttendeesStore.AssertCalled(t, "Store", mock.Anything)
}

func Test_processMessage_ShouldReturnErrorIfUnableToParseMessage(t *testing.T) {
	// Given
	mp := messages.MessageProcessor{AttendeesStore: &attendee.AttendeesStore{}}

	// When
	err := mp.ProcessMessage(events.SQSMessage{Body: ""})

	// Then
	assert.NotNil(t, err)
}
