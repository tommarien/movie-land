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
		slog.Error("main: uncaught error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("failed to parse envvars: %w", err)
	}

	ctx := context.Background()

	ds, err := createDatastore(ctx, cfg)
	if err != nil {
		return err
	}
	defer ds.Close()

	api := api.New(cfg, ds)
	if err = api.Start(ctx); err != nil {
		return fmt.Errorf("failed to start api: %w", err)
	}

	return nil
}

func createDatastore(ctx context.Context, cfg *config.Config) (*datastore.Store, error) {
	opts := pg.PoolOptions{
		ConnString:  cfg.DatabaseUrl,
		PingTimeout: cfg.DatabasePingTimeout,
	}

	pool, err := pg.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	ds := datastore.New(pool)
	return ds, nil
}
