package messages_test

import (
	"github.com/mikebharris/CLAMS/attendee"
	"github.com/stretchr/testify/mock"
)

type MockAttendeesStore struct {
	mock.Mock
}

func (s *MockAttendeesStore) Store(attendee attendee.Attendee) error {
	args := s.Called(attendee)
	return args.Error(0)
}
