package swarmlet

import (
	"context"
	"log"
	"time"

	"github.com/sashabaranov/go-openai"
)

type LLM interface {
	Generate(options LLMOptions, prompt string, messages ...LLMMessage) (string, error)
}

type LLMOptions struct {
	Model       string
	MaxTokens   int
	Temperature float64
}

type LLMMessage struct {
	message string
	role    string
}

type DummyLLM struct{}

func (d *DummyLLM) Generate(propmt string, options LLMOptions) (string, error) {
	time.Sleep(50 * time.Millisecond)
	log.Printf("(DummyLLM) Generated for: \"%s\"", propmt)
	return "Simulated LLM response for: " + propmt, nil
}

type ReverseLLM struct{}

func (d *ReverseLLM) Generate(options LLMOptions, systemPrompt string, messages ...LLMMessage) (string, error) {
	time.Sleep(50 * time.Millisecond)
	output := reverseString(messages[0].message)
	log.Printf("(ReverseLLM) Reversed: \"%s\"", output)

	return output, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := range len(runes) / 2 {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	return string(runes)
}

type OpenAILLM struct {
	ApiKey string
}

func (llm *OpenAILLM) Generate(options LLMOptions, systemMessage string, messages ...LLMMessage) (string, error) {
	client := openai.NewClient(llm.ApiKey)

	llmMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
	}
	for _, m := range messages {
		openAIMessasge := openai.ChatCompletionMessage{
			Role:    m.role,
			Content: m.message,
		}

		llmMessages = append(llmMessages, openAIMessasge)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: llmMessages,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
