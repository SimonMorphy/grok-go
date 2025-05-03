package test

import (
	"context"
	"encoding/json"
	"testing"

	grok "github.com/SimonMorphy/grok-go"
)

// TestCreateChatCompletionWithTools tests the integration of function tools in chat completions
// It verifies that the API can process tool definitions and return tool calls
func TestCreateChatCompletionWithTools(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint
	client.BaseUrl = "https://api.example.com/v1/"

	// Define calculator tool
	calculatorTool := grok.Tool{
		Type: "function",
		Function: grok.Function{
			Name:        "calculator",
			Description: "Performs mathematical calculations",
			Parameters: &grok.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "The expression to calculate, e.g. '2 + 2' or '5 * 3'",
					},
				},
				Required: []string{"expression"},
			},
		},
	}

	// Define weather tool
	weatherTool := grok.Tool{
		Type: "function",
		Function: grok.Function{
			Name:        "get_weather",
			Description: "Gets weather information for a specified location",
			Parameters: &grok.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "City name, e.g. 'Beijing' or 'New York'",
					},
					"unit": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"celsius", "fahrenheit"},
						"description": "Temperature unit, celsius or fahrenheit",
					},
				},
				Required: []string{"location"},
			},
		},
	}

	// Create request with tools
	request := &grok.ChatCompletionRequest{
		Model: "test-model",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "Calculate 25 times 16, then check the weather in Beijing",
			},
		},
		Tools:       []grok.Tool{calculatorTool, weatherTool},
		ToolChoice:  "auto", // Let the model choose appropriate tools
		Temperature: 0.7,
	}

	// Skip actual API call in tests
	t.Skip("Skipping tools test - for demonstration only")

	// In a real test, we would make the API call
	ctx := context.Background()
	response, err := grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		t.Fatalf("Chat completion with tools failed: %v", err)
	}

	// Verify tool calls in the response
	if len(response.Choices) == 0 {
		t.Errorf("No choices returned in response")
		return
	}

	// Test that tool calls are properly formatted
	message := response.Choices[0].Message
	if len(message.ToolCalls) > 0 {
		// Test first tool call
		toolCall := message.ToolCalls[0]

		// Verify basic tool call properties
		if toolCall.Type != "function" {
			t.Errorf("Expected tool type 'function', got '%s'", toolCall.Type)
		}

		if toolCall.Function.Name != "calculator" && toolCall.Function.Name != "get_weather" {
			t.Errorf("Unexpected function name: %s", toolCall.Function.Name)
		}

		// Verify function arguments are valid JSON
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			t.Errorf("Invalid function arguments JSON: %v", err)
		}

		// In a real test, we could handle the tool response and make a follow-up request
	}
}

// TestToolDefinitionStructure tests that tool definitions are structured correctly
// and can be serialized to valid JSON.
func TestToolDefinitionStructure(t *testing.T) {
	// Define a calculator tool
	calculatorTool := grok.Tool{
		Type: "function",
		Function: grok.Function{
			Name:        "calculator",
			Description: "Performs mathematical calculations",
			Parameters: &grok.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "The expression to calculate, e.g. '2 + 2' or '5 * 3'",
					},
				},
				Required: []string{"expression"},
			},
		},
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(calculatorTool)
	if err != nil {
		t.Fatalf("Failed to marshal tool: %v", err)
	}

	// Deserialize to map to check fields
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["type"] != "function" {
		t.Errorf("Expected type 'function', got '%v'", result["type"])
	}

	// Check function field
	function, ok := result["function"].(map[string]interface{})
	if !ok {
		t.Fatalf("Function is not an object")
	}

	if function["name"] != "calculator" {
		t.Errorf("Expected function name 'calculator', got '%v'", function["name"])
	}

	// Check parameters field
	parameters, ok := function["parameters"].(map[string]interface{})
	if !ok {
		t.Fatalf("Parameters is not an object")
	}

	if parameters["type"] != "object" {
		t.Errorf("Expected parameters type 'object', got '%v'", parameters["type"])
	}
}

// TestToolCallSerialization tests the serialization of tool calls
// in assistant messages.
func TestToolCallSerialization(t *testing.T) {
	// Create an assistant message with a tool call
	message := grok.AssistantMessage{
		BaseMessage: grok.BaseMessage{
			Role:    "assistant",
			Content: "I'll help you calculate that.",
		},
		ToolCalls: []grok.ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: struct {
					Arguments string `json:"arguments"`
					Name      string `json:"name"`
				}{
					Name:      "calculator",
					Arguments: `{"expression": "25 * 16"}`,
				},
			},
		},
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Deserialize to map to check fields
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check tool_calls field
	toolCalls, ok := result["tool_calls"].([]interface{})
	if !ok {
		t.Fatalf("tool_calls is not an array")
	}
	if len(toolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(toolCalls))
	}

	// Check first tool call
	toolCall, ok := toolCalls[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Tool call is not an object")
	}

	if toolCall["id"] != "call_123" {
		t.Errorf("Expected tool call id 'call_123', got '%v'", toolCall["id"])
	}
	if toolCall["type"] != "function" {
		t.Errorf("Expected tool call type 'function', got '%v'", toolCall["type"])
	}

	// Check function field
	function, ok := toolCall["function"].(map[string]interface{})
	if !ok {
		t.Fatalf("Function is not an object")
	}

	if function["name"] != "calculator" {
		t.Errorf("Expected function name 'calculator', got '%v'", function["name"])
	}

	// Verify that arguments is a valid JSON string
	args := function["arguments"].(string)
	var argsMap map[string]interface{}
	if err := json.Unmarshal([]byte(args), &argsMap); err != nil {
		t.Errorf("Arguments is not valid JSON: %v", err)
	}

	if argsMap["expression"] != "25 * 16" {
		t.Errorf("Expected expression '25 * 16', got '%v'", argsMap["expression"])
	}
}

// TestToolResponseHandling tests the structure of tool response messages
func TestToolResponseHandling(t *testing.T) {
	// Create a tool response message
	toolMessage := grok.NewToolCallMessage("The result is 400", "call_123")

	// Serialize to JSON
	jsonData, err := json.Marshal(toolMessage)
	if err != nil {
		t.Fatalf("Failed to marshal tool message: %v", err)
	}

	// Deserialize to map to check fields
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check fields
	if result["role"] != "tool" {
		t.Errorf("Expected role 'tool', got '%v'", result["role"])
	}
	if result["content"] != "The result is 400" {
		t.Errorf("Expected content 'The result is 400', got '%v'", result["content"])
	}
	if result["tool_call_id"] != "call_123" {
		t.Errorf("Expected tool_call_id 'call_123', got '%v'", result["tool_call_id"])
	}
}
