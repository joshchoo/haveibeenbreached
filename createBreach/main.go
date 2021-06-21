package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

// Create the DynamoDB client
var svc = dynamodb.New(sess)

func Handler(ctx context.Context) (Response, error) {
	input := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(input)
	if err != nil {
		return Response{StatusCode: 500, Body: "Failed to list DynamoDB tables"}, err
	}

	tableNames := []string{}
	for _, name := range result.TableNames {
		tableNames = append(tableNames, *name)
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": strings.Join(tableNames, "\n"),
	})
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	resp := Response{
		StatusCode:      201,
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
