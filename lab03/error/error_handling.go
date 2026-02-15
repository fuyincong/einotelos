package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

// generateWithRetry 尝试多次调用 ChatModel 的 Generate 方法，直到成功或达到最大重试次数。
func generateWithRetry(ctx context.Context, chatModel *openai.ChatModel, messages []*schema.Message, maxRetries int) (*schema.Message, error) {
	// 记录最后一次错误
	var lastErr error

	// 重试逻辑
	for i := 0; i < maxRetries; i++ {
		response, err := chatModel.Generate(ctx, messages)
		if err == nil {
			return response, nil
		}

		lastErr = err
		log.Printf("尝试 %d/%d 失败: %v", i+1, maxRetries, err)

		// 指数退避
		if i < maxRetries-1 {
			backoff := time.Duration(1<<uint(i)) * time.Second
			log.Printf("等待 %v 后重试...", backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("重试 %d 次后仍然失败: %w", maxRetries, lastErr)
}

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	// 创建智谱 GLM ChatModel 实例（从 .env 读 GLM_API_KEY）
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GLM_API_KEY"),
		Model:   "GLM-4.7",
		BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	messages := []*schema.Message{
		schema.UserMessage("你好"),
	}

	// 带重试的生成
	response, err := generateWithRetry(ctx, chatModel, messages, 3)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatalf("请求超时")
		}
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("成功! 回答: %s\\n", response.Content)
}
