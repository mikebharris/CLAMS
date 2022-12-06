package main

import (
	"attendee-writer/attendee"
	"attendee-writer/messages"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: check that it processes multiple messages
//TODO: check that if one message fails, another is still processed

type failingAttendeeStore struct {
	attendees *[]attendee.Attendee
}

func (s failingAttendeeStore) Store(_ interface{}) error {
	return fmt.Errorf("some storage error")
}

func Test_ShouldPutMessageProcessingFailuresInBatchItemFailures(t *testing.T) {
	// Given
	ctx := context.Background()

	h := Handler{
		MessageProcessor: messages.MessageProcessor{AttendeesStore: failingAttendeeStore{}},
	}

	// When
	aMessage := events.SQSMessage{MessageId: "abcdef"}
	response, _ := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{aMessage}})

	// Then
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
		{ItemIdentifier: aMessage.MessageId},
	}}, response)
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

func Test_processMessage_ShouldPutMessageOnBatchItemFailuresIfUnableToParseBody(t *testing.T) {
	// Given
	h := Handler{
		MessageProcessor: messages.MessageProcessor{},
	}

	// When
	response, _ := h.HandleRequest(context.Background(), events.SQSEvent{Records: []events.SQSMessage{{MessageId: "12345", Body: ""}}})

	// Then
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
		{ItemIdentifier: "12345"},
	}}, response)
}
