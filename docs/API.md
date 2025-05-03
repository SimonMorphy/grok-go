# API Documentation

This document provides a comprehensive reference for the grok-go library API.

## Client

### Creating a Client

```go
// Create a default client
client, err := grok.NewClient(apiKey)

// Create a client with custom options
client, err := grok.NewClientWithOptions(
    apiKey,
    grok.WithBaseURL("https://api.x.ai/v1/"),
    grok.WithTimeout(60 * time.Second),
    grok.WithEndpoint("chat/completions"),
)

// Create specialized clients
chatClient, err := grok.CreateChatCompletionClient(apiKey)
imageClient, err := grok.CreateImageGenerationClient(apiKey)
modelsClient, err := grok.CreateModelsClient(apiKey)
```

### Client Configuration Options

```go
// Set custom base URL
grok.WithBaseURL(baseURL string)

// Set custom timeout
grok.WithTimeout(timeout time.Duration)

// Set custom endpoint
grok.WithEndpoint(endpoint string)

// Set custom HTTP client
grok.WithHTTPClient(httpClient *http.Client)
```

## Chat Completions

### Creating a Chat Completion

```go
// Basic chat completion
request := &grok.ChatCompletionRequest{
    Model: "grok-3",
    Messages: []grok.ChatCompletionMessage{
        {
            Role:    "user",
            Content: "Hello, how are you?",
        },
    },
}

response, err := grok.CreateChatCompletion(ctx, client, request)

// Access the response
content := response.Choices[0].Message.Content
```

### Working with Messages

```go
// Create different types of messages
systemMsg := grok.NewTextMessage("system", "You are a helpful assistant.")
userMsg := grok.NewTextMessage("user", "Hello!")
assistantMsg := grok.NewAssistantMessage("Hi there!")
toolMsg := grok.NewToolCallMessage("Result: 42", "tool_123")

// Create a message collection
messages := grok.NewMessages()
messages.AddMessage(systemMsg)
messages.AddMessage(userMsg)

// Add a name to a message
namedMsg := grok.WithName(userMsg, "user_123")
```

### Streaming Chat Completions

```go
request := &grok.ChatCompletionRequest{
    Model:  "grok-3",
    Messages: []grok.ChatCompletionMessage{
        {
            Role:    "user",
            Content: "Tell me a story",
        },
    },
    Stream: true,
}

stream, err := grok.CreateChatCompletionStream(ctx, client, request)
if err != nil {
    // Handle error
}
defer stream.Close()

// Process the stream
for {
    response, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        // Handle error
        break
    }
    
    // Handle each chunk
    chunk := response.Choices[0].Delta.Content
    fmt.Print(chunk)
}
```

## Function Calling and Tools

### Defining Tools

```go
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
```

### Using Tools in Requests

```go
request := &grok.ChatCompletionRequest{
    Model: "grok-3",
    Messages: []grok.ChatCompletionMessage{
        {
            Role:    "user",
            Content: "Calculate 25 times 16",
        },
    },
    Tools:      []grok.Tool{calculatorTool},
    ToolChoice: "auto", // Let the model choose which tool to use
}

response, err := grok.CreateChatCompletion(ctx, client, request)
```

### Processing Tool Calls

```go
// Check if there are tool calls in the response
if len(response.Choices) > 0 && len(response.Choices[0].Message.ToolCalls) > 0 {
    // Process each tool call
    for _, toolCall := range response.Choices[0].Message.ToolCalls {
        // Execute the tool based on the function name
        var result string
        if toolCall.Function.Name == "calculate" {
            // Parse the arguments
            var args map[string]interface{}
            json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
            
            // Execute the function and get the result
            result = executeCalculator(args)
        }

        // Add the tool result to messages for a follow-up request
        messages := append(request.Messages, grok.ChatCompletionMessage{
            Role:      "assistant",
            Content:   response.Choices[0].Message.Content,
            ToolCalls: response.Choices[0].Message.ToolCalls,
        })

        messages = append(messages, grok.ChatCompletionMessage{
            Role:       "tool",
            Content:    result,
            ToolCallID: toolCall.ID,
        })

        // Create a follow-up request with the tool results
        followupRequest := &grok.ChatCompletionRequest{
            Model:    request.Model,
            Messages: messages,
        }

        // Send the follow-up request
        followupResponse, err := grok.CreateChatCompletion(ctx, client, followupRequest)
    }
}
```

## Image Generation

### Creating Images

```go
request := &grok.ImageGenerationRequest{
    Model:  "grok-3-image",
    Prompt: "A beautiful mountain landscape",
    Size:   "1024x1024",
    N:      1,
}

response, err := grok.CreateImage(ctx, client, request)
```

### Processing Image Responses

```go
// Handle image URLs
if response.Data[0].URL != "" {
    imageURL := response.Data[0].URL
    // Download the image...
}

// Handle Base64-encoded images
if response.Data[0].B64JSON != "" {
    base64Data := response.Data[0].B64JSON
    // Decode and save the image...
}
```

## Error Handling

```go
response, err := grok.CreateChatCompletion(ctx, client, request)
if err != nil {
    // Handle API errors
    if apiErr, ok := err.(*grok.APIError); ok {
        fmt.Printf("API Error: %s, Type: %s\n", apiErr.Message, apiErr.Type)
    } else {
        // Handle network or other errors
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Full Documentation

For more detailed documentation, please refer to the [package documentation](https://pkg.go.dev/github.com/SimonMorphy/grok-go) on pkg.go.dev. 