package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	// 对话历史
	messages := []*schema.Message{
		schema.SystemMessage("你是一个懂得哲学的程序员。"),
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("开始对话（输入 'exit' 退出）：")

	for {
		fmt.Print("\\n你: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "exit" {
			fmt.Println("再见！")
			break
		}

		if userInput == "" {
			continue
		}

		// 添加用户消息
		messages = append(messages, schema.UserMessage(userInput))

		// 生成 AI 响应
		response, err := chatModel.Generate(ctx, messages)
		if err != nil {
			log.Printf("生成失败: %v", err)
			continue
		}

		// 添加 AI 响应到历史
		messages = append(messages, response)

		fmt.Printf("\\nAI: %s\\n", response.Content)
	}
}
