package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDynamoClient struct {
	mock.Mock
}

func (d *MockDynamoClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := d.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func TestAttendees_ShouldPutItemInDynamoDbWhenStoreIsCalled(t *testing.T) {
	// Given
	ctx := context.Background()

	attendee := Attendee{}
	marshalMap, _ := attributevalue.MarshalMap(attendee)

	dynamoClient := MockDynamoClient{}
	dynamoClient.On("PutItem", ctx, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	store := AttendeesStore{Db: &dynamoClient, Table: "some-table"}

	// When
	err := store.Store(ctx, Attendee{})

	// Then
	assert.Nil(t, err)
	dynamoClient.AssertNumberOfCalls(t, "PutItem", 1)
	dynamoClient.AssertCalled(t, "PutItem", ctx, &dynamodb.PutItemInput{
		Item:      marshalMap,
		TableName: aws.String("some-table"),
	})
}

func TestAttendees_ShouldReturnErrorIfUnableToPutItemInDynamoDB(t *testing.T) {
	// Given
	ctx := context.Background()

	dynamoClient := MockDynamoClient{}
	dynamoClient.On("PutItem", ctx, mock.Anything).Return(&dynamodb.PutItemOutput{}, fmt.Errorf("some dynamo error"))

	store := AttendeesStore{Db: &dynamoClient, Table: "some-table"}

	// When
	err := store.Store(ctx, Attendee{})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("putting attendee in datastore: some dynamo error"), err)
	dynamoClient.AssertNumberOfCalls(t, "PutItem", 1)
}
