package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	defaultOpenRouterBaseURL = "https://openrouter.ai/api/v1"
	defaultOpenRouterModel  = "anthropic/claude-3.5-sonnet"
	defaultMaxTokens = 1000
	defaultTemperature float32 = 0.7
)

// OpenRouterProvider implements the Provider interface for OpenRouter
type OpenRouterProvider struct {
	apiKey     string
	baseURL    string
	modelName  string
	httpClient *http.Client
	defaultMaxTokens int
	defaultTemperature float32
}

// OpenRouterRequest represents the request structure for OpenRouter API
type OpenRouterRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

// NewOpenRouterProvider creates a new OpenRouter provider instance
func NewOpenRouterProvider(options *ProviderOptions) (Provider, error) {
	if options.APIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = defaultOpenRouterBaseURL
	}

	modelName := options.ModelName
	if modelName == "" {
		modelName = defaultOpenRouterModel
	}

	maxTokens := defaultMaxTokens
	if options.MaxTokens > 0 {
		maxTokens = options.MaxTokens
	}

	temperature := defaultTemperature
	if options.Temperature > 0 {
		temperature = options.Temperature
	}

	return &OpenRouterProvider{
		apiKey:     options.APIKey,
		baseURL:    baseURL,
		modelName:  modelName,
		httpClient: &http.Client{},
		defaultMaxTokens: maxTokens,
		defaultTemperature: temperature,
	}, nil
}

func (p *OpenRouterProvider) GetName() string {
	return "openrouter"
}

func (p *OpenRouterProvider) GetDefaultModel() string {
	return defaultOpenRouterModel
}

func (p *OpenRouterProvider) ChatStream(ctx context.Context, messages []Message, options *ProviderOptions) (io.ReadCloser, error) {
	// Use provider options if provided, otherwise use defaults
	maxTokens := p.defaultMaxTokens
	temperature := p.defaultTemperature
	if options != nil {
		if options.MaxTokens > 0 {
			maxTokens = options.MaxTokens
		}
		if options.Temperature > 0 {
			temperature = options.Temperature
		}
	}

	reqBody := OpenRouterRequest{
		Model:       p.modelName,
		Messages:    messages,
		Stream:      true,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/rezashahnazar")
	req.Header.Set("X-Title", "Golang AI Web Service")
	req.Header.Set("Accept", "text/event-stream")

	log.Printf("Sending request to OpenRouter: %+v", reqBody)
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("OpenRouter response headers: %+v", resp.Header)
	return resp.Body, nil
}

func init() {
	RegisterProvider("openrouter", NewOpenRouterProvider)
} 