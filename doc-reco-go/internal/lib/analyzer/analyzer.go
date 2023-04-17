package analyzer

import (
	"bytes"
	"doc-reco-go/internal/config"
	"doc-reco-go/internal/constants"
	"encoding/json"
	"fmt"
	"net/http"
)

type NormText struct{}

func (a NormText) Tokenize(text string) ([]string, error) {
	var resBody map[string]interface{}

	postBody, _ := json.Marshal(map[string]interface{}{
		"text":        text,
		"filter":      []interface{}{"lowercase", "asciifolding", "TokenFilterMathStop"},
		"char_filter": []interface{}{"CharFilterBlanks", "HandleChemFormulae", "HandleArcTrig", "splitEquations", "HandleDegree", "HandleLatexContext", "SeperateTextNum", "CharFilterStopWord"},
		"tokenizer":   "TokenizerMathText",
	})
	reqBody := bytes.NewBuffer(postBody)

	esUrl := config.Config.Elasticsearch.Search
	endpoint := fmt.Sprintf("%s/%s/%s", esUrl, constants.EsIndexConstant.SearchQuestion, "_analyze")
	res, err := http.Post(endpoint, "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("%s returned %s", endpoint, res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	var tokens []string
	if tokenSlice, ok := resBody["tokens"].([]interface{}); ok {
		for _, token := range tokenSlice {
			t := token.(map[string]interface{})
			tokens = append(tokens, t["token"].(string))
		}
	}

	return tokens, nil
}
