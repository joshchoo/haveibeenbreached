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
	"github.com/joshuous/haveibeenbreached"
)

type Response events.APIGatewayProxyResponse

type AddAccountEvent struct {
	Accounts       []string
	PathParameters struct {
		BreachName string
	}
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)
var tableName = "Breaches"
var repo = haveibeenbreached.NewRepo(svc)
var addAccountsToBreachHandler = makeAddAccountsToBreachHandler(repo)

func makeAddAccountsToBreachHandler(repo haveibeenbreached.Repo) func(ctx context.Context, event AddAccountEvent) (Response, error) {
	return func(ctx context.Context, event AddAccountEvent) (Response, error) {
		rawAccounts := event.Accounts
		breachName := event.PathParameters.BreachName

		accounts, err := mapToAccount(rawAccounts)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Invalid email: %s", err)}, err
		}

		accounts, err = setAccountBreaches(accounts, breachName)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Invalid email: %s", err)}, err
		}

		items := make([]haveibeenbreached.DBItem, 0, len(accounts))
		for _, a := range accounts {
			items = append(items, a.Item())
		}

		if err = repo.PutItems(items); err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error adding accounts to breach: %s", err)}, err
		}

		numAccounts := len(accounts)
		body, err := json.Marshal(map[string]interface{}{
			"message": fmt.Sprintf("Successfully added/updated %d accounts to the %s breach.", numAccounts, breachName),
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
}

func main() {
	lambda.Start(addAccountsToBreachHandler)
}

func mapToAccount(accounts []string) ([]haveibeenbreached.Account, error) {
	accs := make([]haveibeenbreached.Account, 0, len(accounts))
	var err error

	for _, account := range accounts {
		email, emailErr := haveibeenbreached.NewEmail(account)
		if err != nil {
			err = emailErr
			break
		}
		newAccount := haveibeenbreached.Account{
			Username: email,
			Breaches: make([]string, 0),
		}
		accs = append(accs, newAccount)
	}
	if err != nil {
		return []haveibeenbreached.Account{}, err
	}
	return accs, nil
}

func setAccountBreaches(accounts []haveibeenbreached.Account, breachName string) ([]haveibeenbreached.Account, error) {
	accs := make([]haveibeenbreached.Account, 0, len(accounts))
	var err error

	for _, account := range accounts {
		accountItem := account.Item()
		input := &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"PK": {
					S: aws.String(accountItem.PK),
				},
				"SK": {
					S: aws.String(accountItem.SK),
				},
			},
		}
		result, getItemErr := svc.GetItem(input)
		if getItemErr != nil {
			err = getItemErr
			break
		}
		if result.Item != nil {
			existingAcc := &haveibeenbreached.AccountItem{}
			unmarshalErr := dynamodbattribute.UnmarshalMap(result.Item, existingAcc)
			if unmarshalErr != nil {
				err = unmarshalErr
				break
			}
			var breaches []string
			if contains(existingAcc.Breaches, breachName) {
				breaches = existingAcc.Breaches
			} else {
				breaches = append(existingAcc.Breaches, breachName)
			}
			account.Breaches = breaches
		} else {
			account.Breaches = []string{breachName}
		}
		accs = append(accs, account)
	}
	if err != nil {
		return []haveibeenbreached.Account{}, err
	}
	return accs, nil
}

func contains(arr []string, str string) bool {
	for _, el := range arr {
		if el == str {
			return true
		}
	}
	return false
}
