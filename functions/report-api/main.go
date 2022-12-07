package main

import (
	"clams/attendee"
	"clams/report-api/handler"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
)

func main() {
	cfg := newConfig()

	lambdaHandler := handler.Handler{
		AttendeesStore: &attendee.AttendeesStore{
			Db:    dynamodb.NewFromConfig(cfg),
			Table: os.Getenv("ATTENDEES_TABLE_NAME"),
		},
	}
	lambda.Start(lambdaHandler.HandleRequest)
}

func newConfig() aws.Config {
	dynamoEndpointOverride := os.Getenv("DYNAMO_ENDPOINT_OVERRIDE")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	if len(dynamoEndpointOverride) > 0 {
		defaultEndpointResolver := cfg.EndpointResolverWithOptions
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID && len(dynamoEndpointOverride) > 0 {
				return aws.Endpoint{URL: dynamoEndpointOverride}, nil
			}
			return defaultEndpointResolver.ResolveEndpoint(service, region)
		})
	}

	return cfg
}
