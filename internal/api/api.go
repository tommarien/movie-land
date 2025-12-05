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
)

const gracefulShutDownTimeout = 30 * time.Second

var logger = slog.With("name", "api")

type Api struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Api {
	return &Api{
		cfg: cfg,
	}
}

func (api *Api) Start() error {
	mux := http.NewServeMux()
	AddRoutes(mux)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", api.cfg.Port),
		Handler: mux,
	}

	errChan := make(chan error)
	notify, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("started listening", "port", api.cfg.Port)
		err := svr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("api: could not start: %w", err)

	case <-notify.Done():
		logger.Info("received an os signal to shutdown")
		stop() // stop receiving signals, next ones will be handled by default behavior
	}

	logger.Info("gracefully shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutDownTimeout)
	defer cancel()

	err := svr.Shutdown(ctx)
	switch {
	case errors.Is(err, http.ErrServerClosed):
		logger.Warn("graceful shutdown timed out, forcing exit")
		return svr.Close()
	case err != nil:
		return fmt.Errorf("api: shutdown error: %w", err)
	}

	return nil
}
