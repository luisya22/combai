package swarmlet

import (
	"encoding/json"
	"fmt"
	"log"
)

type AugmentedLLMNode struct {
	BaseNode
	systemPrompt      string
	promptTemplate    string
	LLMOptions        LLMOptions
	tools             []LLMTool
	Children          []WorkflowNode
	maxToolIterations int
}

type AugmentedLLMNodeOption func(*AugmentedLLMNode)

func WithAugmentedID(id string) AugmentedLLMNodeOption {
	return func(node *AugmentedLLMNode) {
		node.nodeID = id
	}
}

func WithAugmentedSystemPrompt(prompt string) AugmentedLLMNodeOption {
	return func(node *AugmentedLLMNode) {
		node.systemPrompt = prompt
	}
}

func WithAugmentedPromptTemplate(template string) AugmentedLLMNodeOption {
	return func(node *AugmentedLLMNode) {
		node.promptTemplate = template
	}
}

func WithAugmentedChildren(children ...WorkflowNode) AugmentedLLMNodeOption {
	return func(node *AugmentedLLMNode) {
		node.Children = append(node.Children, children...)
	}
}

func WithAugmentedLLMOptions(opts LLMOptions) AugmentedLLMNodeOption {
	return func(node *AugmentedLLMNode) {
		node.LLMOptions = opts
	}
}

func WithAugmentedTools(tools ...LLMTool) AugmentedLLMNodeOption {
	return func(node *AugmentedLLMNode) {
		node.tools = append(node.tools, tools...)
	}
}

func NewAugmentedLLMNode(opts ...AugmentedLLMNodeOption) *AugmentedLLMNode {
	node := &AugmentedLLMNode{
		systemPrompt:   DefaultAugmentedSystemPrompt,
		promptTemplate: "%s",
		LLMOptions: LLMOptions{
			Temperature: 0.5,
			MaxTokens:   -1,
		},
		Children:          []WorkflowNode{},
		maxToolIterations: 5,
	}
	for _, opt := range opts {
		opt(node)
	}
	return node
}

func (e *AugmentedLLMNode) Execute(ctx AgentContext, runCtx *RunContext, nodeInput ...string) (string, error) {
	anyInput := make([]any, len(nodeInput))
	for i, v := range nodeInput {
		anyInput[i] = v
	}
	fullPrompt := fmt.Sprintf(e.promptTemplate, anyInput...)
	log.Printf("[LLMCallExecutor] Executing prompt: %s\n", fullPrompt)

	currentUserMessage := LLMMessage{
		role:    "user",
		message: fullPrompt,
	}

	runCtx.AddMessage(e.nodeID, currentUserMessage)
	runCtx.AddInput(e.nodeID, fullPrompt)

	finalResponse := ""

	for i := range e.maxToolIterations {
		log.Printf("[AugmentedLLMNode-%s] Iteration %d: Calling LLM with %d messages and %d tools.\n", e.nodeID, i+1, len(runCtx.MessageHistory), len(e.tools))

		messages, _ := runCtx.GetMessages(e.nodeID)
		llmResponse, err := ctx.LLM.Generate(e.LLMOptions, e.tools, e.systemPrompt, messages...)
		if err != nil {
			runCtx.AddError(e.nodeID, err)
			return "", err
		}

		runCtx.AddMessage(e.nodeID, llmResponse)

		if llmResponse.message != "" {
			log.Printf("[Augmented-%s] LLM responded with content.\n", e.nodeID)
			runCtx.AddOutput(e.nodeID, llmResponse.message)
		} else if llmResponse.toolCalls != nil {
			log.Printf("[AugmentedLLMNode-%s] LLM requested %d tool calls.\n", e.nodeID, len(llmResponse.toolCalls))
		}

		if llmResponse.message != "" && (llmResponse.toolCalls == nil || len(llmResponse.toolCalls) == 0) {
			finalResponse = llmResponse.message
			break
		}

		log.Printf("[AugmentedLLMNode-%s] Executing requested tools...\n", e.nodeID)
		for _, toolCall := range llmResponse.toolCalls {
			toolName := toolCall.function.name
			toolArgsRaw := toolCall.function.arguments
			toolCallID := toolCall.id

			foundTool := false
			for _, registeredTool := range e.tools {
				if registeredTool.Name == toolName {
					foundTool = true

					var argsMap map[string]any
					if err := json.Unmarshal([]byte(toolArgsRaw), &argsMap); err != nil {
						toolErrorMsg := fmt.Sprintf("Error unmarshaling tool arguments for '%s' (ID: %s): %v", toolName, toolCallID, err)
						log.Printf("[AugmentedLLMNode-%s] %s\n", e.nodeID, toolErrorMsg)

						runCtx.AddMessage(e.nodeID, LLMMessage{
							role:       "tool",
							message:    toolErrorMsg,
							toolCallId: toolCallID,
						})
						continue
					}

					toolOutput, err := registeredTool.Executor(argsMap)
					if err != nil {
						toolErrorMsg := fmt.Sprintf("Error executing tool '%s' (ID: %s): %v", toolName, toolCallID, err)
						log.Printf("[AugmentedLLMNode-%s] %s\n", e.nodeID, toolErrorMsg)
						runCtx.AddMessage(e.nodeID, LLMMessage{
							role:       "tool",
							message:    toolErrorMsg,
							toolCallId: toolCallID,
						})

						continue
					}

					log.Printf("[AugmentedLLMNode-%s] Tool '%s' (ID: %s) executed successfully. Output: %s\n",
						e.nodeID, toolName, toolCallID, toolOutput)

					runCtx.AddMessage(e.nodeID, LLMMessage{
						role:       "tool",
						message:    toolOutput,
						toolCallId: toolCallID,
					})
				}
			}

			if !foundTool {
				errorMsg := fmt.Sprintf("LLM Requested unknown tool: '%s' (ID: %s)", toolName, toolCallID)
				log.Printf("[AugmentedLLMNode-%s] %s\n", e.nodeID, errorMsg)
				runCtx.AddMessage(e.nodeID, LLMMessage{
					role:       "tool",
					message:    errorMsg,
					toolCallId: toolCallID,
				})
			}
		}
	}

	if finalResponse == "" {
		finalResponse = "The AI assistant could not fully resolve the request after multiple attempts."
		runCtx.AddError(e.nodeID, fmt.Errorf("max tool iterations reached without a final response"))
		log.Printf("[AugmentedLLMNode-%s] Max tool iterations reached without a final response.\n", e.nodeID)
	}

	for _, cNode := range e.Children {
		_, err := cNode.Execute(ctx, runCtx, finalResponse)
		if err != nil {
			return "", err
		}
	}

	return finalResponse, nil
}

const DefaultAugmentedSystemPrompt = `
You are "Swarmlet-Assistant", an intelligent, helpful, and highly capable AI assistant.
Your primary directive is to understand and fulfill user requests by leveraging the specialized tools at your disposal.

Here are your key operating principles:
1.  **Prioritize Tool Use**: If a user's request can be fulfilled by one or more of your tools, you MUST use the appropriate tool(s) first. Do not attempt to answer questions based on general knowledge if a tool is more relevant or required for accuracy.
2.  **Transparent Tooling**: When you decide to use a tool, acknowledge this intention or briefly explain what tool you are using and why, before providing the final answer (e.g., "I'm checking the knowledge base for that...").
3.  **Synthesize Results**: After executing tools and receiving their outputs, synthesize the information into a clear, concise, and helpful response for the user. Do not just return raw tool output.
4.  **Ask for Clarification**: If a request is ambiguous, lacks necessary parameters for a tool, or you need more information to proceed, politely ask the user clarifying questions.
5.  **Handle Limitations**: If a request is outside your current capabilities or the scope of your available tools, or if a tool execution fails, inform the user gracefully and suggest what you *can* do instead.
6.  **Maintain Context**: Remember and utilize information from previous turns of conversation to provide coherent and relevant responses.
7.  **Be Polite and Professional**: Always maintain a helpful, calm, and professional demeanor throughout the interaction.

You have access to a suite of specialized tools to assist you. Utilize them wisely.
`
