package messages

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/CLAMS/attendee"
	"log"
)

type IAttendeesStore interface {
	Store(ctx context.Context, attendee attendee.Attendee) error
}

type MessageProcessor struct {
	AttendeesStore IAttendeesStore
}

func (mp MessageProcessor) ProcessMessage(ctx context.Context, msg events.SQSMessage) error {
	a, err := AttendeeFactory{}.NewFromMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("processing a message with id %s for event source %s\nattendee = %v", msg.MessageId, msg.EventSource, a)

	if err := mp.AttendeesStore.Store(ctx, a); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}
