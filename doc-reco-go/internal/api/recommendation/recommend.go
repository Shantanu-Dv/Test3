package recommendation

import (
	"doc-reco-go/internal/lib/recommender"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
	"net/http"
)

func Recommend(w http.ResponseWriter, r *http.Request) {

	request := &types.RecommendRequest{}
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

	request.SetDefault(r.URL.Query())

	response, err := recommender.Selector(request).Recommend(request)
	if err != nil {
		utils.SentryCaptureError(r, err)
		utils.SendInternalServerError(w, nil, err.Error())
		return
	}

	// fmt.Print(response.RecommendationList)
	// fmt.Print(response.Confidence)
	utils.SendSuccess(w, response, "success")
	return
}
