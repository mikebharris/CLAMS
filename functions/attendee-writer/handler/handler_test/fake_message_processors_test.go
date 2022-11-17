package handler_test

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
)

type messageProcessorThatFailsToProcessMessage struct{}

func (m messageProcessorThatFailsToProcessMessage) ProcessMessage(msg events.SQSMessage) error {
	return errors.New("")
}
