package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func errorHandler(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	slog.ErrorContext(r.Context(), "error occurred", "message", message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(ErrorMessage{Message: message}); err != nil {
		slog.ErrorContext(r.Context(), "failed to encode response", "error", err)
	}
}
