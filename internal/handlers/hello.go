package handlers

import (
	"net/http"
)

// HelloHandler handles hello-related endpoints
type HelloHandler struct {
	BaseHandler
}

// NewHelloHandler creates a new HelloHandler instance
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// HandleHello handles the GET /hello endpoint
func (h *HelloHandler) HandleHello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	h.RespondSuccess(w, map[string]string{
		"message": "Hello, World!",
	})
} 