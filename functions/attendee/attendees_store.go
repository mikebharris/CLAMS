package attendee

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type IDatastore interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(input *dynamodb.Options)) (*dynamodb.ScanOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type AttendeesStore struct {
	Db    IDatastore
	Table string
}

func (as *AttendeesStore) GetAttendees(authCode string) ([]Attendee, error) {
	if authCode == "" {
		return as.getAllAttendees()
	} else {
		return as.getAttendeesWithAuthCode(authCode)
	}
}

func (as *AttendeesStore) getAttendeesWithAuthCode(authCode string) ([]Attendee, error) {
	record, err := as.Db.GetItem(context.Background(), &dynamodb.GetItemInput{
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

	var attendees []Attendee
	attendee := toAttendee(record.Item)
	if err != nil {
		return nil, fmt.Errorf("marshaling record %v to Attendee{} failed with error: %v", as, err)
	}
	attendees = append(attendees, attendee)
	return attendees, nil
}

func (as *AttendeesStore) getAllAttendees() ([]Attendee, error) {
	records, err := as.Db.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String(as.Table),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching attendees from DynamoDB: %v", err)
	}

	if len(records.Items) == 0 {
		return nil, nil
	}

	var attendees []Attendee
	for _, r := range records.Items {
		attendees = append(attendees, toAttendee(r))
	}

	return attendees, nil
}

func toAttendee(record map[string]types.AttributeValue) Attendee {
	var attendee Attendee
	if err := attributevalue.UnmarshalMap(record, &attendee); err != nil {
		return Attendee{}
	}
	return attendee
}
