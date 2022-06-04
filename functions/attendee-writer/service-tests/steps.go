package service_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"net/http"
	"testing"

	"github.com/cucumber/godog"
)

type steps struct {
	containers   Containers
	DynamoClient DynamoClient
	CtxKey       string
	t            *testing.T
}

type Message struct {
	AuthCode     string
	Name         string
	Email        string
	AmountToPay  int
	AmountPaid   int
	DatePaid     string
	Telephone    string
	ArrivalDay   string
	StayingLate  string
	NumberOfKids int
	Diet         string
}

func (s *steps) startContainers(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	err := s.containers.Start()
	if err != nil {
		fmt.Printf("startContainers error %s", err)
		return ctx, err
	}
	return ctx, nil
}

func (s *steps) setUpDynamoClient(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	localDynamoPort, err := s.containers.GetLocalHostDynamoPort()
	if err != nil {
		fmt.Printf("setUpDynamoClient error %s", err)
		return ctx, err
	}

	s.DynamoClient, err = newDynamoClient("localhost", localDynamoPort)
	if err != nil {
		return ctx, err
	}

	err = s.DynamoClient.createAttendeesTable()
	return ctx, err
}

func (s *steps) stopContainers(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	fmt.Println("Lambda log:")
	readCloser, err := s.containers.GetLambdaLog()
	buf := new(bytes.Buffer)
	buf.ReadFrom(readCloser)
	newStr := buf.String()
	fmt.Println(newStr)

	fmt.Println("Stopping containers")
	newErr := s.containers.Stop()
	if newErr != nil && err == nil {
		err = newErr
	}
	return ctx, err
}

func (s *steps) theAttendeeIsAddedToTheAttendeesDatastore() error {
	attendee, err := s.DynamoClient.getAttendeeByCode("123456")
	if err != nil {
		return err
	}

	assert.Equal(s.t, "123456", attendee.AuthCode)
	assert.Equal(s.t, "Frank Ostrowski", attendee.Name)
	assert.Equal(s.t, "frank.o@gfa.de", attendee.Email)
	assert.Equal(s.t, "123456789", attendee.Telephone)
	assert.Equal(s.t, 1, attendee.NumberOfKids)
	assert.Equal(s.t, "I eat BASIC code for lunch", attendee.Diet)
	assert.Equal(s.t, 5, attendee.NumberOfNights)
	assert.Equal(s.t, 75, attendee.Financials.AmountToPay)
	assert.Equal(s.t, 50, attendee.Financials.AmountPaid)
	assert.Equal(s.t, "28/05/2022", attendee.Financials.DatePaid)
	assert.Equal(s.t, 25, attendee.Financials.AmountDue)
	assert.Equal(s.t, "Yes", attendee.StayingLate)
	assert.Equal(s.t, "Wednesday", attendee.ArrivalDay)

	return nil
}

func (s *steps) theAttendeeIsUpdatedInTheAttendeesDatastore() error {
	attendee, err := s.DynamoClient.getAttendeeByCode("123456")
	if err != nil {
		return err
	}

	assert.Equal(s.t, "123456", attendee.AuthCode)
	assert.Equal(s.t, "Frank Ostrowski", attendee.Name)
	assert.Equal(s.t, "frank.o@gfa.de", attendee.Email)
	assert.Equal(s.t, "123456789", attendee.Telephone)
	assert.Equal(s.t, 1, attendee.NumberOfKids)
	assert.Equal(s.t, "I eat BASIC code for lunch", attendee.Diet)
	assert.Equal(s.t, 4, attendee.NumberOfNights)
	assert.Equal(s.t, 75, attendee.Financials.AmountToPay)
	assert.Equal(s.t, 75, attendee.Financials.AmountPaid)
	assert.Equal(s.t, "29/05/2022", attendee.Financials.DatePaid)
	assert.Equal(s.t, 0, attendee.Financials.AmountDue)
	assert.Equal(s.t, "No", attendee.StayingLate)
	assert.Equal(s.t, "Wednesday", attendee.ArrivalDay)

	return nil
}

func (s *steps) theLambdaIsInvoked(payload Message) error {
	localLambdaInvocationPort, err := s.containers.GetLocalHostLambdaPort()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", localLambdaInvocationPort)

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	request := events.SQSEvent{Records: []events.SQSMessage{{Body: string(body)}}}
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.String()
		return fmt.Errorf("invoking Lambda: %d %s", response.StatusCode, body)
	}

	return nil
}

func (s *steps) theAttendeeWriterIsInvokedWithANewAttendeeRecord() error {
	request := Message{
		Name:         "Frank Ostrowski",
		Email:        "frank.o@gfa.de",
		AuthCode:     "123456",
		AmountToPay:  75,
		AmountPaid:   50,
		DatePaid:     "28/05/2022",
		Telephone:    "123456789",
		ArrivalDay:   "Wednesday",
		Diet:         "I eat BASIC code for lunch",
		StayingLate:  "Yes",
		NumberOfKids: 1,
	}

	return s.theLambdaIsInvoked(request)
}

func (s *steps) theAttendeeWriterIsInvokedWithAnUpdatedAttendeeRecord() error {
	request := Message{
		Name:         "Frank Ostrowski",
		Email:        "frank.o@gfa.de",
		AuthCode:     "123456",
		AmountToPay:  75,
		AmountPaid:   75,
		DatePaid:     "29/05/2022",
		Telephone:    "123456789",
		ArrivalDay:   "Wednesday",
		Diet:         "I eat BASIC code for lunch",
		StayingLate:  "No",
		NumberOfKids: 1,
	}
	return s.theLambdaIsInvoked(request)
}
