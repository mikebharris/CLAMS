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
	Code        string
	Name        string
	Email       string
	Phone       string
	Kids        uint
	Diet        string
	Financials  Financials
	Arrival     string
	Nights      uint
	StayingLate string
	CreatedTime time.Time
}

type Financials struct {
	ToPay    uint   `json:"To Pay"`
	Paid     uint   `json:"Paid"`
	PaidDate string `json:"Paid date"`
	Due      int    `json:"Due"`
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
				AttributeName: aws.String("Code"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Code"),
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

func (d DynamoClient) getAttendeeByCode(code string) (*Attendee, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(attendeesTableName),
		KeyConditionExpression: aws.String("Code = :code"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":code": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s", code)},
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
