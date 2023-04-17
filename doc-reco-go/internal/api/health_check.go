package api

import (
	"doc-reco-go/internal/utils"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	utils.SendSuccess(w, nil, "System is up and running")
}