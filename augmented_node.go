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

// TODO: pass tools to llms and use libraries for now on how to pass tools
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

// TODO:
// Finish the Augmented Node. Test it
// Add other nodes
// Add antrophic llm
// Add Other llms
// Research how to add Retrieval and Memory
