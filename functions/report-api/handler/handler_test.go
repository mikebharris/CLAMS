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

var jsonHeader = map[string]string{
	"Content-Type": "application/json",
}

type MockAttendeesStore struct {
	mock.Mock
}

func (fs MockAttendeesStore) GetAllAttendees(context.Context) ([]attendee.Attendee, error) {
	args := fs.Called()
	return args.Get(0).([]attendee.Attendee), args.Error(1)
}

func Test_shouldReturnReportWhenAttendeesExistInDatastore(t *testing.T) {
	// Given
	ctx := context.Background()

	mockAttendeesStore := MockAttendeesStore{}
	mockAttendeesStore.On("GetAllAttendees").Return(someAttendees(), nil)
	handler := Handler{AttendeesStore: mockAttendeesStore}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: jsonHeader, Body: string(report())}, response)
}

func Test_shouldReturnNoContentWhenThereAreNoAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockAttendeesStore{}
	mockAttendeesDatastore.On("GetAllAttendees").Return([]attendee.Attendee{}, nil)
	handler := Handler{AttendeesStore: &mockAttendeesDatastore}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, response)
}

func Test_shouldReturnErrorWhenUnableToFetchAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockAttendeesStore{}
	mockAttendeesDatastore.On("GetAllAttendees").Return([]attendee.Attendee{}, fmt.Errorf("some error"))
	handler := Handler{AttendeesStore: &mockAttendeesDatastore}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, response)
}

func report() []byte {
	r, _ := json.Marshal(Report{
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
	return r
}

func someAttendees() []attendee.Attendee {
	return []attendee.Attendee{
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
			Name:           "Craig Duffy",
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
}
