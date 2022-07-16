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

type RegisterInterface interface {
	Attendees(ctx context.Context) (*ApiResponse, error)
	AttendeesWithAuthCode(ctx context.Context, authCode string) (*ApiResponse, error)
}

type Handler struct {
	register RegisterInterface
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var attendees *ApiResponse
	var err error

	authCode := request.PathParameters["authCode"]
	if authCode != "" {
		attendees, err = h.register.AttendeesWithAuthCode(ctx, authCode)
	} else {
		attendees, err = h.register.Attendees(ctx)
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
