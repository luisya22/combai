package main

import (
	"log"
	"os"

	"github.com/luisya22/swarmlet"
)

func main() {
	output := swarmlet.NewOutputNode("output", "1", true)
	node1 := swarmlet.NewLLmCallNode(
		swarmlet.WithID("1"),
		swarmlet.WithChildren(output),
		swarmlet.WithSystemPrompt("You are a travel agent recommending places to go"),
	)

	apiKey := os.Getenv("LLM_API_KEY")

	llm := swarmlet.NewOpenAILLM(apiKey, "gpt-4o-mini")
	memory := swarmlet.NewDummyMemory()

	pipeline := swarmlet.NewPipeline("Pipeline", node1, llm, memory)

	stdWriter := os.Stdout

	_, err := pipeline.Run("I want to go to San Juan", "102", stdWriter)
	if err != nil {
		log.Fatal(err)
	}
}
