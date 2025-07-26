package swarmlet

import "fmt"

type OutputNode struct {
	BaseNode
	FromNode string
	Visible  bool
}

func NewOutputNode(id string, fromNode string, visible bool) *OutputNode {
	return &OutputNode{
		BaseNode: BaseNode{
			nodeID: id,
		},
		FromNode: fromNode,
		Visible:  visible,
	}
}

func (n *OutputNode) Execute(ctx AgentContext, runContext *RunContext, input ...string) (string, error) {
	output, ok := runContext.GetOutput(n.FromNode)
	if !ok {
		return "", fmt.Errorf("%s: no output found for node '%s'", n.ID(), n.FromNode)
	}

	if n.Visible && runContext.StreamWriter != nil {
		_, err := runContext.StreamWriter.Write(fmt.Appendln([]byte{}, output))
		if err != nil {
			runContext.AddError(n.ID(), err)
			return "", err
		}
	}

	runContext.AddOutput(n.ID(), output)
	return "", nil
}
