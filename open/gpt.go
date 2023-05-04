package open

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

func InitGPT() {
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "一句话简单介绍一下golang",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
