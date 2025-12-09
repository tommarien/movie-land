package api

import (
	"log/slog"
	"net/http"
)

func handleInternalServerEror(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("unhandled error", "method", r.Method, "url", r.URL, "err", err)
	w.WriteHeader(http.StatusInternalServerError)
}
