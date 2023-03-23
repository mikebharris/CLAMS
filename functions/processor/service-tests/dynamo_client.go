package service_tests

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const workshopSignupsTableName = "workshops"

type DynamoClient struct {
	dynamoDbHandle *dynamodb.Client
}

type WorkshopSignup struct {
	WorkshopSignupId int
	WorkshopTitle    string
	Name             string
	Role             string
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

func (d DynamoClient) createWorkshopsSignupsTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("WorkshopSignupId"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("WorkshopSignupId"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(workshopSignupsTableName),
	}

	_, err := d.dynamoDbHandle.CreateTable(context.Background(), input)
	if err != nil {
		panic(err)
	}
}

func (d DynamoClient) getWorkshopSignup(signupId int) (*WorkshopSignup, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(workshopSignupsTableName),
		KeyConditionExpression: aws.String("WorkshopSignupId = :WorkshopSignupId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":WorkshopSignupId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", signupId)},
		},
	}

	q, err := d.dynamoDbHandle.Query(context.Background(), queryInput)
	if err != nil {
		return nil, err
	}

	var signups []WorkshopSignup
	if err := attributevalue.UnmarshalListOfMaps(q.Items, &signups); err != nil {
		return nil, err
	}

	if len(signups) == 0 {
		return nil, nil
	}

	return &signups[0], nil
}
