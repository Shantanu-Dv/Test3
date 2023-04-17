package types

import (
	"doc-reco-go/internal/utils"
)

type IndexQuestionRequestBody struct {
	QuestionId int                    `json:"question_id"`
	Document   map[string]interface{} `json:"document"`
	Metadata   map[string]interface{} `json:"metadata"`
}

func (i IndexQuestionRequestBody) Validate() error {
	return utils.ValidateRequestStructTag(i)
}