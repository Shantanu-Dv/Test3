package types

import (
	"context"
	"doc-reco-go/internal/constants"
	esQuery "doc-reco-go/internal/lib/es_query"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

type ESRecommender struct {
	EsClient   *elasticsearch.Client
	IndexName  string
	QueryClass esQuery.QueryBuilder
	QueryLimit int
	Message    types.DocMessage
}

func (e ESRecommender) Recommend(doubts *types.RecommendRequest) (types.RecommendResponse, error) {

	websearchText := e.Message.Body

	response := types.RecommendResponse{
		WebsearchText: websearchText,
		ExtractedText: map[string]interface{}{"": websearchText}, // added this way to create entry in DOC
	}

	if doubts.OcrOnly {
		return response, nil
	}

	err := utils.CheckBadQueryError(websearchText, doubts.BadQuery, false)
	if err != nil {
		response.Error = err.Error()
		return response, nil
	}

	//subjectIdsFilter, err := subjects.DocSubject{}.Init(nil, nil, nil).GetSubjectIdFilter(doubts.Session.Klass, doubts.Session.SubjectId)
	//if err != nil {
	//	return response, err
	//}

	searchDict := map[string]interface{}{
		"search_text": e.Message.Body,
	}
	filter := map[string]interface{}{
		//"subject_id":  subjectIdsFilter,
		"query_limit": e.QueryLimit,
	}
	result, err := e.runSearch(searchDict, filter)
	if err != nil {
		return response, err
	}

	response.RecommendationList = formatResult(result, websearchText, constants.QuestionPlainField, doubts.Context)

	return response, nil
}

func (e ESRecommender) runSearch(searchDict map[string]interface{}, filter map[string]interface{}) ([]interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	queryText, err := e.QueryClass.CreateQuery(ctx, searchDict, filter)

	if err != nil {
		return nil, err
	}
	es := e.EsClient

	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(e.IndexName),
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

	resMaxScore, ok := jsonResponse["hits"].(map[string]interface{})["max_score"].(float64)

	if !ok && resMaxScore < 0.5 {
		res := make([]interface{}, 0)
		return res, nil
	}

	return jsonResponse["hits"].(map[string]interface{})["hits"].([]interface{}), nil
}

func formatResult(results []interface{}, websearchText, questionField string, context int) (formattedResult []interface{}) {

	var formattedStruct []types.RecommendSerializer
	websearchText = utils.ReplaceMathsTextToSymbol(websearchText)

	for _, result := range results {
		if res, ok := result.(map[string]interface{}); ok {
			quesId, _ := strconv.Atoi(res["_id"].(string))
			source := res["_source"].(map[string]interface{})
			serializedResult := types.RecommendSerializer{
				QuestionId:     quesId,
				RecommendScore: res["_score"].(float64),
				Source:         source,
				Content:        map[string]interface{}{"question": source[questionField].(string), "question_full": source["question_full"].(string)},
			}

			// fmt.Printf("The type of source[questionField] is %T\n", source[questionField])
			// fmt.Println("The source questionsField : ", source[questionField])

			// fmt.Printf("The type of result is %T\n", result)
			// fmt.Println("The full question is : ", source["question_full"])

			// fmt.Println("Res : ", result)
			if context != constants.RecommendationContextConstant.ConceptSuggest {
				serializedResult.PercentageSimilarity = utils.GetPercentageQuestionMatch(websearchText, source[questionField].(string))
			} //TODO : add surce[questionField].(string) to the serializedResult

			formattedStruct = append(formattedStruct, serializedResult)
		}
	}

	sort.SliceStable(formattedStruct, func(i, j int) bool {
		return formattedStruct[i].PercentageSimilarity > formattedStruct[j].PercentageSimilarity
	})

	for i := range formattedStruct {
		formattedResult = append(formattedResult, formattedStruct[i].PostInit(context).ToDict())
	}

	return
}
