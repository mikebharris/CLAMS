package service_tests

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsClient struct {
	sqsHandle *sqs.Client
}

func newSqsClient(host string, port int) SqsClient {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	endpoint := fmt.Sprintf("http://%s:%d", host, port)
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{URL: endpoint, SigningRegion: "us-east-1"}, nil
	})

	client := sqs.NewFromConfig(cfg)

	return SqsClient{sqsHandle: client}
}

func (s *SqsClient) getQueueUrl() *string {
	output, err := s.sqsHandle.GetQueueUrl(
		context.Background(),
		&sqs.GetQueueUrlInput{
			QueueName: aws.String("db-trigger-queue"),
		},
	)
	if err != nil {
		panic(err)
	}

	return output.QueueUrl
}

func (s *SqsClient) getMessages() []types.Message {
	output, err := s.sqsHandle.ReceiveMessage(
		context.Background(),
		&sqs.ReceiveMessageInput{
			QueueUrl:            s.getQueueUrl(),
			MaxNumberOfMessages: 10,
		},
	)
	if err != nil {
		panic(err)
	}

	return output.Messages
}
