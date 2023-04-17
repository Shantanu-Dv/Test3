package types

import (
	"doc-reco-go/internal/constants"
	"math"
)

type RecommendSerializer struct {
	QuestionId     int
	RecommendScore float64
	Signature      string
	Interaction    int
	Type           string
	PushDownScores []float64
	// QuestionField        string
	Source               map[string]interface{}
	Content              map[string]interface{}
	PercentageSimilarity int
}

func (s RecommendSerializer) PostInit(context int) RecommendSerializer {
	if context == constants.RecommendationContextConstant.ConceptSuggest {
		s.Type = constants.RecommendationType.RecommendConcept
	} else if val, ok := s.Source["answr_app_q_id"]; ok && val != nil {
		s.Type = constants.RecommendationType.AnswrRecommendQuestion
	} else {
		s.Type = constants.RecommendationType.RecommendQuestion
	}

	s.Interaction = constants.RecommendationInteractionConstant.BaseState
	s.RecommendScore = math.Round(s.RecommendScore*1000) / 1000 // rounding off to 3 digit
	return s
}

func (s RecommendSerializer) ToDict() interface{} {
	return map[string]interface{}{
		"_type":                  s.Type,
		"_recommend_score":       s.RecommendScore,
		"_id":                    s.QuestionId,
		"_percentage_similarity": s.PercentageSimilarity,
		"_content":               s.Content,
		// "_questionField":		  s.QuestionField,
		//"_signature":             s.Signature,
		//"_interaction":           s.Interaction,
		//"_pushdown_scores":       s.PushDownScores,
		//"content":                s.Content,
	}
}
