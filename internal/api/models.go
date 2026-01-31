package api

import (
	"time"

	"github.com/tommarien/movie-land/internal/datastore"
	"github.com/tommarien/movie-land/internal/validator"
)

type CreateGenreBody struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (i CreateGenreBody) Validate(v *validator.Validator) {
	if v == nil {
		return
	}

	v.Required("slug", i.Slug)
	v.MaxLength("slug", i.Slug, 40)
	v.Slug("slug", i.Slug)
	v.MaxLength("name", i.Name, 40)
}

type GenreResource struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func NewGenreResource(genre *datastore.Genre) *GenreResource {
	dto := &GenreResource{
		ID:        genre.ID,
		Slug:      genre.Slug,
		CreatedAt: genre.CreatedAt,
	}

	if genre.Name.Valid {
		dto.Name = genre.Name.String
	}
	return dto
}
