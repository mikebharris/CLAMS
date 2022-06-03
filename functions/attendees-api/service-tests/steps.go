package service_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

var responseFromLambda events.APIGatewayProxyResponse

type steps struct {
	containers   Containers
	DynamoClient DynamoClient
	t            *testing.T
}

type ApiResponse struct {
	AuthCode       string
	Name           string
	Email          string
	Telephone      string
	NumberOfKids   uint
	Diet           string
	NumberOfNights uint
	Financials     struct {
		AmountToPay int    `json:"AmountToPay"`
		AmountPaid  int    `json:"AmountPaid"`
		AmountDue   int    `json:"AmountDue"`
		DatePaid    string `json:"DatePaid"`
	} `json:"Financials"`
}

func (s *steps) startContainers(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	err := s.containers.Start()
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (s *steps) setUpDynamoClient(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	localDynamoPort, err := s.containers.GetLocalHostDynamoPort()
	if err != nil {
		return ctx, err
	}

	s.DynamoClient, err = newDynamoClient("localhost", localDynamoPort)
	if err != nil {
		return ctx, err
	}

	err = s.DynamoClient.createAttendeesTable()
	if err != nil {
		return ctx, err
	}

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

func (s *steps) anAttendeeRecordExistsInTheAttendeesDatastore() error {
	err := s.DynamoClient.addAttendee(Attendee{
		AuthCode:     "12345",
		Name:         "Frank",
		Email:        "frank.o@gfa.de",
		Telephone:    "123456789",
		NumberOfKids: 4,
		Diet:         "I eat BASIC code for lunch",
		Financials: Financials{
			AmountToPay: 1024,
			AmountPaid:  512,
			DatePaid:    "10/05/2022",
			AmountDue:   512,
		},
		ArrivalDay:     "Wednesday",
		NumberOfNights: 5,
		StayingLate:    "Yes",
		CreatedTime:    time.Now(),
	})

	return err
}

func (s *steps) theFrontendFetchesTheRecordFromTheAPI() error {
	localLambdaInvocationPort, err := s.containers.GetLocalHostLambdaPort()
	if err != nil {
		fmt.Println(err)
		return err
	}
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", localLambdaInvocationPort)

	params := make(map[string]string)
	params["authCode"] = "12345"
	request := events.APIGatewayProxyRequest{PathParameters: params}
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.String()
		return fmt.Errorf("unexpected response when invoking lambda: %d %s", response.StatusCode, body)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err = json.Unmarshal(body, &responseFromLambda); err != nil {
		fmt.Println(err)
		return fmt.Errorf("unmarshalling proxy response: %s", err)
	}

	return nil
}

func (s *steps) theRecordIsReturned() error {

	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	apiResponse := ApiResponse{}
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &apiResponse); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, "12345", apiResponse.AuthCode)
	assert.Equal(s.t, "Frank", apiResponse.Name)
	assert.Equal(s.t, "frank.o@gfa.de", apiResponse.Email)
	assert.Equal(s.t, "123456789", apiResponse.Telephone)
	assert.Equal(s.t, uint(4), apiResponse.NumberOfKids)
	assert.Equal(s.t, "I eat BASIC code for lunch", apiResponse.Diet)
	assert.Equal(s.t, uint(5), apiResponse.NumberOfNights)
	assert.Equal(s.t, 1024, apiResponse.Financials.AmountToPay)
	assert.Equal(s.t, 512, apiResponse.Financials.AmountPaid)
	assert.Equal(s.t, "10/05/2022", apiResponse.Financials.DatePaid)
	assert.Equal(s.t, 512, apiResponse.Financials.AmountDue)

	return nil
}
