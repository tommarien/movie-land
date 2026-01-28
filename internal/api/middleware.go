package api

import (
	"log/slog"
	"net/http"
)

type HandlerFuncWithError func(w http.ResponseWriter, r *http.Request) error

func withError(h HandlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if httpErr, ok := err.(*HttpStatusError); ok {
				handleHTTPStatusError(w, r, httpErr)
			} else {
				handleInternalServerError(w, r, err)
			}
		}
	}
}
func handleHTTPStatusError(w http.ResponseWriter, r *http.Request, err *HttpStatusError) {
	response := map[string]any{
		"status":  err.StatusCode(),
		"message": err.Message(),
	}

	if len(err.Errors()) > 0 {
		response["errors"] = err.Errors()
	}

	if err := writeJSON(w, r, err.StatusCode(), response, nil); err != nil {
		handleInternalServerError(w, r, err)
	}
}

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
