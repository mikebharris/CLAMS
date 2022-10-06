package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
)

func newAwsConfig() (*aws.Config, error) {
	dynamoEndpointOverride := os.Getenv("DYNAMO_ENDPOINT_OVERRIDE")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, fmt.Errorf("loading config: %v", err)
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
	return &cfg, nil
}
