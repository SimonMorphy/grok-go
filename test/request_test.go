package test

import (
	"encoding/json"
	"testing"

	grok "github.com/SimonMorphy/grok-go"
)

// TestChatCompletionRequestSerialization tests that the chat completion request
// serializes to valid JSON with the expected structure.
func TestChatCompletionRequestSerialization(t *testing.T) {
	// Create a test request
	request := &grok.ChatCompletionRequest{
		Model: "grok-3",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "Hello, how are you?",
			},
		},
		MaxTokens:   100,
		Temperature: 0.7,
		TopP:        0.9,
		Stream:      false,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Deserialize to map to check fields
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["model"] != "grok-3" {
		t.Errorf("Expected model 'grok-3', got '%v'", result["model"])
	}

	// Check messages array
	messages, ok := result["messages"].([]interface{})
	if !ok {
		t.Fatalf("Messages is not an array")
	}
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	// Validate parameter types
	if _, ok := result["max_tokens"].(float64); !ok {
		t.Errorf("max_tokens should be a number, got %T", result["max_tokens"])
	}
	if _, ok := result["temperature"].(float64); !ok {
		t.Errorf("temperature should be a number, got %T", result["temperature"])
	}
}

// TestImageGenerationRequestSerialization tests that the image generation request
// serializes to valid JSON with the expected structure.
func TestImageGenerationRequestSerialization(t *testing.T) {
	// Create a test request
	request := &grok.ImageGenerationRequest{
		Model:  "grok-3-image",
		Prompt: "A serene mountain landscape with a lake",
		N:      1,
		Size:   "512x512",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Deserialize to map to check fields
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["model"] != "grok-3-image" {
		t.Errorf("Expected model 'grok-3-image', got '%v'", result["model"])
	}
	if result["prompt"] != "A serene mountain landscape with a lake" {
		t.Errorf("Expected prompt doesn't match")
	}
	if result["n"].(float64) != 1 {
		t.Errorf("Expected n=1, got %v", result["n"])
	}
}

// TestToolRequestSerialization tests that tool request structures
// serialize correctly to JSON.
func TestToolRequestSerialization(t *testing.T) {
	// Create a test tool
	calculatorTool := grok.Tool{
		Type: "function",
		Function: grok.Function{
			Name:        "calculate",
			Description: "Perform mathematical calculations",
			Parameters: &grok.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "The expression to calculate",
					},
				},
				Required: []string{"expression"},
			},
		},
	}

	// Create a request with the tool
	request := &grok.ChatCompletionRequest{
		Model: "grok-3",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Calculate 25 * 16",
			},
		},
		Tools:      []grok.Tool{calculatorTool},
		ToolChoice: "auto",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Deserialize to map to check fields
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check tools array
	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatalf("Tools is not an array")
	}
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}

	// Check tool choice field
	if result["tool_choice"] != "auto" {
		t.Errorf("Expected tool_choice 'auto', got '%v'", result["tool_choice"])
	}
}
