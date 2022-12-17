package main

import (
	"clams/awscfg"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"os"
	"regexp"
)

const extractionRegex = "^Basic ((?:[A-Za-z\\d+\\/]{4})*(?:[A-Za-z\\d+\\/]{3}=|[A-Za-z\\d+\\/]{2}==)?)$"
const userParameterNameTemplate = "/submissions/%s/metadata-api/user"
const passwordParameterNameTemplate = "/submissions/%s/metadata-api/password"

type SimpleAuthorizerResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

func handleRequest(event events.APIGatewayCustomAuthorizerRequestTypeRequest) (SimpleAuthorizerResponse, error) {
	isAuthorised := false
	authorizationHeader := event.Headers["authorization"]

	compiledRegex := regexp.MustCompile(extractionRegex)
	authMatches := compiledRegex.FindStringSubmatch(authorizationHeader)
	if len(authMatches) == 2 {
		ssmClient := newSsmClient()
		environment := os.Getenv("ENVIRONMENT")
		userParameterName := fmt.Sprintf(userParameterNameTemplate, environment)
		passwordParameterName := fmt.Sprintf(passwordParameterNameTemplate, environment)
		user := getParameterValue(ssmClient, userParameterName)
		password := getParameterValue(ssmClient, passwordParameterName)
		correctUserAndPassword := fmt.Sprintf("%s:%s", user, password)

		givenUserAndPassword := mustDecodeBase64ToString(authMatches[1])
		if givenUserAndPassword == correctUserAndPassword {
			isAuthorised = true
		}
	}
	return SimpleAuthorizerResponse{IsAuthorized: isAuthorised}, nil
}

func mustDecodeBase64ToString(base64Encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		panic(err)
	}
	return string(decoded)
}

func newSsmClient() *ssm.Client {
	awscfg.GetAwsConfig(ssm.ServiceID, os.Getenv("SSM_ENDPOINT_OVERRIDE"))
	return ssm.NewFromConfig(*awscfg.GetAwsConfig(ssm.ServiceID, os.Getenv("SSM_ENDPOINT_OVERRIDE")))
}

func getParameterValue(ssmClient *ssm.Client, parameterName string) string {
	response, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true)},
	)
	if err != nil {
		panic(err)
	}
	return *response.Parameter.Value
}

func main() {
	lambda.Start(handleRequest)
}
