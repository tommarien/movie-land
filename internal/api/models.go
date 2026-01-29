package api

import (
	"github.com/tommarien/movie-land/internal/validator"
)

type CreateGenreRequest struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (i CreateGenreRequest) Validate(v *validator.Validator) {
	if v == nil {
		return
	}

	v.Required("slug", i.Slug)
	v.MaxLength("slug", i.Slug, 40)
	v.Slug("slug", i.Slug)
	v.MaxLength("name", i.Name, 40)
}
