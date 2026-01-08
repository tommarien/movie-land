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

func handleGenreGet(store GenreStore) http.HandlerFunc {
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

func handleGenreIndex(store GenreStore) http.HandlerFunc {
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

func handleGenrePost(store GenreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Slug string `json:"slug"`
			Name string `json:"name"`
		}

		err := readJSON(w, r, &input)
		if err != nil {
			handleBadRequest(w, err.Error(), nil)
			return
		}

		v := validator.New()
		v.Required("slug", input.Slug)
		v.MaxLength("slug", input.Slug, 40)
		v.Slug("slug", input.Slug)
		v.MaxLength("name", input.Name, 40)

		if !v.IsValid() {
			handleBadRequest(w, "", v.GetErrors())
			return
		}

		genre := &datastore.Genre{
			Slug: input.Slug,
		}

		if input.Name != "" {
			genre.Name.String = input.Name
			genre.Name.Valid = true
		}

		err = store.InsertGenre(r.Context(), genre)
		if err != nil {
			if errors.Is(err, datastore.ErrGenreSlugExists) {
				handleConflict(w, "genre with this slug already exists")
				return
			}
			handleInternalServerEror(w, r, err)
			return
		}

		dto := mapGenre(genre)
		err = writeJSON(w, http.StatusCreated, map[string]any{
			"data": dto,
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
