package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/device-ops-dashboard/dashboard"
)

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}
	baseURL := os.Getenv("INFRAI_BASE_URL")
	client := dashboard.NewInfraiClient(baseURL, key, nil)
	channel := "developer-device-status"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.CreateChannel(ctx, channel); err != nil {
		log.Fatalf("create realtime channel: %v", err)
	}

	server := dashboard.NewServer(client, channel)
	httpServer := &http.Server{Addr: ":8080", Handler: server.Handler()}
	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- httpServer.ListenAndServe()
	}()

	log.Printf("device dashboard listening on http://localhost:8080")
	select {
	case <-stopCtx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("shut down dashboard: %v", err)
		}
		shutdownCancel()
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("serve dashboard: %v", err)
		}
	}

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cleanupCancel()
	if err := client.DeleteChannel(cleanupCtx, channel); err != nil {
		log.Fatalf("delete realtime channel: %v", err)
	}
}
