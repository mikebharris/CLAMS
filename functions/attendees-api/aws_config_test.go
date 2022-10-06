package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnDifferentConfigurationWhenDynamoDbEndpointIsOverridden(t *testing.T) {
	// Given
	defaultConfig, _ := newAwsConfig("")
	endpointUrl := "some-overridden-dynamodb-endpoint"
	os.Setenv("DYNAMO_ENDPOINT_OVERRIDE", endpointUrl)

	// When
	configWithEndpointOverride, err := newAwsConfig("")

	// Then
	assert.Nil(t, err)
	assert.NotEqual(t, defaultConfig, configWithEndpointOverride)
}
