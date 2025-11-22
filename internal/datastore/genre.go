package datastore

import (
	"context"
	"errors"
	"time"

	"database/sql"
)

type Genre struct {
	ID        int
	Slug      string
	Name      sql.NullString
	CreatedAt time.Time
}

var (
	ErrGenreSlugExists = errors.New("datastore: genre with this slug already exists")
	ErrGenreNotFound   = errors.New("datastore: genre not found")
)

func (ds *Store) GetGenre(ctx context.Context, ID int) (*Genre, error) {
	var genre Genre

	const qry = `
	SELECT id, slug, name, created_at 
	FROM genres WHERE id=$1`

	err := ds.pool.QueryRow(
		ctx,
		qry,
		ID,
	).Scan(
		&genre.ID,
		&genre.Slug,
		&genre.Name,
		&genre.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGenreNotFound
		}
		return nil, err
	}

	return &genre, nil
}

func (ds *Store) InsertGenre(ctx context.Context, genre *Genre) error {
	if genre == nil {
		return errors.New("datastore: InsertGenre: genre is nil")
	}

	const qry = `
	INSERT INTO genres (slug, name) 
	VALUES ($1, $2) RETURNING id, created_at`

	err := ds.pool.QueryRow(
		ctx,
		qry,
		genre.Slug,
		genre.Name,
	).Scan(
		&genre.ID,
		&genre.CreatedAt,
	)

	if err != nil {
		if isConstraintViolation(err) {
			return ErrGenreSlugExists
		}
		return err
	}

	return nil
}

func (ds *Store) UpdateGenre(ctx context.Context, genre *Genre) error {
	if genre == nil {
		return errors.New("datastore: UpdateGenre: genre is nil")
	}

	const qry = `
	UPDATE genres 
	SET slug = $2, name = $3 
	WHERE id = $1`

	_, err := ds.pool.Exec(
		ctx,
		qry,
		genre.ID,
		genre.Slug,
		genre.Name,
	)

	if err != nil {
		if isConstraintViolation(err) {
			return ErrGenreSlugExists
		}
		return err
	}

	return nil
}
