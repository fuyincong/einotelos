package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	// 示例1: 基础配置
	fmt.Println("=== 示例1: 基础配置 ===")
	basicExample(ctx)

	// 示例2: 高级配置
	fmt.Println("\\n=== 示例2: 高级配置 ===")
	advancedExample(ctx)

	// 示例3: 创意写作配置
	fmt.Println("\\n=== 示例3: 创意写作配置 ===")
	creativeExample(ctx)
}

// 基础配置示例
func basicExample(ctx context.Context) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GLM_API_KEY"),
		Model:   "GLM-4.7",
		BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage("你是一个友好的 AI 助手"),
		schema.UserMessage("用一句话介绍 Eino 框架"),
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("AI 响应: %s\\n", response.Content)
	printTokenUsage(response)
}

// 高级配置示例 - 精确控制输出
func advancedExample(ctx context.Context) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GLM_API_KEY"),
		Model:   "GLM-4.7",
		BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
		Timeout: 30 * time.Second,
		Temperature: floatPtr(0.7),
		TopP:         floatPtr(0.9),
		MaxTokens:    intPtr(500),
		Stop:         []string{"\n\n", "总结:"},
		PresencePenalty:  floatPtr(0.6),
		FrequencyPenalty: floatPtr(0.5),
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage("你是一个专业的技术文档撰写专家"),
		schema.UserMessage("详细介绍 Eino 框架的核心特性，包括架构、组件和优势"),
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("AI 响应: %s\\n", response.Content)
	printTokenUsage(response)
}

// 创意写作配置示例 - 高随机性
func creativeExample(ctx context.Context) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GLM_API_KEY"),
		Model:   "GLM-4.7",
		BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
		Temperature: floatPtr(1.2),
		TopP:         floatPtr(0.95),
		MaxTokens:    intPtr(800),
		PresencePenalty:  floatPtr(0.3),
		FrequencyPenalty: floatPtr(0.3),
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage("你是一个富有创造力的故事作家"),
		schema.UserMessage("创作一个关于 AI 框架变成超级英雄的有趣故事开头"),
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("AI 响应: %s\\n", response.Content)
	printTokenUsage(response)
}

func floatPtr(f float32) *float32 { return &f }
func intPtr(i int) *int           { return &i }

// 打印 Token 使用情况
func printTokenUsage(response *schema.Message) {
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Printf("\\nToken 使用统计:\\n")
		fmt.Printf("  输入 Token: %d\\n", response.ResponseMeta.Usage.PromptTokens)
		fmt.Printf("  输出 Token: %d\\n", response.ResponseMeta.Usage.CompletionTokens)
		fmt.Printf("  总计 Token: %d\\n", response.ResponseMeta.Usage.TotalTokens)
		if response.ResponseMeta.Usage.PromptTokenDetails.CachedTokens > 0 {
			fmt.Printf("  缓存 Token: %d\\n", response.ResponseMeta.Usage.PromptTokenDetails.CachedTokens)
		}
	}
}
