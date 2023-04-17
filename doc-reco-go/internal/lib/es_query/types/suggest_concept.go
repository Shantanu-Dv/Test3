package types

import (
	"bytes"
	"context"
	"doc-reco-go/internal/constants"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"fmt"
)

type ConceptQuery struct{}

func (q ConceptQuery) CreateQuery(_ context.Context, searchDict map[string]interface{}, additionalFilter map[string]interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer

	searchText, ok := searchDict["search_text"].(string)
	if !ok {
		return buf, fmt.Errorf("search text missing")
	}

	queryLimit, ok := additionalFilter["query_limit"]
	if !ok {
		queryLimit = constants.EsDefaultQueryLimit
	}

	searchText = utils.FilterSearchText(searchText)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":                searchText,
							"fields":               []string{"heading"},
							"boost":                1.0,
							"fuzziness":            "AUTO",
							"minimum_should_match": "50%",
						},
					},
					{
						"multi_match": map[string]interface{}{
							"query":                searchText,
							"fields":               []string{"item_type"},
							"boost":                1.0,
							"fuzziness":            "AUTO",
							"minimum_should_match": "50%",
						},
					},
					{
						"multi_match": map[string]interface{}{
							"query":                searchText,
							"fields":               []string{"text_plain"},
							"boost":                1.0,
							"fuzziness":            "AUTO",
							"minimum_should_match": "50%",
						},
					},
				},
			},
		},
		"size": queryLimit,
	}

	err := json.NewEncoder(&buf).Encode(query)

	return buf, err
}
