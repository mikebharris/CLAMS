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
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

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

type IAttendees interface {
	FetchAttendee(ctx context.Context, code string) (*ApiResponse, error)
	FetchAllAttendees(ctx context.Context) (*ApiResponse, error)
}

type Attendees struct {
	Db    attendeesDb
	Table string
}

func (a *Attendees) FetchAttendee(ctx context.Context, authCode string) (*ApiResponse, error) {
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
	attendee, err := a.toAttendee(record.Item)
	if err != nil {
		return nil, fmt.Errorf("marshaling record %v to Attendee{} failed with error: %v", a, err)
	}
	attendees.Attendees = append(attendees.Attendees, attendee)

	return &attendees, nil
}

func (a *Attendees) FetchAllAttendees(ctx context.Context) (*ApiResponse, error) {
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
		attendee, err := a.toAttendee(r)
		if err != nil {
			continue
		}
		attendees.Attendees = append(attendees.Attendees, attendee)
	}

	return &attendees, nil
}

func (a *Attendees) toAttendee(record map[string]types.AttributeValue) (Attendee, error) {
	var attendee Attendee
	if err := attributevalue.UnmarshalMap(record, &attendee); err != nil {
		return Attendee{}, fmt.Errorf("marshaling records %v to Attendee{} failed with error: %v", a, err)
	}
	return attendee, nil
}
