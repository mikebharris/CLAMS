package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type attendeesDb interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type Attendee struct {
	AuthCode       string
	Name           string
	Email          string
	Telephone      string
	NumberOfKids   uint
	Diet           string
	Financials     Financials
	ArrivalDay     string
	NumberOfNights uint
	StayingLate    string
	CreatedTime    time.Time
}

type Financials struct {
	AmountToPay uint   `json:"AmountToPay"`
	AmountPaid  uint   `json:"AmountPaid"`
	AmountDue   int    `json:"AmountDue"`
	DatePaid    string `json:"DatePaid"`
}

type IAttendees interface {
	FetchAttendee(ctx context.Context, code string) (*Attendee, error)
}

type Attendees struct {
	Db    attendeesDb
	Table string
}

func (r *Attendees) FetchAttendee(ctx context.Context, authCode string) (*Attendee, error) {
	record, err := r.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.Table),
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

	var attendee Attendee
	err = attributevalue.UnmarshalMap(record.Item, &attendee)
	if err != nil {
		return nil, fmt.Errorf("marshaling record %v to Attendee{} failed with error: %v", record.Item, err)
	}

	return &attendee, nil
}
