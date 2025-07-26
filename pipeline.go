package swarmlet

import (
	"context"
	"fmt"
	"io"
)

// Holds and executes all the pipeline components
type Pipeline struct {
	Name   string
	Root   WorkflowNode
	LLM    LLM
	Memory Memory
}

func NewPipeline(name string, rootNode WorkflowNode, llm LLM, memory Memory) *Pipeline {
	return &Pipeline{
		Name:   name,
		Root:   rootNode,
		LLM:    llm,
		Memory: memory,
	}
}

func (p *Pipeline) Run(ctx context.Context, initialInput string, runID string, w io.Writer) (finalOutput string, err error) {
	if p.Root == nil {
		return "", fmt.Errorf("pipeline has no root node")
	}

	runContext := NewRunContext(runID, w)

	agentCtx := AgentContext{
		LLM:    p.LLM,
		Memory: p.Memory,
	}

	fmt.Println("1", initialInput)

	err = p.executeNode(ctx, p.Root, initialInput, agentCtx, runContext)
	if err != nil {
		return "", err
	}

	finalOutput = runContext.NodeOutputs[p.Root.ID()]
	return finalOutput, nil
}

func (p *Pipeline) executeNode(
	ctx context.Context,
	node WorkflowNode,
	nodeInput string,
	agentCtx AgentContext,
	runContext *RunContext,
) error {
	_, err := node.Execute(ctx, agentCtx, runContext, nodeInput)
	if err != nil {
		return err
	}

	return nil
}
