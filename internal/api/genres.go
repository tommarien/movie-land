package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/tommarien/movie-land/internal/datastore"
	"github.com/tommarien/movie-land/internal/validator"
)

type GenreDto struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func handleGenreGet(store GenreStore) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := getIntParam(r, "id")
		if err != nil {
			return NewNotFound("genre not found")
		}

		genre, err := store.GetGenre(r.Context(), id)
		if err != nil {
			if errors.Is(err, datastore.ErrGenreNotFound) {
				return NewNotFound("genre not found")
			}
			return err
		}

		return writeJSON(w, r, http.StatusOK, map[string]any{
			"data": mapGenre(genre),
		}, nil)
	}
}

func handleGenreIndex(store GenreStore) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		genres, err := store.ListGenres(r.Context())
		if err != nil {
			return err
		}

		data := make([]*GenreDto, 0, len(genres))
		for _, g := range genres {
			dto := mapGenre(g)
			data = append(data, dto)
		}

		return writeJSON(w, r, http.StatusOK, map[string]any{
			"data": data,
		}, nil)
	}
}

func handleGenrePost(store GenreStore) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		var payload CreateGenreRequest

		err := readJSON(w, r, &payload)
		if err != nil {
			return NewBadRequest(err.Error(), nil)
		}

		v := validator.New()
		payload.Validate(v)

		if !v.IsValid() {
			return NewBadRequest("", v.GetErrors())
		}

		genre := &datastore.Genre{
			Slug: payload.Slug,
		}

		if payload.Name != "" {
			genre.Name.String = payload.Name
			genre.Name.Valid = true
		}

		err = store.InsertGenre(r.Context(), genre)
		if err != nil {
			if errors.Is(err, datastore.ErrGenreSlugExists) {
				return NewConflict("genre with this slug already exists")
			}
			return err
		}

		dto := mapGenre(genre)
		return writeJSON(w, r, http.StatusCreated, map[string]any{
			"data": dto,
		}, nil)
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
