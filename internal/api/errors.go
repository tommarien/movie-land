package api

import (
	"log/slog"
	"net/http"
)

func handleInternalServerEror(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("unhandled error", "method", r.Method, "url", r.URL, "err", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func handleNotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "resource not found"
	}

	writeJSON(w, http.StatusNotFound, map[string]any{
		"status":  http.StatusNotFound,
		"message": message,
	}, nil)
}
