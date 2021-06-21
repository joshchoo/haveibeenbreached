package haveibeenbreached

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tableName = "Breaches"

type Repo struct {
	svc *dynamodb.DynamoDB
}

func NewRepo(svc *dynamodb.DynamoDB) Repo {
	return Repo{svc}
}

func (r Repo) PutItem(item interface{}) error {
	attrVal, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil
	}
	input := &dynamodb.PutItemInput{
		Item:      attrVal,
		TableName: aws.String(tableName),
	}
	_, err = r.svc.PutItem(input)
	return err
}
