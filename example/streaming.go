package example

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// StreamingChatExample demonstrates how to use the streaming API for chat completions
// It shows how to establish a streaming connection and process chunks of the response as they arrive
func StreamingChatExample() {
	// Get API key from environment variable
	apiKey := os.Getenv(APIKeyEnvVar)
	if apiKey == "" {
		log.Printf("Warning: Environment variable %s not set. Using placeholder.", APIKeyEnvVar)
		apiKey = "your-api-key-here" // This will prevent actual API calls
	}

	// Initialize the client with the API key
	client, err := grok.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Set the API base URL and a longer timeout for streaming
	client.BaseUrl = "https://api.x.ai/v1/"
	client.HttpClient.Timeout = 2 * time.Minute

	// Create a chat completion request with streaming enabled
	request := &grok.ChatCompletionRequest{
		Model: "grok-3", // Specify the model to use
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Write a short story about artificial intelligence",
			},
		},
		Stream:      true, // Enable streaming
		Temperature: 0.7,  // Control the randomness of the response
		MaxTokens:   1000, // Limit the response length
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fmt.Println("Establishing streaming connection...")

	// Create the streaming connection
	stream, err := grok.CreateChatCompletionStream(ctx, client, request)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	fmt.Println("Connection established. Streaming response:")
	fmt.Println("--------------------------------------------")

	// Variables to collect response and track tokens
	var completeResponse string
	tokenCount := 0
	startTime := time.Now()

	// Process the stream, reading each chunk as it arrives
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			break
		}
		if err != nil {
			log.Printf("Error receiving stream: %v", err)
			break
		}

		// Extract and display content from the chunk
		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			if content != "" {
				fmt.Print(content) // Print content immediately
				completeResponse += content
				tokenCount++
			}
		}
	}

	// Calculate and display statistics
	elapsed := time.Since(startTime)

	fmt.Println("\n--------------------------------------------")
	fmt.Printf("Stream completed in %.2f seconds\n", elapsed.Seconds())
	fmt.Printf("Received approximately %d tokens\n", tokenCount)
	if tokenCount > 0 {
		fmt.Printf("Average %.2f tokens per second\n", float64(tokenCount)/elapsed.Seconds())
	}

	// In a real application, you might want to use the complete response
	fmt.Printf("Total response length: %d characters\n", len(completeResponse))
}
