package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	config2 "restaurant-management-app/cmd/config"
	router2 "restaurant-management-app/internal/router"
	"syscall"
	"time"
)

const shutdownTimeout = 10 * time.Second

func main() {
	// `server healthcheck` is used by the Docker HEALTHCHECK, since the
	// distroless image has no curl/wget.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(runHealthcheck())
	}

	config, err := config2.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config2.DBConnection(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %v", err)
	}
	defer sqlDB.Close()

	srv := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           router2.NewRouter(db, config.CORSAllowedOrigins),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Bind the port first so a busy port fails immediately with a clear error.
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		sqlDB.Close()
		log.Fatalf("Failed to listen on %s: %v", srv.Addr, err)
	}
	log.Printf("Server listening on %s", srv.Addr)

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	// Wait for Ctrl+C / docker stop, or for the server to fail on its own.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Printf("Server error: %v", err)
			sqlDB.Close()
			os.Exit(1)
		}
		return
	case <-ctx.Done():
	}

	// Stop accepting new requests and give in-flight ones time to finish.
	log.Println("Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}
	log.Println("Server stopped")
}

// runHealthcheck calls the local readiness endpoint and returns the exit code
// (0 = healthy, 1 = unhealthy). It reads only PORT, not the full config.
func runHealthcheck() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/health/ready")
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck failed:", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck failed: status", resp.StatusCode)
		return 1
	}
	return 0
}
