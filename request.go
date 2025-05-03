package grok_go

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SimonMorphy/grok-go/constants"
	"github.com/SimonMorphy/grok-go/constants/errors"
)

// BuildHttpRequest creates an HTTP request with the specified parameters
// It returns the constructed HTTP request or an error if creation fails
// The context is used for request cancellation and logging
func (r Request) BuildHttpRequest(ctx context.Context) (*http.Request, error) {
	// Create a new HTTP request with the provided context
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, r.Url, bufio.NewReader(bytes.NewBuffer(r.Body)))
	if err != nil {
		// Wrap the error in our standard error type
		err = errors.NewWithError(constants.ErrnoBuildRequestError, err)
		slog.ErrorContext(ctx, "failed to build http request", "error", err)
		return nil, err
	}

	// Set standard headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.Auth))

	return request, nil
}
