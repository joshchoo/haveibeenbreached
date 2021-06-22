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

type FindAccountEvent struct {
	PathParameters struct {
		Username string
	}
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)
var repo = haveibeenbreached.NewRepo(svc)
var findAccountHandler = makeFindAccountHandler(repo)

func main() {
	lambda.Start(findAccountHandler)
}

func makeFindAccountHandler(repo haveibeenbreached.Repo) func(ctx context.Context, event FindAccountEvent) (Response, error) {
	return func(ctx context.Context, event FindAccountEvent) (Response, error) {
		rawUsername := event.PathParameters.Username

		username, err := haveibeenbreached.NewEmailAccount(rawUsername)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Invalid username: %s", rawUsername)}, err
		}

		account, err := repo.GetAccount(username)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error finding account: %s", err)}, err
		}
		if account == nil {
			return Response{StatusCode: 404, Body: "No such account."}, err
		}

		body, err := json.Marshal(map[string]interface{}{
			"message": "Found breaches for account.",
			"data": map[string]interface{}{
				"Breaches": account.Breaches,
			},
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
