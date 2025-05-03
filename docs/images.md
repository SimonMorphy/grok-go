# Image Generation Guide

This guide explains how to use the grok-go library to generate images with the X.AI API.

## Overview

X.AI's image generation API allows you to create images from text prompts. The grok-go library provides a simple interface to this API, making it easy to generate images in your Go applications.

## Setting Up

First, create an image generation client:

```go
import (
	"context"
	"fmt"
	"log"
	"os"
	
	grok "github.com/SimonMorphy/grok-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set")
	}
	
	// Create an image generation client
	client, err := grok.CreateImageGenerationClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	
	// Now you can use the client to generate images
}
```

Alternatively, you can use a standard client with the image generation endpoint:

```go
client, err := grok.NewClientWithOptions(
	apiKey,
	grok.WithEndpoint("images/generations"),
)
```

## Basic Image Generation

To generate an image, create an `ImageGenerationRequest` and call `CreateImage`:

```go
// Create a context
ctx := context.Background()

// Create an image generation request
request := &grok.ImageGenerationRequest{
	Model:  "grok-3-image", // Specify the image generation model
	Prompt: "A serene landscape with mountains, a lake, and a sunset sky",
	Size:   "1024x1024",    // Image size (available: 256x256, 512x512, 1024x1024)
	N:      1,              // Number of images to generate
}

// Send the request
response, err := grok.CreateImage(ctx, client, request)
if err != nil {
	log.Fatalf("Image generation failed: %v", err)
}

// Process the response
fmt.Printf("Generated %d images\n", len(response.Data))
for i, imageData := range response.Data {
	fmt.Printf("Image %d: %s\n", i+1, imageData.URL)
}
```

## Request Parameters

The `ImageGenerationRequest` struct supports the following parameters:

| Parameter | Type | Description |
|-----------|------|-------------|
| `Model` | string | The model to use for image generation (e.g., "grok-3-image") |
| `Prompt` | string | A text description of the desired image |
| `N` | int | The number of images to generate (default: 1) |
| `Size` | string | The size of the generated images (e.g., "256x256", "512x512", "1024x1024") |
| `ResponseFormat` | string | The format in which to return the generated images ("url" or "b64_json") |
| `User` | string | A unique identifier for the end-user |

## Handling Responses

The API can return images in two formats: URLs or Base64-encoded JSON. You can specify the format using the `ResponseFormat` parameter.

### URL Responses

By default, the API returns URLs to the generated images:

```go
// Response format defaults to "url"
request := &grok.ImageGenerationRequest{
	Model:  "grok-3-image",
	Prompt: "A serene landscape",
	Size:   "512x512",
}

response, err := grok.CreateImage(ctx, client, request)
if err != nil {
	log.Fatalf("Image generation failed: %v", err)
}

// Process the URL responses
for i, imageData := range response.Data {
	if imageData.URL != "" {
		fmt.Printf("Image %d URL: %s\n", i+1, imageData.URL)
		
		// Download the image (see example below)
	}
}
```

### Base64 Responses

To receive images as Base64-encoded strings:

```go
request := &grok.ImageGenerationRequest{
	Model:          "grok-3-image",
	Prompt:         "A serene landscape",
	Size:           "512x512",
	ResponseFormat: "b64_json",
}

response, err := grok.CreateImage(ctx, client, request)
if err != nil {
	log.Fatalf("Image generation failed: %v", err)
}

// Process the Base64 responses
for i, imageData := range response.Data {
	if imageData.B64JSON != "" {
		fmt.Printf("Received Base64 data for image %d (length: %d)\n", i+1, len(imageData.B64JSON))
		
		// Save the Base64-encoded image (see example below)
	}
}
```

## Downloading and Saving Images

### Downloading from URL

```go
import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Function to download an image from a URL
func downloadImage(url, outputPath string) error {
	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status: %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy image data to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Usage example
outputDir := "generated_images"
if err := os.MkdirAll(outputDir, 0755); err != nil {
	log.Fatalf("Failed to create output directory: %v", err)
}

outputPath := filepath.Join(outputDir, "image.jpg")
err = downloadImage(imageData.URL, outputPath)
if err != nil {
	log.Printf("Failed to download image: %v", err)
} else {
	fmt.Printf("Image saved to: %s\n", outputPath)
}
```

### Saving Base64 Images

```go
import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

// Function to save a Base64-encoded image
func saveBase64Image(base64Data, outputPath string) error {
	// Remove common Base64 prefixes if present
	base64Data = strings.TrimPrefix(base64Data, "data:image/jpeg;base64,")
	base64Data = strings.TrimPrefix(base64Data, "data:image/png;base64,")

	// Decode Base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("base64 decode error: %w", err)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write image data to file
	_, err = out.Write(imageData)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Usage example
outputDir := "generated_images"
if err := os.MkdirAll(outputDir, 0755); err != nil {
	log.Fatalf("Failed to create output directory: %v", err)
}

outputPath := filepath.Join(outputDir, "image.jpg")
err = saveBase64Image(imageData.B64JSON, outputPath)
if err != nil {
	log.Printf("Failed to save Base64 image: %v", err)
} else {
	fmt.Printf("Image saved to: %s\n", outputPath)
}
```

## Complete Example

Here's a complete example that demonstrates image generation with both URL and Base64 response formats:

```go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	grok "github.com/SimonMorphy/grok-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("GROK_API_KEY environment variable not set")
	}

	// Create an image generation client
	client, err := grok.CreateImageGenerationClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Create output directory
	outputDir := "generated_images"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate an image with URL response format
	generateAndSaveImage(ctx, client, "url", outputDir)

	// Generate an image with Base64 response format
	generateAndSaveImage(ctx, client, "b64_json", outputDir)
}

func generateAndSaveImage(ctx context.Context, client *grok.Client, responseFormat, outputDir string) {
	fmt.Printf("\nGenerating image with response format: %s\n", responseFormat)

	// Create an image generation request
	request := &grok.ImageGenerationRequest{
		Model:          "grok-3-image",
		Prompt:         "A beautiful mountain landscape with a lake and sunset",
		Size:           "512x512",
		N:              1,
		ResponseFormat: responseFormat,
	}

	// Send the request
	response, err := grok.CreateImage(ctx, client, request)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	// Process the response
	fmt.Printf("Generated %d images\n", len(response.Data))

	for i, imageData := range response.Data {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("image_%s_%d.jpg", responseFormat, i+1))

		if imageData.URL != "" {
			fmt.Printf("Image URL: %s\n", imageData.URL)
			err = downloadImage(imageData.URL, outputPath)
		} else if imageData.B64JSON != "" {
			fmt.Printf("Received Base64 data (length: %d)\n", len(imageData.B64JSON))
			err = saveBase64Image(imageData.B64JSON, outputPath)
		}

		if err != nil {
			fmt.Printf("Failed to save image: %v\n", err)
		} else {
			fmt.Printf("Image saved to: %s\n", outputPath)
		}
	}
}

func downloadImage(url, outputPath string) error {
	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status: %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy image data to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func saveBase64Image(base64Data, outputPath string) error {
	// Remove common Base64 prefixes if present
	base64Data = strings.TrimPrefix(base64Data, "data:image/jpeg;base64,")
	base64Data = strings.TrimPrefix(base64Data, "data:image/png;base64,")

	// Decode Base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("base64 decode error: %w", err)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write image data to file
	_, err = out.Write(imageData)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
```

## Best Practices

1. **Prompt Engineering**: Be specific and detailed in your prompts for better results.
2. **Error Handling**: Always check for errors when making API calls.
3. **Timeouts**: Set appropriate timeouts for image generation, as it can take longer than text requests.
4. **Rate Limits**: Be aware of API rate limits for image generation.
5. **Image Size**: Choose the appropriate image size based on your needs (larger sizes may take longer to generate). 