package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	grok "github.com/SimonMorphy/grok-go"
)

// ToolsExample demonstrates how to use function tools with the chat completion API.
// It shows how to define tools, create tool-enabled requests, and handle tool call responses.
func ToolsExample() {
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

	// Set the API base URL
	client.BaseUrl = "https://api.x.ai/v1/"

	// Define a calculator tool
	calculatorTool := grok.Tool{
		Type: "function",
		Function: grok.Function{
			Name:        "calculate",
			Description: "Perform mathematical calculations",
			Parameters: &grok.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"operation": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"add", "subtract", "multiply", "divide", "power"},
						"description": "The mathematical operation to perform",
					},
					"x": map[string]interface{}{
						"type":        "number",
						"description": "The first operand",
					},
					"y": map[string]interface{}{
						"type":        "number",
						"description": "The second operand",
					},
				},
				Required: []string{"operation", "x", "y"},
			},
		},
	}

	// Define a weather information tool
	weatherTool := grok.Tool{
		Type: "function",
		Function: grok.Function{
			Name:        "get_weather",
			Description: "Get current weather information for a location",
			Parameters: &grok.FunctionParameters{
				Type: "object",
				Properties: map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The city and country, e.g., 'Tokyo, Japan'",
					},
					"unit": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"celsius", "fahrenheit"},
						"description": "The temperature unit to use",
					},
				},
				Required: []string{"location"},
			},
		},
	}

	// Create a request using tools
	request := &grok.ChatCompletionRequest{
		Model: "grok-3", // Specify the model to use
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "I need to calculate 15.7 to the power of 2.3, and then get the weather for Berlin, Germany",
			},
		},
		Tools:       []grok.Tool{calculatorTool, weatherTool},
		ToolChoice:  "auto", // Let the model choose which tool to use
		Temperature: 0.7,
	}

	// Create a context for the request
	ctx := context.Background()

	fmt.Println("Sending chat completion request with tools...")

	// Send the request to the API
	response, err := grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		log.Fatalf("Chat completion request failed: %v", err)
	}

	// Process the response and handle tool calls
	if len(response.Choices) > 0 {
		message := response.Choices[0].Message

		fmt.Println("\nInitial Response:")
		fmt.Println(message.Content)

		// Handle any tool calls in the response
		if len(message.ToolCalls) > 0 {
			fmt.Println("\nTool Calls:")

			// Store messages for the follow-up request
			messages := append(request.Messages, grok.ChatCompletionMessage{
				Role:      "assistant",
				Content:   message.Content,
				ToolCalls: message.ToolCalls,
			})

			// Process each tool call
			for _, toolCall := range message.ToolCalls {
				fmt.Printf("- Tool: %s\n", toolCall.Function.Name)
				fmt.Printf("  Arguments: %s\n", toolCall.Function.Arguments)

				// Execute the appropriate tool based on the function name
				var toolResult string
				if toolCall.Function.Name == "calculate" {
					toolResult = executeCalculator(toolCall.Function.Arguments)
				} else if toolCall.Function.Name == "get_weather" {
					toolResult = executeWeatherLookup(toolCall.Function.Arguments)
				} else {
					toolResult = fmt.Sprintf("Unknown tool: %s", toolCall.Function.Name)
				}

				// Add the tool result to our messages
				messages = append(messages, grok.ChatCompletionMessage{
					Role:       "tool",
					Content:    toolResult,
					ToolCallID: toolCall.ID,
				})

				fmt.Printf("  Result: %s\n", toolResult)
			}

			// Create a follow-up request with the tool results
			followupRequest := &grok.ChatCompletionRequest{
				Model:       request.Model,
				Messages:    messages,
				Temperature: request.Temperature,
			}

			fmt.Println("\nSending follow-up request with tool results...")

			// Send the follow-up request
			followupResponse, err := grok.CreateChatCompletion(ctx, client, followupRequest)
			if err != nil {
				log.Fatalf("Follow-up request failed: %v", err)
			}

			// Display the final response that incorporates the tool results
			if len(followupResponse.Choices) > 0 {
				fmt.Println("\nFinal Response:")
				fmt.Println(followupResponse.Choices[0].Message.Content)
			}
		}
	} else {
		fmt.Println("No response content received")
	}
}

// executeCalculator performs the mathematical operation specified in the arguments.
// It takes a JSON string and returns a formatted result string.
func executeCalculator(jsonArgs string) string {
	// Parse the JSON arguments
	var args struct {
		Operation string  `json:"operation"`
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
	}

	if err := json.Unmarshal([]byte(jsonArgs), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	// Perform the calculation
	var result float64
	switch args.Operation {
	case "add":
		result = args.X + args.Y
	case "subtract":
		result = args.X - args.Y
	case "multiply":
		result = args.X * args.Y
	case "divide":
		if args.Y == 0 {
			return "Error: Division by zero"
		}
		result = args.X / args.Y
	case "power":
		result = math.Pow(args.X, args.Y)
	default:
		return fmt.Sprintf("Unknown operation: %s", args.Operation)
	}

	return fmt.Sprintf("The result of %g %s %g is %g", args.X, args.Operation, args.Y, result)
}

// executeWeatherLookup simulates looking up weather information.
// In a real application, this would call an actual weather API.
func executeWeatherLookup(jsonArgs string) string {
	// Parse the JSON arguments
	var args struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}

	if err := json.Unmarshal([]byte(jsonArgs), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	// Use default unit if not specified
	unit := args.Unit
	if unit == "" {
		unit = "celsius"
	}

	// In a real application, you would call a weather API here
	// This is a simulated response for demonstration purposes
	temp := 22.5
	if unit == "fahrenheit" {
		temp = temp*9/5 + 32
	}

	condition := "Partly Cloudy"
	humidity := 65
	windSpeed := 10.5

	// Get current time for the response
	now := time.Now().Format("2006-01-02 15:04:05")

	return fmt.Sprintf("Weather for %s at %s:\nTemperature: %.1f°%s\nCondition: %s\nHumidity: %d%%\nWind Speed: %.1f km/h",
		args.Location, now, temp, unit[0:1], condition, humidity, windSpeed)
}
