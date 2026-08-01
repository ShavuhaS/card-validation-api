package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondJSON(w http.ResponseWriter, statusCode int, body any) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		slog.Error("Failed to marshal the response body", "msg", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonBody)
}
