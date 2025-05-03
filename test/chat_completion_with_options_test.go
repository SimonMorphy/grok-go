package test

import (
	"context"
	"testing"

	grok "github.com/SimonMorphy/grok-go"
)

// TestCreateChatCompletionWithOptions tests chat completion with advanced options
// It verifies that the API properly handles custom parameters like temperature and token limits
func TestCreateChatCompletionWithOptions(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint
	client.BaseUrl = "https://api.example.com/v1/"

	// Create a request with advanced options
	request := &grok.ChatCompletionRequest{
		Model: "test-model",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a professional math teacher. Answer math questions concisely.",
			},
			{
				Role:    "user",
				Content: "Explain the Pythagorean theorem with a simple example.",
			},
		},
		// Advanced options
		MaxTokens:        500, // Limit output length
		Temperature:      0.7, // Control creativity
		TopP:             0.9, // Nucleus sampling parameter
		PresencePenalty:  0.0, // Presence penalty
		FrequencyPenalty: 0.5, // Frequency penalty
	}

	// Skip actual API call in tests
	t.Skip("Skipping API call test - for demonstration only")

	// In a real test, we would make the API call
	ctx := context.Background()
	_, err = grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		t.Fatalf("Chat completion with options failed: %v", err)
	}

	// Assert response data in real tests
}

// TestCreateChatCompletionWithJsonMode tests the JSON response format option
// It verifies that the API can return structured JSON data when requested
func TestCreateChatCompletionWithJsonMode(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint
	client.BaseUrl = "https://api.example.com/v1/"

	// Configure JSON response format
	jsonFormat := &grok.ResponseFormat{
		Type: "json_object",
	}

	// Create a request that should return JSON
	request := &grok.ChatCompletionRequest{
		Model: "test-model",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Give me information about 5 science fiction books in JSON format",
			},
		},
		ResponseFormat: jsonFormat,
		Temperature:    0.7,
	}

	// Skip actual API call in tests
	t.Skip("Skipping API call test - for demonstration only")

	// In a real test, we would make the API call and verify JSON format
	ctx := context.Background()
	_, err = grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		t.Fatalf("Chat completion with JSON mode failed: %v", err)
	}

	// Assert JSON response structure in real tests
}
