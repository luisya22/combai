package swarmlet

import "context"

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
	Execute(ctx context.Context, agentContext AgentContext, runContext *RunContext, input ...string) (string, error)
}

type BaseNode struct {
	nodeID   string
	nodeType string
}

func (b *BaseNode) ID() string   { return b.nodeID }
func (b *BaseNode) Type() string { return b.nodeType }

type MemoryAndStreamingConfig struct {
	UseMemory          bool
	MemoryKey          string
	Streaming          bool
	MaxHistoryMessages int
}

type AgenticLLMNode struct {
	InitialPropmtTemplate string
	MaxIterations         int
}
