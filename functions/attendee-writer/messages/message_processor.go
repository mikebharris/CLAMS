package messages

import (
	"attendee-writer/attendee"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

type IAttendeesStore interface {
	Store(ctx context.Context, attendee attendee.Attendee) error
}

type MessageProcessor struct {
	AttendeesStore IAttendeesStore
	Clock          IClock
}

func (mp MessageProcessor) ProcessMessage(ctx context.Context, msg events.SQSMessage) error {
	attendeeFactory := AttendeeFactory{Clock: mp.Clock}
	attendee, err := attendeeFactory.NewFromMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("processing a message with id %s for event source %s\nattendee = %v", msg.MessageId, msg.EventSource, attendee)

	if err := mp.AttendeesStore.Store(ctx, attendee); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}
