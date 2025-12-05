package main

import (
	"fmt"
	"log/slog"
	"os"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/tommarien/movie-land/internal/api"
	"github.com/tommarien/movie-land/internal/config"
	// "github.com/tommarien/movie-land/internal/datastore"
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

	api := api.New(cfg)

	err = api.Start()
	if err != nil {
		return fmt.Errorf("failed to start api: %w", err)
	}

	return nil

	// dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	// if err != nil {
	// 	return err
	// }
	//
	// defer dbpool.Close()
	//
	// ds := datastore.New(dbpool)
	//
	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	// defer cancel()
	//
	// err = ds.Connect(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// fmt.Println("Successfully connected to the database")
	//
	// genre, err := ds.GetGenre(context.Background(), 1)
	// if err != nil {
	// 	return err
	// }
	//
	// fmt.Printf("Genre: %+v\n", genre)
	//
	// return nil
}
