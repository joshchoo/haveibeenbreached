package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

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

	email, err := NewEmail("john@doe.com")
	if err != nil {
		return Response{StatusCode: 400, Body: fmt.Sprintf("Invalid account: %s", err)}, err
	}

	newAccount := Account{
		PK:       email.PartitionKey(),
		SK:       email.SortKey(),
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

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Email struct {
	Domain string
	Alias  string
}

func NewEmail(emailStr string) (Email, error) {
	if !emailRegex.MatchString(emailStr) {
		return Email{}, fmt.Errorf("not a valid email address: %s", emailStr)
	}
	email := strings.Split(emailStr, "@")
	return Email{
		Alias:  email[0],
		Domain: email[1],
	}, nil
}

func (e Email) PartitionKey() string {
	return fmt.Sprintf("EMAIL#%s", e.Domain)
}

func (e Email) SortKey() string {
	return fmt.Sprintf("EMAIL#%s", e.Alias)
}
