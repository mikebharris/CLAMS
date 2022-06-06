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

type AttendeesApiResponse struct {
	Attendees []Attendee `json:"Attendees"`
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

func (s *steps) someAttendeeRecordsExistInTheAttendeesDatastore() error {
	if err := s.DynamoClient.addAttendee(Attendee{
		AuthCode:     "123456",
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
	}); err != nil {
		return err
	}

	return s.DynamoClient.addAttendee(Attendee{
		AuthCode:     "678901",
		Name:         "Zak Mindwarp",
		Email:        "zakm@spangled.net",
		Telephone:    "123456789",
		NumberOfKids: 0,
		Diet:         "I eat LSD for lunch",
		Financials: Financials{
			AmountToPay: 40,
			AmountPaid:  40,
			DatePaid:    "22/05/2022",
			AmountDue:   0,
		},
		ArrivalDay:     "Thursday",
		NumberOfNights: 3,
		StayingLate:    "No",
		CreatedTime:    time.Now(),
	})
}

func (s *steps) theFrontendRequestsASpecificRecordFromTheAPI() error {
	return s.invokeLambdaWithParameters(map[string]string{"authCode": "123456"})
}

func (s *steps) theFrontendRequestsAllRecordsFromTheAPI() error {
	return s.invokeLambdaWithParameters(nil)
}

func (s *steps) invokeLambdaWithParameters(params map[string]string) error {
	localLambdaInvocationPort, err := s.containers.GetLocalHostLambdaPort()
	if err != nil {
		fmt.Println(err)
		return err
	}
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", localLambdaInvocationPort)

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

func (s *steps) aSingleRecordIsReturned() error {
	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	apiResponse := AttendeesApiResponse{}
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &apiResponse); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, 1, len(apiResponse.Attendees))

	assert.Equal(s.t, "123456", apiResponse.Attendees[0].AuthCode)
	assert.Equal(s.t, "Frank", apiResponse.Attendees[0].Name)
	assert.Equal(s.t, "frank.o@gfa.de", apiResponse.Attendees[0].Email)
	assert.Equal(s.t, "123456789", apiResponse.Attendees[0].Telephone)
	assert.Equal(s.t, 4, apiResponse.Attendees[0].NumberOfKids)
	assert.Equal(s.t, "I eat BASIC code for lunch", apiResponse.Attendees[0].Diet)
	assert.Equal(s.t, 5, apiResponse.Attendees[0].NumberOfNights)
	assert.Equal(s.t, 1024, apiResponse.Attendees[0].Financials.AmountToPay)
	assert.Equal(s.t, 512, apiResponse.Attendees[0].Financials.AmountPaid)
	assert.Equal(s.t, "10/05/2022", apiResponse.Attendees[0].Financials.DatePaid)
	assert.Equal(s.t, 512, apiResponse.Attendees[0].Financials.AmountDue)
	assert.Equal(s.t, "Wednesday", apiResponse.Attendees[0].ArrivalDay)
	assert.Equal(s.t, "Yes", apiResponse.Attendees[0].StayingLate)
	assert.IsType(s.t, time.Time{}, apiResponse.Attendees[0].CreatedTime)

	return nil
}

func (s *steps) theRecordsAreReturned() error {
	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	apiResponse := AttendeesApiResponse{}
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &apiResponse); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, 2, len(apiResponse.Attendees))

	assert.Equal(s.t, "678901", apiResponse.Attendees[0].AuthCode)
	assert.Equal(s.t, "Zak Mindwarp", apiResponse.Attendees[0].Name)
	assert.Equal(s.t, "zakm@spangled.net", apiResponse.Attendees[0].Email)
	assert.Equal(s.t, "123456789", apiResponse.Attendees[0].Telephone)
	assert.Equal(s.t, 0, apiResponse.Attendees[0].NumberOfKids)
	assert.Equal(s.t, "I eat LSD for lunch", apiResponse.Attendees[0].Diet)
	assert.Equal(s.t, 3, apiResponse.Attendees[0].NumberOfNights)
	assert.Equal(s.t, 40, apiResponse.Attendees[0].Financials.AmountToPay)
	assert.Equal(s.t, 40, apiResponse.Attendees[0].Financials.AmountPaid)
	assert.Equal(s.t, "22/05/2022", apiResponse.Attendees[0].Financials.DatePaid)
	assert.Equal(s.t, 0, apiResponse.Attendees[0].Financials.AmountDue)
	assert.Equal(s.t, "Thursday", apiResponse.Attendees[0].ArrivalDay)
	assert.Equal(s.t, "No", apiResponse.Attendees[0].StayingLate)
	assert.IsType(s.t, time.Time{}, apiResponse.Attendees[0].CreatedTime)

	assert.Equal(s.t, "123456", apiResponse.Attendees[1].AuthCode)
	assert.Equal(s.t, "Frank", apiResponse.Attendees[1].Name)
	assert.Equal(s.t, "frank.o@gfa.de", apiResponse.Attendees[1].Email)
	assert.Equal(s.t, "123456789", apiResponse.Attendees[1].Telephone)
	assert.Equal(s.t, 4, apiResponse.Attendees[1].NumberOfKids)
	assert.Equal(s.t, "I eat BASIC code for lunch", apiResponse.Attendees[1].Diet)
	assert.Equal(s.t, 5, apiResponse.Attendees[1].NumberOfNights)
	assert.Equal(s.t, 1024, apiResponse.Attendees[1].Financials.AmountToPay)
	assert.Equal(s.t, 512, apiResponse.Attendees[1].Financials.AmountPaid)
	assert.Equal(s.t, "10/05/2022", apiResponse.Attendees[1].Financials.DatePaid)
	assert.Equal(s.t, 512, apiResponse.Attendees[1].Financials.AmountDue)
	assert.Equal(s.t, "Wednesday", apiResponse.Attendees[1].ArrivalDay)
	assert.Equal(s.t, "Yes", apiResponse.Attendees[1].StayingLate)
	assert.IsType(s.t, time.Time{}, apiResponse.Attendees[1].CreatedTime)

	return nil
}
