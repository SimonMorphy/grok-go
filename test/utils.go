package test

import (
	"os"
	"testing"
)

// APIKeyEnvVar is the environment variable name that stores the API key
const (
	APIKeyEnvVar = "GROK_API_KEY"
	TestAPIURL   = "https://api.x.ai/v1/"
)

// GetTestAPIKey retrieves an API key from environment variables
// If the environment variable doesn't exist, returns a test placeholder key
func GetTestAPIKey(t *testing.T) string {
	// Try to get API key from environment variable
	apiKey, exists := os.LookupEnv(APIKeyEnvVar)
	if exists && apiKey != "" {
		t.Logf("Using API key from environment variable %s", APIKeyEnvVar)
		return apiKey
	}

	// If environment variable doesn't exist or is empty, return test placeholder key
	t.Logf("Environment variable %s not set, using test API key", APIKeyEnvVar)
	return "test-api-key-for-validation-only"
}

// SkipIfNoRealAPIKey skips the test if no real API key is provided
// Returns the API key if it exists
func SkipIfNoRealAPIKey(t *testing.T) string {
	apiKey, exists := os.LookupEnv(APIKeyEnvVar)
	if !exists || apiKey == "" {
		t.Skip("Skipping test: No API key provided in environment variable")
	}
	return apiKey
}
