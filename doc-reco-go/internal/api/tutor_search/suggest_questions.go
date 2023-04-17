package tutor_search

import (
	"context"
	"doc-reco-go/internal/constants"
	esQuery "doc-reco-go/internal/lib/es_query"
	queryTypes "doc-reco-go/internal/lib/es_query/types"
	"doc-reco-go/internal/provider/elasticsearch"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SuggestQuestions(w http.ResponseWriter, r *http.Request) {
	print("hi")
	request := &types.SuggestRequest{}
	err := utils.ParseRequestBody(r, request)
	if err != nil {
		utils.SendBadRequestError(w, nil, err.Error())
		return
	}
	err = request.Validate()
	if err != nil {
		utils.SendBadRequestError(w, nil, err.Error())
		return
	}
	request.SetDefault()

	searchDict := map[string]interface{}{
		"search_text": request.Text,
	}
	additionalFilter := map[string]interface{}{
		"query_limit": request.QuerySize,
	}

	results, err := runSearch(searchDict, additionalFilter, constants.EsIndexConstant.SearchQuestion, queryTypes.QueryComposite{})
	if err != nil {
		utils.SentryCaptureError(r, err)
		utils.SendInternalServerError(w, nil, err.Error())
		return
	}

	var response []interface{}
	for _, result := range results {
		res := result.(map[string]interface{})
		response = append(response, res["_source"])
	}
	utils.SendSuccessV2(w, response)
	return
}

func runSearch(searchDict, filter map[string]interface{}, esIndex string, query esQuery.QueryBuilder) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryText, err := query.CreateQuery(ctx, searchDict, filter)
	if err != nil {
		return nil, err
	}
	es := elasticsearch.ESSearchClient
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(esIndex),
		es.Search.WithBody(&queryText),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error fetching data from elasticsearch")
	}

	var jsonResponse map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	//resMaxScore, ok := jsonResponse["hits"].(map[string]interface{})["max_score"].(float64)
	//if !ok && resMaxScore < 0.5 {
	//	res := make([]interface{}, 0)
	//	return res, nil
	//}

	return jsonResponse["hits"].(map[string]interface{})["hits"].([]interface{}), nil
}
