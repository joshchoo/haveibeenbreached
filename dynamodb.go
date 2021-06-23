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

type Itemable interface {
	ToItem() DBItem
}

type Repo struct {
	db *dynamodb.DynamoDB
}

func NewRepo(db *dynamodb.DynamoDB) Repo {
	return Repo{db}
}

func (r Repo) GetAccount(username Username) (*Account, error) {
	pk := username.PartitionKey()
	sk := username.SortKey()
	accountItem := AccountItem{}
	found, err := r.getItem(pk, sk, &accountItem)
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

func (r Repo) GetBreach(breachName string) (*Breach, error) {
	pk := BreachPartitionKey(breachName)
	sk := BreachSortKey(breachName)
	breachItem := BreachItem{}
	found, err := r.getItem(pk, sk, &breachItem)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	breach := breachItem.ToBreach()
	return &breach, nil
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
	result, err := r.db.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item != nil {
		err := dynamodbattribute.UnmarshalMap(result.Item, output)
		return true, err
	}
	return false, nil
}

func (r Repo) PutItem(itemable Itemable) error {
	attrVal, err := dynamodbattribute.MarshalMap(itemable.ToItem())
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      attrVal,
		TableName: aws.String(tableName),
	}
	_, err = r.db.PutItem(input)
	return err
}

func (r Repo) PutItems(items []Itemable) error {
	for _, item := range items {
		if err := r.PutItem(item); err != nil {
			return err
		}
	}
	return nil
}
