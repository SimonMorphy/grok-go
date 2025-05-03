package test

import (
	"testing"
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// TestClientCreation tests client creation with different configuration options
func TestClientCreation(t *testing.T) {
	// Get test API key
	apiKey := GetTestAPIKey(t)

	// Test default client creation
	t.Run("DefaultClient", func(t *testing.T) {
		client, err := grok.NewClient(apiKey)
		if err != nil {
			t.Fatalf("Failed to create default client: %v", err)
		}

		// Verify default values
		if client.BaseUrl != "https://api.x.ai/v1/" {
			t.Errorf("Unexpected default base URL: %s", client.BaseUrl)
		}
		if client.Endpoint != "chat/completions" {
			t.Errorf("Unexpected default endpoint: %s", client.Endpoint)
		}
	})

	// Test client creation with custom options
	t.Run("ClientWithOptions", func(t *testing.T) {
		customTimeout := 60 * time.Second
		customBaseURL := "https://custom-api.example.com/v1/"
		customEndpoint := "custom/endpoint"

		client, err := grok.NewClientWithOptions(
			apiKey,
			grok.WithBaseURL(customBaseURL),
			grok.WithTimeout(customTimeout),
			grok.WithEndpoint(customEndpoint),
		)
		if err != nil {
			t.Fatalf("Failed to create client with options: %v", err)
		}

		// Verify custom values
		if client.BaseUrl != customBaseURL {
			t.Errorf("Expected base URL %s, got %s", customBaseURL, client.BaseUrl)
		}
		if client.Endpoint != customEndpoint {
			t.Errorf("Expected endpoint %s, got %s", customEndpoint, client.Endpoint)
		}
		if client.Timeout != customTimeout {
			t.Errorf("Expected timeout %v, got %v", customTimeout, client.Timeout)
		}
	})

	// Test client creation with invalid API key
	t.Run("InvalidAPIKey", func(t *testing.T) {
		_, err := grok.NewClient("")
		if err == nil {
			t.Error("Expected error with empty API key, got nil")
		}
	})
}

// TestSpecializedClientCreation tests creation of clients for specific purposes
func TestSpecializedClientCreation(t *testing.T) {
	// Get test API key
	apiKey := GetTestAPIKey(t)

	// Test chat completion client
	t.Run("ChatCompletionClient", func(t *testing.T) {
		client, err := grok.CreateChatCompletionClient(apiKey)
		if err != nil {
			t.Fatalf("Failed to create chat completion client: %v", err)
		}

		if client.Endpoint != "chat/completions" {
			t.Errorf("Expected endpoint 'chat/completions', got '%s'", client.Endpoint)
		}
	})

	// Test image generation client
	t.Run("ImageGenerationClient", func(t *testing.T) {
		client, err := grok.CreateImageGenerationClient(apiKey)
		if err != nil {
			t.Fatalf("Failed to create image generation client: %v", err)
		}

		if client.Endpoint != "images/generations" {
			t.Errorf("Expected endpoint 'images/generations', got '%s'", client.Endpoint)
		}
	})

	// Test models client
	t.Run("ModelsClient", func(t *testing.T) {
		client, err := grok.CreateModelsClient(apiKey)
		if err != nil {
			t.Fatalf("Failed to create models client: %v", err)
		}

		if client.Endpoint != "models" {
			t.Errorf("Expected endpoint 'models', got '%s'", client.Endpoint)
		}
	})
}

// TestClientValidation tests client configuration validation
func TestClientValidation(t *testing.T) {
	// Test API key validation
	t.Run("APIKeyValidation", func(t *testing.T) {
		client := &grok.Client{
			BaseUrl:  "https://api.x.ai/v1/",
			Endpoint: "chat/completions",
			Timeout:  30 * time.Second,
		}

		// Missing API key
		err := grok.ValidateClient(client)
		if err == nil {
			t.Error("Expected error for missing API key, got nil")
		}

		// Add API key
		client.ApiKey = "test-api-key"
		err = grok.ValidateClient(client)
		if err != nil {
			t.Errorf("Unexpected error after adding API key: %v", err)
		}
	})

	// Test URL validation
	t.Run("URLValidation", func(t *testing.T) {
		client := &grok.Client{
			ApiKey:   "test-api-key",
			Endpoint: "chat/completions",
			Timeout:  30 * time.Second,
		}

		// Missing base URL
		err := grok.ValidateClient(client)
		if err == nil {
			t.Error("Expected error for missing base URL, got nil")
		}

		// Invalid base URL format
		client.BaseUrl = "invalid-url"
		err = grok.ValidateClient(client)
		if err == nil {
			t.Error("Expected error for invalid base URL format, got nil")
		}

		// Add valid base URL
		client.BaseUrl = "https://api.x.ai/v1/"
		err = grok.ValidateClient(client)
		if err != nil {
			t.Errorf("Unexpected error after adding valid base URL: %v", err)
		}
	})
}
