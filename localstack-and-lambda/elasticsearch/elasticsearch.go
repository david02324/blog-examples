package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"localstack-and-lambda/model"
	"log"
	"time"
)

type ESClient interface {
	BulkCreate(memos *[]model.Memo) error
	GetItem(index string, id string) map[string]any
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
	return NewESClientWithAddress("http://172.17.0.1:9200")
}

func NewESClientWithAddress(address string) (ESClient, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{address}})
	if err != nil {
		return nil, fmt.Errorf("error creating es client: %s", err)
	}

	return &ESClientImpl{es}, nil
}

func (c ESClientImpl) BulkCreate(memos *[]model.Memo) error {
	var packets []byte
	for _, memo := range *memos {
		index := ESIndex{Index: GetMemoESIndex(), Id: memo.Id}
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

	_, err := c.client.Bulk(bytes.NewReader(packets))
	return err
}

func (c ESClientImpl) GetItem(index string, id string) map[string]any {
	res, err := c.client.Get(index, id)
	if err != nil {
		log.Fatalf("es GetItem error\n%v", err)
	}
	doc := make(map[string]any)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		log.Fatalf("buffer read error\n%v", err)
	}

	err = json.Unmarshal(buf.Bytes(), &doc)
	if err != nil {
		log.Fatalf("json unmarshal error\n%v", err)
	}

	return doc
}

func GetMemoESIndex() string {
	t := time.Now()
	return fmt.Sprintf("memos-%d-%02d", t.Year(), t.Month())
}
