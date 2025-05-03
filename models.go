package grok_go

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// ListModels fetches available models from the X.AI API
func ListModels(ctx context.Context, client *Client) (*ModelListResponse, error) {
	// Save current endpoint and restore it after the call
	originalEndpoint := client.Endpoint
	client.Endpoint = "models"
	defer func() {
		client.Endpoint = originalEndpoint
	}()

	// Create HTTP request
	httpRequest := Request{
		Auth: client.ApiKey,
		Url:  client.Url(),
		Body: nil, // No body needed for GET request
	}

	// Build the HTTP request
	req, err := httpRequest.BuildHttpRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Change method to GET
	req.Method = http.MethodGet

	// Execute the request
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
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
	var response ModelListResponse
	if err = json.Unmarshal(body, &response); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetModelInfo fetches information about a specific model
func GetModelInfo(ctx context.Context, client *Client, modelID string) (*ModelInfo, error) {
	if modelID == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	// Save current endpoint and restore it after the call
	originalEndpoint := client.Endpoint
	client.Endpoint = fmt.Sprintf("models/%s", modelID)
	defer func() {
		client.Endpoint = originalEndpoint
	}()

	// Create HTTP request
	httpRequest := Request{
		Auth: client.ApiKey,
		Url:  client.Url(),
		Body: nil, // No body needed for GET request
	}

	// Build the HTTP request
	req, err := httpRequest.BuildHttpRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Change method to GET
	req.Method = http.MethodGet

	// Execute the request
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
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
	var modelInfo ModelInfo
	if err = json.Unmarshal(body, &modelInfo); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal response", "error", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &modelInfo, nil
}
