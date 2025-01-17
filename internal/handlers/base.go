package handlers

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct{}

// RespondJSON sends a JSON response with the given status code
func (h *BaseHandler) RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// RespondError sends an error response with the given status code
func (h *BaseHandler) RespondError(w http.ResponseWriter, status int, message string) {
	h.RespondJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

// RespondSuccess sends a success response with the given data
func (h *BaseHandler) RespondSuccess(w http.ResponseWriter, data interface{}) {
	h.RespondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
} 