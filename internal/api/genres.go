package api

import (
	"context"
	"net/http"
	"time"

	"github.com/tommarien/movie-land/internal/datastore"
)

type genreStore interface {
	ListGenres(ctx context.Context) ([]*datastore.Genre, error)
}

type GenreDto struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func handleGenreIndex(store genreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbGenres, err := store.ListGenres(r.Context())

		if err != nil {
			handleInternalServerEror(w, r, err)
			return
		}

		data := make([]*GenreDto, 0, len(dbGenres))
		for _, genre := range dbGenres {
			dto := mapGenre(genre)
			data = append(data, dto)
		}

		err = encode(w, http.StatusOK, map[string]any{
			"data": data,
		})

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
