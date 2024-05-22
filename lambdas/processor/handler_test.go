package main

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_handler_handleRequestReturnsErrorWhenRequestEmpty(t *testing.T) {
	// given
	handler := handler{}

	// when
	response, err := handler.handleRequest(context.Background(), events.SQSEvent{})

	// then
	assert.NotNil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure(nil)}, response)
}

func Test_handler_handleRequestPutsMessageInBatchItemFailuresWhenMessageHasNoBody(t *testing.T) {
	// given
	handler := handler{}
	records := []events.SQSMessage{{MessageId: "123"}}

	// when
	response, err := handler.handleRequest(context.Background(), events.SQSEvent{Records: records})

	// then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{{ItemIdentifier: "123"}}}, response)
}

func Test_handler_handleRequestPutsMessageInBatchItemFailuresWhenBodyIsNotWellFormed(t *testing.T) {
	// given
	handler := handler{}
	records := []events.SQSMessage{{MessageId: "123", Body: ""}}

	// when
	response, err := handler.handleRequest(context.Background(), events.SQSEvent{Records: records})

	// then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{{ItemIdentifier: "123"}}}, response)
}

func Test_handler_handleRequestPutsMessageInBatchItemFailuresWhenRecordDoesNotExistInDatabase(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening stub repository connexion: %v", err)
	}
	defer db.Close()
	emptyRow := sqlmock.NewRows([]string{"title", "name", "role_name"})
	mock.ExpectQuery("^select (.*)$").WithArgs(42).WillReturnRows(emptyRow)

	handler := handler{dbConx: db}

	m := message{WorkshopSignupId: 42}
	body, _ := json.Marshal(m)
	records := []events.SQSMessage{{MessageId: "123", Body: string(body)}}

	// when
	response, err := handler.handleRequest(context.Background(), events.SQSEvent{Records: records})

	// then
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{{ItemIdentifier: "123"}}}, response)
}

type SpyingDynDs struct {
	recordThatWasStored *WorkshopSignupRecord
}

func (s SpyingDynDs) Store(thing interface{}) error {
	record := thing.(WorkshopSignupRecord)
	*s.recordThatWasStored = record
	return nil
}

func Test_handler_handleRequestReturnsNoErrorAddsToDynamoAndWritesFileToS3WhenThereAreRecordsInTheDatabase(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("opening stub repository connexion: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"title", "name", "role_name"}).
		AddRow("Some Workshop Title", "Some Person Name", "Some role")
	sql := "select w.title, concat(p.forename, ' ', p.surname) as name, r.role_name " +
		"from people p join workshop_signups s on s.people_id = p.id " +
		"join workshops w on s.workshop_id = w.id " +
		"join roles r on s.role_id = r.id " +
		"where s.id = $1"
	mock.ExpectQuery(sql).WithArgs(42).WillReturnRows(rows)

	var recordThatWasStored WorkshopSignupRecord
	handler := handler{
		dbConx:    db,
		datastore: SpyingDynDs{recordThatWasStored: &recordThatWasStored},
	}

	body, _ := json.Marshal(message{WorkshopSignupId: 42})
	records := []events.SQSMessage{{MessageId: "123", Body: string(body)}}

	// When
	response, err := handler.handleRequest(context.Background(), events.SQSEvent{Records: records})

	// Then

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// no errors nor failed message deliveries
	assert.Nil(t, err)
	assert.Equal(t, events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure(nil)}, response)

	// record written to datastore
	assert.Equal(t, "Some Workshop Title", recordThatWasStored.WorkshopTitle)
	assert.Equal(t, "Some Person Name", recordThatWasStored.Name)
	assert.Equal(t, "Some role", recordThatWasStored.Role)
	assert.Equal(t, 42, recordThatWasStored.WorkshopSignupId)
}
