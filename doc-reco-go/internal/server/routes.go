package server

import (
	"doc-reco-go/internal/api"
	"doc-reco-go/internal/api/indexing"
	"doc-reco-go/internal/api/recommendation"
	"doc-reco-go/internal/api/tutor_search"
	"encoding/json"
	"net/http"
)

func (server *Server) InitializeRoutes() {
	r := server.Router

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello from recommendation service"))
	})

	//health check
	r.HandleFunc("/v2/health-check/", api.HealthCheck).Methods(http.MethodGet)

	// some service calls with trailing slash, some without. TODO: keep only one route
	r.HandleFunc("/v2/recommend", recommendation.Recommend).Methods(http.MethodPost)
	r.HandleFunc("/v2/recommend/", recommendation.Recommend).Methods(http.MethodPost)

	r.HandleFunc("/recommendation-used/", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"data": "Successfully logged"})
	})

	r.HandleFunc("/suggest/questions", tutor_search.SuggestQuestions).Methods(http.MethodPost)

	r.HandleFunc("/suggest/concepts", tutor_search.SuggestConcepts).Methods(http.MethodPost)
	// index questions
	r.HandleFunc("/indexcontent/question/",indexing.IndexQuestions).Methods(http.MethodPost)

}
