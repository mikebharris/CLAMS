package main

import (
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
			af := &AttendeeFactory{}
			if got := af.computeNights(tt.args.arrival, tt.args.stayingLate); got != tt.want {
				t.Errorf("computeNights() = %v, want %v", got, tt.want)
			}
		})
	}
}
