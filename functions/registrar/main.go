package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
	"registrar/storage"
)

const (
	awsRegion = "us-east-1"
)

func main() {
	cfg := newConfig()

	lambdaHandler := handler{
		attendees: &storage.Attendees{
			Db:    dynamodb.NewFromConfig(cfg),
			Table: os.Getenv("ATTENDEES_TABLE_NAME"),
		},
	}

	lambda.Start(lambdaHandler.handleRequest)
}

func newConfig() aws.Config {
	dynamoEndpointOverride := os.Getenv("DYNAMO_ENDPOINT_OVERRIDE")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
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
