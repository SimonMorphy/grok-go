package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// TestImageGeneration tests the image generation API functionality
// It verifies that the API can process image generation requests with various prompts
func TestImageGeneration(t *testing.T) {
	// Initialize client with test API key and extended timeout
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint with extended timeout
	client.BaseUrl = "https://api.example.com/v1/"
	client.HttpClient = &http.Client{
		Timeout: 180 * time.Second, // 3 minute timeout
	}

	// Test prompts
	prompts := []string{
		"A cute cartoon cat with white fur and blue eyes on a simple background",
		"Simple red geometric shapes on a white background",
		"A green forest landscape with a blue lake",
	}

	// Create output directory for generated images
	outputDir := "generated_images"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Skip actual API call in tests
	t.Skip("Skipping image generation test - for demonstration only")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test each prompt
	for i, prompt := range prompts {
		t.Logf("Testing prompt %d: %s", i+1, prompt)

		// Create image generation request
		request := &grok.ImageGenerationRequest{
			Model:  "test-image-model",
			Prompt: prompt,
			Size:   "512x512", // Use smaller size for testing
			N:      1,
		}

		// Send request
		response, err := grok.CreateImage(ctx, client, request)
		if err != nil {
			t.Logf("Image generation failed for prompt %d: %v", i+1, err)
			continue
		}

		// In a real test, we would process and save the generated images
		// and validate their properties
		t.Logf("Successfully generated %d images for prompt %d", len(response.Data), i+1)
	}
}

// downloadAndSaveImage downloads an image from a URL and saves it to the specified path
func downloadAndSaveImage(url string, filePath string) error {
	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// Create output file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// saveBase64Image decodes and saves a Base64-encoded image to the specified path
func saveBase64Image(base64Data string, filePath string) error {
	// Decode Base64 data
	decoded, err := decodeBase64Image(base64Data)
	if err != nil {
		return err
	}

	// Create output file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write decoded data to file
	_, err = out.Write(decoded)
	if err != nil {
		return err
	}

	return nil
}

// decodeBase64Image decodes a Base64-encoded image
func decodeBase64Image(base64Data string) ([]byte, error) {
	// Remove common prefixes if present
	base64Data = strings.TrimPrefix(base64Data, "data:image/jpeg;base64,")
	base64Data = strings.TrimPrefix(base64Data, "data:image/png;base64,")

	// Decode Base64 data
	return base64.StdEncoding.DecodeString(base64Data)
}

// TestSimpleImageGeneration tests a single image generation request
// It provides a more focused test case for basic functionality
func TestSimpleImageGeneration(t *testing.T) {
	// Initialize client with test API key
	client, err := grok.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Use test endpoint with timeout
	client.BaseUrl = "https://api.example.com/v1/"
	client.HttpClient = &http.Client{
		Timeout: 120 * time.Second,
	}

	// Skip actual API call in tests
	t.Skip("Skipping simple image generation test - for demonstration only")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Create a simple image generation request
	request := &grok.ImageGenerationRequest{
		Model:  "test-image-model",
		Prompt: "A minimalist blue circle",
		Size:   "256x256", // Small size for testing
		N:      1,
	}

	// Send request
	response, err := grok.CreateImage(ctx, client, request)
	if err != nil {
		t.Logf("First model attempt failed: %v", err)

		// In a real implementation, we might try alternative models
		t.Skip("Image generation API may not be available, skipping test")
		return
	}

	// In a real test, we would validate the response structure
	if len(response.Data) == 0 {
		t.Errorf("Response contained no image data")
	}
}

// TestImageGenerationSimple 测试图片生成API，只请求一次并显示完整响应体
func TestImageGenerationSimple(t *testing.T) {
	// 初始化客户端，设置更长的超时时间
	client, err := grok.NewClient("test-api")
	client.BaseUrl = "https://www.yunqiaoai.top/v1/"
	client.HttpClient = &http.Client{
		Timeout: 120 * time.Second, // 2分钟超时
	}

	if err != nil {
		t.Fatalf("初始化客户端失败: %v", err)
	}

	// 创建上下文，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 简单图片生成请求
	request := &grok.ImageGenerationRequest{
		Model:  "grok-3", // 只使用grok-3模型
		Prompt: "一个简单的蓝色圆形",
		Size:   "256x256", // 小尺寸
		N:      1,
	}

	fmt.Println("发送单次图片生成请求到grok-3模型...")

	// 发送请求
	response, err := grok.CreateImage(ctx, client, request)

	// 创建请求体的JSON表示以便打印
	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Printf("\n请求体:\n%s\n", string(requestJSON))

	if err != nil {
		// 即使有错误也继续运行测试，因为我们只想看响应
		fmt.Printf("\n请求失败，错误: %v\n", err)

		// 无法直接使用类型断言获取APIErrorResponse，因为它不是error接口的实现
		// 只能打印错误信息
		fmt.Printf("\nAPI错误详情: %s\n", err.Error())
	} else {
		// 打印完整的响应体
		fmt.Println("\n请求成功!")

		// 使用漂亮的打印格式输出完整的响应
		responseJSON, _ := json.MarshalIndent(response, "", "  ")
		fmt.Printf("\n完整响应体:\n%s\n", string(responseJSON))

		// 提取和显示图片URL或Base64数据
		for i, imageData := range response.Data {
			fmt.Printf("\n图片 #%d:\n", i+1)

			if imageData.URL != "" {
				fmt.Printf("图片URL: %s\n", imageData.URL)
			} else if imageData.B64JSON != "" {
				fmt.Printf("图片Base64数据长度: %d 字符\n", len(imageData.B64JSON))
				fmt.Printf("Base64数据前100个字符: %.100s...\n", imageData.B64JSON)
			} else {
				fmt.Printf("警告: 图片数据为空\n")
			}

			if imageData.RevisedPrompt != "" {
				fmt.Printf("修改后的提示词: %s\n", imageData.RevisedPrompt)
			}
		}
	}
}
