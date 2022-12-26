package service_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"time"

	"testing"
)

type steps struct {
	containers   Containers
	auroraClient AuroraClient
	t            *testing.T
}

func (s *steps) startContainers() {
	s.containers.start()
	fmt.Println("Giving the containers a chance to start before running tests")
	time.Sleep(5 * time.Second)
}

func (s *steps) setUpAuroraClient() {
	s.auroraClient.port = s.containers.getAuroraPort()
	s.auroraClient.host = "localhost"
	s.auroraClient.dbconx = s.auroraClient.connectToDatabase()
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
	s.auroraClient.createDatabaseEntries()
}

type message struct {
	WorkshopSignupId int
	WorkshopId       int
	PeopleId         int
	RoleId           int
}

func (s *steps) theProcessorLambdaIsInvoked() error {
	request := message{
		WorkshopSignupId: 1,
		WorkshopId:       1,
		PeopleId:         2,
		RoleId:           1,
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
