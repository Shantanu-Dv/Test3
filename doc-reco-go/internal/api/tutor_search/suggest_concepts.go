package tutor_search

import (
	"doc-reco-go/internal/constants"
	queryTypes "doc-reco-go/internal/lib/es_query/types"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
	"net/http"
)

// /suggest/concepts TODO: Remove this endpoint based on usage
func SuggestConcepts(w http.ResponseWriter, r *http.Request) {
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

	results, err := runSearch(searchDict, additionalFilter, constants.EsIndexConstant.SearchConcept, queryTypes.ConceptQuery{})
	if err != nil {
		utils.SentryCaptureError(r, err)
		utils.SendInternalServerError(w, nil, err.Error())
		return
	}

	var response []interface{}
	for i := range results {
		if res, ok := results[i].(map[string]interface{}); ok {
			response = append(response, res["_source"])
		}
	}
	utils.SendSuccessV2(w, response)
	return
}
