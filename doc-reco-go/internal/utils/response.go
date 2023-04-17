package utils

import (
	"encoding/json"
	"net/http"
)

func SendSuccess(w http.ResponseWriter, data interface{}, message string) {
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"data":        data,
		"message":     message,
		"status_code": http.StatusOK,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SendSuccessV2(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SendBadRequestError(w http.ResponseWriter, data interface{}, message string) {
	w.WriteHeader(http.StatusBadRequest)
	response := map[string]interface{}{
		"data":    data,
		"message": message,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SendInternalServerError(w http.ResponseWriter, data interface{}, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	response := map[string]interface{}{
		"data":    data,
		"message": message,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
