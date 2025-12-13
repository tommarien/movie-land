package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/tommarien/movie-land/internal/api"
	"github.com/tommarien/movie-land/internal/config"
	"github.com/tommarien/movie-land/internal/datastore"
	"github.com/tommarien/movie-land/internal/pg"
)

func main() {
	if err := run(); err != nil {
		slog.Error("uncaught error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("failed to parse envvars: %w", err)
	}

	ctx := context.Background()

	pool, err := pg.Connect(ctx, pg.PoolOptions{
		ConnString:  cfg.DatabaseUrl,
		PingTimeout: cfg.DatabasePingTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	defer pool.Close()

	store := datastore.New(pool)

	api := api.New(cfg, store)
	if err = api.Start(ctx); err != nil {
		return fmt.Errorf("failed to start api: %w", err)
	}

	return nil
}
