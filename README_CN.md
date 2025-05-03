# Grok-Go: X.AI API 的 Go 客户端库

[![Go Reference](https://pkg.go.dev/badge/github.com/SimonMorphy/grok-go.svg)](https://pkg.go.dev/github.com/SimonMorphy/grok-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/SimonMorphy/grok-go)](https://goreportcard.com/report/github.com/SimonMorphy/grok-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

*[English](README.md)*

这是一个轻量级、类型安全的 X.AI API Go 客户端库，可让您在 Go 应用程序中轻松集成和使用 X.AI 的 Grok 模型。

## 特性

- **完整 API 覆盖**：支持所有 X.AI API 端点
- **类型安全**：所有请求和响应均有原生 Go 类型支持
- **流式响应**：高效处理流式 API 响应
- **函数调用**：支持函数调用和工具功能
- **图像生成**：使用 Grok 模型生成图像
- **错误处理**：全面的错误处理机制和详细的错误信息

## 安装

```bash
go get github.com/SimonMorphy/grok-go
```

需要 Go 1.18 或更高版本。

## 快速开始

### 认证

推荐通过环境变量设置 API 密钥：

```bash
export GROK_API_KEY="你的-API-密钥"
```

或者在代码中直接提供（不推荐在生产环境中使用）：

```go
client, err := grok.NewClient("你的-API-密钥") // 生产环境不推荐
```

### 基础聊天示例

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	
	grok "github.com/SimonMorphy/grok-go"
)

func main() {
	// 从环境变量获取 API 密钥
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("未设置 GROK_API_KEY 环境变量")
	}
	
	// 初始化客户端
	client, err := grok.NewClient(apiKey)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	
	// 创建聊天完成请求
	request := &grok.ChatCompletionRequest{
		Model: "grok-3",
		Messages: []grok.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "什么是人工智能？",
			},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}
	
	// 发送请求
	ctx := context.Background()
	response, err := grok.CreateChatCompletion(ctx, client, request)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	
	// 打印响应
	if len(response.Choices) > 0 {
		fmt.Println(response.Choices[0].Message.Content)
	}
}
```

## 高级用法

### 函数调用

```go
// 定义一个计算器工具
calculatorTool := grok.Tool{
    Type: "function",
    Function: grok.Function{
        Name:        "calculate",
        Description: "执行数学计算",
        Parameters: &grok.FunctionParameters{
            Type: "object",
            Properties: map[string]interface{}{
                "operation": map[string]interface{}{
                    "type":        "string",
                    "enum":        []string{"add", "subtract", "multiply", "divide"},
                    "description": "要执行的数学运算",
                },
                "x": map[string]interface{}{
                    "type":        "number",
                    "description": "第一个操作数",
                },
                "y": map[string]interface{}{
                    "type":        "number",
                    "description": "第二个操作数",
                },
            },
            Required: []string{"operation", "x", "y"},
        },
    },
}

// 创建使用该工具的请求
request := &grok.ChatCompletionRequest{
    Model: "grok-3",
    Messages: []grok.ChatCompletionMessage{
        {
            Role:    "user",
            Content: "计算 25 乘以 16",
        },
    },
    Tools:      []grok.Tool{calculatorTool},
    ToolChoice: "auto",
}
```

### 流式响应

```go
request.Stream = true
stream, err := grok.CreateChatCompletionStream(ctx, client, request)
if err != nil {
    log.Fatalf("创建流失败: %v", err)
}
defer stream.Close()

// 处理流
for {
    response, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Printf("流错误: %v", err)
        break
    }
    
    // 处理数据块
    if len(response.Choices) > 0 {
        chunk := response.Choices[0].Delta.Content
        if chunk != "" {
            fmt.Print(chunk)
        }
    }
}
```

### 图像生成

```go
imageClient, err := grok.CreateImageGenerationClient(apiKey)
if err != nil {
    log.Fatalf("创建客户端失败: %v", err)
}

request := &grok.ImageGenerationRequest{
    Model:  "grok-3-image",
    Prompt: "一个未来主义城市，有飞行汽车和高耸的摩天大楼",
    Size:   "1024x1024",
    N:      1,
}

response, err := grok.CreateImage(ctx, imageClient, request)
if err != nil {
    log.Fatalf("图像生成失败: %v", err)
}

// 处理图像 URL 或 Base64 数据
for i, imageData := range response.Data {
    if imageData.URL != "" {
        fmt.Printf("图像 %d URL: %s\n", i+1, imageData.URL)
        // 下载图像...
    } else if imageData.B64JSON != "" {
        fmt.Printf("图像 %d 作为 Base64 数据接收\n", i+1)
        // 保存 Base64 图像...
    }
}
```

## 示例

完整的可运行示例可以在 [`example`](example/) 目录中找到：

- [基础聊天完成](example/chat.go)
- [流式响应](example/streaming.go)
- [函数调用](example/tools.go)
- [图像生成](example/image.go)

## 文档

- [包文档](https://pkg.go.dev/github.com/SimonMorphy/grok-go)
- [API 参考](docs/API.md)
- [图像生成指南](docs/images.md)
- [X.AI API 参考](https://platform.x.ai/docs/api-reference)

## 测试

运行所有测试（集成测试需要 API 密钥）：

```bash
export GROK_API_KEY="你的-API-密钥"
go test ./...
```

仅运行单元测试（不需要 API 密钥）：

```bash
go test ./... -short
```

## 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 贡献

欢迎贡献！您可以：

- 报告 bug
- 请求新功能
- 提交拉取请求

请确保您的代码通过所有测试并遵循 Go 最佳实践。 