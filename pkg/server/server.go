package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server represents the HTTP server and its dependencies
type Server struct {
	router     *http.ServeMux
	logger     *zap.Logger
	middleware []func(http.Handler) http.Handler
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer(logger *zap.Logger) *Server {
	return &Server{
		router: http.NewServeMux(),
		logger: logger,
	}
}

// Use adds middleware to the server
func (s *Server) Use(middleware ...func(http.Handler) http.Handler) {
	s.middleware = append(s.middleware, middleware...)
}

// Handle registers a new route with the server
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.router.Handle(pattern, handler)
}

// HandleFunc registers a new route with a handler function
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	var handler http.Handler = s.router
	// Apply middleware in reverse order
	for i := len(s.middleware) - 1; i >= 0; i-- {
		handler = s.middleware[i](handler)
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080" // Fallback default
	}

	addr := fmt.Sprintf(":%s", port)
	s.logger.Info("Server configuration", 
		zap.String("address", addr),
		zap.Int("read_timeout", viper.GetInt("READ_TIMEOUT")),
		zap.Int("write_timeout", viper.GetInt("WRITE_TIMEOUT")),
		zap.Int("idle_timeout", viper.GetInt("IDLE_TIMEOUT")),
	)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(viper.GetInt("READ_TIMEOUT")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("WRITE_TIMEOUT")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("IDLE_TIMEOUT")) * time.Second,
	}

	s.logger.Info("Server starting", zap.String("port", port))
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
} 