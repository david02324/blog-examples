package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"localstack-and-lambda/elasticsearch"
	"localstack-and-lambda/model"
)

func HandleRequest(request events.DynamoDBEvent) {
	var memos []model.Memo

	fmt.Println(memos)
	for _, record := range request.Records {
		if record.EventName == "INSERT" {
			image := record.Change.NewImage
			memo := model.NewMemoFromDDBImage(image)

			memos = append(memos, *memo)
		}
	}

	client, err := elasticsearch.NewESClient()
	if err != nil {
		fmt.Println(fmt.Errorf("elasticsearch client error %w", err))
	}

	err = client.BulkCreate(&memos)
	if err != nil {
		fmt.Println(fmt.Errorf("elasticsearch bulk error %w", err))
	}
}

func main() {
	lambda.Start(HandleRequest)
}
