package database

import (
	"context"
	"errors"
	"time"
)

type Genre struct {
	ID        int
	Slug      string
	Name      string
	CreatedAt time.Time
}

var ErrGenreSlugExists = errors.New("genre with this slug already exists")

func (db *Database) GetGenre(ctx context.Context, ID int) (*Genre, error) {
	var genre Genre

	const sql = `
	SELECT id, slug, name, created_at 
	FROM genres WHERE id=$1`

	err := db.pool.
		QueryRow(ctx, sql, ID).
		Scan(&genre.ID, &genre.Slug, &genre.Name, &genre.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &genre, nil
}

func (db *Database) InsertGenre(ctx context.Context, genre *Genre) error {
	if genre == nil {
		return errors.New("database: InsertGenre: genre is nil")
	}

	const sql = `
	INSERT INTO genres (slug, name) 
	VALUES ($1, $2) RETURNING id, created_at`

	err := db.pool.
		QueryRow(ctx, sql, genre.Slug, genre.Name).
		Scan(&genre.ID, &genre.CreatedAt)

	if err != nil {
		if isConstraintViolation(err) {
			return ErrGenreSlugExists
		}
		return err
	}

	return nil
}

func (db *Database) UpdateGenre(ctx context.Context, genre *Genre) error {
	if genre == nil {
		return errors.New("database: UpdateGenre: genre is nil")
	}

	const sql = `
	UPDATE genres 
	SET slug = $2, name = $3 
	WHERE id = $1`

	_, err := db.pool.Exec(ctx, sql, genre.ID, genre.Slug, genre.Name)
	if err != nil {
		if isConstraintViolation(err) {
			return ErrGenreSlugExists
		}
		return err
	}

	return nil
}
