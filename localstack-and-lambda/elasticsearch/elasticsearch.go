package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"localstack-and-lambda/model"
	"time"
)

type ESClient interface {
	BulkCreate(memos *[]model.Memo) error
}

type ESClientImpl struct {
	client *elasticsearch.Client
}

type ESIndex struct {
	Id    string `json:"_id"`
	Index string `json:"_index"`
}

type ESCreateIndex struct {
	Index ESIndex `json:"create"`
}

func NewESClient() (ESClient, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{"http://172.17.0.1:9200"}})
	if err != nil {
		return nil, fmt.Errorf("error creating es client: %s", err)
	}

	return &ESClientImpl{es}, nil
}

func (c ESClientImpl) BulkCreate(memos *[]model.Memo) error {
	var packets []byte
	for _, memo := range *memos {
		index := ESIndex{Index: GetMessageESIndex(), Id: memo.Id}
		createIndex := ESCreateIndex{index}

		indexBytes, err := json.Marshal(createIndex)
		if err != nil {
			return err
		}

		docBytes, err := json.Marshal(memo)
		if err != nil {
			return err
		}

		packets = append(packets, indexBytes...)
		packets = append(packets, "\n"...)
		packets = append(packets, docBytes...)
		packets = append(packets, "\n"...)
	}

	fmt.Println(string(packets))
	_, err := c.client.Bulk(bytes.NewReader(packets))
	return err
}

func GetMessageESIndex() string {
	t := time.Now()
	return fmt.Sprintf("memos-%d-%02d", t.Year(), t.Month())
}
