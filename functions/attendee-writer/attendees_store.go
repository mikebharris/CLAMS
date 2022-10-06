package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type IDynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type AttendeesStore struct {
	Db    IDynamoClient
	Table string
}

func (a AttendeesStore) Store(ctx context.Context, attendee Attendee) error {
	marshalMap, _ := attributevalue.MarshalMap(attendee)
	_, err := a.Db.PutItem(ctx,
		&dynamodb.PutItemInput{
			Item:      marshalMap,
			TableName: aws.String(a.Table),
		})

	if err != nil {
		return fmt.Errorf("putting attendee in datastore: %v", err)
	}
	return nil
}
