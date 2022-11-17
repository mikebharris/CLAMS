package handler_test

import (
	"attendee-writer/handler"
	"attendee-writer/messages"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: check that it processes multiple messages
//TODO: check that if one message fails, another is still processed

func Test_ShouldPutMessageProcessingFailuresInBatchItemFailures(t *testing.T) {
	// Given
	ctx := context.Background()

	h := handler.Handler{
		MessageProcessor: messageProcessorThatFailsToProcessMessage{},
	}

	// When
	aMessage := events.SQSMessage{MessageId: "abcdef"}
	request, _ := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{aMessage}})

	// Then
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
		{ItemIdentifier: aMessage.MessageId},
	}}, request)
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
