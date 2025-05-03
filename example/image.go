package example

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
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// ImageGenerationExample demonstrates how to use the image generation API.
// It shows how to create image generation requests and handle the response.
func ImageGenerationExample() {
	// Get API key from environment variable
	apiKey := os.Getenv(APIKeyEnvVar)
	if apiKey == "" {
		log.Printf("Warning: Environment variable %s not set. Using placeholder.", APIKeyEnvVar)
		apiKey = "your-api-key-here" // This will prevent actual API calls
	}

	// Initialize the client with the API key
	client, err := grok.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Set the API base URL and a longer timeout for image generation
	client.BaseUrl = "https://api.x.ai/v1/"
	client.HttpClient = &http.Client{
		Timeout: 3 * time.Minute, // Longer timeout for image generation
	}

	// Create a simple image generation request
	request := &grok.ImageGenerationRequest{
		Model:  "grok-3-image", // Specify the image generation model
		Prompt: "A serene landscape with mountains, a lake, and a sunset sky",
		Size:   "512x512", // Image size
		N:      1,         // Number of images to generate
	}

	fmt.Println("Sending image generation request...")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Send the request to generate images
	response, err := grok.CreateImage(ctx, client, request)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	fmt.Printf("Successfully generated %d images\n", len(response.Data))

	// Create output directory if it doesn't exist
	outputDir := "generated_images"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Process each generated image
	for i, imageData := range response.Data {
		fmt.Printf("\nImage #%d:\n", i+1)

		// Handle image response based on format (URL or Base64)
		if imageData.URL != "" {
			// API returned an image URL
			fmt.Printf("Image URL: %s\n", imageData.URL)

			// Download and save the image
			outputPath := filepath.Join(outputDir, fmt.Sprintf("image_%d.jpg", i+1))
			err = downloadImage(imageData.URL, outputPath)
			if err != nil {
				fmt.Printf("Failed to download image: %v\n", err)
			} else {
				fmt.Printf("Image saved to: %s\n", outputPath)
			}
		} else if imageData.B64JSON != "" {
			// API returned Base64-encoded image data
			fmt.Printf("Received Base64-encoded image (length: %d characters)\n", len(imageData.B64JSON))

			// Save the Base64-encoded image
			outputPath := filepath.Join(outputDir, fmt.Sprintf("image_%d.jpg", i+1))
			err = saveBase64Image(imageData.B64JSON, outputPath)
			if err != nil {
				fmt.Printf("Failed to save Base64 image: %v\n", err)
			} else {
				fmt.Printf("Image saved to: %s\n", outputPath)
			}
		} else {
			fmt.Println("No image data received")
		}

		// Display any revised prompt information
		if imageData.RevisedPrompt != "" {
			fmt.Printf("Revised prompt: %s\n", imageData.RevisedPrompt)
		}
	}
}

// downloadImage downloads an image from a URL and saves it to the specified path.
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

// saveBase64Image decodes and saves a Base64-encoded image to the specified path.
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
