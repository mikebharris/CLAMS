package messages_test

import (
	"clams/attendee"
	"clams/attendee-writer/messages"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AttendeeFactory_computeNights(t *testing.T) {
	type args struct {
		arrival     string
		stayingLate string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Should report 4 days if arriving Wednesday and not staying late",
			args: args{
				arrival:     "Wednesday AM",
				stayingLate: "No",
			},
			want: 4,
		},
		{
			name: "Should report 3 days if arriving Thursday and not staying late",
			args: args{
				arrival:     "Thursday PM",
				stayingLate: "No",
			},
			want: 3,
		},
		{
			name: "Should report 2 days if arriving Friday and not staying late",
			args: args{
				arrival:     "Friday AM",
				stayingLate: "No",
			},
			want: 2,
		},
		{
			name: "Should report 1 day if arriving Saturday and not staying late",
			args: args{
				arrival:     "Saturday",
				stayingLate: "No",
			},
			want: 1,
		},
		{
			name: "Should report 3 days if arriving Friday and staying late",
			args: args{
				arrival:     "Friday",
				stayingLate: "Yes",
			},
			want: 3,
		},
		{
			name: "Should default to five days if unknown number of nights",
			args: args{
				arrival:     "Miércoles",
				stayingLate: "No",
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			boxForSpyToPutStoredAttendeesIn := &[]attendee.Attendee{}

			mp := messages.MessageProcessor{AttendeesStore: spyingAttendeesStore{boxForSpyToPutStoredAttendeesIn}}
			message := messages.Message{ArrivalDay: tt.args.arrival, StayingLate: tt.args.stayingLate}

			// When
			_ = mp.ProcessMessage(message)

			// Then
			x := *boxForSpyToPutStoredAttendeesIn
			assert.Equal(t, tt.args.arrival, x[0].ArrivalDay)
			assert.Equal(t, tt.args.stayingLate, x[0].StayingLate)
			assert.Equal(t, tt.want, x[0].NumberOfNights)
		})
	}
}
