package types

import (
	"bytes"
	"context"
	"doc-reco-go/internal/constants"
	"doc-reco-go/internal/lib/analyzer"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

type AutoQuery struct {
	priority      []string
	equationRegex *regexp.Regexp
	probStatRegex *regexp.Regexp
	maxTextWords  int
	lowWordCount  int

	queryLimit int
}

func (q *AutoQuery) init() {
	q.priority = []string{"mathpix", "vision", "text"}
	//q.equationRegex, _ = regexp.Compile(`(\\mathrm\\{[A-Z0-9^}]{1,4}\\})`)
	//q.probStatRegex, _ = regexp.Compile(`[Pp]\s?\([A-Z](\s?\|\s?[A-Z])?\)`)
	q.maxTextWords = 193 // should allow the query to execute under 4 sec
	q.lowWordCount = 15  // don't trigger shingles for large text
}

func (q AutoQuery) CreateQuery(_ context.Context, searchDict map[string]interface{}, additionalFilter map[string]interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer

	queryLimit, ok := additionalFilter["query_limit"].(int)
	if !ok {
		queryLimit = constants.EsDefaultQueryLimit
	}
	q.queryLimit = queryLimit

	q.init()
	query, err := q.getESQuery(searchDict)

	if err != nil {
		return buf, err
	}

	err = json.NewEncoder(&buf).Encode(query)
	return buf, err
}

func (q AutoQuery) getESQuery(searchDict map[string]interface{}) (map[string]interface{}, error) {
	searchText := q.getText(searchDict)

	if text, err := strconv.Unquote(searchText); err == nil {
		searchText = text
	}
	searchText = utils.LimitWords(searchText, q.maxTextWords)
	//searchText = utils.FilterSearchText(searchText)

	visionText := ""
	if vTxt, ok := searchDict["vision"]; ok {
		asciiString := strconv.QuoteToASCII(vTxt.(string))
		searchDict["vision"] = asciiString
		visionText = utils.LimitWords(asciiString, q.maxTextWords)
		if text, err := strconv.Unquote(visionText); err == nil {
			visionText = text
		}
	}
	textToken, latexToken, err := utils.SeparateMath(searchText)
	if err != nil {
		return nil, err
	}

	latexFrac := float32(len(strings.Join(latexToken, " "))) / float32(len(strings.Join(textToken, " "))+1)

	var queryString interface{}
	if len(latexToken) > 0 {
		if len(latexToken) > 2 || latexFrac < 0.80 {
			queryString = q.getLatexAndTextQuery(searchText, visionText)
		} else if len(strings.Join(latexToken, " ")) < 10 {
			queryString = q.getTextQuery(searchText)
		} else {
			queryString, err = q.getLatexQuery(searchText, strings.Join(latexToken, " "))
			if err != nil {
				return nil, err
			}
		}
	} else {
		queryString = q.getTextQuery(searchText)
	}

	return map[string]interface{}{"query": queryString, "size": q.queryLimit}, nil
}

func (q AutoQuery) getText(searchDict map[string]interface{}) string {
	for _, priority := range q.priority {
		if text, ok := searchDict[priority]; ok && text != "" {
			if priority == constants.OcrService.Mathpix && !q.useMathpix(searchDict) {
				if _, ok = searchDict["vision"]; ok {
					continue
				}
			}
			return searchDict[priority].(string)
		}
	}
	return ""
}

func (q AutoQuery) useMathpix(searchDict map[string]interface{}) bool {
	/*
		Function to decide whether to use mathpix or not
		:param mathpix_text:
		:return: bool to decide whether to use mathpix or not
	*/
	mathpixText, ok := searchDict["mathpix"].(string)
	if !ok {
		return false
	}
	visionText := searchDict["vision"].(string)
	mathpixTokens, _ := analyzer.NormText{}.Tokenize(mathpixText)

	visionTokens, _ := analyzer.NormText{}.Tokenize(visionText)

	if len(strings.Join(mathpixTokens, " "))-len(strings.Join(visionTokens, " ")) < 5 {
		return false
	}
	return true
}

func (q AutoQuery) getLatexAndTextQuery(latexText, plainText string) interface{} {
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"multi_match": map[string]interface{}{
						"query":                latexText,
						"fields":               []string{"question_full"},
						"boost":                0.4,
						"fuzziness":            "AUTO",
						"minimum_should_match": "50%",
					},
				},
			},
			"should": []map[string]interface{}{
				{
					"multi_match": map[string]interface{}{
						"query":                plainText,
						"fields":               []string{"question_plain_full"},
						"boost":                0.6,
						"fuzziness":            "AUTO",
						"minimum_should_match": "45%",
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
	}
}

func (q AutoQuery) getTextQuery(text string) interface{} {

	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"multi_match": map[string]interface{}{
						"query":                text,
						"fields":               []string{"question_plain_full"},
						"boost":                0.7,
						"fuzziness":            "AUTO",
						"minimum_should_match": "55%",
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
	}

	if utils.UseShingles(text, q.lowWordCount) {
		query["bool"].(map[string]interface{})["should"] = []map[string]interface{}{
			{
				"multi_match": map[string]interface{}{
					"query":                text,
					"fields":               []string{"question_plain_full.shingles"},
					"boost":                0.3,
					"fuzziness":            "AUTO",
					"minimum_should_match": "65%",
				},
			},
		}
	}

	return query
}

func (q AutoQuery) getLatexQuery(fullText, latexText string) (interface{}, error) {
	tokens, err := analyzer.NormText{}.Tokenize(latexText)
	if err != nil {
		return nil, err
	}

	tokenLength, scale := len(tokens), 3

	query := map[string]interface{}{
		"function_score": map[string]interface{}{
			"functions": []map[string]interface{}{
				{
					"gauss": map[string]interface{}{
						"question_full_latex.length": map[string]interface{}{
							"origin": tokenLength,
							"scale":  scale,
						},
					},
				},
			},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"multi_match": map[string]interface{}{
								"query":                latexText,
								"fields":               []string{"question_full_latex.shingles"},
								"boost":                1,
								"fuzziness":            "AUTO",
								"minimum_should_match": "50%",
							},
						},
						{
							"multi_match": map[string]interface{}{
								"query":                fullText,
								"fields":               []string{"question_full.shingles"},
								"boost":                1,
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
		},
	}

	return query, nil
}
