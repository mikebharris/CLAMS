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
	Code        string
	Name        string
	Email       string
	Phone       string
	Kids        uint
	Diet        string
	Financials  Financials
	Arrival     string
	Nights      uint
	StayingLate string
	CreatedTime time.Time
}

type Financials struct {
	ToPay    uint   `json:"ToPay"`
	Paid     uint   `json:"Paid"`
	Due      int    `json:"Due"`
	PaidDate string `json:"PaidDate"`
}

type IAttendees interface {
	FetchAttendee(ctx context.Context, code string) (*Attendee, error)
}

type Attendees struct {
	Db    attendeesDb
	Table string
}

func (r *Attendees) FetchAttendee(ctx context.Context, code string) (*Attendee, error) {
	record, err := r.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.Table),
		Key: map[string]types.AttributeValue{
			"Code": &types.AttributeValueMemberS{Value: code},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fetching attendee %s from DynamoDB: %v", code, err)
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
