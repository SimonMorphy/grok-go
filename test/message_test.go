package test

import (
	"encoding/json"
	"os"
	"testing"

	grok "github.com/SimonMorphy/grok-go"
)

// TestMessageInterface tests if message types correctly implement the Message interface
func TestMessageInterface(t *testing.T) {
	// Verify that all message types implement the Message interface
	var _ grok.Message = grok.BaseMessage{}
	var _ grok.Message = grok.AssistantMessage{}
	var _ grok.Message = grok.ToolCallMessage{}

	// Test interface methods return expected values
	msg := grok.NewTextMessage("user", "test message")
	if msg.GetRole() != "user" {
		t.Errorf("GetRole() returned unexpected value: %s", msg.GetRole())
	}
	if msg.GetContent() != "test message" {
		t.Errorf("GetContent() returned unexpected value: %v", msg.GetContent())
	}
	if msg.GetName() != "" {
		t.Errorf("GetName() returned unexpected value: %s", msg.GetName())
	}
}

// TestMessageCreation tests the creation of various message types
func TestMessageCreation(t *testing.T) {
	t.Run("TextMessage", func(t *testing.T) {
		msg := grok.NewTextMessage("system", "test message")
		if msg.Role != "system" || msg.Content != "test message" {
			t.Errorf("Text message creation failed: %+v", msg)
		}
	})

	t.Run("AssistantMessage", func(t *testing.T) {
		msg := grok.NewAssistantMessage("I am an assistant")
		if msg.Role != "assistant" || msg.Content != "I am an assistant" {
			t.Errorf("Assistant message creation failed: %+v", msg)
		}
	})

	t.Run("ToolCallMessage", func(t *testing.T) {
		msg := grok.NewToolCallMessage("tool result", "tool_123")
		if msg.Role != "tool" || msg.Content != "tool result" || msg.ToolCallID != "tool_123" {
			t.Errorf("Tool call message creation failed: %+v", msg)
		}
	})

	t.Run("MultiContentMessage", func(t *testing.T) {
		// Create a temporary file for testing
		tmpFile, err := os.CreateTemp("", "test-*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		items := []any{"text content", *tmpFile}
		msg := grok.NewMultiContentMessage("user", items)
		if msg.Role != "user" || len(msg.Content.([]any)) != 2 {
			t.Errorf("Multi-content message creation failed: %+v", msg)
		}
	})

	t.Run("WithName", func(t *testing.T) {
		msg := grok.NewTextMessage("user", "message")
		namedMsg := grok.WithName(msg, "test_user")
		if namedMsg.Name != "test_user" || namedMsg.Content != "message" {
			t.Errorf("Named message creation failed: %+v", namedMsg)
		}
	})
}

// TestMessageSerialization tests JSON serialization of various message types
func TestMessageSerialization(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		message  grok.Message
		expected map[string]interface{}
	}{
		{
			name:    "System Message",
			message: grok.NewTextMessage("system", "You are a helpful assistant"),
			expected: map[string]interface{}{
				"role":    "system",
				"content": "You are a helpful assistant",
			},
		},
		{
			name:    "User Message with Name",
			message: grok.WithName(grok.NewTextMessage("user", "Hello"), "user_123"),
			expected: map[string]interface{}{
				"role":    "user",
				"content": "Hello",
				"name":    "user_123",
			},
		},
		{
			name: "Assistant Message with Tool Calls",
			message: grok.AssistantMessage{
				BaseMessage: grok.BaseMessage{
					Role:    "assistant",
					Content: nil,
				},
				ToolCalls: []grok.ToolCall{
					{
						ID:   "call_123",
						Type: "function",
						Function: struct {
							Arguments string `json:"arguments"`
							Name      string `json:"name"`
						}{
							Name:      "get_weather",
							Arguments: `{"location": "Beijing"}`,
						},
					},
				},
			},
			expected: map[string]interface{}{
				"role": "assistant",
				"tool_calls": []map[string]interface{}{
					{
						"id":   "call_123",
						"type": "function",
						"function": map[string]interface{}{
							"name":      "get_weather",
							"arguments": `{"location": "Beijing"}`,
						},
					},
				},
			},
		},
		{
			name:    "Tool Message",
			message: grok.NewToolCallMessage("Weather is sunny", "call_123"),
			expected: map[string]interface{}{
				"role":         "tool",
				"content":      "Weather is sunny",
				"tool_call_id": "call_123",
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize message to JSON
			jsonData, err := json.Marshal(tt.message)
			if err != nil {
				t.Fatalf("Serialization failed: %v", err)
			}

			// Parse JSON to map
			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Fatalf("JSON parsing failed: %v", err)
			}

			// Verify fields
			for key, expectedValue := range tt.expected {
				actualValue, exists := result[key]
				if !exists {
					t.Errorf("Missing expected field '%s'", key)
					continue
				}

				// Special handling for tool_calls field
				if key == "tool_calls" {
					verifyToolCalls(t, expectedValue, actualValue)
				} else if actualValue != expectedValue {
					t.Errorf("Field '%s': Expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestMessagesCollection tests message collection functionality
func TestMessagesCollection(t *testing.T) {
	messages := grok.NewMessages()

	// Add different types of messages
	messages.AddMessage(grok.NewTextMessage("system", "You are a helpful assistant"))
	messages.AddMessage(grok.WithName(grok.NewTextMessage("user", "Hello"), "user_123"))
	messages.AddMessage(grok.NewAssistantMessage("Hi there!"))

	// Verify message count
	if len(messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(messages))
	}

	// Verify message order and roles
	expectedRoles := []string{"system", "user", "assistant"}
	for i, role := range expectedRoles {
		if messages[i].GetRole() != role {
			t.Errorf("Message %d: Expected role '%s', got '%s'", i, role, messages[i].GetRole())
		}
	}

	// Test serialization
	jsonData, err := json.Marshal(messages)
	if err != nil {
		t.Fatalf("Failed to serialize message array: %v", err)
	}

	// Verify serialized structure
	var result []map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(result))
	}
}

// verifyToolCalls is a helper function to compare tool calls in serialized messages
func verifyToolCalls(t *testing.T, expected, actual interface{}) {
	expectedCalls := expected.([]map[string]interface{})
	actualCalls := actual.([]interface{})

	if len(expectedCalls) != len(actualCalls) {
		t.Errorf("Expected %d tool calls, got %d", len(expectedCalls), len(actualCalls))
		return
	}

	for i, expectedCall := range expectedCalls {
		actualCall := actualCalls[i].(map[string]interface{})
		for callKey, callValue := range expectedCall {
			if callKey == "function" {
				expectedFunc := callValue.(map[string]interface{})
				actualFunc := actualCall[callKey].(map[string]interface{})
				for funcKey, funcValue := range expectedFunc {
					if actualFunc[funcKey] != funcValue {
						t.Errorf("Tool call %d function: Expected %s=%v, got %v",
							i, funcKey, funcValue, actualFunc[funcKey])
					}
				}
			} else if actualCall[callKey] != callValue {
				t.Errorf("Tool call %d: Expected %s=%v, got %v",
					i, callKey, callValue, actualCall[callKey])
			}
		}
	}
}
