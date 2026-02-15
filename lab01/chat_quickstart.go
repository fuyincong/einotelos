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
	// 从项目根目录的 .env 加载环境变量
	_ = godotenv.Load()

	// 1. 创建上下文
	ctx := context.Background()

	// 2. 创建 ChatModel 实例（智谱 GLM，OpenAI 兼容 API）
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GLM_API_KEY"),
		Model:   "GLM-4.7",
		BaseURL: "https://open.bigmodel.cn/api/coding/paas/v4",
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	// 3. 准备发送聊天请求
	messages := []*schema.Message{
		schema.SystemMessage("你是一个知识渊博的篮球解说员"), // 系统消息，设定对话背景
		schema.UserMessage("科比为什么被人叫做牢大？"),     // 用户消息，提出问题
	}

	// 4. 调用模型生成响应
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}

	// 5. 输出结果
	fmt.Printf("AI 响应: %s\n", response.Content)

	// 6. 输出 token 使用情况（可选项）
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Printf("\nToken 使用统计:\n")
		fmt.Printf("  输入 Token: %d\n", response.ResponseMeta.Usage.PromptTokens)
		fmt.Printf("  输出 Token: %d\n", response.ResponseMeta.Usage.CompletionTokens)
		fmt.Printf("  总计 Token: %d\n", response.ResponseMeta.Usage.TotalTokens)
	}
}
