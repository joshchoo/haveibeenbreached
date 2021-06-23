package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
var db = dynamodb.New(sess)
var repo = haveibeenbreached.NewRepo(db)
var addAccountsToBreachHandler = makeAddAccountsToBreachHandler(repo)

func makeAddAccountsToBreachHandler(repo haveibeenbreached.Repo) func(ctx context.Context, event AddAccountEvent) (Response, error) {
	return func(ctx context.Context, event AddAccountEvent) (Response, error) {
		rawAccounts := event.Accounts
		breachName := event.PathParameters.BreachName

		breach, err := repo.GetBreach(breachName)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error occurred while fetching breach: %s", err)}, err
		}
		if breach == nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Cannot add accounts to non-existent breach: %s", breachName)}, nil
		}

		accounts, err := parseAccounts(rawAccounts)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Invalid email: %s", err)}, err
		}

		accounts, err = setAccountBreaches(accounts, breachName)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error setting breaches on account: %s", err)}, err
		}

		items := mapAccountToItemable(accounts)

		if err = repo.PutItems(items); err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error adding accounts to breach: %s", err)}, err
		}

		body, err := json.Marshal(map[string]interface{}{
			"message": fmt.Sprintf("Successfully added/updated %d accounts to the %s breach.", len(accounts), breachName),
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

func parseAccounts(accounts []string) ([]haveibeenbreached.Account, error) {
	accs := make([]haveibeenbreached.Account, 0, len(accounts))
	for _, account := range accounts {
		email, err := haveibeenbreached.NewEmailAccount(account)
		if err != nil {
			return []haveibeenbreached.Account{}, err
		}
		newAccount := haveibeenbreached.Account{
			Username: email,
			Breaches: make([]string, 0),
		}
		accs = append(accs, newAccount)
	}
	return accs, nil
}

func setAccountBreaches(accounts []haveibeenbreached.Account, breachName string) ([]haveibeenbreached.Account, error) {
	accs := make([]haveibeenbreached.Account, 0, len(accounts))
	for _, account := range accounts {
		foundAccount, err := repo.GetAccount(account.Username)
		if err != nil {
			return []haveibeenbreached.Account{}, err
		}
		var breaches []string
		if foundAccount == nil {
			breaches = []string{breachName}
		} else {
			if contains(foundAccount.Breaches, breachName) {
				breaches = foundAccount.Breaches
			} else {
				breaches = append(foundAccount.Breaches, breachName)
			}
		}
		account.Breaches = breaches
		accs = append(accs, account)
	}
	return accs, nil
}

func mapAccountToItemable(accounts []haveibeenbreached.Account) []haveibeenbreached.Itemable {
	items := make([]haveibeenbreached.Itemable, 0, len(accounts))
	for _, a := range accounts {
		items = append(items, a)
	}
	return items
}

func contains(arr []string, str string) bool {
	for _, el := range arr {
		if el == str {
			return true
		}
	}
	return false
}
