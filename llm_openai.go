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

func (llm *openAILLM) Generate(ctx context.Context, options LLMOptions, tools []LLMTool, systemMessage string, messages ...LLMMessage) (LLMMessage, error) {
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
		if m.role == openai.ChatMessageRoleTool {
			openAIMessasge.ToolCallID = m.toolCallId
		}

		if m.role == openai.ChatMessageRoleAssistant && len(m.toolCalls) > 0 {
			for _, tc := range m.toolCalls {
				openAIToolCall := openai.ToolCall{
					Index: tc.index,
					ID:    tc.id,
					Type:  openai.ToolType(tc.toolType),
					Function: openai.FunctionCall{
						Name:      tc.function.name,
						Arguments: tc.function.arguments,
					},
				}

				openAIMessasge.ToolCalls = append(openAIMessasge.ToolCalls, openAIToolCall)
			}
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

	req := openai.ChatCompletionRequest{
		Model:    llm.model,
		Messages: llmMessages,
		Tools:    llmTools,
	}

	if options.Temperature > 0 {
		req.Temperature = options.Temperature
	}

	if options.MaxTokens > 0 {
		req.MaxTokens = options.MaxTokens
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return LLMMessage{}, err
	}

	responseMessage := LLMMessage{
		role:    resp.Choices[0].Message.Role,
		message: resp.Choices[0].Message.Content,
	}

	if resp.Choices[0].Message.ToolCalls != nil && len(resp.Choices[0].Message.ToolCalls) > 0 {
		for _, tc := range resp.Choices[0].Message.ToolCalls {
			toolCall := LLMToolCall{
				index:    tc.Index,
				id:       tc.ID,
				toolType: string(tc.Type),
				function: LLMFunctionCall{
					name:      tc.Function.Name,
					arguments: tc.Function.Arguments,
				},
			}

			responseMessage.toolCalls = append(responseMessage.toolCalls, toolCall)
		}
	}

	return responseMessage, nil
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

// TODO: When responding with tool you need to add the tool call to the messages and then add the response using toolid and also add to messages
