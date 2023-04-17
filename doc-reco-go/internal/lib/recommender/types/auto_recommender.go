package types

import (
	"context"
	"doc-reco-go/internal/constants"
	esQuery "doc-reco-go/internal/lib/es_query"
	"doc-reco-go/internal/provider/ocr/mathpix"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/getsentry/sentry-go"
	"log"
	"strings"
	"time"
)

type AutoRecommender struct {
	EsClient   *elasticsearch.Client
	IndexName  string
	QueryClass esQuery.QueryBuilder
	QueryLimit int
	Message    types.DocMessage

	context context.Context
}

func (a AutoRecommender) Recommend(doubts *types.RecommendRequest) (types.RecommendResponse, error) {

	var ok bool
	var cancel context.CancelFunc

	a.context, cancel = context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	ocrDict := a.detectText(a.Message.ImageUrl)

	ocrService := constants.OcrService.Mathpix
	quesField := constants.QuestionPlainField
	ocrConfidence := ocrDict["confidence"].(map[string]interface{})

	websearchText := ocrDict["mathpixAscii"].(string)
	if websearchText == "" {
		websearchText = ocrDict["mathpix"].(string)
		quesField = constants.QuestionField
	}

	if websearchText == "" {
		if websearchText, ok = ocrDict["vision"].(string); ok && websearchText != "" {
			ocrService = constants.OcrService.Vision
		}
	}

	extractedText := make(map[string]interface{})
	if mathpixText, ok := ocrDict["mathpix"].(string); ok && mathpixText != "" {
		extractedText["mathpix"] = mathpixText
	}
	if visionText, ok := ocrDict["vision"].(string); ok && visionText != "" {
		extractedText["google_vision"] = visionText
	}

	response := types.RecommendResponse{
		ExtractedText: extractedText,
		WebsearchText: websearchText,
		Confidence:    ocrConfidence,
		Service:       ocrService,
	}

	if doubts.OcrOnly {
		return response, nil
	}

	err := utils.CheckBadQueryError(websearchText, doubts.BadQuery, true)
	if err != nil {
		response.Error = err.Error()
		return response, nil
	}

	filter := map[string]interface{}{
		"query_limit": a.QueryLimit,
	}

	res, err := a.runSearch(ocrDict, filter)
	if err != nil {
		log.Println("AUTO RECOMMENDER: ", err)
	}

	response.RecommendationList = formatResult(res, websearchText, quesField, doubts.Context)

	return response, nil
}

func (a AutoRecommender) runSearch(searchDict map[string]interface{}, filter map[string]interface{}) ([]interface{}, error) {
	queryText, err := a.QueryClass.CreateQuery(a.context, searchDict, filter)
	if err != nil {
		return nil, err
	}
	es := a.EsClient

	res, err := es.Search(
		es.Search.WithContext(a.context),
		es.Search.WithIndex(a.IndexName),
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

func (a AutoRecommender) detectText(imageUrl string) map[string]interface{} {

	searchMap := make(map[string]interface{})
	var mathpixLatex, visionText, mathpixAscii string
	var mathpixConfidence, visionConfidence float64

	for i := 0; i < 3; i++ {
		// don't loop once deadline exceeded
		if a.context.Err() == context.DeadlineExceeded {
			break
		}

		searchMap = a.extractOcrText(imageUrl)
		if tmp, ok := searchMap["mathpix"].(map[string]interface{}); ok {

			if err, ok := tmp["error"]; ok {
				// err: Content not found
				fmt.Println("MATHPIX: ", err)
				break
			}

			mathpixLatex = tmp["latex"].(string)
			mathpixAscii = tmp["asciimath"].(string)
			mathpixConfidence = tmp["confidence"].(float64)
		}
		if tmp, ok := searchMap["vision"].(map[string]interface{}); ok {
			visionText = tmp["text"].(string)
			visionConfidence = tmp["confidence"].(float64)
		}

		if mathpixLatex != "" || visionText != "" {
			break
		}
	}
	if mathpixLatex != "" && visionText != "" && float32(len(mathpixLatex))/float32(len(visionText)+1) < 0.8 {
		mathpixLatex = ""
	}

	return map[string]interface{}{
		"vision":       strings.Replace(visionText, "\n", " ", -1),
		"mathpixAscii": strings.Replace(mathpixAscii, "\n", " ", -1),
		"mathpix":      strings.Replace(mathpixLatex, "\n", " ", -1),
		"confidence":   map[string]interface{}{"mathpix": mathpixConfidence, "google_vision": visionConfidence},
	}
}

func (a AutoRecommender) extractOcrText(imageUrl string) map[string]interface{} {
	text := make(chan map[string]interface{})
	searchMap := make(map[string]interface{})

	go func(localHub *sentry.Hub) {
		t, e := mathpix.ExtractText(imageUrl, a.context)
		// not logging `deadline exceeded` error on sentry
		if e != nil && a.context.Err() != context.DeadlineExceeded {
			localHub.CaptureMessage(fmt.Sprintf("%s: %s", "mathpix", e.Error()))
		}
		text <- t
	}(sentry.CurrentHub().Clone())

	select {
	case m := <-text:
		for k, v := range m {
			searchMap[k] = v
		}
	}
	return searchMap
}
