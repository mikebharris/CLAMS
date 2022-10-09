package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"time"
)

type IAttendeesStore interface {
	Store(ctx context.Context, attendee Attendee) error
}

type IClock interface {
	Now() time.Time
}

type MessageProcessor struct {
	attendeesStore IAttendeesStore
	clock          IClock
}

func (mp MessageProcessor) processMessage(ctx context.Context, msg events.SQSMessage) error {
	attendeeFactory := AttendeeFactory{mp.clock}
	attendee, err := attendeeFactory.NewFromMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("processing a message with id %s for event source %s\nattendee = %v", msg.MessageId, msg.EventSource, attendee)

	if err := mp.attendeesStore.Store(ctx, attendee); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}
