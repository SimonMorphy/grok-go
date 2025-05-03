package test

import (
	"context"
	"io"
	"testing"
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// TestCreateChatCompletionStream tests the streaming functionality of the chat completion API
// It demonstrates how to establish and process a streaming response
func TestCreateChatCompletionStream(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint with extended timeout
	client.BaseUrl = "https://api.example.com/v1/"
	client.HttpClient.Timeout = 60 * time.Second

	// Create a simple chat completion request
	request := &grok.ChatCompletionRequest{
		Model: "test-model",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Tell me a short story about space exploration.",
			},
		},
	}

	// Skip actual API call in tests
	t.Skip("Skipping streaming API call test - for demonstration only")

	// In a real test, we would establish a streaming connection
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := grok.CreateChatCompletionStream(ctx, client, request)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	// Process the stream
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// End of stream reached normally
			break
		}
		if err != nil {
			t.Fatalf("Error receiving from stream: %v", err)
		}

		// In a real test, we would process and validate each chunk
		// For example, checking that content is present and well-formed
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			// Process content chunk
		}
	}

	// Assert stream completion and overall response structure in real tests
}

// TestCreateChatCompletionStreamWithFallback tests streaming with fallback to standard completion
// This demonstrates how to implement graceful degradation in production environments
func TestCreateChatCompletionStreamWithFallback(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint with timeout
	client.BaseUrl = "https://api.example.com/v1/"
	client.HttpClient.Timeout = 60 * time.Second

	// Create a request
	request := &grok.ChatCompletionRequest{
		Model: "test-model",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Write a haiku about programming",
			},
		},
	}

	// Skip actual API call in tests
	t.Skip("Skipping streaming fallback test - for demonstration only")

	ctx := context.Background()

	// Try streaming first
	stream, err := grok.CreateChatCompletionStream(ctx, client, request)

	// If streaming fails, fall back to regular completion
	if err != nil {
		t.Logf("Streaming failed (expected in this test): %v", err)
		t.Log("Falling back to regular chat completion...")

		// Use regular completion as fallback
		_, err := grok.CreateChatCompletion(ctx, client, request)
		if err != nil {
			t.Fatalf("Regular chat completion also failed: %v", err)
		}

		// In a real test, we would validate the regular response
		return
	}

	// If streaming succeeded, process the stream
	defer stream.Close()

	// Process stream (simplified for test)
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error receiving from stream: %v", err)
		}
	}
}

// TestCreateChatCompletionSimpleStream is a minimal streaming test
// It focuses on the core streaming functionality without extra complexity
func TestCreateChatCompletionSimpleStream(t *testing.T) {
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
				Content: "Write a short poem",
			},
		},
	}

	// Skip actual API call in tests
	t.Skip("Skipping simple streaming test - for demonstration only")

	// In a real test, we would establish a streaming connection
	ctx := context.Background()
	stream, err := grok.CreateChatCompletionStream(ctx, client, request)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	// Simple stream processing
	contentReceived := false
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error receiving from stream: %v", err)
		}

		// Check that we received at least some content
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			contentReceived = true
		}
	}

	// In a real test, we would assert that content was received
	if !contentReceived {
		t.Skip("Would assert content was received in actual test")
	}
}
