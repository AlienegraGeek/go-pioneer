package open

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

func InitGPT(content string) (error, string) {
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err, "error"
	}

	fmt.Println("gpt response:", resp.Choices[0].Message.Content)
	return nil, resp.Choices[0].Message.Content
}
