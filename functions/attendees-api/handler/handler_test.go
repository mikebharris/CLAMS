package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/CLAMS/attendee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

type MockRegister struct {
	mock.Mock
}

func (r *MockRegister) GetAllAttendees() ([]attendee.Attendee, error) {
	args := r.Called()
	return args.Get(0).([]attendee.Attendee), args.Error(1)
}

func (r *MockRegister) GetAttendeesWithAuthCode(authCode string) ([]attendee.Attendee, error) {
	args := r.Called(authCode)
	return args.Get(0).([]attendee.Attendee), args.Error(1)
}

func Test_shouldReturnAllAttendeesWhenNoAuthCodeProvided(t *testing.T) {
	// Given
	mockRegister := MockRegister{}

	attendees := []attendee.Attendee{
		{
			AuthCode:       "12345",
			Name:           "Bob Storey-Day",
			Email:          "",
			Telephone:      "",
			NumberOfKids:   1,
			Diet:           "",
			Financials:     attendee.Financials{},
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
			Financials:     attendee.Financials{},
			ArrivalDay:     "",
			NumberOfNights: 3,
			StayingLate:    "",
			CreatedTime:    time.Time{},
		},
	}

	ctx := context.Background()
	mockRegister.On("GetAllAttendees").Return(attendees, nil)
	handler := Handler{AttendeesStore: &mockRegister}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"Attendees\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)

	mockRegister.AssertCalled(t, "GetAllAttendees")
}

func Test_shouldReturnSingleAttendeesWhenAuthCodeProvided(t *testing.T) {
	// Given
	mockRegister := MockRegister{}

	attendees := []attendee.Attendee{
		{
			AuthCode:       "12345",
			Name:           "Bob Storey-Day",
			Email:          "",
			Telephone:      "",
			NumberOfKids:   1,
			Diet:           "",
			Financials:     attendee.Financials{},
			ArrivalDay:     "",
			NumberOfNights: 3,
			StayingLate:    "",
			CreatedTime:    time.Time{},
		},
	}

	ctx := context.Background()
	mockRegister.On("GetAttendeesWithAuthCode", "12345").Return(attendees, nil)
	handler := Handler{AttendeesStore: &mockRegister}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{PathParameters: map[string]string{"authCode": "12345"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"Attendees\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)

	mockRegister.AssertCalled(t, "GetAttendeesWithAuthCode", "12345")
}

func Test_shouldReturnNoContentWhenThereAreNoAttendees(t *testing.T) {
	// Given
	mockRegister := MockRegister{}
	ctx := context.Background()
	mockRegister.On("GetAllAttendees").Return([]attendee.Attendee{}, nil)
	handler := Handler{AttendeesStore: &mockRegister}

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
	mockAttendeesDatastore.On("GetAllAttendees").Return([]attendee.Attendee{}, fmt.Errorf("some error"))
	handler := Handler{AttendeesStore: &mockAttendeesDatastore}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
}
