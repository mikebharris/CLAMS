package handler

import (
	"clams/attendee"
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type IAttendeesStore interface {
	GetAttendees(authCode string) ([]attendee.Attendee, error)
}

type ApiResponse struct {
	Attendees []attendee.Attendee `json:"Attendees"`
}

type Handler struct {
	AttendeesStore IAttendeesStore
}

func (h Handler) HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.APIGatewayProxyResponse, error) {
	attendees, err := h.AttendeesStore.GetAttendees(h.getAuthCodeFromPath(request.RawPath))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(attendees) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
	}

	var responseBody []byte

	if strings.Contains(request.RawPath, "report") {
		responseBody = doReport(attendees)
	} else {
		responseBody = doOtherThing(attendees)
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(responseBody)}, nil

}

func (h Handler) getAuthCodeFromPath(path string) string {
	r := regexp.MustCompile("/attendees/([A-Za-z0-9]{6})")
	submatch := r.FindStringSubmatch(path)
	return submatch[1]
}

func doOtherThing(attendees []attendee.Attendee) []byte {
	m, _ := json.Marshal(ApiResponse{Attendees: attendees})
	return m
}

func doReport(attendees []attendee.Attendee) []byte {

	report := Report{
		TotalAttendees: len(attendees),
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
	return m
}
