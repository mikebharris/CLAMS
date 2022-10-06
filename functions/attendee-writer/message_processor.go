package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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

func (mp MessageProcessor) processMessage(ctx context.Context, message events.SQSMessage) error {
	attendeeFactory := AttendeeFactory{mp.clock}
	attendee, err := attendeeFactory.NewFromMessage(message)
	if err != nil {
		return err
	}

	if err := mp.attendeesStore.Store(ctx, attendee); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}
