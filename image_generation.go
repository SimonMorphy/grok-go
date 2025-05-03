package grok_go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SimonMorphy/grok-go/constants"
	"github.com/SimonMorphy/grok-go/constants/errors"
)

// CreateImage sends an image generation request to the X.AI API
// It returns an ImageGenerationResponse object and an error if one occurred
// Context is used for request cancellation and logging
func CreateImage(ctx context.Context, client *Client, request *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// Save current endpoint and restore it after the call
	originalEndpoint := client.Endpoint
	client.Endpoint = "images/generations"
	defer func() {
		client.Endpoint = originalEndpoint
	}()

	// Validate the request
	if err := ValidateImageGenerationRequest(request); err != nil {
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
	var response ImageGenerationResponse
	if err = json.Unmarshal(body, &response); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal response", "error", err)
		return nil, errors.NewWithError(constants.ErrnoUnmarshalError, err)
	}

	return &response, nil
}
