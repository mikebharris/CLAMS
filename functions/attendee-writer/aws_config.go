package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
)

var loadDefaultConfig = config.LoadDefaultConfig

func newAwsConfig(awsRegion string) (*aws.Config, error) {
	dynamoEndpointOverride := os.Getenv("DYNAMO_ENDPOINT_OVERRIDE")

	cfg, err := loadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, fmt.Errorf("loading config: %v", err)
	}

	if len(dynamoEndpointOverride) == 0 {
		return &cfg, nil
	}

	cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{URL: dynamoEndpointOverride}, nil
		}
		return cfg.EndpointResolver.ResolveEndpoint(service, region)
	})
	return &cfg, nil
}
