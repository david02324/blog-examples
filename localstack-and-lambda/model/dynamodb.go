package model

import "github.com/aws/aws-lambda-go/events"

func NewMemoFromDDBImage(image map[string]events.DynamoDBAttributeValue) *Memo {
	id := image["id"].String()
	name := image["name"].String()
	body := image["body"].String()

	return &Memo{Id: id, Name: name, Body: body}
}
