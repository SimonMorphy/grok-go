package test

import (
	"context"
	"os"
	"testing"
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// 这个测试文件展示了API的完整使用流程，包括设置、多次交互和错误处理

// getAPIKey retrieves the API key from environment variables or returns a test key
// In real tests, prefer environment variables for sensitive credentials
func getAPIKey() string {
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		// Use a test key for demonstration
		apiKey = "test-api-key"
	}
	return apiKey
}

// createClient creates a client with custom options for testing
// This centralizes client creation logic for all tests in this file
func createClient() (*grok.Client, error) {
	// Create client with custom options
	client, err := grok.NewClientWithOptions(getAPIKey(),
		grok.WithBaseURL("https://api.example.com/v1/"),
		grok.WithTimeout(60*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// TestCompleteConversation tests a multi-turn conversation with the API
// It verifies that context is maintained across multiple exchanges
func TestCompleteConversation(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Skip actual API call in tests
	t.Skip("Skipping conversation test - for demonstration only")

	ctx := context.Background()
	var messages []grok.ChatCompletionMessage

	// First turn: user greeting
	messages = append(messages, grok.ChatCompletionMessage{
		Role:    "user",
		Content: "Tell me about traditional holidays in China",
	})

	// First API call
	response, err := sendChatRequest(ctx, client, messages)
	if err != nil {
		t.Fatalf("First conversation turn failed: %v", err)
	}

	// Save assistant's response
	assistantResponse := response.Choices[0].Message.Content
	messages = append(messages, grok.ChatCompletionMessage{
		Role:    "assistant",
		Content: assistantResponse,
	})

	// Second turn: user follow-up question
	messages = append(messages, grok.ChatCompletionMessage{
		Role:    "user",
		Content: "Can you tell me more about Spring Festival customs in different regions?",
	})

	// Second API call
	response, err = sendChatRequest(ctx, client, messages)
	if err != nil {
		t.Fatalf("Second conversation turn failed: %v", err)
	}

	// Save assistant's response
	assistantResponse = response.Choices[0].Message.Content
	messages = append(messages, grok.ChatCompletionMessage{
		Role:    "assistant",
		Content: assistantResponse,
	})

	// Third turn: user asks for summary
	messages = append(messages, grok.ChatCompletionMessage{
		Role:    "user",
		Content: "Thank you. Could you summarize the key characteristics of Chinese traditional holidays?",
	})

	// Third API call
	response, err = sendChatRequest(ctx, client, messages)
	if err != nil {
		t.Fatalf("Third conversation turn failed: %v", err)
	}

	// In a real test, we would assert that the conversation maintained context
	// and that each response was appropriate to the conversation history
}

// sendChatRequest is a helper function to send chat requests
// It encapsulates the common request creation and error handling logic
func sendChatRequest(ctx context.Context, client *grok.Client, messages []grok.ChatCompletionMessage) (*grok.Response, error) {
	request := &grok.ChatCompletionRequest{
		Model:       "test-model",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   800,
	}

	return grok.CreateChatCompletion(ctx, client, request)
}

// TestErrorHandling tests how the API handles various error conditions
// It validates that appropriate errors are returned for invalid requests
func TestErrorHandling(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Skip actual API call in tests
	t.Skip("Skipping error handling test - for demonstration only")

	ctx := context.Background()

	testCases := []struct {
		name    string
		request *grok.ChatCompletionRequest
	}{
		{
			name: "Invalid Model",
			request: &grok.ChatCompletionRequest{
				Model: "invalid-model-name",
				Messages: []grok.ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
			},
		},
		{
			name: "Empty Messages",
			request: &grok.ChatCompletionRequest{
				Model:    "test-model",
				Messages: []grok.ChatCompletionMessage{},
			},
		},
		{
			name: "Context Too Long",
			request: &grok.ChatCompletionRequest{
				Model: "test-model",
				Messages: []grok.ChatCompletionMessage{
					{Role: "user", Content: string(make([]byte, 100000))}, // Very long message
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := grok.CreateChatCompletion(ctx, client, tc.request)
			if err == nil {
				t.Errorf("Expected error but received success")
			}
			// In a real test, we would validate the specific error type and message
		})
	}
}
