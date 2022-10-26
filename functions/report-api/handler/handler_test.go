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

type MockAttendeesStore struct {
	mock.Mock
}

func (fs MockAttendeesStore) GetAllAttendees(ctx context.Context) ([]attendee.Attendee, error) {
	args := fs.Called(ctx)
	return args.Get(0).([]attendee.Attendee), args.Error(1)
}

func Test_shouldReturnReportWhenAttendeesExistInDatastore(t *testing.T) {
	// Given
	mockAttendeesStore := MockAttendeesStore{}

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
	mockAttendeesStore.On("GetAllAttendees", ctx).Return(attendees, nil)
	handler := Handler{AttendeesStore: mockAttendeesStore}

	// When
	response, err := handler.HandleRequest(ctx, events.APIGatewayProxyRequest{})
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
	mockAttendeesDatastore := MockAttendeesStore{}
	mockAttendeesDatastore.On("GetAllAttendees", mock.Anything).Return([]attendee.Attendee{}, nil)
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
	mockAttendeesDatastore.On("GetAllAttendees", mock.Anything).Return([]attendee.Attendee{}, fmt.Errorf("some error"))
	handler := Handler{AttendeesStore: &mockAttendeesDatastore}

	// When
	response, err := handler.HandleRequest(context.Background(), events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, response)
}
