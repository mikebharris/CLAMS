package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"time"
)

type attendeesDb interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
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
	Store(ctx context.Context, attendee Attendee) error
}

type Attendees struct {
	Db    attendeesDb
	Table string
}

func (a *Attendees) Store(ctx context.Context, attendee Attendee) error {
	marshalMap, err := attributevalue.MarshalMap(attendee)
	if err != nil {
		return err
	}
	_, err = a.Db.PutItem(ctx,
		&dynamodb.PutItemInput{
			Item:      marshalMap,
			TableName: aws.String(a.Table),
		})
	return err
}
