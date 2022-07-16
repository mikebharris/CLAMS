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

func (r *MockRegister) Attendees(ctx context.Context) (*ApiResponse, error) {
	args := r.Called()
	return args.Get(0).(*ApiResponse), args.Error(1)
}

func (r *MockRegister) AttendeesWithAuthCode(ctx context.Context, authCode string) (*ApiResponse, error) {
	args := r.Called()
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

	mockRegister.On("Attendees").Return(&ApiResponse{Attendees: attendees}, nil)
	handler := Handler{register: &mockRegister}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{})
	fmt.Println(response)

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"AttendeesWithAuthCode\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)
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

	mockRegister.On("AttendeesWithAuthCode").Return(&ApiResponse{Attendees: attendees}, nil)
	handler := Handler{register: &mockRegister}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{PathParameters: map[string]string{"authCode": "12345"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"AttendeesWithAuthCode\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)
}

func Test_shouldReturnNoContentWhenThereAreNoAttendees(t *testing.T) {
	// Given
	mockRegister := MockRegister{}
	mockRegister.On("Attendees").Return(&ApiResponse{Attendees: []Attendee{}}, nil)
	handler := Handler{register: &mockRegister}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, response)
}

func Test_shouldReturnErrorWhenUnableToFetchAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockRegister{}
	mockAttendeesDatastore.On("Attendees").Return(&ApiResponse{}, fmt.Errorf("some error"))
	handler := Handler{register: &mockAttendeesDatastore}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
}
