package main

import (
	"dd/agents"
	"fmt"
	"log"
	"strings"

	"github.com/openai/openai-go"
)

func main() {
	// TODO: show the code of the agents

	// STEP 1: Creation of the agents
	Lyralei, err := agents.GetLyraleiAgent("ai/qwen2.5:latest")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	Thorin, err := agents.GetThorinAgent("ai/qwen2.5:latest")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// STEP 2: add a user message to the list of messages
	Lyralei.Params.Messages = append(
		Lyralei.Params.Messages,
		openai.UserMessage("Qui est tu?"),
	)

	Thorin.Params.Messages = append( // NOTE: with append, we can keep the conversationnal memory
		Thorin.Params.Messages,
		openai.UserMessage("Qui es-tu?"),
	)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" Lyralei is speaking...")
	fmt.Println(strings.Repeat("=", 50))

	// STEP 3: ask the question to Lyralei == completion
	answer, _ := Lyralei.ChatCompletionStream(func(self *agents.TinyAgent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	// STEP 4: add the answer to the conversation memory
	// And resume the conversation with a new question
	Lyralei.Params.Messages = append(
		Lyralei.Params.Messages,
		openai.AssistantMessage(answer), // IMPORTANT: avoid to answering the same question twice
		openai.UserMessage("Quelles sont tes motivations?"),
	)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" Thorin is speaking...")
	fmt.Println(strings.Repeat("=", 50))

	// 
	Thorin.ChatCompletionStream(func(self *agents.TinyAgent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" Lyralei is speaking...")
	fmt.Println(strings.Repeat("=", 50))

	Lyralei.ChatCompletionStream(func(self *agents.TinyAgent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

}
