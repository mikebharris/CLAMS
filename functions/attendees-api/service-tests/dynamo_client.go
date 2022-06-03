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
	NumberOfKids   uint
	Diet           string
	Financials     Financials
	ArrivalDay     string
	NumberOfNights uint
	StayingLate    string
	CreatedTime    time.Time
}

type Financials struct {
	AmountToPay uint   `json:"To Pay"`
	AmountPaid  uint   `json:"AmountPaid"`
	DatePaid    string `json:"DatePaid"`
	AmountDue   int    `json:"AmountDue"`
}

func newDynamoClient(host string, port int) (DynamoClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return DynamoClient{}, err
	}

	endpoint := fmt.Sprintf("http://%s:%d", host, port)
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{URL: endpoint}, nil
	})

	return DynamoClient{dynamodb.NewFromConfig(cfg)}, nil
}

func (d DynamoClient) createAttendeesTable() error {
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
	return err
}

func (d DynamoClient) addAttendee(entry Attendee) error {
	marshalMap, err := attributevalue.MarshalMap(entry)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      marshalMap,
		TableName: aws.String(attendeesTableName),
	}
	_, err = d.dynamoDbHandle.PutItem(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}
