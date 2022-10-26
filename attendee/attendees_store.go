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

func (as *AttendeesStore) GetAttendeesWithAuthCode(ctx context.Context, authCode string) ([]Attendee, error) {
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

	var attendees []Attendee
	attendee := as.toAttendee(record.Item)
	if err != nil {
		return nil, fmt.Errorf("marshaling record %v to Attendee{} failed with error: %v", as, err)
	}
	attendees = append(attendees, attendee)
	return attendees, nil
}

func (as *AttendeesStore) GetAllAttendees(ctx context.Context) ([]Attendee, error) {
	records, err := as.Db.Scan(ctx, &dynamodb.ScanInput{
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
		attendees = append(attendees, as.toAttendee(r))
	}

	return attendees, nil
}

func (as *AttendeesStore) toAttendee(record map[string]types.AttributeValue) Attendee {
	var attendee Attendee
	if err := attributevalue.UnmarshalMap(record, &attendee); err != nil {
		return Attendee{}
	}
	return attendee
}

func (as *AttendeesStore) Store(ctx context.Context, attendee Attendee) error {
	marshalMap, _ := attributevalue.MarshalMap(attendee)
	_, err := as.Db.PutItem(ctx,
		&dynamodb.PutItemInput{
			Item:      marshalMap,
			TableName: aws.String(as.Table),
		})

	if err != nil {
		return fmt.Errorf("putting attendee in datastore: %v", err)
	}
	return nil
}
