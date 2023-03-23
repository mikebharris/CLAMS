package service_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"

	"testing"
)

var signupIdForGrace int
var signupIdForMike int

type steps struct {
	containers   Containers
	auroraClient AuroraClient
	dynamoClient DynamoClient
	t            *testing.T
}

func (s *steps) startContainers() {
	s.containers.start()
	fmt.Println("Giving the containers a chance to start before running tests")
	time.Sleep(10 * time.Second)
}

func (s *steps) setUpAuroraClient() {
	s.auroraClient.port = s.containers.getAuroraPort()
	s.auroraClient.host = "localhost"
	s.auroraClient.dbconx = s.auroraClient.connectToDatabase()
}

func (s *steps) setUpDynamoClient() {
	s.dynamoClient = newDynamoClient("localhost", s.containers.getDynamoPort())
	s.dynamoClient.createWorkshopsSignupsTable()
}

func (s *steps) stopContainers() {
	fmt.Println("Lambda log:")
	readCloser := s.containers.getLambdaLog()
	buf := new(bytes.Buffer)
	buf.ReadFrom(readCloser)
	newStr := buf.String()
	fmt.Println(newStr)

	fmt.Println("Stopping containers")
	s.containers.stop()
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
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", s.containers.getLambdaPort())

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

func (s *steps) theWorkshopSignupsDatastoreIsUpdated() {
	signup, err := s.dynamoClient.getWorkshopSignup(signupIdForMike)
	if err != nil {
		panic(err)
	}

	assert.NotNil(s.t, signup, "No signup found")
	assert.Equal(s.t, "My Exciting Workshop on COBOL", signup.WorkshopTitle)
	assert.Equal(s.t, "Mike Harris", signup.Name)
	assert.Equal(s.t, "attendee", signup.Role)
}
