package messages_test

import (
	"clams/attendee"
	"fmt"
)

type spyingAttendeesStore struct {
	attendees *[]attendee.Attendee
}

func (s spyingAttendeesStore) Store(a interface{}) error {
	*s.attendees = append(*s.attendees, a.(attendee.Attendee))
	return nil
}

type failingAttendeeStore struct {
	attendees *[]attendee.Attendee
}

func (s failingAttendeeStore) Store(_ interface{}) error {
	return fmt.Errorf("some storage error")
}
