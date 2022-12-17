package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

type SimpleAuthorizerResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

type steps struct {
	containers     Containers
	responseStatus int
	responseBody   string
	t              *testing.T
}

func (s *steps) startContainers() {
	if err := s.containers.Start(); err != nil {
		panic(err)
	}
}

func (s *steps) stopContainers() {
	fmt.Println("Lambda log:")
	readCloser := s.containers.GetLambdaLog()
	buf := new(bytes.Buffer)
	buf.ReadFrom(readCloser)
	newStr := buf.String()
	fmt.Println(newStr)

	fmt.Println("Stopping containers")
	if err := s.containers.Stop(); err != nil {
		panic(err)
	}
}

func (s *steps) theLambdaIsInvokedWithValidCredentials() error {
	err := s.theLambdaIsInvokedWithCredentials("user", "password")
	return err
}

func (s *steps) aSuccessResponseIsReturned() error {
	s.assertResponse(true)
	return nil
}

func (s *steps) theLambdaIsInvokedWithInvalidCredentials() error {
	err := s.theLambdaIsInvokedWithCredentials("fgsgd", "dfgsd")
	return err
}

func (s *steps) aFailureResponseIsReturned() error {
	s.assertResponse(false)
	return nil
}

func (s *steps) theLambdaIsInvokedWithCredentials(user string, password string) error {
	localLambdaInvocationPort, err := s.containers.GetLocalHostLambdaPort()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://localhost:%d/2015-03-31/functions/myfunction/invocations", localLambdaInvocationPort)
	response, err := http.Post(url, "application/json",
		strings.NewReader("{\"headers\": {\"authorization\":\"Basic dXNlcjpwYXNzd29yZA==\"}}"))
	if err != nil {
		return err
	}
	s.responseStatus = response.StatusCode
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	s.responseBody = buf.String()

	return nil
}

func (s *steps) assertResponse(expectedAuthorised bool) error {
	assert.Equal(s.t, 200, s.responseStatus)
	simpleAuthorizerResponse := SimpleAuthorizerResponse{}
	err := json.Unmarshal([]byte(s.responseBody), simpleAuthorizerResponse)
	if err != nil {
		return err
	}
	assert.Equal(s.t, expectedAuthorised, simpleAuthorizerResponse.IsAuthorized)
	return nil
}
