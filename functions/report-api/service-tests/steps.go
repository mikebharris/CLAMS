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
	Attendees []Attendee `json:"Attendee"`
}

type Day string

type HeadCount struct {
	Day   Day
	Count int
}

type ReportApiResponse struct {
	TotalAttendees        int
	TotalKids             int
	TotalNightsCamped     int
	TotalCampingCharge    int
	TotalPaid             int
	TotalToPay            int
	TotalIncome           int
	AveragePaidByAttendee int
	DailyHeadCounts       []HeadCount
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
		NumberOfKids: 1,
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

func (s *steps) theFrontEndRequestsTheStatsFromTheReportAPI() error {
	localLambdaInvocationPort, err := s.containers.GetLocalHostLambdaPort()
	if err != nil {
		fmt.Println(err)
		return err
	}
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", localLambdaInvocationPort)

	request := events.APIGatewayProxyRequest{}
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

func (s *steps) someStatisticsAboutTheEventAreReturned() error {
	assert.Equal(s.t, http.StatusOK, responseFromLambda.StatusCode)

	apiResponse := ReportApiResponse{}
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &apiResponse); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, 2, apiResponse.TotalAttendees)
	assert.Equal(s.t, 8, apiResponse.TotalNightsCamped)
	assert.Equal(s.t, 80*100, apiResponse.TotalCampingCharge)
	assert.Equal(s.t, 552, apiResponse.TotalPaid)
	assert.Equal(s.t, 512, apiResponse.TotalToPay)
	assert.Equal(s.t, 1064, apiResponse.TotalIncome)
	assert.Equal(s.t, 5, apiResponse.TotalKids)

	//assert.Equal(s.t, "Wednesday", apiResponse.DailyHeadCounts[0].Day)
	//assert.Equal(s.t, 1, apiResponse.DailyHeadCounts[0].Count)
	//assert.Equal(s.t, "Thursday", apiResponse.DailyHeadCounts[1].Day)
	//assert.Equal(s.t, 2, apiResponse.DailyHeadCounts[1].Count)

	return nil
}
