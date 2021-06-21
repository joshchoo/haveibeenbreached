package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Response events.APIGatewayProxyResponse

type BreachEvent struct {
	BreachName  string `json:"BreachName"`
	Title       string `json:"Title"`
	Domain      string `json:"Domain"`
	Description string `json:"Description"`
	BreachDate  string `json:"BreachDate"` // YYYY-MM-DD
}

type Breach struct {
	PK          string
	SK          string
	BreachName  string
	Title       string
	Domain      string
	Description string
	BreachDate  time.Time
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)
var tableName = "Breaches"

func Handler(ctx context.Context, event BreachEvent) (Response, error) {
	breachDate, err := time.Parse("2006-01-02", event.BreachDate)
	if err != nil {
		return Response{StatusCode: 400, Body: fmt.Sprintf("Error parsing BreachDate: %s", err)}, err
	}

	newBreach := Breach{
		PK:          partitionKey(event.BreachName),
		SK:          sortKey(event.BreachName),
		BreachName:  event.BreachName,
		Title:       event.Title,
		Domain:      event.Domain,
		Description: event.Description,
		BreachDate:  breachDate,
	}

	attrVal, err := dynamodbattribute.MarshalMap(newBreach)
	if err != nil {
		return Response{StatusCode: 400, Body: fmt.Sprintf("Error marshalling new Breach: %s", err)}, err
	}
	input := &dynamodb.PutItemInput{
		Item:      attrVal,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		return Response{StatusCode: 400, Body: fmt.Sprintf("Error putting Breach: %s", err)}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": fmt.Sprintf("Successfully added breach: %s", newBreach.BreachName),
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
	return fmt.Sprintf("BREACH#%s", key)
}

func sortKey(key string) string {
	return fmt.Sprintf("BREACH#%s", key)
}
