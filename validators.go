package grok_go

import (
	"github.com/SimonMorphy/grok-go/constants"
	"github.com/SimonMorphy/grok-go/constants/errors"
)

// ValidateClient validates client configuration
func ValidateClient(client *Client) error {
	if client.ApiKey == "" {
		return errors.New(constants.ErrnoInvalidApiKey)
	}
	if client.BaseUrl == "" {
		return errors.New(constants.ErrnoInvalidBaseUrl)
	}
	if client.Endpoint == "" {
		return errors.New(constants.ErrnoInvalidEndpoint)
	}
	return nil
}

// ValidateChatCompletionRequest validates a chat completion request
func ValidateChatCompletionRequest(request *ChatCompletionRequest) error {
	if request.Model == "" {
		return errors.NewWithMsgF(constants.ErrnoInvalidValidationError, "model is required")
	}
	if len(request.Messages) == 0 {
		return errors.NewWithMsgF(constants.ErrnoInvalidValidationError, "messages are required")
	}
	return nil
}

// ValidateImageGenerationRequest validates an image generation request
func ValidateImageGenerationRequest(request *ImageGenerationRequest) error {
	if request.Model == "" {
		return errors.NewWithMsgF(constants.ErrnoInvalidValidationError, "model is required")
	}
	if request.Prompt == "" {
		return errors.NewWithMsgF(constants.ErrnoInvalidValidationError, "prompt is required")
	}
	return nil
}

// ValidateLanguageCode validates a language code
func ValidateLanguageCode(languageCode string) error {
	if languageCode == "" {
		return errors.NewWithMsgF(constants.ErrnoInvalidValidationError, "language code is required")
	}
	return nil
}
