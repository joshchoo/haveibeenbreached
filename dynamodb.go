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

func (r Repo) GetAccount(partitionKey, sortKey string) (*Account, error) {
	accountItem := AccountItem{}
	found, err := r.getItem(partitionKey, sortKey, &accountItem)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	account, err := accountItem.ToAccount()
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r Repo) getItem(partitionKey string, sortKey string, output interface{}) (found bool, err error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(partitionKey),
			},
			"SK": {
				S: aws.String(sortKey),
			},
		},
	}
	result, err := r.svc.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item != nil {
		err := dynamodbattribute.UnmarshalMap(result.Item, output)
		return true, err
	}
	return false, nil
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
