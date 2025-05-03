package grok_go

import (
	"net/http"
	"time"

	"github.com/SimonMorphy/grok-go/constants"
	"github.com/SimonMorphy/grok-go/constants/errors"
)

// Default configuration values
const (
	DefaultBaseURL  = "https://api.x.ai/v1/"
	DefaultTimeout  = 30 * time.Second
	DefaultEndpoint = "chat/completions"
)

// ClientOption defines a function type that can modify a Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseUrl = baseURL
	}
}

// WithTimeout sets a custom timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

// WithEndpoint sets a custom API endpoint
func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.Endpoint = endpoint
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HttpClient = httpClient
	}
}

// NewClient creates a new client with default configuration
// It takes an API key and returns a configured client or an error
func NewClient(apiKey string) (*Client, error) {
	return NewClientWithOptions(apiKey)
}

// NewClientWithOptions creates a new client with custom options
// It takes an API key and optional configuration options
// Returns a configured client or an error if validation fails
func NewClientWithOptions(apiKey string, options ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New(constants.ErrnoInvalidApiKey)
	}

	// Create client with default values
	client := &Client{
		ApiKey:   apiKey,
		BaseUrl:  DefaultBaseURL,
		Timeout:  DefaultTimeout,
		Endpoint: DefaultEndpoint,
	}

	// Apply custom options
	for _, option := range options {
		option(client)
	}

	// Set default HTTP client if not provided
	if client.HttpClient == nil {
		client.HttpClient = &http.Client{
			Timeout: client.Timeout,
		}
	}

	// Validate client configuration
	if err := ValidateClient(client); err != nil {
		return nil, err
	}

	return client, nil
}

// CreateChatCompletionClient creates a client specifically configured for chat completions
// It sets the endpoint to "chat/completions" and applies any additional options
func CreateChatCompletionClient(apiKey string, options ...ClientOption) (*Client, error) {
	// Force the chat completions endpoint
	return NewClientWithOptions(
		apiKey,
		append(options, WithEndpoint("chat/completions"))...,
	)
}

// CreateImageGenerationClient creates a client specifically configured for image generation
// It sets the endpoint to "images/generations" and applies any additional options
func CreateImageGenerationClient(apiKey string, options ...ClientOption) (*Client, error) {
	// Force the image generations endpoint
	return NewClientWithOptions(
		apiKey,
		append(options, WithEndpoint("images/generations"))...,
	)
}

// CreateModelsClient creates a client specifically configured for listing models
// It sets the endpoint to "models" and applies any additional options
func CreateModelsClient(apiKey string, options ...ClientOption) (*Client, error) {
	// Force the models endpoint
	return NewClientWithOptions(
		apiKey,
		append(options, WithEndpoint("models"))...,
	)
}
