package db

import "localstack-and-lambda/model"

type Client interface {
	InsertMemo(memo model.Memo) error
}

func NewClient() Client {
	return newPostgresqlClient()
}
