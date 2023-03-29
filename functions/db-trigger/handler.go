package main

import (
	"database/sql"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"golang.org/x/net/context"
)

type ISqsClient interface {
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type handler struct {
	dbConx     *sql.DB
	sqsService ISqsClient
}

type TriggerNotification struct {
	WorkshopSignupId int
	WorkshopId       int
	PeopleId         int
	RoleId           int
}

func (h handler) handleRequest(_ context.Context) (events.LambdaFunctionURLResponse, error) {
	queueName := "db-trigger-queue"
	queueUrl, _ := h.sqsService.GetQueueUrl(
		context.Background(),
		&sqs.GetQueueUrlInput{
			QueueName: &queueName,
		},
	)

	r := repository{dbConx: h.dbConx}
	notifications := r.getTriggerNotifications()

	for _, n := range notifications {
		h.sqsService.SendMessage(
			context.Background(),
			&sqs.SendMessageInput{
				QueueUrl:    aws.String(*queueUrl.QueueUrl),
				MessageBody: aws.String(n),
			},
		)
	}

	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}
