package types

import (
	"bytes"
	"context"
	"doc-reco-go/internal/constants"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

type QueryComposite struct {
	maxTextWords int
}

func (q *QueryComposite) init() {
	q.maxTextWords = 100 // to reduce ES timeout and decrease latency
}

func (q QueryComposite) CreateQuery(_ context.Context, searchDict map[string]interface{}, additionalFilter map[string]interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer

	searchText, ok := searchDict["search_text"].(string)
	if !ok {
		return buf, fmt.Errorf("search text missing")
	}

	queryLimit, ok := additionalFilter["query_limit"]
	if !ok {
		queryLimit = constants.EsDefaultQueryLimit
	}

	q.init()
	searchText = utils.LimitWords(searchText, q.maxTextWords)
	searchText = utils.FilterSearchText(searchText)

	if text, err := strconv.Unquote(searchText); err == nil {
		searchText = text
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":                searchText,
							"fields":               []string{"question_plain_full"},
							"boost":                0.7,
							"fuzziness":            "AUTO",
							"minimum_should_match": "55%",
						},
					},
				},
				"should": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":                searchText,
							"fields":               []string{"question_plain_full.shingles"},
							"boost":                0.3,
							"fuzziness":            "AUTO",
							"minimum_should_match": "50%",
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"question_status": "published",
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"question_style": "passage",
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
