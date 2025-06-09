package main

import (
	"context"
	"dd/agents"
	"dd/helpers"
	"dd/ui"
	"os/exec"

	"fmt"
	"log"
	"strings"

	"github.com/openai/openai-go"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	ctx := context.Background()

	Zephyr, err1 := agents.GetZephyrAgent("ai/qwen2.5:3B-F16")
	// IMPORTANT: the model must support the tools
	// NOTE: ai/qwen2.5:3B-F16 is pretty good
	// NOTE: ai/qwen3:latest try to always answer the question
	Thorin, err2 := agents.GetThorinAgent("ai/qwen2.5:latest")
	Lyralei, err3 := agents.GetLyraleiAgent("ai/qwen2.5:latest")
	Aldric, err4 := agents.GetAldricAgent("ai/qwen2.5:latest")
	Grash, err5 := agents.GetGrashAgent("ai/qwen2.5:latest")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		log.Fatalln("😡:", err1, err2, err3, err4, err5)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	ui.Println(ui.Green, "Bienvenue dans le monde de Dungeons & Dragons !")
	fmt.Println(strings.Repeat("=", 50))

	characters := map[string]*agents.TinyAgent{
		"aldric":  Aldric,
		"grash":   Grash,
		"lyralei": Lyralei,
		"thorin":  Thorin,
		"zephyr":  Zephyr,
	}
	selectedCharacter, exists := characters["lyralei"]
	if !exists {
		log.Fatalln("😡 character not found")
	}
	toolsMaster := characters["zephyr"]
	if !exists {
		log.Fatalln("😡 character not found")
	}

	// BEGIN: MCP SECTION STEP 1: Initialize MCP server
	cmd := exec.Command(
		"docker",
		"run",
		"-i",
		"--rm",
		"mcp-dd:demo",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("😡 Failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("😡 Failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("😡 Failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	clientTransport := stdio.NewStdioServerTransportWithIO(stdout, stdin)

	// STEP 2: Create a new MCP client
	mcpClient := mcp_golang.NewClient(clientTransport)

	if _, err := mcpClient.Initialize(ctx); err != nil {
		log.Fatalf("😡 Failed to initialize client: %v", err)
	}

	// STEP 3: Get the list of the available MCP tools
	mcpTools, err := mcpClient.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("😡 Failed to list tools: %v", err)
	}
	// STEP 4: IMPORTANT: Convert the MCP tools to OpenAI tools
	openAITools := helpers.ConvertToOpenAITools(mcpTools)

	// STEP 5: Register the tools to the tools master agent
	toolsMaster.Params.Tools = openAITools
	
	// END: MCP SECTION

	/*
		je souhaite parler avec un nain
		je veux discuter avec une elfe
		j'ai une question sur la magie
	*/

	for {
		question, _ := ui.Input(
			"#660707",
			fmt.Sprintf("%s [%s] que puis-je faire pour toi ? ",
				selectedCharacter.Avatar, selectedCharacter.Name),
		)
		if question == "bye" {
			break
		}

		ui.Println(ui.Green, "⏳ checking...")

		// STEP 1: tools detection
		// Create messages list with the question
		toolsMaster.Params.Messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}
		// Run the tools detection (completion)
		detectedToolCalls, _ := toolsMaster.ToolsCompletion()

		// STEP 2: BEGIN: of tools execution
		if len(detectedToolCalls) > 0 {
			for _, toolCall := range detectedToolCalls {
				// Display the detected tool call
				ui.Println(ui.Blue, "💡 tool detection:", toolCall.Function.Name, toolCall.Function.Arguments)

				args := helpers.ExtractArgsFromJSONString(toolCall.Function.Arguments)

				switch toolCall.Function.Name {
				// TOOL 1:
				case "choisir_un_personnage_par_son_espece":
					//selectedCharacter = characters[tools.ChooseCharacterBySpecies(args)]

					// NOTE: Call the MCP tool with the arguments
					toolResponse, _ := mcpClient.CallTool(ctx, toolCall.Function.Name, args)
					characterName := toolResponse.Content[0].TextContent.Text
					selectedCharacter = characters[characterName]


				// TOOL 2:
				case "detecter_le_vrai_sujet_du_message_utilisateur":
					//selectedCharacter = characters[tools.ChooseCharacterFromTopic(args)]
					
					// NOTE: Call the MCP tool with the arguments
					toolResponse, _ := mcpClient.CallTool(ctx, toolCall.Function.Name, args)
					characterName := toolResponse.Content[0].TextContent.Text
					selectedCharacter = characters[characterName]

				default:
					ui.Println(ui.Red, "❌ Error: unknown tool", toolCall.Function.Name)

				}
			}
		} // END: of tools execution
		// Reset the messages for the tools master
		toolsMaster.Params.Messages = []openai.ChatCompletionMessageParamUnion{}

		// STEP 3: chat with the selected character / chat completion
		// Add the question to the messages
		selectedCharacter.Params.Messages = append(selectedCharacter.Params.Messages,
			openai.UserMessage(question),
		)

		// TODO:
		// display who is speaking...
		// add a color to the struct agent
		ui.Println(ui.Magenta, "[", selectedCharacter.Avatar, "]", selectedCharacter.Name, "is speaking...")

		// Run the chat completion
		answer, _ := selectedCharacter.ChatCompletionStream(func(self *agents.TinyAgent, content string, err error) error {
			ui.Print(selectedCharacter.Color, content)
			return nil
		})

		// IMPORTANT: avoid to answering the same question twice
		selectedCharacter.Params.Messages = append(
			selectedCharacter.Params.Messages,
			openai.AssistantMessage(answer),
		)

	}

}
