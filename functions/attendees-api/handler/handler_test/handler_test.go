package handler_test

import (
	"clams/attendee"
	"clams/attendees-api/handler"
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

func (r *MockRegister) GetAttendees(authCode string) ([]attendee.Attendee, error) {
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
	mockRegister.On("GetAttendees", "").Return(attendees, nil)
	h := handler.Handler{AttendeesStore: &mockRegister}

	// When
	response, err := h.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"Attendees\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)

	mockRegister.AssertCalled(t, "GetAttendees", "")
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
	mockRegister.On("GetAttendees", "12345").Return(attendees, nil)
	h := handler.Handler{AttendeesStore: &mockRegister}

	// When
	response, err := h.HandleRequest(ctx, events.APIGatewayProxyRequest{PathParameters: map[string]string{"authCode": "12345"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}
	m, _ := json.Marshal(attendees)
	body := fmt.Sprintf("{\"Attendees\":%s}", m)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: body}, response)
	assert.Nil(t, err)

	mockRegister.AssertCalled(t, "GetAttendees", "12345")
}

func Test_shouldReturnNoContentWhenThereAreNoAttendees(t *testing.T) {
	// Given
	mockRegister := MockRegister{}
	ctx := context.Background()
	mockRegister.On("GetAttendees", "").Return([]attendee.Attendee{}, nil)
	h := handler.Handler{AttendeesStore: &mockRegister}

	// When
	response, err := h.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, response)
}

func Test_shouldReturnErrorWhenUnableToFetchAttendees(t *testing.T) {
	// Given
	mockAttendeesDatastore := MockRegister{}
	ctx := context.Background()
	mockAttendeesDatastore.On("GetAttendees", "").Return([]attendee.Attendee{}, fmt.Errorf("some error"))
	h := handler.Handler{AttendeesStore: &mockAttendeesDatastore}

	// When
	response, err := h.HandleRequest(ctx, events.APIGatewayProxyRequest{})

	// Then
	assert.Equal(t, fmt.Errorf("some error"), err)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
}

var jsonHeader = map[string]string{
	"Content-Type": "application/json",
}

func Test_shouldReturnReportWhenAttendeesExistInDatastore(t *testing.T) {
	// Given
	ctx := context.Background()

	mockAttendeesDatastore := MockRegister{}
	mockAttendeesDatastore.On("GetAttendees", "").Return(someAttendees(), nil)
	h := handler.Handler{AttendeesStore: &mockAttendeesDatastore}

	// When
	response, err := h.HandleRequest(ctx, events.APIGatewayProxyRequest{Path: "/report"})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: jsonHeader, Body: string(report())}, response)
}

func report() []byte {
	r, _ := json.Marshal(handler.Report{
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
