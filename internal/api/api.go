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
	cfg   *config.Config
	store *datastore.Store
}

func New(cfg *config.Config, store *datastore.Store) *Api {
	return &Api{
		cfg:   cfg,
		store: store,
	}
}

func (api *Api) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	RegisterRoutes(mux, api.store)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", api.cfg.Port),
		Handler: mux,
	}

	errChan := make(chan error)
	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
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

	case <-shutdownCtx.Done():
		slog.Info("received a signal to shutdown")
		stop() // stop receiving signals, next ones will be handled by default behavior
	}

	slog.Info("gracefully shutting down")

	// Important to use a new context here,
	// as the shutdownCtx is already done when a signal is received
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutDownTimeout)
	defer cancel()

	err := svr.Shutdown(ctx)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		slog.Warn("graceful shutdown timed out, forcing exit")
		return svr.Close()
	case err != nil:
		return fmt.Errorf("api: shutdown error: %w", err)
	}

	return nil
}
