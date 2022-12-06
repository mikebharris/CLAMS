package main

import (
	"attendee-writer/messages"
	"errors"
)

type messageProcessorThatFailsToProcessMessage struct{}

func (m messageProcessorThatFailsToProcessMessage) ProcessMessage(message messages.Message) error {
	return errors.New("")
}
