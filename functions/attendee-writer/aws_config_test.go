package main

import (
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
