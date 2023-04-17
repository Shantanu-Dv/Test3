package utils

import (
	"bytes"
	"doc-reco-go/internal/config"
	"doc-reco-go/internal/constants"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func createESIndex() error {
	var err error

	indexName := constants.EsIndexConstant.AutoSuggestionQuestion
	esUrl := config.Config.Elasticsearch.Search
	mappingFile, err := os.Open("../static/test.json")

	if err != nil {
		fmt.Println(err)
	}

	byteValueOfMap, _ := ioutil.ReadAll(mappingFile)
	client := &http.Client{}
	reqBody := bytes.NewBuffer(byteValueOfMap)
	req, err := http.NewRequest(http.MethodPut, esUrl+indexName, reqBody)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != 200 {
		log.Println(resp)
		return errors.New("error creating ES index")
	}
	log.Println("Index created successfully!")
	return nil
}
