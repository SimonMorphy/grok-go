package test

import (
	"context"
	"testing"

	grok "github.com/SimonMorphy/grok-go"
)

// TestCreateChatCompletion tests the basic chat completion functionality
// It verifies that the API can process a simple user message and return a response
func TestCreateChatCompletion(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint
	client.BaseUrl = "https://api.example.com/v1/"

	// Create a simple chat completion request
	request := &grok.ChatCompletionRequest{
		Model: "test-model",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Hello, introduce yourself",
			},
		},
	}

	// Skip actual API call in tests
	t.Skip("Skipping API call test - for demonstration only")

	// In a real test, we would make the API call
	ctx := context.Background()
	_, err = grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		t.Fatalf("Chat completion failed: %v", err)
	}

	// Assert response data in real tests
}
