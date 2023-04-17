package constants

type recommendationContext struct {
	PreChat         int
	ChatOngoing     int
	PostChat        int
	ConceptSuggest  int
	AnswrWebSuggest int
}

type docDoubtMessage struct {
	Text         string
	Image        string
	TextAndImage string
}

type recommendationInteraction struct {
	BaseState      int
	ShownToUser    int
	UserViews      int
	UserLikes      int
	UserClosesChat int
}

type esIndex struct {
	SearchQuestion         string
	SearchConcept          string
	AutoSuggestionQuestion string
}

type badQueryErrorMessage struct {
	InsufficientOcr     string
	UnsupportedLanguage string
}

type ocrService struct {
	Mathpix string
	Vision  string
}

type recommendationType struct {
	RecommendConcept       string
	RecommendQuestion      string
	AnswrRecommendQuestion string
}

type apiSource struct {
	AutoSuggestionSource string
	SearchSource         string
	BertSearchSource     string
}
