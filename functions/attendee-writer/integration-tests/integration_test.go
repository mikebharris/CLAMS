package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cucumber/godog"
	"github.com/mikebharris/testcontainernetwork-go"
	"github.com/mikebharris/testcontainernetwork-go/clients"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	dynamoDbHostname   = "dynamodb"
	dynamoDbPort       = 8000
	attendeesTableName = "attendees"
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
			ctx.Step(`^the Attendee Writer is invoked with an attendee record from BAMS to be processed$`, steps.theAttendeeWriterIsInvokedWithANewAttendeeRecord)
			ctx.Step(`^the Attendee Writer is invoked with an updated attendee record from BAMS to be processed$`, steps.theAttendeeWriterIsInvokedWithAnUpdatedAttendeeRecord)
			ctx.Step(`^an attendee record is added to CLAMS$`, steps.theAttendeeIsAddedToTheAttendeesDatastore)
			ctx.Step(`^the attendee record is updated in CLAMS$`, steps.theAttendeeIsUpdatedInTheAttendeesDatastore)
		},
		Options: &godog.Options{
			StopOnFailure: true,
			Strict:        true,
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
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
	_ = s.networkOfContainers.StartWithDelay(5 * time.Second)
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

	dynamoDbClient, err := clients.DynamoDbClient{}.New(s.dynamoDbContainer.MappedPort())
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

func (s *steps) theAttendeeWriterIsInvokedWithANewAttendeeRecord() {
	s.theLambdaIsInvoked(Payload{
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
	})
}

func (s *steps) theAttendeeWriterIsInvokedWithAnUpdatedAttendeeRecord() {
	s.theLambdaIsInvoked(Payload{
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
	})
}

func (s *steps) theLambdaIsInvoked(payload Payload) {
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	request := events.SQSEvent{Records: []events.SQSMessage{{Body: string(body)}}}
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("marshalling request: %v", err)
	}

	response, err := http.Post(s.lambdaContainer.InvocationUrl(), "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		log.Fatalf("triggering lambda: %v", err)
	}

	if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(response.Body); err != nil {
			panic(err)
		}
		body := buf.String()
		log.Fatalf("invoking Lambda: %d %s", response.StatusCode, body)
	}
}

func (s *steps) theAttendeeIsAddedToTheAttendeesDatastore() error {
	attendee := s.getAttendeeByCode("123456")

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
	attendee := s.getAttendeeByCode("123456")

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

func (s *steps) getAttendeeByCode(authCode string) *Attendee {
	dynamoDbClient, err := clients.DynamoDbClient{}.New(s.dynamoDbContainer.MappedPort())
	if err != nil {
		log.Fatalf("creating DynamoDB client: %v", err)
	}

	itemsInTable, err := dynamoDbClient.GetItemsInTable(attendeesTableName)
	if err != nil {
		log.Fatalf("getting items in table %s: %v", attendeesTableName, err)
	}

	var attendees []Attendee
	if err := attributevalue.UnmarshalListOfMaps(itemsInTable, &attendees); err != nil {
		log.Fatalf("unmarshalling list of attendees: %v", err)
	}

	for _, a := range attendees {
		if a.AuthCode == authCode {
			return &a
		}
	}

	return nil
}
