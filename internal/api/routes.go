package api

import (
	"net/http"

	"github.com/tommarien/movie-land/internal/datastore"
)

func registerRoutes(mux *http.ServeMux, store *datastore.Store) {
	mux.HandleFunc("GET /healtz", handleHealtzIndex)

	mux.HandleFunc("GET /api/v1/genres", handleGenreIndex(store))
	mux.HandleFunc("GET /api/v1/genres/{id}", handleGenreGet(store))
	mux.HandleFunc("POST /api/v1/genres", handleGenrePost(store))
}
