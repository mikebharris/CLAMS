package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

type Attendees struct {
	Db    *dynamodb.Client
	Table string
}

func (a *Attendees) Store(ctx context.Context, attendee Attendee) error {
	marshalMap, _ := attributevalue.MarshalMap(attendee)
	_, err := a.Db.PutItem(ctx,
		&dynamodb.PutItemInput{
			Item:      marshalMap,
			TableName: aws.String(a.Table),
		})
	return err
}
