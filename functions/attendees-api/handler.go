package main

import (
	"attendees-api/storage"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type Handler struct {
	attendees storage.Attendees
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var attendees *storage.ApiResponse
	var err error

	authCode := request.PathParameters["authCode"]
	if authCode != "" {
		attendees, err = h.attendees.FetchAttendee(ctx, authCode)
	} else {
		attendees, err = h.attendees.FetchAllAttendees(ctx)
	}
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(attendees.Attendees) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	m, err := json.Marshal(attendees)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
