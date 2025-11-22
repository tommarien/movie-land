package datastore_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tommarien/movie-land/internal/datastore"
)

func connect(t *testing.T) *pgxpool.Pool {
	t.Helper()

	databaseUrl := os.Getenv("DATABASE_TEST_URL")
	if databaseUrl == "" {
		t.Fatal("DATABASE_TEST_URL not set")
	}

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { dbpool.Close() })

	return dbpool
}

func createGenre(t *testing.T, dbpool *pgxpool.Pool) int {
	t.Helper()

	var id int

	err := dbpool.QueryRow(
		context.Background(),
		`INSERT INTO genres (slug, name) VALUES ($1, $2) RETURNING id`,
		"drama", "Drama",
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to insert genre: %v", err)
	}

	t.Cleanup(func() {
		_, err := dbpool.Exec(context.Background(), `DELETE FROM genres`)
		if err != nil {
			t.Fatal(err)
		}
	})

	return id
}

func TestGetGenre(t *testing.T) {
	pool := connect(t)
	ds := datastore.New(pool)

	t.Run("returns ErrGenreNotFound if genre does not exist", func(t *testing.T) {
		_, err := ds.GetGenre(context.Background(), 1)
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, datastore.ErrGenreNotFound) {
			t.Errorf("expected error to be ErrGenreNotFound, got %v", err)
		}
	})

	t.Run("returns the genre if it exists", func(t *testing.T) {
		genreId := createGenre(t, pool)

		genre, err := ds.GetGenre(context.Background(), genreId)
		if err != nil {
			t.Fatalf("failed to get genre: %v", err)
		}

		if genre.ID != genreId {
			t.Errorf("expected genre.ID to be %d, got %d", genreId, genre.ID)
		}

		if genre.Slug != "drama" {
			t.Errorf("expected genre.Slug to be drama, got %s", genre.Slug)
		}

		if genre.Name.String != "Drama" {
			t.Errorf("expected genre.Name to be Drama, got %v", genre.Name)
		}

		if genre.CreatedAt.IsZero() {
			t.Error("expected genre.CreatedAt to be set")
		}

		if time.Until(genre.CreatedAt).Abs() > time.Second {
			t.Errorf("expected genre.CreatedAt to be close to now, got %v", genre.CreatedAt)
		}
	})
}

func TestInsertGenre(t *testing.T) {
	pool := connect(t)
	ds := datastore.New(pool)

	t.Run("inserts a new genre", func(t *testing.T) {
		genre := &datastore.Genre{
			Slug: "comedy",
			Name: sql.NullString{String: "Comedy", Valid: true},
		}

		err := ds.InsertGenre(context.Background(), genre)
		if err != nil {
			t.Fatalf("failed to insert genre: %v", err)
		}

		if genre.ID == 0 {
			t.Error("expected genre.ID to be set")
		}

		if time.Until(genre.CreatedAt).Abs() > time.Second {
			t.Errorf("expected genre.CreatedAt to be close to now, got %v", genre.CreatedAt)
		}

		t.Cleanup(func() {
			_, err := pool.Exec(context.Background(), `DELETE FROM genres`)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("returns error when inserting duplicate slug", func(t *testing.T) {
		createGenre(t, pool)

		genre1 := &datastore.Genre{
			Slug: "drama",
		}

		err := ds.InsertGenre(context.Background(), genre1)
		if !errors.Is(err, datastore.ErrGenreSlugExists) {
			t.Fatalf("expected ErrGenreSlugExists, got %v", err)
		}
	})
}
