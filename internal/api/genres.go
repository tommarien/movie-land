package api

import (
	"errors"
	"net/http"

	"github.com/tommarien/movie-land/internal/datastore"
	"github.com/tommarien/movie-land/internal/validator"
)

func handleGenreGet(store GenreStore) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := getIntParam(r, "id")
		if err != nil {
			return NewNotFoundStatusError("genre not found")
		}

		genre, err := store.GetGenre(r.Context(), id)
		if err != nil {
			if errors.Is(err, datastore.ErrGenreNotFound) {
				return NewNotFoundStatusError("genre not found")
			}
			return err
		}

		return writeJSON(w, r, http.StatusOK, map[string]any{
			"data": NewGenreResource(genre),
		}, nil)
	}
}

func handleGenreIndex(store GenreStore) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		genres, err := store.ListGenres(r.Context())
		if err != nil {
			return err
		}

		data := make([]*GenreResource, 0, len(genres))
		for _, g := range genres {
			data = append(data, NewGenreResource(g))
		}

		return writeJSON(w, r, http.StatusOK, map[string]any{
			"data": data,
		}, nil)
	}
}

func handleGenrePost(store GenreStore) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		var payload CreateGenreBody

		err := readJSON(w, r, &payload)
		if err != nil {
			return NewBadRequestStatusError(err.Error(), nil)
		}

		v := validator.New()
		payload.Validate(v)

		if !v.IsValid() {
			return NewBadRequestStatusError("", v.GetErrors())
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
				return NewConflictStatusError("genre with this slug already exists")
			}
			return err
		}

		return writeJSON(w, r, http.StatusCreated, map[string]any{
			"data": NewGenreResource(genre),
		}, nil)
	}
}
