package types

import (
	"doc-reco-go/internal/constants"
	"doc-reco-go/internal/utils"
	"fmt"
	"net/url"
)

type RecommendRequest struct {
	Context int          `json:"context"`
	Session DocSession   `json:"session"`
	Message []DocMessage `json:"messages" validate:"required"`

	BadQuery bool
	OcrOnly  bool
	Source   string
}

type RecommendResponse struct {
	ExtractedText      map[string]interface{} `json:"extracted_text"`
	WebsearchText      string                 `json:"websearch_text"`
	RecommendationList []interface{}          `json:"recommendation_list"`
	Confidence         map[string]interface{} `json:"confidence"`
	Service            string                 `json:"service"`
	Error              string                 `json:"error"`
}

type DocSession struct {
	//SessionId         int       `json:"session_id"`
	//UserProfileId     int       `json:"user_profile_id"`
	//SubjectId int `json:"subject_id"`
	//StartedOn         float32 `json:"started_on"`
	//LastUpdatedOn     float32 `json:"last_updated_on"`
	//PlatformStartedOn float32 `json:"platform_started_on"`
	//State             string       `json:"state" validate:"required"`
	//IsDoubtSolved     bool      `json:"is_doubt_solved"`
	//NBroadcasts       int       `json:"n_broadcasts"`
	//Klass int `json:"klass"`
	//KlassInt          int       `json:"klass_int"`
	//Rating            int       `json:"rating"`
	//EndedOn           float32 `json:"ended_on"`
	//IsClosed          bool      `json:"is_closed"`
	//CloseInitiatedBy  string    `json:"close_initiated_by"`
	//TutorId           int       `json:"tutor_id"`
	//TutorType         string    `json:"tutor_type"`
	SearchSource string `json:"search_source"`
	//ResultSet         []string  `json:"result_set"`
	//ExtraFilters      []string  `json:"extra_filters"`
}
type DocMessage struct {
	MessageType string `json:"message_type"`
	//MessageId     int       `json:"message_id"`
	//SentById      int       `json:"sent_by_id"`
	//SentOn        float32 `json:"sent_on"`
	//LastUpdatedOn float32 `json:"last_updated_on"`
	//UserType      string    `json:"user_type"`
	//Platform      string    `json:"platform"`
	Body          string `json:"body"`
	ImageUrl      string `json:"image_url" validate:"url"`
	WebsearchText string `json:"websearch_text"`
}

func (r *RecommendRequest) SetDefault(qp url.Values) {
	//if r.Session.Rating == 0 {
	//	r.Session.Rating = -1
	//}
	if r.Session.SearchSource == "" {
		r.Session.SearchSource = constants.DocSearchSource
	}

	//r.Session.ResultSet = []string{"_id"}

	if badQuery, ok := qp["bad_query"]; ok && badQuery[0] == constants.True {
		r.BadQuery = true
	}

	if ocrOnly, ok := qp["ocr_only"]; ok && ocrOnly[0] == constants.True {
		r.OcrOnly = true
	}

	if source, ok := qp["source"]; ok && source[0] == constants.ApiSource.AutoSuggestionSource {
		r.Source = constants.ApiSource.AutoSuggestionSource
	} else if ok && source[0] == constants.ApiSource.BertSearchSource {
		r.Source = constants.ApiSource.BertSearchSource
	} else {
		r.Source = constants.ApiSource.SearchSource
	}

}

func (r RecommendRequest) Validate() error {
	if len(r.Message) == 0 {
		return fmt.Errorf("empty message list recieced")
	}
	return utils.ValidateRequestStructTag(r)
}
