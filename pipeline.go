package combai

import (
	"fmt"
)

// Holds and executes all the pipeline components

type Pipeline struct {
	Name string
	Root WorkflowNode
	// Add base LLM - each node can have a different llm
	// Memory
}

type NodeType int

const (
	LLM_CALL NodeType = iota
	GATE
	ROUTER
	ORCHESTRATOR
	EVALUATOR
)

type WorkflowNode interface {
	ID() string
	Type() NodeType
	ExecuteLogic(input string) (string, error)
	GetChildren() []WorkflowNode
}

type NodeExecutor interface {
	Execute(nodeInput string) error
}

type RunContext struct {
	RunID       string
	NodeInputs  map[string]string
	NodeOutputs map[string]string
	NodeErrors  map[string]error
}

func NewRunContext(runID string) *RunContext {
	return &RunContext{
		RunID:       runID,
		NodeInputs:  make(map[string]string),
		NodeOutputs: make(map[string]string),
		NodeErrors:  make(map[string]error),
	}
}

func (p *Pipeline) Run(initialInput string, runID string) (finalOutput string, err error) {
	if p.Root == nil {
		return "", fmt.Errorf("pipeline has no root node")
	}

	runContext := NewRunContext(runID)

	err = p.executeNode(p.Root, initialInput, runContext)
	if err != nil {
		return "", err
	}

	finalOutput = runContext.NodeOutputs[p.Root.ID()]
	return finalOutput, nil
}

func (p *Pipeline) executeNode(
	node WorkflowNode,
	nodeInput string,
	runtContext *RunContext,
) error {
	output, err := node.ExecuteLogic(nodeInput)
	if err != nil {
		return err
	}

	for _, child := range node.GetChildren() {
		if err := p.executeNode(child, output, runtContext); err != nil {
			return err
		}
	}

	return nil
}
