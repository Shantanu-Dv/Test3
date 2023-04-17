package types

import (
	"doc-reco-go/internal/constants"
	"doc-reco-go/internal/utils"
)

type SuggestRequest struct {
	Text string `json:"text" validate:"required"`
	//SessionId int    `json:"session_id"`
	QuerySize int `json:"size"`
}

func (t *SuggestRequest) SetDefault() {

	if t.QuerySize == 0 || t.QuerySize > constants.EsDefaultQueryLimit {
		t.QuerySize = constants.EsDefaultQueryLimit
	}
}

func (t SuggestRequest) Validate() error {
	return utils.ValidateRequestStructTag(t)
}
