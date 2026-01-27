package api

import (
	"log/slog"
	"net/http"
)

func handleInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("unhandled error", "method", r.Method, "url", r.URL, "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)

	_, writeErr := w.Write([]byte(`{"status":500,"message":"internal server error"}`))
	if writeErr != nil {
		slog.Error("handleInternalServerError: write response", "method", r.Method, "url", r.URL, "err", writeErr)
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "resource not found"
	}

	statusCode := http.StatusNotFound

	writeJSON(w, r, statusCode, map[string]any{
		"status":  statusCode,
		"message": message,
	}, nil)
}

func handleConflict(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "conflict"
	}

	statusCode := http.StatusConflict

	writeJSON(w, r, statusCode, map[string]any{
		"status":  statusCode,
		"message": message,
	}, nil)
}

func handleBadRequest(w http.ResponseWriter, r *http.Request, message string, errors []string) {
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

	writeJSON(w, r, http.StatusBadRequest, response, nil)
}
