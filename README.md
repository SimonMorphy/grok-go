# Grok-Go: Go Client Library for X.AI API

This is a Go client library for the X.AI API, allowing you to easily integrate Grok AI models into your Go applications.

## Features

- Complete API coverage for X.AI endpoints
- Support for chat completions with the Grok model
- Streaming responses support
- Function calling and tool use capabilities
- Image generation support
- Type-safe request and response handling
- Comprehensive error handling

## Installation

```bash
go get github.com/SimonMorphy/grok-go
```

## Authentication

To use the X.AI API, you need an API key. You can set it as an environment variable:

```bash
export GROK_API_KEY="your-api-key-here"
```

Or provide it directly in your code (not recommended for production):

```go
client, err := grok.NewClient("your-api-key-here")
```

## Quick Start

### Chat Completion Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	
	grok "github.com/SimonMorphy/grok-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set")
	}
	
	// Initialize client
	client, err := grok.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	
	// Create chat completion request
	request := &grok.ChatCompletionRequest{
		Model: "grok-3",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "What is artificial intelligence?",
			},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}
	
	// Send request
	ctx := context.Background()
	response, err := grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	
	// Print response
	if len(response.Choices) > 0 {
		fmt.Println(response.Choices[0].Message.Content)
	}
}
```

## Documentation

For detailed documentation and more examples, check out the following:

- [API Documentation](docs/API.md)
- [Example Code](example/)
- [Chat Completion Guide](docs/chat.md)
- [Streaming Guide](docs/streaming.md)
- [Tool Calling Guide](docs/tools.md)
- [Image Generation Guide](docs/images.md)

## Examples

The `example` directory contains complete, runnable examples demonstrating various API features:

- Basic chat completion: [example/chat.go](example/chat.go)
- Streaming responses: [example/streaming.go](example/streaming.go)
- Function calling: [example/tools.go](example/tools.go)
- Image generation: [example/image.go](example/image.go)

## Testing

Run the tests with:

```bash
export GROK_API_KEY="your-api-key-here"  # Optional for integration tests
go test ./...
```

For unit tests only (which don't require an API key):

```bash
go test ./... -short
```

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 