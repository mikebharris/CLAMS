package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
)

type messageProcessorThatFailsToProcessMessage struct{}

func (m messageProcessorThatFailsToProcessMessage) ProcessMessage(_ events.SQSMessage) error {
	return errors.New("")
}
