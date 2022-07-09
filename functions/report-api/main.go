package main

import (
	"context"
	"os"
	"report-api/storage"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg := newConfig()

	lambdaHandler := Handler{
		attendeesDatastore: &storage.DynamoAttendeesDataStore{
			Db:    dynamodb.NewFromConfig(cfg),
			Table: os.Getenv("ATTENDEES_TABLE_NAME"),
		},
	}
	lambda.Start(lambdaHandler.Handle)
}

func newConfig() aws.Config {
	dynamoEndpointOverride := os.Getenv("DYNAMO_ENDPOINT_OVERRIDE")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	if len(dynamoEndpointOverride) > 0 {
		defaultEndpointResolver := cfg.EndpointResolver
		cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID && len(dynamoEndpointOverride) > 0 {
				return aws.Endpoint{URL: dynamoEndpointOverride}, nil
			}
			return defaultEndpointResolver.ResolveEndpoint(service, region)
		})
	}

	return cfg
}
