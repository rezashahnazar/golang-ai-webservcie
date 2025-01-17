package ai

import (
	"context"
	"io"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ProviderOptions contains common configuration options for AI providers
type ProviderOptions struct {
	APIKey     string
	BaseURL    string
	ModelName  string
	MaxTokens  int
	Temperature float32
}

// Provider defines the interface that all AI providers must implement
type Provider interface {
	// ChatStream streams the chat completion response
	ChatStream(ctx context.Context, messages []Message, options *ProviderOptions) (io.ReadCloser, error)
	
	// GetName returns the provider name
	GetName() string
	
	// GetDefaultModel returns the default model for this provider
	GetDefaultModel() string
}

// Factory is a function type that creates new Provider instances
type Factory func(options *ProviderOptions) (Provider, error)

// registry stores provider factories by name
var registry = make(map[string]Factory)

// RegisterProvider registers a new provider factory
func RegisterProvider(name string, factory Factory) {
	registry[name] = factory
}

// GetProvider returns a provider instance by name
func GetProvider(name string, options *ProviderOptions) (Provider, error) {
	factory, exists := registry[name]
	if !exists {
		return nil, ErrProviderNotFound{Name: name}
	}
	return factory(options)
}

// ErrProviderNotFound is returned when a provider is not found in the registry
type ErrProviderNotFound struct {
	Name string
}

func (e ErrProviderNotFound) Error() string {
	return "ai provider not found: " + e.Name
} 