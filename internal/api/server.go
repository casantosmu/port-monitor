package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Server(ctx context.Context, port string, routes http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		<-ctx.Done()

		log.Printf("[server] shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(shutdownCtx)
	}()

	log.Printf("[server] listening on http://localhost:%s", port)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Println("[server] stopped")
	return nil
}
