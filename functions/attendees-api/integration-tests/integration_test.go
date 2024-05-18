package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mikebharris/testcontainernetwork-go"
	"github.com/mikebharris/testcontainernetwork-go/clients"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	var steps steps
	steps.t = t

	suite := godog.TestSuite{
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(steps.startContainerNetwork)
			ctx.BeforeSuite(steps.initialiseDynamoDb)
			ctx.AfterSuite(steps.stopContainerNetwork)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^some attendee records exist in CLAMS$`, steps.someAttendeeRecordsExistInTheAttendeesDatastore)
			ctx.Step(`^the front-end requests a specific attendee record from the endpoint$`, steps.theFrontendRequestsASpecificRecordFromTheEndpoint)
			ctx.Step(`^the record is returned$`, steps.aSingleRecordIsReturned)

			ctx.Step(`^the front-end requests all records from the endpoint$`, steps.theFrontendRequestsAllRecordsFromTheEndpoint)
			ctx.Step(`^all available records are returned$`, steps.theRecordsAreReturned)

			ctx.Step(`^the front-end requests the stats from the report endpoint$`, steps.theFrontEndRequestsTheStatsFromTheReportEndpoint)
			ctx.Step(`^some statistics about the event are returned$`, steps.someStatisticsAboutTheEventAreReturned)

		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type steps struct {
	t                   *testing.T
	networkOfContainers testcontainernetwork.NetworkOfDockerContainers
	lambdaContainer     testcontainernetwork.LambdaDockerContainer
	dynamoDbContainer   testcontainernetwork.DynamoDbDockerContainer
}

var responseFromLambda events.APIGatewayProxyResponse

const (
	dynamoDbHostname   = "dynamodb"
	dynamoDbPort       = 8000
	attendeesTableName = "attendees"
)

func (s *steps) startContainerNetwork() {
	s.dynamoDbContainer = testcontainernetwork.DynamoDbDockerContainer{
		Config: testcontainernetwork.DynamoDbDockerContainerConfig{
			Hostname: dynamoDbHostname,
			Port:     dynamoDbPort,
		},
	}

	s.lambdaContainer = testcontainernetwork.LambdaDockerContainer{
		Config: testcontainernetwork.LambdaDockerContainerConfig{
			Hostname:   "lambda",
			Executable: "../main",
			Environment: map[string]string{
				"DYNAMO_ENDPOINT_OVERRIDE": fmt.Sprintf("http://%s:%d", dynamoDbHostname, dynamoDbPort),
				"ATTENDEES_TABLE_NAME":     attendeesTableName,
			},
		},
	}

	s.networkOfContainers =
		testcontainernetwork.NetworkOfDockerContainers{}.
			WithDockerContainer(&s.lambdaContainer).
			WithDockerContainer(&s.dynamoDbContainer)
	_ = s.networkOfContainers.StartWithDelay(2 * time.Second)
}

func (s *steps) initialiseDynamoDb() {
	i := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("AuthCode"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("AuthCode"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(attendeesTableName),
	}

	dynamoDbClient, err := clients.DynamoDbClient{}.New("localhost", s.dynamoDbContainer.MappedPort())
	if err != nil {
		log.Fatalf("creating DynamoDB client: %v", err)
	}

	if err = dynamoDbClient.CreateTable(i); err != nil {
		log.Fatalf("creating table: %v", err)
	}
}

func (s *steps) stopContainerNetwork() {
	if err := s.networkOfContainers.Stop(); err != nil {
		log.Fatalf("stopping docker containers: %v", err)
	}
}

func (s *steps) someAttendeeRecordsExistInTheAttendeesDatastore() {
	dynamoDbClient, err := clients.DynamoDbClient{}.New("localhost", s.dynamoDbContainer.MappedPort())
	if err != nil {
		log.Fatalf("creating DynamoDB client: %v", err)
	}

	if err := dynamoDbClient.PutObject(attendeesTableName, Attendee{
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
		log.Fatalf("adding attendee: %v", err)
	}

	if err := dynamoDbClient.PutObject(attendeesTableName, Attendee{
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
	}); err != nil {
		log.Fatalf("adding attendee: %v", err)
	}
}

func (s *steps) theFrontendRequestsASpecificRecordFromTheEndpoint() {
	s.invokeLambdaUsingRequest(events.APIGatewayProxyRequest{PathParameters: map[string]string{"authCode": "123456"}})
}

func (s *steps) theFrontendRequestsAllRecordsFromTheEndpoint() {
	s.invokeLambdaUsingRequest(events.APIGatewayProxyRequest{})
}

func (s *steps) invokeLambdaUsingRequest(request events.APIGatewayProxyRequest) {
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("marshalling lambda request %v", err)
	}
	response, err := http.Post(s.lambdaContainer.InvocationUrl(), "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		log.Fatalf("triggering lambda: %v", err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("invoking Lambda: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		log.Fatalf("reading response body: %v", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &responseFromLambda); err != nil {
		log.Fatalf("unmarshalling response: %v", err)
	}
}

func (s *steps) aSingleRecordIsReturned() error {
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
	apiResponse := AttendeesApiResponse{}
	if err := json.Unmarshal([]byte(responseFromLambda.Body), &apiResponse); err != nil {
		return fmt.Errorf("unmarshalling result: %s", err)
	}

	assert.Equal(s.t, 2, len(apiResponse.Attendees))

	assert.Equal(s.t, "678901", apiResponse.Attendees[0].AuthCode)
	assert.Equal(s.t, "Zak Mindwarp", apiResponse.Attendees[0].Name)
	assert.Equal(s.t, "zakm@spangled.net", apiResponse.Attendees[0].Email)
	assert.Equal(s.t, "123456789", apiResponse.Attendees[0].Telephone)
	assert.Equal(s.t, 1, apiResponse.Attendees[0].NumberOfKids)
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

func (s *steps) theFrontEndRequestsTheStatsFromTheReportEndpoint() {
	s.invokeLambdaUsingRequest(events.APIGatewayProxyRequest{Path: "/report"})
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

	return nil
}
