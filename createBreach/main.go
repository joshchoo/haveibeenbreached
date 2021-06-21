package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joshuous/haveibeenbreached"
)

type Response events.APIGatewayProxyResponse

/// BreachEvent is the payload passed to the handler.
type BreachEvent struct {
	BreachName  string
	Title       string
	Domain      string
	Description string
	BreachDate  string // YYYY-MM-DD
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)
var repo = haveibeenbreached.NewRepo(svc)
var createBreachHandler = makeCreateBreachHandler(repo)
var timeLayout = "2006-01-02"

func main() {
	lambda.Start(createBreachHandler)
}

func makeCreateBreachHandler(repo haveibeenbreached.Repo) func(ctx context.Context, event BreachEvent) (Response, error) {
	return func(ctx context.Context, event BreachEvent) (Response, error) {
		breachDate, err := time.Parse(timeLayout, event.BreachDate)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error parsing BreachDate: %s", err)}, err
		}

		breach := haveibeenbreached.Breach{
			BreachName:  event.BreachName,
			Title:       event.Title,
			Domain:      event.Domain,
			Description: event.Description,
			BreachDate:  breachDate,
		}

		if err = repo.PutItem(breach.Item()); err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Error putting Breach: %s", err)}, err
		}

		body, err := json.Marshal(map[string]interface{}{
			"message": fmt.Sprintf("Successfully added breach: %s", breach.BreachName),
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
