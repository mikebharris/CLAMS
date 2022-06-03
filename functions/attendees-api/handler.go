package main

import (
	"attendees-api/storage"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	attendees storage.IAttendees
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	authCode := request.PathParameters["authCode"]
	if authCode == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotImplemented}, nil
	}

	attendee, err := h.attendees.FetchAttendee(ctx, authCode)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	if attendee == nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	m, _ := json.Marshal(attendee)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
