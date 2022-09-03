package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"localstack-and-lambda/elasticsearch"
	"log"
	"testing"
	"time"
)

func insertToDDB(item map[string]string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://localhost:4566"),
			Region:   aws.String("ap-northeast-2"),
		},
	}))

	ddbItem := convertToDDBItem(item)
	client := dynamodb.New(sess)
	_, err := client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("memos"),
		Item:      ddbItem,
	})

	if err != nil {
		log.Fatalf("ddb upsert item error\n%v", err)
	}
}

func convertToDDBItem(item map[string]string) map[string]*dynamodb.AttributeValue {
	ddbItem := make(map[string]*dynamodb.AttributeValue)

	for k, v := range item {
		ddbItem[k] = &dynamodb.AttributeValue{S: aws.String(v)}
	}

	return ddbItem
}

func TestCreate(t *testing.T) {
	dummyId := uuid.New().String()
	memo := map[string]string{"id": dummyId, "name": "dummyName", "body": "dummyBody"}
	index := elasticsearch.GetMemoESIndex()
	client, _ := elasticsearch.NewESClientWithAddress("http://localhost:9200")

	insertToDDB(memo)
	time.Sleep(time.Second * 5)

	doc := client.GetItem(index, dummyId)
	sourceDoc := doc["_source"].(map[string]any)

	assert.True(t, doc["found"].(bool))
	assert.Equal(t, memo["id"], doc["_id"])
	assert.Equal(t, memo["name"], sourceDoc["name"])
	assert.Equal(t, memo["body"], sourceDoc["body"])
}
