package handler

import (
	"context"
	"encoding/json"
	"github.com/mikebharris/CLAMS/attendee"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type IAttendeesStore interface {
	GetAllAttendees(ctx context.Context) ([]attendee.Attendee, error)
	GetAttendeesWithAuthCode(ctx context.Context, authCode string) ([]attendee.Attendee, error)
}

type ApiResponse struct {
	Attendees []attendee.Attendee `json:"Attendees"`
}

type Handler struct {
	AttendeesStore IAttendeesStore
}

func (h Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var attendees []attendee.Attendee
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

	if len(attendees) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	m, _ := json.Marshal(ApiResponse{Attendees: attendees})

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, nil
}
