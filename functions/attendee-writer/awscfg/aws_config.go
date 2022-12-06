package awscfg

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	awsRegion = "us-east-1"
)

var loadDefaultConfig = config.LoadDefaultConfig

func awsConfigForBespokeServiceEndpoint(awsRegion string, awsService string, endpoint string) (*aws.Config, error) {
	cfg, err := loadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, fmt.Errorf("loading config: %v", err)
	}

	cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == awsService {
			return aws.Endpoint{URL: endpoint}, nil
		}
		return cfg.EndpointResolverWithOptions.ResolveEndpoint(service, region)
	})
	return &cfg, nil
}

func GetAwsConfig(id string, endpoint string) *aws.Config {
	var awsConfig *aws.Config
	if len(endpoint) == 0 {
		cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
		awsConfig = &cfg
	} else {
		awsConfig, _ = awsConfigForBespokeServiceEndpoint(awsRegion, id, endpoint)
	}
	return awsConfig
}
