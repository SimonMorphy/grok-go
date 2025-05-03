package example

import (
	"context"
	"fmt"
	"log"
	"os"

	grok "github.com/SimonMorphy/grok-go"
)

// APIKeyEnvVar is the environment variable name for the API key
const APIKeyEnvVar = "GROK_API_KEY"

// ChatExample demonstrates a basic chat completion request using the X.AI API
// It shows how to initialize a client, create a request, and process the response
func ChatExample() {
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

	// Set the API base URL
	client.BaseUrl = "https://api.x.ai/v1/"

	// Create a simple chat completion request
	request := &grok.ChatCompletionRequest{
		Model: "grok-3", // Specify the model to use
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Explain quantum computing in simple terms",
			},
		},
		Temperature: 0.7, // Control the randomness of the response
		MaxTokens:   500, // Limit the response length
	}

	// Create a context for the request
	ctx := context.Background()

	// Send the request to the API
	response, err := grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		log.Fatalf("Chat completion request failed: %v", err)
	}

	// Process and display the response
	if len(response.Choices) > 0 {
		fmt.Println("AI Response:")
		fmt.Println(response.Choices[0].Message.Content)

		// Print token usage statistics
		fmt.Printf("\nToken usage - Prompt: %d, Completion: %d, Total: %d\n",
			response.Usage.PromptTokens,
			response.Usage.CompletionTokens,
			response.Usage.TotalTokens)
	} else {
		fmt.Println("No response content received")
	}
}
