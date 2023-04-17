package mathpix

import (
	"bytes"
	"context"
	"doc-reco-go/internal/config"
	"doc-reco-go/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var reCngDelimiter = regexp.MustCompile(`(?:\\)(?:\(|\)|\[|\])`)

var delimiterMap = map[string]interface{}{
	"latex":     map[string]string{"left": "$$", "right": "$$"},
	"asciimath": map[string]string{"left": "'", "right": "'"},
}

func ExtractText(imageUrl string, ctx context.Context) (map[string]interface{}, error) {

	if imageUrl == "" {
		return nil, fmt.Errorf("img url is missing")
	}
	res, err := makeRequest(imageUrl, ctx)
	if err != nil {
		return nil, err
	}

	if e, ok := res["error"].(string); ok {
		return map[string]interface{}{
			"mathpix": map[string]interface{}{
				"latex": "", "asciimath": "", "confidence": 0, "error": e,
			},
		}, nil
	}

	resText, ok := "", false
	if resText, ok = res["text"].(string); !ok {
		resText = ""
	}
	latexText := reCngDelimiter.ReplaceAllString(resText, "$$$$")

	var resData []interface{}
	if resData, ok = res["data"].([]interface{}); !ok {
		resData = []interface{}{}
	}

	outFormat := make(map[string]bool)
	for _, item := range resData {
		t := item.(map[string]interface{})["type"].(string)
		if _, ok := outFormat[t]; !ok {
			outFormat[t] = true
		}
	}
	delimiters := []map[string]interface{}{{"left": "$$", "right": "$$"}}
	splits := utils.SplitWithDelimiters(latexText, delimiters)

	outputs := make(map[string][]string)
	for k, _ := range outFormat {
		outputs[k] = []string{}
	}

	for _, item := range splits {
		if item["type"].(string) == "text" {
			if len(item["data"].(string)) != 0 {
				for k, _ := range outputs {
					outputs[k] = append(outputs[k], strings.TrimSpace(item["data"].(string)))
				}
			}
		} else if item["type"].(string) == "math" {
			for i := 0; i < len(outFormat); i++ {
				d := resData[i]
				s1, s2, s3 := "", "", ""
				t := d.(map[string]interface{})["type"].(string)
				s1, ok := delimiterMap[t].(map[string]string)["left"]
				if !ok {
					s1 = "$$"
				}
				s2 = d.(map[string]interface{})["value"].(string)
				s3, ok = delimiterMap[t].(map[string]string)["right"]
				if !ok {
					s1 = "$$"
				}
				outputs[t] = append(outputs[t], fmt.Sprintf("%s%s%s", s1, s2, s3))
			}
		}
	}
	result := make(map[string]string)
	for k, v := range outputs {
		result[k] = strings.Join(v, " ")
	}

	var confidence interface{}
	if confidence, ok = res["confidence"]; !ok || confidence == nil {
		confidence = res["latex_confidence"]
	}

	ret := map[string]interface{}{
		"mathpix": map[string]interface{}{
			"latex": latexText, "asciimath": result["asciimath"], "confidence": confidence,
		},
	}

	return ret, nil
}

func makeRequest(imageUrl string, ctx context.Context) (map[string]interface{}, error) {

	postBody, err := json.Marshal(map[string]interface{}{
		"src":     imageUrl,
		"formats": []string{"text", "data"},
		"data_options": map[string]bool{
			"include_asciimath": true,
			"include_latex":     true,
		},
	})
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(postBody)

	client := http.Client{}
	conf := config.Config.Mathpix
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.ApiUrl, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("app_id", conf.AppId)
	req.Header.Set("app_key", conf.AppKey)
	req.Header.Set("content-type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", r.Status)
	}
	defer r.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if e, ok := res["error"]; ok {
		return map[string]interface{}{"error": e}, nil
	}
	return res, nil
}
