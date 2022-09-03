package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"localstack-and-lambda/db"
	"localstack-and-lambda/model"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ApiResponse := events.APIGatewayProxyResponse{}

	switch request.HTTPMethod {
	case "GET":
		ApiResponse = events.APIGatewayProxyResponse{Body: "Hi there!", StatusCode: 200}

	case "POST":
		memo := model.Memo{}

		if err := json.Unmarshal([]byte(request.Body), &memo); err != nil {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Invalid memo format", StatusCode: 500}
			break
		}

		if err := db.NewClient().InsertMemo(memo); err != nil {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Insert memo failed", StatusCode: 500}
			break
		}

		ApiResponse = events.APIGatewayProxyResponse{Body: "Success to insert memo", StatusCode: 200}
	}

	return ApiResponse, nil
}

func main() {
	lambda.Start(HandleRequest)
}
