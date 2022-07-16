package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type ApiResponse struct {
	Attendees []Attendee `json:"Attendees"`
}

type Attendee struct {
	AuthCode       string
	Name           string
	Email          string
	Telephone      string
	NumberOfKids   int
	Diet           string
	Financials     Financials
	ArrivalDay     string
	NumberOfNights int
	StayingLate    string
	CreatedTime    time.Time
}

type Financials struct {
	AmountToPay int    `json:"AmountToPay"`
	AmountPaid  int    `json:"AmountPaid"`
	AmountDue   int    `json:"AmountDue"`
	DatePaid    string `json:"DatePaid"`
}

type DatastoreInterface interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(input *dynamodb.Options)) (*dynamodb.ScanOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type Register struct {
	Db    DatastoreInterface
	Table string
}

func (a *Register) AttendeesWithAuthCode(ctx context.Context, authCode string) (*ApiResponse, error) {
	record, err := a.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(a.Table),
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
	attendee := a.toAttendee(record.Item)
	if err != nil {
		return nil, fmt.Errorf("marshaling record %v to Attendee{} failed with error: %v", a, err)
	}
	attendees.Attendees = append(attendees.Attendees, attendee)

	return &attendees, nil
}

func (a *Register) Attendees(ctx context.Context) (*ApiResponse, error) {
	records, err := a.Db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(a.Table),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching attendees from DynamoDB: %v", err)
	}

	if len(records.Items) == 0 {
		return nil, nil
	}

	var attendees ApiResponse
	for _, r := range records.Items {
		attendees.Attendees = append(attendees.Attendees, a.toAttendee(r))
	}

	return &attendees, nil
}

func (a *Register) toAttendee(record map[string]types.AttributeValue) Attendee {
	var attendee Attendee
	if err := attributevalue.UnmarshalMap(record, &attendee); err != nil {
		return Attendee{}
	}
	return attendee
}
