package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rezashahnazar/golang-ai-webservcie/internal/ai"
)

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Messages    []ai.Message `json:"messages"`
	Provider    string       `json:"provider"`
	Model       string       `json:"model,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature float32     `json:"temperature,omitempty"`
}

// ChatHandler handles AI chat requests
type ChatHandler struct {
	BaseHandler
	apiKey string
}

// NewChatHandler creates a new ChatHandler instance
func NewChatHandler(apiKey string) *ChatHandler {
	return &ChatHandler{
		apiKey: apiKey,
	}
}

// Handle processes chat requests and streams responses
func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		h.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Messages) == 0 {
		h.RespondError(w, http.StatusBadRequest, "Messages cannot be empty")
		return
	}

	log.Printf("Received chat request: %+v", req)

	// Set default provider if not specified
	if req.Provider == "" {
		req.Provider = "openrouter"
	}

	// Create provider options from request
	options := &ai.ProviderOptions{
		APIKey:      h.apiKey,
		ModelName:   req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// Get AI provider
	provider, err := ai.GetProvider(req.Provider, options)
	if err != nil {
		h.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to get provider: %v", err))
		return
	}

	// Set response headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable buffering in Nginx

	// Stream the response
	stream, err := provider.ChatStream(r.Context(), req.Messages, options)
	if err != nil {
		log.Printf("Error getting stream: %v", err)
		h.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stream chat: %v", err))
		return
	}
	defer stream.Close()

	// Read and process the stream line by line
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		log.Printf("Received line: %s", line)

		// Skip non-data lines
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Remove "data: " prefix and handle special cases
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return
		}

		// Try to parse the data as JSON
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			log.Printf("Error parsing chunk: %v", err)
			continue
		}

		// Extract content from the chunk
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			content := chunk.Choices[0].Delta.Content
			log.Printf("Sending content: %s", content)
			
			// Format as SSE
			_, err := fmt.Fprintf(w, "data: %s\n\n", content)
			if err != nil {
				log.Printf("Error writing response: %v", err)
				return
			}
			
			// Flush if supported
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
	}
} 