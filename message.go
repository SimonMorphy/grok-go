package grok_go

// Message defines the basic interface for all message types
// All message implementations must provide methods to access content, name, and role
type Message interface {
	// GetContent returns the content of the message
	GetContent() any

	// GetName returns the name of the message author
	GetName() string

	// GetRole returns the role of the message (e.g., "system", "user", "assistant", "tool")
	GetRole() string
}

// BaseMessage implements the basic structure for Message interface
// It provides the common fields and methods required by the Message interface
type BaseMessage struct {
	Content any    `json:"content"`        // Content of the message (can be string or structured content)
	Name    string `json:"name,omitempty"` // Optional name of the message author
	Role    string `json:"role"`           // Role of the message sender
}

// GetContent returns the content of the message
func (m BaseMessage) GetContent() any {
	return m.Content
}

// GetName returns the name of the message
func (m BaseMessage) GetName() string {
	return m.Name
}

// GetRole returns the role of the message
func (m BaseMessage) GetRole() string {
	return m.Role
}

// ToolCallMessage represents a message with a tool call
// Used for responses from tool executions
type ToolCallMessage struct {
	BaseMessage
	ToolCallID string `json:"tool_call_id,omitempty"` // ID of the tool call this message responds to
}

// AssistantMessage represents a message from the assistant
// May contain tool calls that the assistant wants to make
type AssistantMessage struct {
	BaseMessage
	ToolCalls []ToolCall `json:"tool_calls,omitempty"` // Tool calls the assistant wants to make
}

// ToolCall represents a tool call in the message
// Contains information about a function to be called
type ToolCall struct {
	Function struct {
		Arguments string `json:"arguments"` // Arguments to the function as a JSON string
		Name      string `json:"name"`      // Name of the function to call
	} `json:"function"`
	ID    string `json:"id"`              // Unique identifier for this tool call
	Index int    `json:"index,omitempty"` // Index of this tool call in the sequence
	Type  string `json:"type,omitempty"`  // Type of the tool call (usually "function")
}

// WithName adds a name to the message
// Returns a new message with the name field set
func WithName(msg BaseMessage, name string) BaseMessage {
	msg.Name = name
	return msg
}

// NewSystemMessage creates a new system message
// System messages provide instructions or context to the model
func NewSystemMessage(content any) BaseMessage {
	return BaseMessage{
		Content: content,
		Role:    "system",
	}
}

// NewUserMessage creates a new user message
// User messages represent input from the end user
func NewUserMessage(content any) BaseMessage {
	return BaseMessage{
		Content: content,
		Role:    "user",
	}
}

// NewAssistantMessage creates a new assistant message
// Assistant messages represent responses from the AI assistant
func NewAssistantMessage(content any) AssistantMessage {
	return AssistantMessage{
		BaseMessage: BaseMessage{
			Content: content,
			Role:    "assistant",
		},
	}
}

// NewToolCallMessage creates a new tool call message
// Tool call messages represent responses from tools that were called
func NewToolCallMessage(content any, toolCallID string) ToolCallMessage {
	return ToolCallMessage{
		BaseMessage: BaseMessage{
			Content: content,
			Role:    "tool",
		},
		ToolCallID: toolCallID,
	}
}

// Messages represents an array of messages
// Used to build a conversation history
type Messages []Message

// AddMessage adds a message to the array
// Appends the message to the end of the conversation
func (m *Messages) AddMessage(msg Message) {
	*m = append(*m, msg)
}

// NewMessages creates a new empty message array
// Initializes a conversation with zero messages
func NewMessages() Messages {
	return make(Messages, 0)
}

// NewTextMessage creates a new text message with the specified role
// Generic function to create messages with simple text content
func NewTextMessage(role string, text string) BaseMessage {
	return BaseMessage{
		Content: text,
		Role:    role,
	}
}

// NewMultiContentMessage creates a new message with multiple content items
// Used when a message contains a mix of content types (e.g., text and images)
func NewMultiContentMessage(role string, items []any) BaseMessage {
	return BaseMessage{
		Content: items,
		Role:    role,
	}
}

// NewMessageArray creates a new empty message array
// Alias for NewMessages for backward compatibility
func NewMessageArray() Messages {
	return make(Messages, 0)
}
