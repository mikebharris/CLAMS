package service_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type steps struct {
	t              *testing.T
	containers     Containers
	databaseClient DatabaseClient
	sqsClient      SqsClient
}

func (s *steps) startContainers() {
	s.containers.start()
	fmt.Println("Giving the containers a chance to start before running tests")
	time.Sleep(10 * time.Second)
}

func (s *steps) setUpDatabaseClient() {
	s.databaseClient.port = s.containers.getDatabasePort()
	s.databaseClient.host = "localhost"
	s.databaseClient.dbConx = s.databaseClient.connectToDatabase()
}

func (s *steps) setUpSqsClient() {
	s.sqsClient = newSqsClient("localhost", s.containers.getSqsPort())
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

type TriggerNotification struct {
	WorkshopSignupId int
	WorkshopId       int
	PeopleId         int
	RoleId           int
}

func (s *steps) thereAreDatabaseTriggerNotifications() {
	s.databaseClient.insertTriggerNotification(TriggerNotification{
		WorkshopSignupId: 1,
		WorkshopId:       2,
		PeopleId:         3,
		RoleId:           4,
	})
	s.databaseClient.insertTriggerNotification(TriggerNotification{
		WorkshopSignupId: 9,
		WorkshopId:       8,
		PeopleId:         7,
		RoleId:           6,
	})
}

func (s *steps) theDbTriggerLambdaIsInvoked() error {
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", s.containers.getLambdaPort())

	response, err := http.Post(url, "application/json", nil)
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

func (s *steps) theMessagesArePlacedOnTheQueue() {
	messages := s.sqsClient.getMessages()

	if len(messages) == 0 {
		panic(fmt.Errorf("expected at least one message in the SQS queue"))
	}

	assert.Len(s.t, messages, 2)

	var tm TriggerNotification
	json.Unmarshal([]byte(*messages[0].Body), &tm)
	assert.Equal(s.t, 1, tm.WorkshopSignupId)
	assert.Equal(s.t, 2, tm.WorkshopId)
	assert.Equal(s.t, 3, tm.PeopleId)
	assert.Equal(s.t, 4, tm.RoleId)

	json.Unmarshal([]byte(*messages[1].Body), &tm)
	assert.Equal(s.t, 9, tm.WorkshopSignupId)
	assert.Equal(s.t, 8, tm.WorkshopId)
	assert.Equal(s.t, 7, tm.PeopleId)
	assert.Equal(s.t, 6, tm.RoleId)
}

func (s *steps) theNotificationsAreRemovedFromTheDatabase() {
	assert.Equal(s.t, 0, s.databaseClient.countOfNotifications())
}

func (s *steps) closeDatabaseConnection() {
	s.databaseClient.closeDatabaseConnexion()
}
