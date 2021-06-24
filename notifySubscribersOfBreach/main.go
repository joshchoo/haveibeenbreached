package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joshuous/haveibeenbreached"
)

type Response events.APIGatewayProxyResponse

type NotifySubscribersEvent struct {
	BreachName string
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)
var r = haveibeenbreached.NewRepo(db)
var notifySubscribersOfBreachHandler = makeNotifySubscribersOfBreachHandler(r)

func main() {
	lambda.Start(notifySubscribersOfBreachHandler)
}

func makeNotifySubscribersOfBreachHandler(repo haveibeenbreached.Repo) func(ctx context.Context, event NotifySubscribersEvent) (Response, error) {
	return func(ctx context.Context, event NotifySubscribersEvent) (Response, error) {
		breachName := event.BreachName
		subscribers, err := findSubscribers(breachName, repo)
		if err != nil {
			return Response{StatusCode: 400}, err
		}

		for _, subscriber := range subscribers {
			if err := sendEmail(subscriber.Email); err != nil {
				// TODO: send to Dead Letter Queue
				return Response{StatusCode: 400}, err
			}
		}

		body, err := json.Marshal(map[string]interface{}{
			"message": fmt.Sprintf("Notified subscribers of breach: %s", breachName),
		})
		if err != nil {
			return Response{StatusCode: 400}, err
		}

		resp := Response{
			StatusCode:      202,
			IsBase64Encoded: false,
			Body:            string(body),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return resp, nil
	}
}

func findSubscribers(breachName string, repo haveibeenbreached.Repo) ([]haveibeenbreached.Subscriber, error) {
	// TODO: Get all subscribers
	// Perform Scan operation for Type=Subscriber
	// allSubscribers := make([]haveibeenbreached.Subscriber, 0)

	// TODO: Get all breached accounts for the breach
	breach, err := repo.GetBreach(breachName)
	if err != nil {
		return []haveibeenbreached.Subscriber{}, err
	}
	if breach == nil {
		return []haveibeenbreached.Subscriber{}, fmt.Errorf("breach does not exist: %s", breachName)
	}

	// TODO: Filter subscribers who are in the breach
	// Select shortest(allSubscribers, breachedAccounts), make into Set datastructure, iterate the other array and look for matches, then append matches to subscribersInBreach
	subscribersInBreach := make([]haveibeenbreached.Subscriber, 0)
	// For prototyping, just notify all breached accounts
	for _, acc := range breach.BreachedAccounts {
		subscriber, err := haveibeenbreached.NewSubscriber(acc)
		if err == nil {
			subscribersInBreach = append(subscribersInBreach, subscriber)
		}
	}

	return subscribersInBreach, nil
}

func sendEmail(email string) error {
	// TODO: Ideally we would send an email with SES / other services
	log.Printf("Successfully sent email to: %s\n", email)
	return nil
}
