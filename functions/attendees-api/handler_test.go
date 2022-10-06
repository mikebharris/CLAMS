package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

type MockRegister struct {
	mock.Mock
}

func (r *MockRegister) GetAllAttendees(ctx context.Context) (*ApiResponse, error) {
	args := r.Called(ctx)
	return args.Get(0).(*ApiResponse), args.Error(1)
}

func (r *MockRegister) GetAttendeesWithAuthCode(ctx context.Context, authCode string) (*ApiResponse, error) {
	args := r.Called(ctx, authCode)
	return args.Get(0).(*ApiResponse), args.Error(1)
}

func Test_shouldReturnAllAttendeesWhenNoAuthCodeProvided(t *testing.T) {
	// Given
	mockRegister := MockRegister{}

	attendees := []Attendee{
		{
			AuthCode:       "12345",
			Name:           "Bob Storey-Day",
			Email:          "",
			Telephone:      "",
			NumberOfKids:   1,
			Diet:           "",
			Financials:     Financials{},
			ArrivalDay:     "",
			NumberOfNights: 3,
			StayingLate:    "",
			CreatedTime:    time.Time{},
		},
		{
			AuthCode:       "23456",
			Name:           "Craig",
			Email:          "",
			Telephone:      "",
			NumberOfKids:   1,
			Diet:           "",
			Financials:     Financials{},
			ArrivalDay:     "",
			NumberOfNights: 3,
			StayingLate:    "",
			CreatedTime:    time.Time{},
		},
	}

	ctx := context.Background()
	mockRegister.On("GetAllAttendees", ctx).Return(&ApiResponse{Attendees: attendees}, nil)
	handler := Handler{attendeesStore: &mockRegister}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"Attendees\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)

	mockRegister.AssertCalled(t, "GetAllAttendees", ctx)
}

func Test_shouldReturnSingleAttendeesWhenAuthCodeProvided(t *testing.T) {
	// Given
	mockRegister := MockRegister{}

	attendees := []Attendee{
		{
			AuthCode:       "12345",
			Name:           "Bob Storey-Day",
			Email:          "",
			Telephone:      "",
			NumberOfKids:   1,
			Diet:           "",
			Financials:     Financials{},
			ArrivalDay:     "",
			NumberOfNights: 3,
			StayingLate:    "",
			CreatedTime:    time.Time{},
		},
	}

	ctx := context.Background()
	mockRegister.On("GetAttendeesWithAuthCode", ctx, "12345").Return(&ApiResponse{Attendees: attendees}, nil)
	handler := Handler{attendeesStore: &mockRegister}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{PathParameters: map[string]string{"authCode": "12345"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"Attendees\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)

	mockRegister.AssertCalled(t, "GetAttendeesWithAuthCode", ctx, "12345")
}

func Test_shouldReturnNoContentWhenThereAreNoAttendees(t *testing.T) {
	// Given
	mockRegister := MockRegister{}
	ctx := context.Background()
	mockRegister.On("GetAllAttendees", ctx).Return(&ApiResponse{Attendees: []Attendee{}}, nil)
	handler := Handler{attendeesStore: &mockRegister}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, response)
}

func Test_shouldReturnErrorWhenUnableToFetchAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockRegister{}
	ctx := context.Background()
	mockAttendeesDatastore.On("GetAllAttendees", ctx).Return(&ApiResponse{}, fmt.Errorf("some error"))
	handler := Handler{attendeesStore: &mockAttendeesDatastore}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
}
