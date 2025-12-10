package datastore

import (
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueConstraintViolationCode = "23505"

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	if pool == nil {
		panic("store: pool is nil")
	}

	return &Store{pool: pool}
}

func isConstraintViolation(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueConstraintViolationCode {
		return pgErr.ConstraintName
	}
	return ""
}

func (ds *Store) Close() {
	slog.Info("closing database connection pool")
	ds.pool.Close()
}
