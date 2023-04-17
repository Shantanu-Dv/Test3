package recommender

import (
	"doc-reco-go/internal/constants"
	"doc-reco-go/internal/lib/es_query"
	queryTypes "doc-reco-go/internal/lib/es_query/types"
	recoTypes "doc-reco-go/internal/lib/recommender/types"
	"doc-reco-go/internal/provider/elasticsearch"
	"doc-reco-go/internal/types"
	"doc-reco-go/internal/utils"
)

type SelectRecommender interface {
	Recommend(doubts *types.RecommendRequest) (types.RecommendResponse, error)
}

func Selector(doubts *types.RecommendRequest) SelectRecommender {
	messageIndex := utils.GetDoubtMessageId(doubts.Context)
	var message types.DocMessage
	if messageIndex != -1 {
		message = doubts.Message[messageIndex]
	} else {
		message = doubts.Message[len(doubts.Message)-1]
	}
	msgType := message.MessageType

	if doubts.Source == constants.ApiSource.AutoSuggestionSource {
		return &recoTypes.ESRecommender{
			EsClient:   elasticsearch.ESAutoSuggestionClient,
			IndexName:  constants.EsIndexConstant.AutoSuggestionQuestion,
			QueryClass: queryTypes.AutoSuggestQuery{},
			QueryLimit: 10,
			Message:    message,
		}
	}

	if doubts.Context == constants.RecommendationContextConstant.ConceptSuggest {
		return &recoTypes.ESRecommender{
			EsClient:   elasticsearch.ESSearchClient,
			IndexName:  constants.EsIndexConstant.SearchConcept,
			QueryClass: queryTypes.ConceptQuery{},
			QueryLimit: 10,
			Message:    message,
		}
	}

	if msgType == constants.DocDoubtMessageConstants.Image || msgType == constants.DocDoubtMessageConstants.TextAndImage {

		var queryClass es_query.QueryBuilder = queryTypes.AutoQuery{}
		//if doubts.Source == constants.ApiSource.BertSearchSource {
		//	queryClass = queryTypes.BertQuery{}
		//}

		return &recoTypes.AutoRecommender{
			EsClient:   elasticsearch.ESSearchClient,
			IndexName:  constants.EsIndexConstant.SearchQuestion,
			QueryClass: queryClass,
			QueryLimit: 10,
			Message:    message,
		}
	}

	var queryClass es_query.QueryBuilder = queryTypes.PlainTextQuery{}
	//if doubts.Source == constants.ApiSource.BertSearchSource {
	//	queryClass = queryTypes.BertQuery{}
	//}
	return &recoTypes.ESRecommender{
		EsClient:   elasticsearch.ESSearchClient,
		IndexName:  constants.EsIndexConstant.SearchQuestion,
		QueryClass: queryClass,
		QueryLimit: 10,
		Message:    message,
	}
}
