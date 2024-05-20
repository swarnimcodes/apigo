package utils

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(response)
}

func SendMessageResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := MessageResponse{Message: message}
	json.NewEncoder(w).Encode(response)
}
