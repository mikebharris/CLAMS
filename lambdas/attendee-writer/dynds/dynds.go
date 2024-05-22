package dynds

import (
	"clams/awscfg"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDatastore struct {
	Table    string
	Endpoint string
	Region   string
	dbClient *dynamodb.Client
}

func (d *DynamoDatastore) Init() {
	awsConfig := awscfg.GetAwsConfig(dynamodb.ServiceID, d.Endpoint, d.Region)
	d.dbClient = dynamodb.NewFromConfig(*awsConfig)
}

func (d *DynamoDatastore) Store(thing interface{}) error {
	marshalMap, _ := attributevalue.MarshalMap(thing)
	_, err := d.dbClient.PutItem(context.Background(),
		&dynamodb.PutItemInput{
			Item:      marshalMap,
			TableName: aws.String(d.Table),
		})

	if err != nil {
		return fmt.Errorf("putting thing in datastore: %v", err)
	}
	return nil
}
