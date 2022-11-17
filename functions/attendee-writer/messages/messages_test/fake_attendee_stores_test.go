package messages_test

import (
	"fmt"
	"github.com/mikebharris/CLAMS/attendee"
)

type spyingAttendeesStore struct {
	attendees *[]attendee.Attendee
}

func (s spyingAttendeesStore) Store(attendee attendee.Attendee) error {
	*s.attendees = append(*s.attendees, attendee)
	return nil
}

type failingAttendeeStore struct {
	attendees *[]attendee.Attendee
}

func (s failingAttendeeStore) Store(attendee attendee.Attendee) error {
	return fmt.Errorf("some storage error")
}
