package main

import (
	"log"
	"os"

	"github.com/luisya22/swarmlet"
)

func main() {
	output := swarmlet.NewOutputNode("output", "3", true)
	node1 := swarmlet.NewLLmCallNode(
		swarmlet.WithID("1"),
		swarmlet.WithChildren(output),
		swarmlet.WithSystemPrompt("You are a temperature expert. I will give you a temperature in Celsius and you will return in Farenheit. Return just a simple string with the temperature."),
	)
	node2 := swarmlet.NewLLmCallNode(
		swarmlet.WithID("2"),
		swarmlet.WithChildren(node1),
		swarmlet.WithSystemPrompt("You are a temperature expert. Give me plain string of average temperature in this city in Celsius."),
	)
	node3 := swarmlet.NewLLmCallNode(
		swarmlet.WithID("3"),
		swarmlet.WithChildren(node2),
		swarmlet.WithSystemPrompt("You are a reverser agent. Return plain message string of the reversed input"),
	)

	llm := swarmlet.OpenAILLM{
		ApiKey: os.Getenv("LLM_API_KEY"),
	}
	memory := swarmlet.NewDummyMemory()

	pipeline := swarmlet.NewPipeline("Pipeline", node3, &llm, memory)

	stdWriter := os.Stdout

	// San Juan,PR
	_, err := pipeline.Run("RP, nauJ naS", "102", stdWriter)
	if err != nil {
		log.Fatal(err)
	}
}
