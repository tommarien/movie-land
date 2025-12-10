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

func storeGenre(t *testing.T, dbpool *pgxpool.Pool, genre *datastore.Genre) int {
	t.Helper()

	if genre == nil {
		genre = &datastore.Genre{
			Slug: "drama",
			Name: sql.NullString{String: "Drama", Valid: true},
		}
	}

	var id int

	err := dbpool.QueryRow(
		context.Background(),
		`INSERT INTO genres (slug, name) VALUES ($1, $2) RETURNING id`,
		genre.Slug, genre.Name,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to insert genre: %v", err)
	}

	return id
}

func removeAllGenres(t *testing.T, dbpool *pgxpool.Pool) {
	_, err := dbpool.Exec(context.Background(), `DELETE FROM genres`)
	if err != nil {
		t.Fatal(err)
	}
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
		genreId := storeGenre(t, pool, nil)
		defer removeAllGenres(t, pool)

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

		if time.Until(genre.CreatedAt).Abs() > 2*time.Second {
			t.Errorf("expected genre.CreatedAt to be close to now, got %v", genre.CreatedAt)
		}
	})
}

func TestInsertGenre(t *testing.T) {
	pool := connect(t)
	ds := datastore.New(pool)

	t.Run("inserts a new genre", func(t *testing.T) {
		defer removeAllGenres(t, pool)

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

		if time.Until(genre.CreatedAt).Abs() > 2*time.Second {
			t.Errorf("expected genre.CreatedAt to be close to now, got %v", genre.CreatedAt)
		}
	})

	t.Run("returns error when inserting duplicate slug", func(t *testing.T) {
		storeGenre(t, pool, nil)
		defer removeAllGenres(t, pool)

		genre1 := &datastore.Genre{
			Slug: "drama",
		}

		err := ds.InsertGenre(context.Background(), genre1)
		if !errors.Is(err, datastore.ErrGenreSlugExists) {
			t.Fatalf("expected ErrGenreSlugExists, got %v", err)
		}
	})

	t.Run("returns error when inserting nil genre", func(t *testing.T) {
		defer removeAllGenres(t, pool)

		err := ds.InsertGenre(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if want := "store: InsertGenre: genre is nil"; err.Error() != want {
			t.Fatalf("expected error message %q, got %q", want, err.Error())
		}
	})
}

func TestUpdateGenre(t *testing.T) {
	pool := connect(t)
	ds := datastore.New(pool)

	t.Run("updates an existing genre", func(t *testing.T) {
		genreId := storeGenre(t, pool, nil)
		defer removeAllGenres(t, pool)

		genre, err := ds.GetGenre(context.Background(), genreId)
		if err != nil {
			t.Fatalf("failed to get genre: %v", err)
		}

		genre.Name = sql.NullString{String: "Updated Drama", Valid: true}

		err = ds.UpdateGenre(context.Background(), genre)
		if err != nil {
			t.Fatalf("failed to update genre: %v", err)
		}

		updatedGenre, err := ds.GetGenre(context.Background(), genreId)
		if err != nil {
			t.Fatalf("failed to get updated genre: %v", err)
		}

		if updatedGenre.Name.String != "Updated Drama" {
			t.Errorf("expected updated genre name to be 'Updated Drama', got %s",
				updatedGenre.Name.String)
		}
	})

	t.Run("returns error when updating to duplicate slug", func(t *testing.T) {
		defer removeAllGenres(t, pool)

		dramaId := storeGenre(t, pool, nil)
		_ = storeGenre(t, pool, &datastore.Genre{
			Slug: "comedy",
			Name: sql.NullString{String: "Comedy", Valid: true},
		})

		drama, err := ds.GetGenre(context.Background(), dramaId)
		if err != nil {
			t.Fatalf("failed to get drama genre: %v", err)
		}

		drama.Slug = "comedy"

		err = ds.UpdateGenre(context.Background(), drama)
		if !errors.Is(err, datastore.ErrGenreSlugExists) {
			t.Fatalf("expected ErrGenreSlugExists, got %v", err)
		}
	})

	t.Run("returns error when updating non-existent genre", func(t *testing.T) {
		defer removeAllGenres(t, pool)

		genre := &datastore.Genre{
			ID:   99999,
			Slug: "non-existent",
		}

		err := ds.UpdateGenre(context.Background(), genre)
		if !errors.Is(err, datastore.ErrGenreNotFound) {
			t.Fatalf("expected ErrGenreNotFound, got %v", err)
		}
	})

	t.Run("returns error when updating nil genre", func(t *testing.T) {
		err := ds.UpdateGenre(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if want := "store: UpdateGenre: genre is nil"; err.Error() != want {
			t.Fatalf("expected error message %q, got %q", want, err.Error())
		}
	})
}

func TestListGenres(t *testing.T) {
	pool := connect(t)
	ds := datastore.New(pool)

	t.Run("lists all genres", func(t *testing.T) {
		defer removeAllGenres(t, pool)

		storeGenre(t, pool, &datastore.Genre{
			Slug: "drama",
			Name: sql.NullString{String: "Drama", Valid: true},
		})
		storeGenre(t, pool, &datastore.Genre{
			Slug: "comedy",
			Name: sql.NullString{String: "Comedy", Valid: true},
		})

		genres, err := ds.ListGenres(context.Background())
		if err != nil {
			t.Fatalf("failed to list genres: %v", err)
		}

		if len(genres) != 2 {
			t.Fatalf("expected 2 genres, got %d", len(genres))
		}
	})

	t.Run("returns empty list when no genres exist", func(t *testing.T) {
		pool := connect(t)
		ds := datastore.New(pool)
		defer removeAllGenres(t, pool)

		genres, err := ds.ListGenres(context.Background())
		if err != nil {
			t.Fatalf("failed to list genres: %v", err)
		}

		if len(genres) != 0 {
			t.Fatalf("expected 0 genres, got %d", len(genres))
		}
	})
}
