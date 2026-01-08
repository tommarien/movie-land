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

	statusCode := http.StatusNotFound

	writeJSON(w, statusCode, map[string]any{
		"status":  statusCode,
		"message": message,
	}, nil)
}

func handleConflict(w http.ResponseWriter, message string) {
	if message == "" {
		message = "conflict"
	}

	statusCode := http.StatusConflict

	writeJSON(w, statusCode, map[string]any{
		"status":  statusCode,
		"message": message,
	}, nil)
}

func handleBadRequest(w http.ResponseWriter, message string, errors []string) {
	if message == "" {
		message = "bad request"
	}

	response := map[string]any{
		"status":  http.StatusBadRequest,
		"message": message,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	writeJSON(w, http.StatusBadRequest, response, nil)
}
