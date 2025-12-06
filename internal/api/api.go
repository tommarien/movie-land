package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/tommarien/movie-land/internal/config"
	"github.com/tommarien/movie-land/internal/datastore"
)

const gracefulShutDownTimeout = 30 * time.Second

type Api struct {
	cfg *config.Config
	ds  *datastore.Store
}

func New(cfg *config.Config, ds *datastore.Store) *Api {
	return &Api{
		cfg: cfg,
		ds:  ds,
	}
}

func (api *Api) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	AddRoutes(mux)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", api.cfg.Port),
		Handler: mux,
	}

	errChan := make(chan error)
	notify, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("started listening", "port", api.cfg.Port)
		err := svr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("api: could not start: %w", err)

	case <-notify.Done():
		slog.Info("received a signal to shutdown")
		stop() // stop receiving signals, next ones will be handled by default behavior
	}

	slog.Info("gracefully shutting down")

	ctx, cancel := context.WithTimeout(ctx, gracefulShutDownTimeout)
	defer cancel()

	err := svr.Shutdown(ctx)
	switch {
	case errors.Is(err, http.ErrServerClosed):
		slog.Warn("graceful shutdown timed out, forcing exit")
		return svr.Close()
	case err != nil:
		return fmt.Errorf("api: shutdown error: %w", err)
	}

	return nil
}
