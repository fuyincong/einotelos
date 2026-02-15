package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	// 创建 ChatModel（智谱 GLM，从 .env 读 GLM_API_KEY）
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GLM_API_KEY"),
		Model:   "GLM-4.7",
		BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	// 构建消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个懂得哲学的程序员。"),
		schema.UserMessage("什么是存在主义？"),
	}

	// 生成响应
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("回答:\\n%s\\n", response.Content)
}
