package api

import (
	"log/slog"
	"net/http"
)

func handleInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("unhandled error", "method", r.Method, "url", r.URL, "err", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func handleNotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "resource not found"
	}

	statusCode := http.StatusNotFound

	err := writeJSON(w, statusCode, map[string]any{
		"status":  statusCode,
		"message": message,
	}, nil)
	if err != nil {
		slog.Error("handleNotFound: writeJSON", "err", err)
	}
}

func handleConflict(w http.ResponseWriter, message string) {
	if message == "" {
		message = "conflict"
	}

	statusCode := http.StatusConflict

	err := writeJSON(w, statusCode, map[string]any{
		"status":  statusCode,
		"message": message,
	}, nil)
	if err != nil {
		slog.Error("handleConflict: writeJSON", "err", err)
	}
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

	err := writeJSON(w, http.StatusBadRequest, response, nil)
	if err != nil {
		slog.Error("handleBadRequest: writeJSON", "err", err)
	}
}
