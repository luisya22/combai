package swarmlet

import (
	"fmt"
	"log"
)

type LLMCallNode struct {
	BaseNode
	SystemPrompt   string
	PromptTemplate string
	LLMOptions     LLMOptions
	LLMTools       []LLMTool
	Children       []WorkflowNode
}

func (e *LLMCallNode) Execute(ctx AgentContext, runContext *RunContext, nodeInput ...string) (string, error) {
	anyInput := make([]any, len(nodeInput))
	for i, v := range nodeInput {
		anyInput[i] = v
	}
	fullPrompt := fmt.Sprintf(e.PromptTemplate, anyInput...)
	log.Printf("[LLMCallExecutor] Executing prompt: %s\n", fullPrompt)

	runContext.AddInput(e.nodeID, fullPrompt)

	llmMessage := LLMMessage{
		role:    "user",
		message: fullPrompt,
	}

	output, err := ctx.LLM.Generate(e.LLMOptions, e.LLMTools, e.SystemPrompt, llmMessage)
	if err != nil {
		runContext.AddError(e.nodeID, err)
		return "", err
	}

	runContext.AddOutput(e.nodeID, output)

	for _, cNode := range e.Children {
		_, err := cNode.Execute(ctx, runContext, output)
		if err != nil {
			return "", err
		}
	}

	return output, err
}

type LLMCallOption func(*LLMCallNode)

func WithID(id string) LLMCallOption {
	return func(node *LLMCallNode) {
		node.nodeID = id
	}
}

func WithSystemPrompt(prompt string) LLMCallOption {
	return func(node *LLMCallNode) {
		node.SystemPrompt = prompt
	}
}

func WithPropmtTemplate(prompt string) LLMCallOption {
	return func(node *LLMCallNode) {
		node.PromptTemplate = prompt
	}
}

func WithChildren(children ...WorkflowNode) LLMCallOption {
	return func(node *LLMCallNode) {
		node.Children = append(node.Children, children...)
	}
}

func WithLLMOptions(opts LLMOptions) LLMCallOption {
	return func(node *LLMCallNode) {
		node.LLMOptions = opts
	}
}

func NewLLmCallNode(opts ...LLMCallOption) *LLMCallNode {
	node := &LLMCallNode{
		PromptTemplate: "%s",
		LLMOptions: LLMOptions{
			Temperature: 0.5,
			MaxTokens:   -1,
		},
		Children: []WorkflowNode{},
	}
	for _, opt := range opts {
		opt(node)
	}
	return node
}
