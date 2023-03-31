package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	context2 "golang.org/x/net/context"
	"testing"
)

type message struct {
	WorkshopSignupId int
	WorkshopId       int
	PeopleId         int
	RoleId           int
}

type spyingSqsClient struct {
	messagesThatWereSent *[]message
	queueThatWasSentTo   *string
}

func (s spyingSqsClient) GetQueueUrl(_ context.Context, params *sqs.GetQueueUrlInput, _ ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
	*s.queueThatWasSentTo = *params.QueueName
	url := "a-queue-url"
	return &sqs.GetQueueUrlOutput{QueueUrl: &url}, nil
}

func (s spyingSqsClient) SendMessage(_ context.Context, params *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	var msg message
	json.Unmarshal([]byte(*params.MessageBody), &msg)
	*s.messagesThatWereSent = append(*s.messagesThatWereSent, msg)
	return nil, nil
}

func Test_handler_ShouldSendTriggerNotificationsToEventQueueAndDeleteThemFromDatabase(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("opening stub repository connexion: %v", err)
	}
	defer db.Close()

	mock.MatchExpectationsInOrder(true)

	rows := sqlmock.NewRows([]string{"id", "message"}).
		AddRow(1, "{\"WorkshopSignupId\":1, \"WorkshopId\":3, \"PeopleId\":3, \"RoleId\":7}").
		AddRow(2, "{\"WorkshopSignupId\":1, \"WorkshopId\":8, \"PeopleId\":6, \"RoleId\":1}")

	mock.ExpectQuery("select id, message from trigger_notifications").WillReturnRows(rows)

	mock.ExpectExec("delete from trigger_notifications where id = 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from trigger_notifications where id = 2").WillReturnResult(sqlmock.NewResult(1, 1))

	var messagesReceived []message
	var queueThatWasSentTo string
	handler := handler{dbConx: db, sqsService: &spyingSqsClient{queueThatWasSentTo: &queueThatWasSentTo, messagesThatWereSent: &messagesReceived}}

	// When
	response, err := handler.handleRequest(context.Background())

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: 200}, response)
	assert.Equal(t, "db-trigger-queue", queueThatWasSentTo)
	assert.Len(t, messagesReceived, 2)
	assert.Equal(t, message{
		WorkshopSignupId: 1,
		WorkshopId:       3,
		PeopleId:         3,
		RoleId:           7,
	}, messagesReceived[0])
	assert.Equal(t, message{
		WorkshopSignupId: 1,
		WorkshopId:       8,
		PeopleId:         6,
		RoleId:           1,
	}, messagesReceived[1])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_handler_ShouldReturnErrorWhenUnableToFetchTriggerNotifications(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("opening stub repository connexion: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("select id, message from trigger_notifications").WillReturnError(errors.New("some db error"))

	var messagesReceived []message
	var queueThatWasSentTo string
	handler := handler{dbConx: db, sqsService: &spyingSqsClient{queueThatWasSentTo: &queueThatWasSentTo, messagesThatWereSent: &messagesReceived}}

	// When
	response, err := handler.handleRequest(context.Background())

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, "fetching trigger notifications: some db error", err.Error())
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: 500}, response)
	assert.Len(t, messagesReceived, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

type sqsClientOneCantSendMessagesTo struct{}

func (s sqsClientOneCantSendMessagesTo) GetQueueUrl(_ context.Context, params *sqs.GetQueueUrlInput, _ ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
	url := "a-queue-url"
	return &sqs.GetQueueUrlOutput{QueueUrl: &url}, nil
}

func (s sqsClientOneCantSendMessagesTo) SendMessage(ctx context2.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return &sqs.SendMessageOutput{}, errors.New("some aws error")
}

func Test_handler_ShouldReturnErrorWhenUnableToPostEventsToQueue(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("opening stub repository connexion: %v", err)
	}
	defer db.Close()

	mock.MatchExpectationsInOrder(true)

	rows := sqlmock.NewRows([]string{"id", "message"}).
		AddRow(1, "{\"WorkshopSignupId\":1, \"WorkshopId\":3, \"PeopleId\":3, \"RoleId\":7}")
	mock.ExpectQuery("select id, message from trigger_notifications").WillReturnRows(rows)

	handler := handler{dbConx: db, sqsService: &sqsClientOneCantSendMessagesTo{}}

	// When
	response, err := handler.handleRequest(context.Background())

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, "sending message to queue db-trigger-queue: some aws error", err.Error())
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: 500}, response)
	assert.NoError(t, mock.ExpectationsWereMet())
}

type sqsClientOneCantGetAQueueUrlFor struct{}

func (s sqsClientOneCantGetAQueueUrlFor) GetQueueUrl(_ context.Context, params *sqs.GetQueueUrlInput, _ ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
	return &sqs.GetQueueUrlOutput{}, errors.New("some aws error")
}

func (s sqsClientOneCantGetAQueueUrlFor) SendMessage(ctx context2.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	panic("not implemented")
}

func Test_handler_ShouldReturnErrorWhenUnableToGetQueueUrl(t *testing.T) {
	// Given
	handler := handler{sqsService: &sqsClientOneCantGetAQueueUrlFor{}}

	// When
	response, err := handler.handleRequest(context.Background())

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, "getting queue url for db-trigger-queue: some aws error", err.Error())
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: 500}, response)
}

func Test_handler_ShouldReturnErrorWhenUnableToDeleteRow(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("opening stub repository connexion: %v", err)
	}
	defer db.Close()

	mock.MatchExpectationsInOrder(true)

	rows := sqlmock.NewRows([]string{"id", "message"}).
		AddRow(1, "{\"WorkshopSignupId\":1, \"WorkshopId\":3, \"PeopleId\":3, \"RoleId\":7}")

	mock.ExpectQuery("select id, message from trigger_notifications").WillReturnRows(rows)
	mock.ExpectExec("delete from trigger_notifications where id = 1").WillReturnError(errors.New("some db error"))

	var messagesReceived []message
	var queueThatWasSentTo string
	handler := handler{dbConx: db, sqsService: &spyingSqsClient{queueThatWasSentTo: &queueThatWasSentTo, messagesThatWereSent: &messagesReceived}}

	// When
	response, err := handler.handleRequest(context.Background())

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, "deleting trigger notification 1: some db error", err.Error())
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: 500}, response)
	assert.Len(t, messagesReceived, 1)
	assert.Equal(t, message{
		WorkshopSignupId: 1,
		WorkshopId:       3,
		PeopleId:         3,
		RoleId:           7,
	}, messagesReceived[0])

	assert.NoError(t, mock.ExpectationsWereMet())
}
