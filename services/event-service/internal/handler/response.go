package handler

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code      string      `json:"Code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, r *http.Request, status int, code, msg string, details interface{}) {
	respondJSON(w, status, APIError{
		Code:      code,
		Message:   msg,
		RequestID: GetRequestID(r.Context()),
		Details:   details,
	})
}
