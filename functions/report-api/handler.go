package main

import (
	"context"
	"encoding/json"
	"net/http"
	"report-api/storage"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type Handler struct {
	attendees storage.AttendeesDataStore
}

type Day string

type HeadCount struct {
	Day   Day
	Count int
}

type Report struct {
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

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	attendees, err := h.attendees.FetchAllAttendees(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	numberOfAttendees := len(attendees)
	if numberOfAttendees == 0 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}

	report := Report{
		TotalAttendees: numberOfAttendees,
	}

	for _, a := range attendees {
		report.TotalNightsCamped += a.NumberOfNights
		report.TotalCampingCharge += 10 * a.NumberOfNights * 100
		report.TotalToPay += a.Financials.AmountDue
		report.TotalIncome += a.Financials.AmountToPay
		report.TotalPaid += a.Financials.AmountPaid
		report.TotalKids += a.NumberOfKids
	}

	m, err := json.Marshal(report)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
