package service_tests

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const attendeesTableName = "attendees"

type DynamoClient struct {
	dynamoDbHandle *dynamodb.Client
}

type Attendee struct {
	AuthCode       string
	Name           string
	Email          string
	Telephone      string
	NumberOfKids   int
	Diet           string
	Financials     Financials
	ArrivalDay     string
	NumberOfNights int
	StayingLate    string
	CreatedTime    time.Time
}

type Financials struct {
	AmountToPay int    `json:"AmountToPay"`
	AmountPaid  int    `json:"AmountPaid"`
	AmountDue   int    `json:"AmountDue"`
	DatePaid    string `json:"DatePaid"`
}

func newDynamoClient(host string, port int) DynamoClient {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	endpoint := fmt.Sprintf("http://%s:%d", host, port)
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{URL: endpoint}, nil
	})

	return DynamoClient{dynamodb.NewFromConfig(cfg)}
}

func (d DynamoClient) createAttendeesTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("AuthCode"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("AuthCode"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(attendeesTableName),
	}

	_, err := d.dynamoDbHandle.CreateTable(context.Background(), input)
	if err != nil {
		panic(err)
	}
}

func (d DynamoClient) getAttendeeByCode(authCode string) (*Attendee, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(attendeesTableName),
		KeyConditionExpression: aws.String("AuthCode = :authCode"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":authCode": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s", authCode)},
		},
	}

	q, err := d.dynamoDbHandle.Query(context.Background(), queryInput)
	if err != nil {
		return nil, err
	}

	var attendees []Attendee
	if err := attributevalue.UnmarshalListOfMaps(q.Items, &attendees); err != nil {
		return nil, err
	}

	return &attendees[0], nil
}
