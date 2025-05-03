package grok_go

import (
	"net/http"
	"time"
)

// Url returns the full API URL by combining base URL and endpoint
// This is used for constructing API request URLs
func (client *Client) Url() string {
	return client.BaseUrl + client.Endpoint
}

// apiKey sets the API key for the client
// Returns the client for method chaining
func (client *Client) apiKey(apikey string) *Client {
	client.ApiKey = apikey
	return client
}

// baseUrl sets the base URL for the client
// Returns the client for method chaining
func (client *Client) baseUrl(baseUrl string) *Client {
	client.BaseUrl = baseUrl
	return client
}

// timeout sets the timeout duration for the client
// Returns the client for method chaining
func (client *Client) timeout(timeout time.Duration) *Client {
	client.Timeout = timeout
	return client
}

// endpoint sets the API endpoint for the client
// Returns the client for method chaining
func (client *Client) endpoint(endpoint string) *Client {
	client.Endpoint = endpoint
	return client
}

// httpClient sets the HTTP client for API requests
// Returns the client for method chaining
func (client *Client) httpClient(httpClient *http.Client) *Client {
	client.HttpClient = httpClient
	return client
}

// build validates and finalizes the client configuration
// Returns the configured client or an error if validation fails
func (client *Client) build() (*Client, error) {
	// Validate client configuration
	if err := ValidateClient(client); err != nil {
		return nil, err
	}

	// Set default HTTP client if not provided
	if client.HttpClient == nil {
		client.HttpClient = &http.Client{Timeout: client.Timeout}
	}

	return client, nil
}
