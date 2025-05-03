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

// ListLanguages fetches a list of supported languages from the X.AI API
// It returns a LanguageListResponse object and an error if one occurred
// Context is used for request cancellation and logging
func ListLanguages(ctx context.Context, client *Client) (*LanguageListResponse, error) {
	// Save current endpoint and restore it after the call
	originalEndpoint := client.Endpoint
	client.Endpoint = "languages"
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
	var response LanguageListResponse
	if err = json.Unmarshal(body, &response); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal response", "error", err)
		return nil, errors.NewWithError(constants.ErrnoUnmarshalError, err)
	}

	return &response, nil
}

// GetLanguageInfo fetches information about a specific language
// It returns a LanguageInfo object for the requested language code
// Returns an error if the language is not found or if the API request fails
func GetLanguageInfo(ctx context.Context, client *Client, languageCode string) (*LanguageInfo, error) {
	// Validate language code
	if err := ValidateLanguageCode(languageCode); err != nil {
		return nil, err
	}

	// Get all languages
	languages, err := ListLanguages(ctx, client)
	if err != nil {
		return nil, err
	}

	// Find the requested language
	for _, lang := range languages.Data {
		if lang.Code == languageCode {
			return &lang, nil
		}
	}

	return nil, errors.NewWithMsgF(constants.ErrnoInvalidValidationError, "language with code '%s' not found", languageCode)
}
