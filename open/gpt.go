package open

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

const AUTHTOKEN = "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"

func FetchGPT(content string) (string, error) {
	client := openai.NewClient(AUTHTOKEN)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "g res error", err
	}
	//fmt.Println("g response:", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}

func FetchContextGPT() (string, error) {
	client := openai.NewClient(AUTHTOKEN)
	messages := make([]openai.ChatCompletionMessage, 0)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	content := resp.Choices[0].Message.Content
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})
	fmt.Println(content)
	return content, nil
}

var MsgTree map[string]openai.ChatCompletionMessage

func AddNewMsg(treeId string, content string) {
	MsgTree[treeId] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
}
