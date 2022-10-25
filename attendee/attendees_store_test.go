package attendee

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDatastore struct {
	mock.Mock
}

func (dc *MockDatastore) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	args := dc.Called()
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (dc *MockDatastore) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := dc.Called()
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func Test_shouldReturnAttendees(t *testing.T) {
	mockDynamoClient := MockDatastore{}
	mockDynamoClient.
		On("Scan", mock.Anything).
		Return(&dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{{"AuthCode": &types.AttributeValueMemberS{Value: "12345"}}},
		}, nil)

	datastore := AttendeesStore{
		Db:    &mockDynamoClient,
		Table: "some-table",
	}

	// When
	attendees, err := datastore.GetAllAttendees(context.Background())

	// Then
	assert.Nil(t, err)
	assert.Equal(t, &ApiResponse{Attendees: []Attendee{{
		AuthCode: "12345",
	}}}, attendees)
}

func Test_shouldReturnNoAttendeesWhenUnableToScanDynamoDB(t *testing.T) {
	// Given
	mockDynamoClient := MockDatastore{}
	mockDynamoClient.
		On("Scan", mock.Anything).
		Return(&dynamodb.ScanOutput{}, fmt.Errorf("some dynamo error"))

	datastore := AttendeesStore{Db: &mockDynamoClient}

	// When
	response, err := datastore.GetAllAttendees(context.Background())

	// Then
	assert.Equal(t, fmt.Errorf("fetching attendees from DynamoDB: some dynamo error"), err)
	assert.Nil(t, response)
}

func Test_shouldReturnNoAttendeesWhenThereAreNoneInTheDatastore(t *testing.T) {
	// Given
	mockDynamoClient := MockDatastore{}
	mockDynamoClient.
		On("Scan", mock.Anything).
		Return(&dynamodb.ScanOutput{}, nil)

	datastore := AttendeesStore{Db: &mockDynamoClient}

	// When
	response, err := datastore.GetAllAttendees(context.Background())

	// Then
	assert.Nil(t, err)
	assert.Nil(t, response)
}
