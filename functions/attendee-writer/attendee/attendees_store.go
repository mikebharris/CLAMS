package attendee

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"time"
)

type AttendeesStore struct {
	Db    *dynamodb.Client
	Table string
}

func (as *AttendeesStore) Store(attendee Attendee) error {
	attendee.CreatedTime = time.Now()
	marshalMap, _ := attributevalue.MarshalMap(attendee)
	_, err := as.Db.PutItem(context.Background(),
		&dynamodb.PutItemInput{
			Item:      marshalMap,
			TableName: aws.String(as.Table),
		})

	if err != nil {
		return fmt.Errorf("putting attendee in datastore: %v", err)
	}
	return nil
}
