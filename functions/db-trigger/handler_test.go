package main

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
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

func Test_handler_Should_SendTriggerNotificationsToEventQueueAndDeleteThemFromDatabase(t *testing.T) {
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

	mock.ExpectExec("delete from trigger_notifications where id = 1")
	mock.ExpectExec("delete from trigger_notifications where id = 2")

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
