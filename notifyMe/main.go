package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joshuous/haveibeenbreached"
)

type Response events.APIGatewayProxyResponse

/*
TODO:
- validate email
- add email to SQS
*/

type NotifyMeEvent struct {
	Email string
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var db = dynamodb.New(sess)
var r = haveibeenbreached.NewRepo(db)
var sqsQueue = sqs.New(sess)
var mq = haveibeenbreached.NewMessageQueue(sqsQueue)
var notifyMeHandler = makeNotifyMeHandler(r, mq)
var subscriptionsQueueName = os.Getenv("SUBSCRIPTIONS_QUEUE")

func main() {
	lambda.Start(notifyMeHandler)
}

func makeNotifyMeHandler(repo haveibeenbreached.Repo, queue haveibeenbreached.MessageQueue) func(ctx context.Context, event NotifyMeEvent) (Response, error) {
	return func(ctx context.Context, event NotifyMeEvent) (Response, error) {
		rawEmail := event.Email
		subscriber, err := haveibeenbreached.NewSubscriber(rawEmail)
		if err != nil {
			return Response{StatusCode: 400, Body: fmt.Sprintf("Invalid email: %s", rawEmail)}, err
		}

		err = queue.SendMessage(haveibeenbreached.SendMessageInput{
			MessageBody: subscriber.Email,
			QueueName:   subscriptionsQueueName,
		})
		if err != nil {
			return Response{StatusCode: 400}, err
		}

		body, err := json.Marshal(map[string]interface{}{
			"message": "Processing breach notification subscription.",
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
