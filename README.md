# Grok-Go: Go Client Library for X.AI API

[![Go Reference](https://pkg.go.dev/badge/github.com/SimonMorphy/grok-go.svg)](https://pkg.go.dev/github.com/SimonMorphy/grok-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonMorphy/grok-go)](https://goreportcard.com/report/github.com/SimonMorphy/grok-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

*Read this in [中文](README_CN.md)*

A lightweight, type-safe Go client library for the X.AI API. This library enables seamless integration with X.AI's Grok models in Go applications.

## Features

- **Complete API Coverage**: Access all X.AI API endpoints
- **Strong Typing**: First-class Go types for all requests and responses
- **Streaming Support**: Handle streaming responses efficiently
- **Function Calling**: Utilize function calling and tool capabilities
- **Image Generation**: Generate images with the Grok model
- **Error Handling**: Comprehensive error handling with detailed error messages

## Installation

```bash
go get github.com/SimonMorphy/grok-go
```

Requires Go 1.18 or later.

## Quick Start

### Authentication

Set your API key as an environment variable (recommended):

```bash
export GROK_API_KEY="your-api-key-here"
```

Or provide it directly in code:

```go
client, err := grok.NewClient("your-api-key-here") // Not recommended for production
```

### Basic Chat Example

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

## Advanced Usage

### Function Calling

```go
// Define a calculator tool
calculatorTool := grok.Tool{
    Type: "function",
    Function: grok.Function{
        Name:        "calculate",
        Description: "Perform mathematical calculations",
        Parameters: &grok.FunctionParameters{
            Type: "object",
            Properties: map[string]interface{}{
                "operation": map[string]interface{}{
                    "type":        "string",
                    "enum":        []string{"add", "subtract", "multiply", "divide"},
                    "description": "The mathematical operation to perform",
                },
                "x": map[string]interface{}{
                    "type":        "number",
                    "description": "The first operand",
                },
                "y": map[string]interface{}{
                    "type":        "number",
                    "description": "The second operand",
                },
            },
            Required: []string{"operation", "x", "y"},
        },
    },
}

// Create a request using the tool
request := &grok.ChatCompletionRequest{
    Model: "grok-3",
    Messages: []grok.ChatCompletionMessage{
        {
            Role:    "user",
            Content: "Calculate 25 times 16",
        },
    },
    Tools:      []grok.Tool{calculatorTool},
    ToolChoice: "auto",
}
```

### Streaming Responses

```go
request.Stream = true
stream, err := grok.CreateChatCompletionStream(ctx, client, request)
if err != nil {
    log.Fatalf("Failed to create stream: %v", err)
}
defer stream.Close()

// Process the stream
for {
    response, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Printf("Stream error: %v", err)
        break
    }
    
    // Handle chunk
    if len(response.Choices) > 0 {
        chunk := response.Choices[0].Delta.Content
        if chunk != "" {
            fmt.Print(chunk)
        }
    }
}
```

### Image Generation

```go
imageClient, err := grok.CreateImageGenerationClient(apiKey)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}

request := &grok.ImageGenerationRequest{
    Model:  "grok-3-image",
    Prompt: "A futuristic city with flying cars and tall skyscrapers",
    Size:   "1024x1024",
    N:      1,
}

response, err := grok.CreateImage(ctx, imageClient, request)
if err != nil {
    log.Fatalf("Image generation failed: %v", err)
}

// Process the image URL or Base64 data
for i, imageData := range response.Data {
    if imageData.URL != "" {
        fmt.Printf("Image %d URL: %s\n", i+1, imageData.URL)
        // Download the image...
    } else if imageData.B64JSON != "" {
        fmt.Printf("Image %d received as Base64 data\n", i+1)
        // Save the Base64 image...
    }
}
```

## Examples

For complete, runnable examples, check the [`example`](example/) directory:

- [Basic Chat Completion](example/chat.go)
- [Streaming Responses](example/streaming.go)
- [Function Calling](example/tools.go)
- [Image Generation](example/image.go)

## Documentation

- [Package Documentation](https://pkg.go.dev/github.com/SimonMorphy/grok-go)
- [API Reference](docs/API.md)
- [Image Generation Guide](docs/images.md)
- [X.AI API Reference](https://platform.x.ai/docs/api-reference)

## Testing

Run all tests (requires API key for integration tests):

```bash
export GROK_API_KEY="your-api-key-here"
go test ./...
```

Run unit tests only (no API key required):

```bash
go test ./... -short
```

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Feel free to:

- Report bugs
- Request features
- Submit pull requests

Please ensure your code passes all tests and follows Go best practices. 