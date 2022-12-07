package messages

import (
	"clams/attendee"
	"strings"
	"time"
)

type AttendeeFactory struct {
}

func (af AttendeeFactory) newFromMessage(message Message) (attendee.Attendee, error) {
	a := attendee.Attendee{
		AuthCode:     message.AuthCode,
		Name:         message.Name,
		Email:        message.Email,
		Telephone:    message.Telephone,
		NumberOfKids: message.NumberOfKids,
		Diet:         message.Diet,
		Financials: attendee.Financials{
			AmountToPay: message.AmountToPay,
			AmountPaid:  message.AmountPaid,
			AmountDue:   message.AmountToPay - message.AmountPaid,
			DatePaid:    message.DatePaid,
		},
		ArrivalDay:     message.ArrivalDay,
		NumberOfNights: af.computeNights(message.ArrivalDay, message.StayingLate),
		StayingLate:    message.StayingLate,
		CreatedTime:    time.Now(),
	}
	return a, nil
}

func (af AttendeeFactory) computeNights(arrival string, stayingLate string) int {
	var nights int

	if strings.Contains(arrival, "Wednesday") {
		nights = 4
	} else if strings.Contains(arrival, "Thursday") {
		nights = 3
	} else if strings.Contains(arrival, "Friday") {
		nights = 2
	} else if strings.Contains(arrival, "Saturday") {
		nights = 1
	} else {
		nights = 5
	}

	if stayingLate == "Yes" {
		nights += 1
	}

	return nights
}
