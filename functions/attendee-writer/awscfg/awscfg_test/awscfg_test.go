package aws_config_test

import (
	"attendee-writer/awscfg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnDifferentConfigurationWhenDynamoDbEndpointIsOverridden(t *testing.T) {
	// Given
	defaultConfig := awscfg.GetAwsConfig("", "")

	// When
	configWithEndpointOverride := awscfg.GetAwsConfig("", "some-overridden-dynamodb-endpoint")

	// Then
	assert.NotEqual(t, defaultConfig, configWithEndpointOverride)
}

func TestShouldReturnDefaultConfigurationWhenDynamoDbEndpointIsNotOverridden(t *testing.T) {
	// Given
	// When
	configWithEndpointOverride := awscfg.GetAwsConfig("", "")

	// Then
	assert.NotNil(t, configWithEndpointOverride)
}
