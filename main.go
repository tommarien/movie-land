package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tommarien/movie-land/internal/config"
	"github.com/tommarien/movie-land/internal/datastore"
)

func main() {
	fmt.Println("Welcome to Movieland")

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("cfg: %+v\n", cfg)

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	ds := datastore.New(dbpool)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	err = ds.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")

	genre, err := ds.GetGenre(context.Background(), 1)
	if err != nil {
		log.Fatalf("err: %#v", err)
	}

	fmt.Printf("Genre: %+v\n", genre)
}
