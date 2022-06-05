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

type AttendeesDataStore struct {
	Db    *dynamodb.Client
	Table string
}

func (a *AttendeesDataStore) FetchAllAttendees(ctx context.Context) ([]Attendee, error) {
	records, err := a.Db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(a.Table),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching attendees from DynamoDB: %v", err)
	}

	if len(records.Items) == 0 {
		return nil, nil
	}

	var attendees []Attendee
	for _, r := range records.Items {
		attendee, err := a.toAttendee(r)
		if err != nil {
			continue
		}
		attendees = append(attendees, attendee)
	}

	return attendees, nil
}

func (a *AttendeesDataStore) toAttendee(record map[string]types.AttributeValue) (Attendee, error) {
	var attendee Attendee
	if err := attributevalue.UnmarshalMap(record, &attendee); err != nil {
		return Attendee{}, fmt.Errorf("marshaling records %v to Attendee{} failed with error: %v", a, err)
	}
	return attendee, nil
}
