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

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rezashahnazar/golang-ai-webservcie/internal/handlers"
	"github.com/rezashahnazar/golang-ai-webservcie/internal/middleware"
	"github.com/rezashahnazar/golang-ai-webservcie/pkg/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	// Load .env file if it exists
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: error reading config file: %v", err)
		}
	} else {
		log.Printf("Successfully loaded .env file")
	}

	// Bind environment variables
	viper.SetEnvPrefix("")  // No prefix for env vars
	viper.AutomaticEnv()    // Automatically bind env vars

	// Initialize configuration with defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("READ_TIMEOUT", 15)
	viper.SetDefault("WRITE_TIMEOUT", 15)
	viper.SetDefault("IDLE_TIMEOUT", 60)
	viper.SetDefault("RATE_LIMIT_RPS", 100)
	viper.SetDefault("OPENROUTER_API_KEY", "")

	// Debug: Print all configuration values
	log.Printf("Configuration values:")
	log.Printf("PORT: %s", viper.GetString("PORT"))
	log.Printf("OPENROUTER_API_KEY length: %d", len(viper.GetString("OPENROUTER_API_KEY")))

	// Validate required configuration
	apiKey := viper.GetString("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}
}

func main() {
	// Configure the application to use all available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize server
	srv := server.NewServer(logger)

	// Add middleware
	srv.Use(
		middleware.Recovery(logger),
		middleware.Logging(logger),
		middleware.SecurityHeaders,
		middleware.Metrics,
		middleware.RateLimit(viper.GetInt("RATE_LIMIT_RPS")),
	)

	// Initialize handlers
	helloHandler := handlers.NewHelloHandler()
	chatHandler := handlers.NewChatHandler(viper.GetString("OPENROUTER_API_KEY"))

	// Register routes
	srv.HandleFunc("/hello", helloHandler.HandleHello)
	srv.HandleFunc("/v1/chat", chatHandler.Handle)
	srv.Handle("/metrics", promhttp.Handler())
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		helloHandler.RespondSuccess(w, map[string]string{"status": "healthy"})
	})

	// Start server
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
