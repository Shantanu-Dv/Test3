package vision

import (
	"bytes"
	"doc-reco-go/internal/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


var visionFeatureTypeMap = map[string][]map[string]string{
	"Text": {map[string]string{"type": "TEXT_DETECTION"}},
	"DocText": {map[string]string{"type": "DOCUMENT_TEXT_DETECTION"}},
}

func ExtractText(imageUrl, feature string) (map[string]interface{}, error) {
	if feature == "" {
		feature = "Text"
	}
	if imageUrl == "" {
		return nil, fmt.Errorf("img url is missing")
	}

	base64Encoding, err := ToBase64(imageUrl)
	if err != nil {
		return nil, err
	}

	postBody, err := json.Marshal(map[string][]interface{}{
		"requests": {
			map[string]interface{}{
				"image":         map[string]interface{}{"content": base64Encoding},
				"features":      visionFeatureTypeMap[feature],
				"image_context": map[string]interface{}{"language_hints": []string{"en"}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(postBody)
	r, err := http.Post(config.Config.GoogleVision.ApiUrl, "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s",r.Status)
	}

	var res map[string][]interface{}
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	textAnnotations := res["responses"][0].(map[string]interface{})["textAnnotations"].([]interface{})[0]

	ret := map[string]interface{}{
		"vision": map[string]string{"text": textAnnotations.(map[string]interface{})["description"].(string)},
	}
	return ret, nil
}

func ToBase64(imageUrl string) (string, error) {
	img, err := http.Get(imageUrl)
	if err != nil {
		return "", err
	}
	defer img.Body.Close()

	b, err := ioutil.ReadAll(img.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}