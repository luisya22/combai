package swarmlet

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type openAILLM struct {
	apiKey string
	model  string
}

func NewOpenAILLM(apiKey string, model string) *openAILLM {
	return &openAILLM{
		apiKey: apiKey,
		model:  model,
	}
}

func (llm *openAILLM) Generate(options LLMOptions, tools []LLMTool, systemMessage string, messages ...LLMMessage) (string, error) {
	client := openai.NewClient(llm.apiKey)

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

	llmTools := []openai.Tool{}
	for _, t := range tools {
		properties := getOpenAIParams(t.Params)

		params := jsonschema.Definition{
			Type:       jsonschema.Object,
			Properties: properties,
		}

		f := openai.FunctionDefinition{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  params,
		}

		llmTools = append(llmTools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &f,
		})
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    llm.model,
			Messages: llmMessages,
			Tools:    llmTools,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func getOpenAIParams(params map[string]LLMToolFieldProperty) map[string]jsonschema.Definition {
	llmParams := make(map[string]jsonschema.Definition, len(params))
	for n, p := range params {
		llmParams[n] = jsonschema.Definition{
			Type:        jsonschema.DataType(p.Type),
			Description: p.Description,
			Enum:        p.Enum,
		}
	}

	return llmParams
}
