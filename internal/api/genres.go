package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tommarien/movie-land/internal/datastore"
)

type genreStore interface {
	ListGenres(ctx context.Context) ([]*datastore.Genre, error)
	GetGenre(ctx context.Context, ID int) (*datastore.Genre, error)
}

type GenreDto struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func handleGenreGet(store genreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIntParam(r, "id")
		if err != nil {
			handleNotFound(w, "genre not found")
			return
		}

		genre, err := store.GetGenre(r.Context(), id)
		if err != nil {
			if errors.Is(err, datastore.ErrGenreNotFound) {
				handleNotFound(w, "genre not found")
				return
			}
			handleInternalServerEror(w, r, err)
			return
		}

		err = writeJSON(w, http.StatusOK, map[string]any{
			"data": mapGenre(genre),
		}, nil)

		if err != nil {
			handleInternalServerEror(w, r, err)
			return
		}
	}
}

func handleGenreIndex(store genreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genres, err := store.ListGenres(r.Context())
		if err != nil {
			handleInternalServerEror(w, r, err)
			return
		}

		data := make([]*GenreDto, 0, len(genres))
		for _, g := range genres {
			dto := mapGenre(g)
			data = append(data, dto)
		}

		err = writeJSON(w, http.StatusOK, map[string]any{
			"data": data,
		}, nil)

		if err != nil {
			handleInternalServerEror(w, r, err)
			return
		}
	}
}

func mapGenre(genre *datastore.Genre) *GenreDto {
	dto := &GenreDto{
		ID:        genre.ID,
		Slug:      genre.Slug,
		CreatedAt: genre.CreatedAt,
	}

	if genre.Name.Valid {
		dto.Name = genre.Name.String
	}
	return dto
}
