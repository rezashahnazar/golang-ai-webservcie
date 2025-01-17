// Package main provides a simple HTTP server with graceful shutdown capabilities
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// helloHandler responds with a JSON message "Hello, World!"
// It sets the appropriate content type header for JSON responses
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello, World!"}`))
}

func main() {
	// Configure the application to use all available CPU cores
	// This helps with concurrent request handling
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize router and define HTTP endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	// Health check endpoint for container orchestration and monitoring
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Configure the HTTP server with timeouts for security and reliability
	server := &http.Server{
		Addr:         ":8080",           // Listen on all interfaces, port 8080
		Handler:      mux,               // Use our routes
		ReadTimeout:  15 * time.Second,  // Maximum duration for reading entire request
		WriteTimeout: 15 * time.Second,  // Maximum duration for writing response
		IdleTimeout:  60 * time.Second,  // Maximum duration for keep-alive connections
	}

	// Graceful shutdown handler
	// Run in a separate goroutine to not block the main server
	go func() {
		// Create channel to listen for interrupt signals
		sigChan := make(chan os.Signal, 1)
		// Register for SIGINT (Ctrl+C) and SIGTERM signals
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		// Wait for interrupt signal
		<-sigChan
		
		// Create a deadline for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Start the server
	log.Printf("Server starting on port 8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
