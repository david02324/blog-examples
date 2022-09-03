package db

import (
	"database/sql"
	"fmt"
	"localstack-and-lambda/model"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "Don't look!"
	DB_NAME     = "postgres"
)

type PostgresqlClient struct {
	db *sql.DB
}

func newPostgresqlClient() PostgresqlClient {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return PostgresqlClient{db: db}
}

func (c PostgresqlClient) InsertMemo(memo model.Memo) error {
	result, err := c.db.Exec("INSERT INTO memo (name, body) VALUES(%s, %s)", memo.Name, memo.Body)
	if err != nil {
		return err
	}

	if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}
