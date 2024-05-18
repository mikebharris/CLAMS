package integration_tests

import "time"

type Attendee struct {
	AuthCode       string
	Name           string
	Email          string
	Telephone      string
	NumberOfKids   int
	Diet           string
	Financials     Financials
	ArrivalDay     string
	NumberOfNights int
	StayingLate    string
	CreatedTime    time.Time
}

type Financials struct {
	AmountToPay int    `json:"AmountToPay"`
	AmountPaid  int    `json:"AmountPaid"`
	DatePaid    string `json:"DatePaid"`
	AmountDue   int    `json:"AmountDue"`
}

type AttendeesApiResponse struct {
	Attendees []Attendee `json:"Attendees"`
}

type Day string

type HeadCount struct {
	Day   Day
	Count int
}

type ReportApiResponse struct {
	TotalAttendees        int
	TotalKids             int
	TotalNightsCamped     int
	TotalCampingCharge    int
	TotalPaid             int
	TotalToPay            int
	TotalIncome           int
	AveragePaidByAttendee int
	DailyHeadCounts       []HeadCount
}
