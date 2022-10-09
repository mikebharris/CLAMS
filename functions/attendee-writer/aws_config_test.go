package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnDifferentConfigurationWhenDynamoDbEndpointIsOverridden(t *testing.T) {
	// Given
	defaultConfig, _ := newAwsConfig("us-east-1")
	endpointUrl := "some-overridden-dynamodb-endpoint"
	os.Setenv("DYNAMO_ENDPOINT_OVERRIDE", endpointUrl)

	// When
	configWithEndpointOverride, err := newAwsConfig("us-east-1")

	// Then
	assert.Nil(t, err)
	assert.NotEqual(t, defaultConfig, configWithEndpointOverride)
}

func TestShouldReturnDefaultConfigurationWhenDynamoDbEndpointIsNotOverridden(t *testing.T) {
	// Given
	// When
	configWithEndpointOverride, err := newAwsConfig("us-east-1")

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, configWithEndpointOverride)
}

func TestShouldReturnErrorWhenLoadingConfigurationReturnsError(t *testing.T) {
	// Given
	loadDefaultConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
		return aws.Config{}, fmt.Errorf("some error")
	}

	// When
	configWithEndpointOverride, err := newAwsConfig("us-east-1")

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("loading config: some error"), err)
	assert.Nil(t, configWithEndpointOverride)
}
