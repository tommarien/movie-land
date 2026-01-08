package api

import (
	"context"
	"net/http"

	"github.com/tommarien/movie-land/internal/datastore"
)

type GenreStore interface {
	ListGenres(ctx context.Context) ([]*datastore.Genre, error)
	GetGenre(ctx context.Context, ID int) (*datastore.Genre, error)
	InsertGenre(ctx context.Context, genre *datastore.Genre) error
}

func registerRoutes(
	mux *http.ServeMux,
	genreStore GenreStore) {
	mux.HandleFunc("GET /healtz", handleHealtzIndex)

	mux.HandleFunc("GET /api/v1/genres", handleGenreIndex(genreStore))
	mux.HandleFunc("GET /api/v1/genres/{id}", handleGenreGet(genreStore))
	mux.HandleFunc("POST /api/v1/genres", handleGenrePost(genreStore))
}
