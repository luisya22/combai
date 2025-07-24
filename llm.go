package swarmlet

import (
	"log"
	"time"
)

type LLM interface {
	Generate(options LLMOptions, tools []LLMTool, prompt string, messages ...LLMMessage) (LLMMessage, error)
}

type LLMOptions struct {
	Model       string
	MaxTokens   int
	Temperature float32
}

// TODO: API to pass tools to LLM, each node could have an individual LLM
// Need to pass the model to the llm
type LLMTool struct {
	Name        string
	Description string
	Params      map[string]LLMToolFieldProperty
	Executor    func(map[string]any) (string, error)
}

type LLMToolFieldProperty struct {
	Type        string
	Description string
	Enum        []string
}

type LLMMessage struct {
	message    string
	role       string
	toolCallId string
	toolCalls  []LLMToolCall
}

type LLMToolCall struct {
	index    *int
	id       string
	toolType string
	function LLMFunctionCall
}

type LLMFunctionCall struct {
	name      string
	arguments string
}

type DummyLLM struct{}

func (d *DummyLLM) Generate(propmt string, options LLMOptions) (string, error) {
	time.Sleep(50 * time.Millisecond)
	log.Printf("(DummyLLM) Generated for: \"%s\"", propmt)
	return "Simulated LLM response for: " + propmt, nil
}

type ReverseLLM struct{}

func (d *ReverseLLM) Generate(options LLMOptions, tools []LLMTool, systemPrompt string, messages ...LLMMessage) (string, error) {
	time.Sleep(50 * time.Millisecond)
	output := reverseString(messages[0].message)
	log.Printf("(ReverseLLM) Reversed: \"%s\"", output)

	return output, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := range len(runes) / 2 {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	return string(runes)
}
