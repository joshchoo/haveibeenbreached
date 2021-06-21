package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Response events.APIGatewayProxyResponse

type AddAccountEvent struct {
	Accounts       []string
	PathParameters struct {
		BreachName string
	}
}

type Account struct {
	PK       string
	SK       string
	Type     string
	Account  string
	Breaches []string
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)
var tableName = "Breaches"
var entityType = "Account"

func Handler(ctx context.Context, event AddAccountEvent) (Response, error) {
	accounts := event.Accounts
	breachName := event.PathParameters.BreachName

	newAccount := Account{
		PK:       partitionKey("@doe.com"),
		SK:       sortKey("john"),
		Type:     entityType,
		Account:  "john@doe.com",
		Breaches: []string{breachName},
	}

	attrVal, err := dynamodbattribute.MarshalMap(newAccount)
	if err != nil {
		return Response{StatusCode: 400, Body: fmt.Sprintf("Error marshalling new Account: %s", err)}, err
	}
	input := &dynamodb.PutItemInput{
		Item:      attrVal,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		return Response{StatusCode: 400, Body: fmt.Sprintf("Error adding Account to breach: %s", err)}, err
	}

	numAccounts := len(accounts)
	body, err := json.Marshal(map[string]interface{}{
		"message": fmt.Sprintf("Successfully added %d accounts to the %s breach.", numAccounts, breachName),
	})
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func partitionKey(key string) string {
	return fmt.Sprintf("EMAIL#%s", key)
}

func sortKey(key string) string {
	return fmt.Sprintf("EMAIL#%s", key)
}
