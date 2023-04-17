package constants

const (
	DocSearchSource     = "doc"
	EsDefaultQueryLimit = 10

	True = "true"

	QuestionField      = "question"
	QuestionPlainField = "question_plain"
)

var RecommendationContextConstant = recommendationContext{
	PreChat:         1,
	ChatOngoing:     2,
	PostChat:        3,
	ConceptSuggest:  4,
	AnswrWebSuggest: 5,
}

var DocDoubtMessageConstants = docDoubtMessage{
	Text:         "text",
	Image:        "image",
	TextAndImage: "text and image",
}

var RecommendationInteractionConstant = recommendationInteraction{
	BaseState:      0,
	ShownToUser:    1,
	UserViews:      2,
	UserLikes:      3,
	UserClosesChat: 4,
}

var EsIndexConstant = esIndex{
	SearchQuestion:         "question_v2",
	SearchConcept:          "concept",
	AutoSuggestionQuestion: "question_v6",
}

var BadQueryErrorMessage = badQueryErrorMessage{
	InsufficientOcr:     "INSUFFICIENT_OCR",
	UnsupportedLanguage: "UNSUPPORTED_LANGUAGE",
}

var OcrService = ocrService{
	Mathpix: "mathpix",
	Vision:  "vision",
}

var RecommendationType = recommendationType{
	RecommendQuestion:      "RecommendQuestion",
	AnswrRecommendQuestion: "AnswrRecommendQuestion",
	RecommendConcept:       "RecommendConcept",
}

var ApiSource = apiSource{
	AutoSuggestionSource: "auto_suggestion",
	SearchSource:         "search",
	BertSearchSource:     "bert_search",
}
