package attendees_test

import (
	"attendee-writer/attendee"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDynamoClient struct {
	mock.Mock
}

func (dc *MockDynamoClient) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	args := dc.Called()
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (dc *MockDynamoClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := dc.Called()
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (d *MockDynamoClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := d.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func Test_shouldReturnAttendees(t *testing.T) {
	mockDynamoClient := MockDynamoClient{}
	mockDynamoClient.
		On("Scan", mock.Anything).
		Return(&dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{{"AuthCode": &types.AttributeValueMemberS{Value: "12345"}}},
		}, nil)

	datastore := attendee.AttendeesStore{
		Db:    &mockDynamoClient,
		Table: "some-table",
	}

	// When
	attendees, err := datastore.GetAllAttendees()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []attendee.Attendee{{
		AuthCode: "12345",
	}}, attendees)
}

func Test_shouldReturnNoAttendeesWhenUnableToScanDynamoDB(t *testing.T) {
	// Given
	mockDynamoClient := MockDynamoClient{}
	mockDynamoClient.
		On("Scan", mock.Anything).
		Return(&dynamodb.ScanOutput{}, fmt.Errorf("some dynamo error"))

	datastore := attendee.AttendeesStore{Db: &mockDynamoClient}

	// When
	response, err := datastore.GetAllAttendees()

	// Then
	assert.Equal(t, fmt.Errorf("fetching attendees from DynamoDB: some dynamo error"), err)
	assert.Nil(t, response)
}

func Test_shouldReturnNoAttendeesWhenThereAreNoneInTheDatastore(t *testing.T) {
	// Given
	mockDynamoClient := MockDynamoClient{}
	mockDynamoClient.
		On("Scan", mock.Anything).
		Return(&dynamodb.ScanOutput{}, nil)

	datastore := attendee.AttendeesStore{Db: &mockDynamoClient}

	// When
	response, err := datastore.GetAllAttendees()

	// Then
	assert.Nil(t, err)
	assert.Nil(t, response)
}

type SpyingDynamoClient struct {
	a *attendee.Attendee
}

func (s SpyingDynamoClient) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(input *dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return nil, nil
}

func (s SpyingDynamoClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return nil, nil
}

func (s SpyingDynamoClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	attributevalue.UnmarshalMap(params.Item, s.a)
	return nil, nil
}

func TestAttendees_ShouldPutItemInDynamoDbWhenStoreIsCalled(t *testing.T) {
	// Given
	var anAttendee attendee.Attendee
	store := attendee.AttendeesStore{Db: &SpyingDynamoClient{&anAttendee}, Table: "some-table"}

	// When
	err := store.Store(attendee.Attendee{AuthCode: "12345", Name: "Frank Spencer"})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "Frank Spencer", anAttendee.Name)
	assert.Equal(t, "12345", anAttendee.AuthCode)
}

func TestAttendees_ShouldReturnErrorIfUnableToPutItemInDynamoDB(t *testing.T) {
	// Given
	dynamoClient := MockDynamoClient{}
	dynamoClient.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, fmt.Errorf("some dynamo error"))

	store := attendee.AttendeesStore{Db: &dynamoClient, Table: "some-table"}

	// When
	err := store.Store(attendee.Attendee{})

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("putting attendee in datastore: some dynamo error"), err)
	dynamoClient.AssertNumberOfCalls(t, "PutItem", 1)
}
