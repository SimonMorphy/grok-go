package grok_go

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SimonMorphy/grok-go/constants"
	"github.com/SimonMorphy/grok-go/constants/errors"
)

// CreateChatCompletion sends a chat completion request to the X.AI API
// It returns a Response object and an error if one occurred
// Context is used for request cancellation and logging
func CreateChatCompletion(ctx context.Context, client *Client, request *ChatCompletionRequest) (*Response, error) {
	// Validate the request
	if err := ValidateChatCompletionRequest(request); err != nil {
		return nil, err
	}

	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal request", "error", err)
		return nil, errors.NewWithError(constants.ErrnoMarshalError, err)
	}

	// Create HTTP request
	httpRequest := Request{
		Auth: client.ApiKey,
		Url:  client.Url(),
		Body: requestBody,
	}

	// Build and send the HTTP request
	req, err := httpRequest.BuildHttpRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Execute the request
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send request", "error", err)
		return nil, errors.NewWithError(constants.ErrnoBuildRequestError, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read response body", "error", err)
		return nil, errors.NewWithError(constants.ErrnoInvalidRequestError, err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var apiError APIErrorResponse
		if err = json.Unmarshal(body, &apiError); err == nil {
			slog.ErrorContext(ctx, "API error",
				"status", resp.StatusCode,
				"error_type", apiError.Error.Type,
				"error_message", apiError.Error.Message)
			return nil, fmt.Errorf("API error: %s: %s", apiError.Error.Type, apiError.Error.Message)
		}
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal response", "error", err)
		return nil, errors.NewWithError(constants.ErrnoUnmarshalError, err)
	}

	return &response, nil
}

// CreateChatCompletionStream creates a streaming connection for chat completions
// It returns a StreamReader that can be used to read streaming responses
// Context is used for request cancellation and logging
func CreateChatCompletionStream(ctx context.Context, client *Client, request *ChatCompletionRequest) (*StreamReader, error) {
	// Ensure streaming is enabled
	request.Stream = true

	// Validate the request
	if err := ValidateChatCompletionRequest(request); err != nil {
		return nil, err
	}

	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal request", "error", err)
		return nil, errors.NewWithError(constants.ErrnoMarshalError, err)
	}

	// Create HTTP request
	httpRequest := Request{
		Auth: client.ApiKey,
		Url:  client.Url(),
		Body: requestBody,
	}

	// Build the HTTP request
	req, err := httpRequest.BuildHttpRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Set header for streaming response
	req.Header.Set("Accept", "text/event-stream")

	// Execute the request
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send request", "error", err)
		return nil, errors.NewWithError(constants.ErrnoBuildRequestError, err)
	}

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var apiError APIErrorResponse
		if err = json.Unmarshal(body, &apiError); err == nil {
			slog.ErrorContext(ctx, "API error",
				"status", resp.StatusCode,
				"error_type", apiError.Error.Type,
				"error_message", apiError.Error.Message)
			return nil, fmt.Errorf("API error: %s: %s", apiError.Error.Type, apiError.Error.Message)
		}

		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	// Create a stream reader to handle the SSE response
	return &StreamReader{
		reader:     bufio.NewReader(resp.Body),
		response:   resp,
		connection: resp.Body,
	}, nil
}

// StreamReader handles reading from a streaming response
type StreamReader struct {
	reader     *bufio.Reader
	response   *http.Response
	connection io.Closer
}

// Recv reads the next event from the stream
// It returns a StreamResponse object and an error if one occurred
// Returns io.EOF when the stream ends
func (s *StreamReader) Recv() (*StreamResponse, error) {
	// Read lines until we get a data line
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		line = bytes.TrimSpace(line)

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Check for end of stream
		if bytes.Equal(line, []byte("data: [DONE]")) {
			return nil, io.EOF
		}

		// Parse data lines
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))
			var response StreamResponse
			if err = json.Unmarshal(data, &response); err != nil {
				slog.Error("failed to unmarshal stream response", "error", err)
				continue
			}
			return &response, nil
		}
	}
}

// Close closes the stream connection
func (s *StreamReader) Close() error {
	return s.connection.Close()
}
