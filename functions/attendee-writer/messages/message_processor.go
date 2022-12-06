package messages

import (
	"fmt"
)

type IAttendeesStore interface {
	Store(attendee interface{}) error
}

type MessageProcessor struct {
	AttendeesStore IAttendeesStore
}

func (mp MessageProcessor) ProcessMessage(message Message) error {
	a, err := AttendeeFactory{}.newFromMessage(message)
	if err != nil {
		return err
	}

	if err := mp.AttendeesStore.Store(a); err != nil {
		return fmt.Errorf("storing attendee in datastore: %v", err)
	}

	return nil
}
