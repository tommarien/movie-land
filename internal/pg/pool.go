package pg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolOptions struct {
	ConnString  string
	PingTimeout time.Duration
}

func Connect(ctx context.Context, opts PoolOptions) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, opts.ConnString)
	if err != nil {
		return nil, fmt.Errorf("pool: could not create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, opts.PingTimeout)
	defer cancel()

	slog.Info("pinging database to verify connection", "timeout", opts.PingTimeout)
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pool: could not ping database: %w", err)
	}

	return pool, nil
}
