package api

import (
	"net/http"

	"github.com/tommarien/movie-land/internal/datastore"
)

func RegisterRoutes(mux *http.ServeMux, store *datastore.Store) {
	mux.HandleFunc("GET /healtz", handleHealtzIndex)

	mux.HandleFunc("GET /api/v1/genres", handleGenreIndex(store))
}
