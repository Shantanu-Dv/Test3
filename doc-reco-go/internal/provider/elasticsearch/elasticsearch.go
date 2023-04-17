package elasticsearch

import (
	"bytes"
	"context"
	"doc-reco-go/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	ESAutoSuggestionClient *elasticsearch.Client // elastic search client
	ESSearchClient         *elasticsearch.Client
)

func InitializeEsClients() error {
	autoSuggestionCfg := elasticsearch.Config{
		Addresses: []string{config.Config.Elasticsearch.AutoSuggestion},
	}

	searchCfg := elasticsearch.Config{
		Addresses: []string{config.Config.Elasticsearch.Search},
	}

	var err error

	ESAutoSuggestionClient, err = elasticsearch.NewClient(autoSuggestionCfg)
	if err != nil {
		return err
	}

	ESSearchClient, err = elasticsearch.NewClient(searchCfg)
	if err != nil {
		return err
	}

	return nil
}

func AddToES(client *elasticsearch.Client, body map[string]interface{}, id string, index string) error {
	var err error
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyParam := strings.NewReader(bytes.NewBuffer(bodyJson).String())

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bodyParam,
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		msg := fmt.Sprintf("Error getting response: %s", err)
		return errors.New(msg)
	}

	defer res.Body.Close()

	if res.IsError() {
		msg := fmt.Sprintf("[%s] Error indexing document ID", res.Status())
		return errors.New(msg)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		msg := fmt.Sprintf("Error parsing the response body: %s", err)
		return errors.New(msg)
	}
	return nil
}
