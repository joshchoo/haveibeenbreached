package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joshuous/haveibeenbreached"
)

type SubscriptionEvent events.SQSEvent

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)
var r = haveibeenbreached.NewRepo(db)
var sendSubscriptionEmailHandler = makeSendSubscriptionEmailHandler(r)

func main() {
	lambda.Start(sendSubscriptionEmailHandler)
}

func makeSendSubscriptionEmailHandler(repo haveibeenbreached.Repo) func(ctx context.Context, event SubscriptionEvent) error {
	return func(ctx context.Context, event SubscriptionEvent) error {
		subscribers := getSubscribers(event.Records)

		for _, subscriber := range subscribers {
			if err := repo.PutItem(subscriber); err != nil {
				// TODO: send to Dead Letter Queue
				return err
			}
			if err := sendEmail(subscriber.Email); err != nil {
				// TODO: send to Dead Letter Queue
				return err
			}
		}
		return nil
	}
}

func getSubscribers(records []events.SQSMessage) []haveibeenbreached.Subscriber {
	subscribers := make([]haveibeenbreached.Subscriber, 0)
	for _, record := range records {
		email := record.Body
		subscriber, err := haveibeenbreached.NewSubscriber(email)
		if err == nil {
			subscribers = append(subscribers, subscriber)
		}
	}
	return subscribers
}

func sendEmail(email string) error {
	// Ideally we would send an email with SES / other services
	log.Printf("Successfully sent email to: %s\n", email)
	return nil
}
