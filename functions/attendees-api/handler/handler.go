package handler

import (
	"attendees-api/attendee"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type IAttendeesStore interface {
	GetAllAttendees(ctx context.Context) (*attendee.ApiResponse, error)
	GetAttendeesWithAuthCode(ctx context.Context, authCode string) (*attendee.ApiResponse, error)
}

type Handler struct {
	AttendeesStore IAttendeesStore
}

func (h Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var attendees *attendee.ApiResponse
	var err error

	authCode := request.PathParameters["authCode"]
	if authCode != "" {
		attendees, err = h.AttendeesStore.GetAttendeesWithAuthCode(ctx, authCode)
	} else {
		attendees, err = h.AttendeesStore.GetAllAttendees(ctx)
	}

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(attendees.Attendees) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	m, _ := json.Marshal(attendees)

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
