package haveibeenbreached

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tableName = "Breaches"

type DBItem interface {
	isDBItem() bool
}

type Repo struct {
	svc *dynamodb.DynamoDB
}

func NewRepo(svc *dynamodb.DynamoDB) Repo {
	return Repo{svc}
}

func (r Repo) PutItem(item DBItem) error {
	attrVal, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      attrVal,
		TableName: aws.String(tableName),
	}
	_, err = r.svc.PutItem(input)
	return err
}

func (r Repo) PutItems(items []DBItem) error {
	for _, item := range items {
		if err := r.PutItem(item); err != nil {
			return err
		}
	}
	return nil
}