package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dxta-dev/app/internal/outbound-proxy/host"
)

type ApiKeyDetails struct {
	Host   int
	ApiKey string
}

type RateLimitInfo struct {
	RetryAfter int64
}

var rateLimitMap sync.Map
var client = &http.Client{}

func proxyHandler(hostId int) http.HandlerFunc {

	var h host.Host

	switch hostId {
	case host.GITHUB:
		h = host.NewGitHubHost()
	default:
		panic("")
	}

	return func(w http.ResponseWriter, req *http.Request) {
		_, err := h.UnwrapRequest(req)
		if err != nil {
		}

		var resp *http.Response

		resp, err = client.Do(req)
		if err != nil {
		}

		_, err = h.UnwrapResponse(resp)
		if err != nil {
		}

		if resp.StatusCode == http.StatusTooManyRequests {
		}

	}
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "1337"
	}

	http.HandleFunc("/github", proxyHandler(host.GITHUB))
	http.HandleFunc("/gitlab", proxyHandler(host.GITLAB))

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 1337 %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
