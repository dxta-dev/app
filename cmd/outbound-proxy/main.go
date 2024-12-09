package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dxta-dev/app/internal/assert"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx := context.Background()

	port := os.Getenv("PORT")

	if port == "" {
		port = "1337"
	}

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			assert.Never(fmt.Sprintf("Could not listen on port %v %v", port, err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		assert.Never(fmt.Sprintf("Server forced to shutdown: %v", err))
	}
}
