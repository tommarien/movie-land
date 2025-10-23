package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueConstraintViolationCode = "23505"

type Database struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func (db *Database) Connect(ctx context.Context) error {
	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database: connect: %w", err)
	}
	return nil
}

func isConstraintViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueConstraintViolationCode
}
