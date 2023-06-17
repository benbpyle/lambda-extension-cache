package clients

import (
	"cache-layer/models"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

// NewDynamoDBClient inits a DynamoDB session to be used throughout the services
func NewDynamoDBClient() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		log.Fatal(err)
	}

	return dynamodb.New(sess)
}

type DbRepository struct {
	client *dynamodb.DynamoDB
}

func NewDbRepository(client *dynamodb.DynamoDB) *DbRepository {
	return &DbRepository{
		client: client,
	}
}

func (d *DbRepository) ReadItem(ctx context.Context, id string) (*models.Model, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("CacheSample"),
	}

	result, err := d.client.GetItemWithContext(ctx, input)

	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, err
	}

	up := &models.Model{}
	err = dynamodbattribute.UnmarshalMap(result.Item, up)
	return up, err
}
