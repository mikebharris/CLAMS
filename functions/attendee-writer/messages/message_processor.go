package messages

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/CLAMS/attendee"
	"log"
)

type IAttendeesStore interface {
	Store(attendee attendee.Attendee) error
}

type MessageProcessor struct {
	AttendeesStore IAttendeesStore
}

func (mp MessageProcessor) ProcessMessage(msg events.SQSMessage) error {
	a, err := AttendeeFactory{}.newFromMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("processing a message with id %s for event source %s\nattendee = %v", msg.MessageId, msg.EventSource, a)

	if err := mp.AttendeesStore.Store(a); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}
