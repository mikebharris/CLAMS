package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"report-api/attendee"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type IAttendeesStore interface {
	GetAllAttendees(ctx context.Context) ([]attendee.Attendee, error)
}

type Handler struct {
	AttendeesStore IAttendeesStore
}

func (h Handler) HandleRequest(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	attendees, err := h.AttendeesStore.GetAllAttendees(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
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

	m, _ := json.Marshal(report)

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
