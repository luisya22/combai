package swarmlet

type AugmentedLLMNode struct {
	BaseNode
	systemPrompt   string
	promptTemplate string
	LLMOptions     LLMOptions
	Children       []WorkflowNode
}

// func (e *AgenticLLMNode) Execute(ctx AgentContext, runCtx *RunContext, nodeInput ...string) (string, error) {
//
// }

// task: answer user queries, perform actions using available tools, retrieve relevant information, and leverage conversational memory
// goal: is to provide the most accurate, comprohensive, and helpful response possible

var systemPrompt = `
	You are an intelligent assistante designed to %s. Your goal %s.

	Here are the capabilities you can use:

	<capabilities>
		%s
	</capabilities>

	<conversation_history>
		%s
	</conversation_history>


`
