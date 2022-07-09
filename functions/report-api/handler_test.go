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

func (fs *MockRegister) Attendees(ctx context.Context) ([]Attendee, error) {
	args := fs.Called()
	return args.Get(0).([]Attendee), args.Error(1)
}

func Test_shouldReturnReportWhenAttendeesExistInDatastore(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockRegister{}

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

	mockAttendeesDatastore.On("Attendees").Return(attendees, nil)
	handler := Handler{register: &mockAttendeesDatastore}

	// When
	response, err := handler.Handle(context.Background(), events.APIGatewayProxyRequest{})
	fmt.Println(response)

	// Then
	assert.Nil(t, err)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	report, _ := json.Marshal(Report{
		TotalAttendees:        2,
		TotalKids:             2,
		TotalNightsCamped:     6,
		TotalCampingCharge:    6000,
		TotalPaid:             0,
		TotalToPay:            0,
		TotalIncome:           0,
		AveragePaidByAttendee: 0,
		DailyHeadCounts:       nil,
	})

	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(report)}, response)
}

func Test_shouldReturnNoContentWhenThereAreNoAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockRegister{}
	mockAttendeesDatastore.On("Attendees").Return([]Attendee{}, nil)
	handler := Handler{register: &mockAttendeesDatastore}

	// When
	response, err := handler.Handle(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, response)
}

func Test_shouldReturnErrorWhenUnableToFetchAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockRegister{}
	mockAttendeesDatastore.On("Attendees").Return([]Attendee{}, fmt.Errorf("some error"))
	handler := Handler{register: &mockAttendeesDatastore}

	// When
	response, err := handler.Handle(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, response)
}
