package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type Handler struct {
	attendees Attendees
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var attendees *ApiResponse
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

	m, _ := json.Marshal(attendees)

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
