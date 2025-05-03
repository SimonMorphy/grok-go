package grok_go

import (
	"net/http"
	"time"
)

// Client Represents Http Client
type Client struct {
	ApiKey     string        `json:"api_key"`
	BaseUrl    string        `json:"base_url"`
	Timeout    time.Duration `json:"timeout"`
	Endpoint   string        `json:"endpoint"`
	HttpClient *http.Client  `json:"-"`
}

// Request A valid Http request to invoke remote api
type Request struct {
	Auth string
	Url  string
	Body []byte
}

// ChatCompletionRequest Request structure compliant with X.AI API specifications
type ChatCompletionRequest struct {
	// Required parameters
	Model    string                  `json:"model"`    // Model identifier, e.g. "grok-3-beta"
	Messages []ChatCompletionMessage `json:"messages"` // Array of conversation messages

	// Optional streaming parameters
	Stream        bool `json:"stream,omitempty"` // Whether to stream the response
	StreamOptions struct {
		IncludeUsage bool `json:"include_usage,omitempty"` // Whether to include token usage stats in the stream
	} `json:"stream_options,omitempty"`
	Deferred bool `json:"deferred,omitempty"` // Whether to return a deferred response

	// Optional control parameters
	MaxTokens           int                `json:"max_tokens,omitempty"`           // Maximum number of tokens to generate
	Temperature         float64            `json:"temperature,omitempty"`          // Sampling temperature (0-2)
	TopP                float64            `json:"top_p,omitempty"`                // Nucleus sampling parameter
	TopK                int                `json:"top_k,omitempty"`                // Top-k sampling parameter
	Stop                []string           `json:"stop,omitempty"`                 // Sequences where the API will stop generating
	PresencePenalty     float64            `json:"presence_penalty,omitempty"`     // Penalty for new tokens based on presence in text
	FrequencyPenalty    float64            `json:"frequency_penalty,omitempty"`    // Penalty for new tokens based on frequency in text
	LogitBias           map[string]float64 `json:"logit_bias,omitempty"`           // Modify likelihood of specified tokens
	Logprobs            bool               `json:"logprobs,omitempty"`             // Whether to return log probabilities
	TopLogprobs         int                `json:"top_logprobs,omitempty"`         // How many log probabilities to return
	Seed                int64              `json:"seed,omitempty"`                 // Random seed for deterministic results
	ResponseFormat      *ResponseFormat    `json:"response_format,omitempty"`      // Format of the response content
	Tools               []Tool             `json:"tools,omitempty"`                // Tools the model may call
	ToolChoice          interface{}        `json:"tool_choice,omitempty"`          // Controls which tool is called by the model
	User                string             `json:"user,omitempty"`                 // A unique identifier for the end-user
	JSONMode            bool               `json:"json_mode,omitempty"`            // Always output valid JSON
	PrefixTokens        []int              `json:"prefix_tokens,omitempty"`        // Tokens to prepend to generated text
	SkipParameterCheck  bool               `json:"-"`                              // Do not validate parameters
	BackendType         string             `json:"backend_type,omitempty"`         // Backend processing preference
	AdaptivePrompt      bool               `json:"adaptive_prompt,omitempty"`      // Whether to adapt the system prompt based on conversation
	ManagedCredentials  map[string]string  `json:"managed_credentials,omitempty"`  // Credentials for API integrations
	HTTPBasedPlugins    interface{}        `json:"http_based_plugins,omitempty"`   // Custom plugins based on HTTP endpoints
	ConversationContext string             `json:"conversation_context,omitempty"` // Additional context for conversation understanding
}

// ChatCompletionMessage represents a message in the chat completion request
type ChatCompletionMessage struct {
	Role       string        `json:"role"`                 // Role of the message sender (system, user, assistant, tool)
	Content    string        `json:"content"`              // The content of the message
	Name       string        `json:"name,omitempty"`       // The name of the author of this message
	ToolCalls  []APIToolCall `json:"tool_calls,omitempty"` // Tool calls made in this message
	ToolCallID string        `json:"tool_call_id,omitempty"`
	Prefix     bool          `json:"prefix,omitempty"` // Whether this is a prefix message for completion
}

// ResponseFormat specifies the format of the model's response
type ResponseFormat struct {
	Type string `json:"type"` // The format type, e.g., "text" or "json_object"
}

// Tool represents a function that the model may generate JSON inputs for
type Tool struct {
	Type     string   `json:"type"`     // The type of tool, usually "function"
	Function Function `json:"function"` // The function to be called
}

// Function represents a function specification
type Function struct {
	Name        string              `json:"name"`                  // The name of the function
	Description string              `json:"description,omitempty"` // A description of what the function does
	Parameters  *FunctionParameters `json:"parameters,omitempty"`  // The parameters the function accepts
}

// FunctionParameters defines the parameters that a function accepts
type FunctionParameters struct {
	Type       string                 `json:"type"` // The type of the parameters, usually "object"
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"` // Which parameters are required
}

// APIToolCall represents a call to a function made by the model
type APIToolCall struct {
	ID       string       `json:"id"`       // Unique ID of the tool call
	Type     string       `json:"type"`     // Type of the tool, usually "function"
	Function FunctionCall `json:"function"` // Function call details
}

// FunctionCall represents the details of a function call
type FunctionCall struct {
	Name      string `json:"name"`      // Name of the function being called
	Arguments string `json:"arguments"` // Arguments to pass to the function, as a JSON string
}

// Response structure for X.AI API chat completion response
type Response struct {
	ID                string           `json:"id"`                           // Unique identifier for the response
	Object            string           `json:"object"`                       // Object type, usually "chat.completion"
	Created           int64            `json:"created"`                      // Unix timestamp when the response was created
	Model             string           `json:"model"`                        // Model used for the completion
	Choices           []ResponseChoice `json:"choices"`                      // Array of completion choices
	SystemFingerprint string           `json:"system_fingerprint,omitempty"` // Fingerprint of the system configuration
	Usage             Usage            `json:"usage,omitempty"`              // Token usage information
}

// ResponseChoice represents a single completion choice in the response
type ResponseChoice struct {
	Index        int             `json:"index"`              // Index of this choice
	Message      ResponseMessage `json:"message"`            // The message content
	LogProbs     *LogProbs       `json:"logprobs,omitempty"` // Log probabilities if requested
	FinishReason string          `json:"finish_reason"`      // Reason why generation stopped
}

// ResponseMessage represents the message in a response choice
type ResponseMessage struct {
	Role      string        `json:"role"`                 // Role of the message sender, usually "assistant"
	Content   string        `json:"content"`              // The content of the message
	ToolCalls []APIToolCall `json:"tool_calls,omitempty"` // Tool calls made in this message
}

// LogProbs contains log probability information for token generation
type LogProbs struct {
	Content []TokenLogProb `json:"content"` // Log probability information for content tokens
}

// TokenLogProb holds log probability information for a single token
type TokenLogProb struct {
	Token       string             `json:"token"`        // The token string
	LogProb     float64            `json:"logprob"`      // Log probability of this token
	TopLogProbs map[string]float64 `json:"top_logprobs"` // Log probabilities of top alternatives
	Bytes       []int              `json:"bytes"`        // Token byte representation
}

// Usage provides token usage statistics for the request and response
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // Number of tokens in the prompt
	CompletionTokens int `json:"completion_tokens"` // Number of tokens in the completion
	TotalTokens      int `json:"total_tokens"`      // Total tokens used (prompt + completion)
}

// StreamResponse structure for streaming responses
type StreamResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
	Choices           []StreamChoice `json:"choices"`
	Usage             *Usage         `json:"usage,omitempty"`
}

// StreamChoice represents a chunk of the stream response
type StreamChoice struct {
	Index        int           `json:"index"`
	Delta        StreamMessage `json:"delta"`
	LogProbs     *LogProbs     `json:"logprobs,omitempty"`
	FinishReason string        `json:"finish_reason,omitempty"`
}

// StreamMessage represents a message chunk in the stream
type StreamMessage struct {
	Role      string        `json:"role,omitempty"`
	Content   string        `json:"content,omitempty"`
	ToolCalls []APIToolCall `json:"tool_calls,omitempty"`
}

// APIErrorResponse defines the structure for API error responses
type APIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

// ModelListResponse represents a list of available models
type ModelListResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// ModelInfo represents information about a single model
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ImageGenerationRequest defines the request structure for image generation
type ImageGenerationRequest struct {
	Model          string `json:"model"`                     // ID of the model to use
	Prompt         string `json:"prompt"`                    // Text description of the desired image
	N              int    `json:"n,omitempty"`               // Number of images to generate
	Size           string `json:"size,omitempty"`            // Size of the generated images
	Quality        string `json:"quality,omitempty"`         // Quality of the generated images
	Style          string `json:"style,omitempty"`           // Style of the generated images
	ResponseFormat string `json:"response_format,omitempty"` // Format of the response
	User           string `json:"user,omitempty"`            // User identifier
}

// ImageGenerationResponse defines the response structure from image generation
type ImageGenerationResponse struct {
	Created int64                 `json:"created"`
	Data    []ImageGenerationData `json:"data"`
}

// ImageGenerationData represents a single generated image
type ImageGenerationData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// LanguageListResponse represents a list of available languages for translation
type LanguageListResponse struct {
	Object string         `json:"object"`
	Data   []LanguageInfo `json:"data"`
}

// LanguageInfo represents information about a supported language
type LanguageInfo struct {
	Code       string `json:"code"`        // ISO language code
	Name       string `json:"name"`        // Language name
	NativeName string `json:"native_name"` // Language name in the native script
	Supported  bool   `json:"supported"`   // Whether translation is supported
}
