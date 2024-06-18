package integration_tests

import (
	"bytes"
	"database/sql"
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
	"os"
	"path"
	"testing"
	"time"
)

func TestFeatures(t *testing.T) {
	var steps steps
	steps.t = t
	suite := godog.TestSuite{
		TestSuiteInitializer: func(ctx *godog.TestSuiteContext) {
			ctx.BeforeSuite(steps.startContainerNetwork)
			ctx.BeforeSuite(steps.setUpAuroraClient)
			ctx.BeforeSuite(steps.initialiseDynamoDb)
			ctx.AfterSuite(steps.stopContainerNetwork)
		},
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			ctx.Step(`^workshop signup records exist in the database$`, steps.theWorkshopSignupRequestExistsInTheDatabase)
			ctx.Step(`^the workshop signup processor receives a notification$`, steps.theProcessorLambdaIsInvoked)
			ctx.Step(`^the workshops signups datastore is updated$`, steps.theWorkshopSignupsDatastoreIsUpdated)
		},
		Options: &godog.Options{
			StopOnFailure: true,
			Strict:        true,
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t, // Testing instance that will run subtests.
		},
	}

	if run := suite.Run(); run != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

var signupIdForGrace int
var signupIdForMike int

type steps struct {
	auroraClient              AuroraClient
	lambdaContainer           testcontainernetwork.LambdaDockerContainer
	dynamoDbContainer         testcontainernetwork.DynamoDbDockerContainer
	auroraContainer           testcontainernetwork.PostgresDockerContainer
	flywayContainer           testcontainernetwork.FlywayDockerContainer
	networkOfDockerContainers testcontainernetwork.NetworkOfDockerContainers
	t                         *testing.T
	dbConx                    *sql.DB
}

const (
	dynamoDbHostname         = "dynamodb"
	dynamoDbPort             = 8000
	workshopSignupsTableName = "workshops"
)

func (s *steps) startContainerNetwork() {
	wd, _ := os.Getwd()
	s.dynamoDbContainer = testcontainernetwork.DynamoDbDockerContainer{
		Config: testcontainernetwork.DynamoDbDockerContainerConfig{
			Hostname: dynamoDbHostname,
			Port:     dynamoDbPort,
		},
	}

	s.auroraContainer = testcontainernetwork.PostgresDockerContainer{
		Config: testcontainernetwork.PostgresDockerContainerConfig{
			Hostname: "postgres",
			Port:     5432,
			Environment: map[string]string{
				"POSTGRES_PASSWORD": "d0ntHackM3",
				"POSTGRES_USER":     "hacktivista",
				"POSTGRES_DB":       "hacktionlab",
			},
		},
	}

	s.flywayContainer = testcontainernetwork.FlywayDockerContainer{
		Config: testcontainernetwork.FlywayDockerContainerConfig{
			Hostname:        "flyway",
			ConfigFilesPath: path.Join(wd, "flyway/conf"),
			SqlFilesPath:    path.Join(wd, "../../../flyway/sql"),
		},
	}

	s.lambdaContainer = testcontainernetwork.LambdaDockerContainer{
		Config: testcontainernetwork.LambdaDockerContainerConfig{
			Hostname:   "lambda",
			Executable: "../main",
			Environment: map[string]string{
				"DATABASE_URL": fmt.Sprintf(
					"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
					"postgres", 5432, "hacktivista", "d0ntHackM3", "hacktionlab", "hacktionlab_workshops",
				),
				"WORKSHOP_SIGNUPS_TABLE_NAME": workshopSignupsTableName,
				"DYNAMO_ENDPOINT_OVERRIDE":    "http://dynamo:8000",
			},
		},
	}

	s.networkOfDockerContainers =
		testcontainernetwork.NetworkOfDockerContainers{}.
			WithDockerContainer(&s.lambdaContainer).
			WithDockerContainer(&s.dynamoDbContainer).
			WithDockerContainer(&s.auroraContainer).
			WithDockerContainer(&s.flywayContainer)
	_ = s.networkOfDockerContainers.StartWithDelay(5 * time.Second)
}

func (s *steps) stopContainerNetwork() {
	if err := s.networkOfDockerContainers.Stop(); err != nil {
		log.Fatalf("stopping docker containers: %v", err)
	}
}

func (s *steps) initialiseDynamoDb() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("WorkshopSignupId"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("WorkshopSignupId"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(workshopSignupsTableName),
	}

	dynamoDbClient, err := clients.DynamoDbClient{}.New("localhost", s.dynamoDbContainer.MappedPort())
	if err != nil {
		log.Fatalf("creating DynamoDB client: %v", err)
	}

	if err = dynamoDbClient.CreateTable(input); err != nil {
		log.Fatalf("creating table: %v", err)
	}
}

func (s *steps) setUpAuroraClient() {
	s.auroraClient.port = s.auroraContainer.MappedPort()
	s.auroraClient.host = "localhost"
	s.auroraClient.dbConx = s.auroraClient.connectToDatabase()
}

func (s *steps) theWorkshopSignupRequestExistsInTheDatabase() {
	s.auroraClient.createRoles()
	workshopId := s.auroraClient.createWorkshop("My Exciting Workshop on COBOL")
	facilitatorId := s.auroraClient.createPerson("Grace", "Hopper", "g.hopper@codasyl.mil")
	attendeeId := s.auroraClient.createPerson("Mike", "Harris", "mike@cobolenthusiasts.biz")

	signupIdForGrace = s.auroraClient.createWorkshopSignup(workshopId, facilitatorId, facilitatorRoleId)
	signupIdForMike = s.auroraClient.createWorkshopSignup(workshopId, attendeeId, attendeeRoleId)
}

type message struct {
	WorkshopSignupId int
}

func (s *steps) theProcessorLambdaIsInvoked() error {
	request := message{
		WorkshopSignupId: signupIdForMike,
	}

	return s.theLambdaIsInvoked(request)
}

func (s *steps) theLambdaIsInvoked(payload message) error {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("marshalling payload: %v", err)
	}

	request := events.SQSEvent{Records: []events.SQSMessage{{Body: string(body)}}}
	requestJsonBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("marshalling request: %v", err)
	}

	response, err := http.Post(s.lambdaContainer.InvocationUrl(), "application/json", bytes.NewReader(requestJsonBytes))
	if err != nil {
		log.Fatalf("invoking Lambda: %v", err)
	}
	if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(response.Body)
		body := buf.String()
		return fmt.Errorf("invoking Lambda: %d %s", response.StatusCode, body)
	}

	return nil
}

func (s *steps) theWorkshopSignupsDatastoreIsUpdated() {
	dynamoDbClient, err := clients.DynamoDbClient{}.New("localhost", s.dynamoDbContainer.MappedPort())
	if err != nil {
		log.Fatalf("creating DynamoDB client: %v", err)
	}

	results, err := dynamoDbClient.GetItemsInTable(workshopSignupsTableName)
	if err != nil {
		log.Fatalf("getting items in table: %v", err)
	}

	var signups []struct {
		WorkshopSignupId int
		WorkshopTitle    string
		Name             string
		Role             string
	}

	err = attributevalue.UnmarshalListOfMaps(results, &signups)
	if err != nil {
		log.Fatalf("unmarshalling list of maps: %v", err)
	}

	for _, signup := range signups {
		if signup.WorkshopSignupId == signupIdForMike {
			assert.Equal(s.t, "My Exciting Workshop on COBOL", signup.WorkshopTitle)
			assert.Equal(s.t, "Mike Harris", signup.Name)
			assert.Equal(s.t, "attendee", signup.Role)
		}
	}
}
