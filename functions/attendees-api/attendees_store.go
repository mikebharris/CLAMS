package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ApiResponse struct {
	Attendees []Attendee `json:"Attendees"`
}

type DatastoreInterface interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(input *dynamodb.Options)) (*dynamodb.ScanOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type AttendeesStore struct {
	Db    DatastoreInterface
	Table string
}

func (as *AttendeesStore) GetAttendeesWithAuthCode(ctx context.Context, authCode string) (*ApiResponse, error) {
	record, err := as.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(as.Table),
		Key: map[string]types.AttributeValue{
			"AuthCode": &types.AttributeValueMemberS{Value: authCode},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fetching attendee %s from DynamoDB: %v", authCode, err)
	}

	if record.Item == nil {
		return nil, nil
	}

	var attendees ApiResponse
	attendee := as.toAttendee(record.Item)
	if err != nil {
		return nil, fmt.Errorf("marshaling record %v to Attendee{} failed with error: %v", as, err)
	}
	attendees.Attendees = append(attendees.Attendees, attendee)

	return &attendees, nil
}

func (as *AttendeesStore) GetAllAttendees(ctx context.Context) (*ApiResponse, error) {
	records, err := as.Db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(as.Table),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching attendees from DynamoDB: %v", err)
	}

	if len(records.Items) == 0 {
		return nil, nil
	}

	var attendees ApiResponse
	for _, r := range records.Items {
		attendees.Attendees = append(attendees.Attendees, as.toAttendee(r))
	}

	return &attendees, nil
}

func (as *AttendeesStore) toAttendee(record map[string]types.AttributeValue) Attendee {
	var attendee Attendee
	if err := attributevalue.UnmarshalMap(record, &attendee); err != nil {
		return Attendee{}
	}
	return attendee
}
