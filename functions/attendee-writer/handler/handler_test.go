package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockMessageProcessor struct {
	mock.Mock
}

func (mp *MockMessageProcessor) ProcessMessage(ctx context.Context, message events.SQSMessage) error {
	args := mp.Called(ctx, message)
	return args.Error(0)
}

func Test_handleRequest_ShouldProcessSingleMessage(t *testing.T) {
	// Given
	ctx := context.Background()

	msg := events.SQSMessage{MessageId: "abcdef"}

	mp := MockMessageProcessor{}
	mp.On("ProcessMessage", ctx, msg).Return(nil)

	h := Handler{
		MessageProcessor: &mp,
	}

	// When
	request, err := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{msg}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: nil}, request)
	mp.AssertNumberOfCalls(t, "ProcessMessage", 1)
}

func Test_handleRequest_ShouldProcessMultipleMessages(t *testing.T) {
	// Given
	ctx := context.Background()

	msg1 := events.SQSMessage{MessageId: "abcdef"}
	msg2 := events.SQSMessage{MessageId: "123456"}

	mp := MockMessageProcessor{}
	mp.On("ProcessMessage", ctx, msg1).Return(nil)
	mp.On("ProcessMessage", ctx, msg2).Return(nil)

	h := Handler{
		MessageProcessor: &mp,
	}

	// When
	request, err := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{msg1, msg2}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: nil}, request)
	mp.AssertNumberOfCalls(t, "ProcessMessage", 2)
}

func Test_handleRequest_ShouldReturnSliceOfFailedMessages(t *testing.T) {
	// Given
	ctx := context.Background()

	msg1 := events.SQSMessage{MessageId: "fedcba"}
	msg2 := events.SQSMessage{MessageId: "654321"}

	mp := MockMessageProcessor{}
	mp.On("ProcessMessage", ctx, msg1).Return(nil)
	mp.On("ProcessMessage", ctx, msg2).Return(errors.New("cannot process message"))

	h := Handler{
		MessageProcessor: &mp,
	}

	// When
	request, err := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{msg1, msg2}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{{ItemIdentifier: msg2.MessageId}}}, request)
	mp.AssertNumberOfCalls(t, "ProcessMessage", 2)
}

func Test_handleRequest_ShouldReturnErrorIfThereSqsEventContainsNoMessages(t *testing.T) {
	// Given
	ctx := context.Background()

	msg := events.SQSMessage{MessageId: "abcdef"}

	mp := MockMessageProcessor{}
	mp.On("ProcessMessage", ctx, msg).Return(nil)

	h := Handler{
		MessageProcessor: &mp,
	}

	// When
	_, err := h.HandleRequest(ctx, events.SQSEvent{Records: []events.SQSMessage{}})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("sqs event contained no records"), err)
	mp.AssertNotCalled(t, "ProcessMessage")
}
