package indexing

import (
	// "context"
	"doc-reco-go/internal/constants"

	// bertEncoder "doc-reco-go/internal/provider/bert_encoder"
	"doc-reco-go/internal/provider/elasticsearch"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	// "time"

	esPkg "github.com/elastic/go-elasticsearch/v7"
	"github.com/getsentry/sentry-go"
)

func IndexQuestions(w http.ResponseWriter, r *http.Request) {
	request := &types.IndexQuestionRequestBody{}

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

	// print(request)
	questionId := strconv.Itoa(request.QuestionId)

	preprocessedDoc, err := preprocessDocument(request.Document)
	// preprocessedDoc := request.Document
	print(preprocessedDoc)
	if err != nil {
		utils.SentryCaptureError(r, err)
	}

	type IndexDoc struct {
		index string
		doc   map[string]interface{}
	}
	esIndexPairs := map[*esPkg.Client]IndexDoc{
		elasticsearch.ESSearchClient: {
			index: constants.EsIndexConstant.SearchQuestion,
			doc:   preprocessedDoc,
		},
		elasticsearch.ESAutoSuggestionClient: {
			index: constants.EsIndexConstant.AutoSuggestionQuestion,
			doc:   removeQuestionVectorField(preprocessedDoc),
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(esIndexPairs))

	for client, indexDoc := range esIndexPairs {

		go func(localHub *sentry.Hub, client *esPkg.Client, indexDoc IndexDoc) {
			defer wg.Done()
			err = elasticsearch.AddToES(client, indexDoc.doc, questionId, indexDoc.index)
			if err != nil {
				localHub.CaptureMessage(fmt.Sprintf("IndexQuestion: %v", err))
			}
		}(sentry.CurrentHub().Clone(), client, indexDoc)
	}

	wg.Wait()

	response := "Question indexed successfully"
	utils.SendSuccess(w, response, "success")
}

func preprocessDocument(document map[string]interface{}) (map[string]interface{}, error) {

	preprocessedDoc := make(map[string]interface{})

	var questionStr string
	var questionChoices []string
	var columns []string
	var questionAssertions []string

	// question (string)
	if question, ok := document["question"]; ok {
		questionStr = question.(string)
	} else {
		return preprocessedDoc, errors.New("question is empty")
	}

	//choices (arr of dictionary)
	if choices, ok := document["choices"].([]interface{}); ok && choices != nil {
		for _, val := range choices {
			if ch, present := val.(map[string]interface{})["choice"]; present {
				questionChoices = append(questionChoices, ch.(string))
			}
		}
	}
	// columns (to be confirmed currently set as array of string)
	if cols, ok := document["mx_l2"].([]interface{}); ok && cols != nil {
		for _, val := range cols {
			columns = append(columns, val.(string))
		}
	}

	if cols, ok := document["mx_l1"].([]interface{}); ok && cols != nil {
		for _, val := range cols {
			columns = append(columns, val.(string))
		}
	}

	// assertion - reason (string)
	if assertion, ok := document["assertion"]; ok {
		questionAssertions = append(questionAssertions, assertion.(string))
		questionAssertions = append(questionAssertions, document["reason"].(string))
	}

	questionArr := utils.ChainArr(questionAssertions, questionStr, questionChoices, columns)
	questionJointText := strings.Join(questionArr, " \n ")
	_, latexText, _ := utils.SeparateMath(questionJointText)
	latexFormatText := strings.Join(latexText, " ")

	document["question"] = utils.StripHtmlTags(questionStr)
	document["question_plain"] = utils.LatexToPlainText(utils.StripHtmlTags(questionStr))
	document["question_full"] = utils.StripHtmlTags(questionJointText)
	document["question_plain_full"] = utils.LatexToPlainText(utils.StripHtmlTags(questionJointText))
	document["question_full_latex"] = latexFormatText

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()
	// queryVector, err := bertEncoder.Encoder.Encode(ctx, document["question_plain_full"].(string))
	// if err != nil {
	// 	return document, err
	// }

	// document["question_vector"] = queryVector
	document["question_vector"] = "none"
	return document, nil
}

func removeQuestionVectorField(doc map[string]interface{}) map[string]interface{} {

	res := make(map[string]interface{})
	for k, v := range doc {
		if k != "question_vector" {
			res[k] = v
		}
	}
	return res
}
