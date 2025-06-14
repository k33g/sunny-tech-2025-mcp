package main

import (
	"dd/agents"
	"dd/helpers"
	"fmt"
	"log"
	"strings"

	"github.com/openai/openai-go"
)

func main() {

	//Zephyr, err := agents.GetZephyrAgent("ai/qwen2.5:3B-F16")
	Zephyr, err := agents.GetZephyrAgent("ai/qwen2.5:latest")
	//Zephyr, err := agents.GetZephyrAgent("ignaciolopezluna020/watt-tool:8B-Q4_K_M")
	//Zephyr, err := agents.GetZephyrAgent("ignaciolopezluna020/llama-xlam:8B-Q4_K_M")
	//Zephyr, err := agents.GetZephyrAgent("k33g/qwen2.5:0.5b-instruct-q8_0")


	if err != nil {
		log.Fatalln("😡:", err)
	}

	Zephyr.Params.Messages = append(
		Zephyr.Params.Messages,
		openai.UserMessage(`
			je souhaite parler avec un nain

			je veux discuter avec une elfe

			j'ai une question sur la magie
		`),
	)

	// IMPORTANT: show code of Zephyr + tools catalog + tools completion

	detectedToolCalls, _ := Zephyr.ToolsCompletion()

	if len(detectedToolCalls) == 0 {
		fmt.Println("😡 No function call detected")
		fmt.Println()
		return
	}
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" Zephyr is speaking...")
	fmt.Println(strings.Repeat("=", 50))

	for _, toolCall := range detectedToolCalls {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println(toolCall.Function.Name, toolCall.Function.Arguments)
		fmt.Println(helpers.ToolCallsToJSONString([]openai.ChatCompletionMessageToolCall{toolCall}))
	}

}
