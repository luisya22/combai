package main

import (
	"context"
	"log"
	"os"

	"github.com/luisya22/swarmlet"
)

func main() {

	tools := []swarmlet.LLMTool{
		{
			Name:        "get_temperature",
			Description: "Get temperature from any country",
			Params: map[string]swarmlet.LLMToolFieldProperty{
				"name": {
					Type:        "string",
					Description: "API to get temperature from any country",
				},
			},
			Executor: func(args map[string]any) (string, error) {
				return "400 F", nil
			},
		},
	}

	output := swarmlet.NewOutputNode("output", "1", true)

	node1 := swarmlet.NewAugmentedLLMNode(
		swarmlet.WithAugmentedID("1"),
		swarmlet.WithAugmentedChildren(output),
		swarmlet.WithAugmentedTools(tools...),
	)

	apiKey := os.Getenv("LLM_API_KEY")

	llm := swarmlet.NewOpenAILLM(apiKey, "gpt-4o-mini")

	memory := swarmlet.NewDummyMemory()

	pipeline := swarmlet.NewPipeline("Pipeline", node1, llm, memory)

	stdWriter := os.Stdout

	// San Juan,PR
	_, err := pipeline.Run(context.Background(), "What is the temperature in Nebraska", "102", stdWriter)
	if err != nil {
		log.Fatal(err)
	}
}
