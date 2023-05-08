package open

import (
	"AlienegraGeek/go-pioneer/util"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
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

func FetchContextGPT(talker string, pro string) (string, error) {
	client := openai.NewClient(AUTHTOKEN)
	messages := make([]openai.ChatCompletionMessage, 0)
	//historyMsg := make([]openai.ChatCompletionMessage, 0)
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("-> ")
	//text, _ := reader.ReadString('\n')
	//// convert CRLF to LF
	//text = strings.Replace(text, "\n", "", -1)

	// 获取全局缓存实例
	c := util.GetCacheInstance()

	// 从缓存中获取 key 为 "talker" 的历史消息记录
	if historyMsg, found := c.Get(talker); found {
		fmt.Println("talker:", historyMsg)
		historyMsg = append(historyMsg.([]openai.ChatCompletionMessage), openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: pro,
		})
		// 将字符串切片转换为数组类型
		//copy(messages, historyMsg.([]openai.ChatCompletionMessage))
		messages = historyMsg.([]openai.ChatCompletionMessage)
	} else {
		fmt.Println("talker not found")
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: pro,
		})
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			//Model:    openai.GPT3Dot5Turbo,
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
			User:     talker,
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
	// 存储一个字符串类型的缓存数据，key 为 "用户名"，value 为 "消息"，过期时间为默认的 5 分钟
	c.Set(talker, messages, 0)
	return content, nil
}

var MsgTree map[string]openai.ChatCompletionMessage

func AddNewMsg(treeId string, content string) {
	MsgTree[treeId] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
}
