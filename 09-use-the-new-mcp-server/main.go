package main

import (
	"context"
	"dd/agents"
	"dd/helpers"
	"dd/ui"
	"os"

	"fmt"
	"log"
	"strings"

	//mcp_golang "github.com/metoro-io/mcp-golang"
	//mcp_http "github.com/metoro-io/mcp-golang/transport/http"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to load .env file:", err)
	}
	// Ensure MODEL_RUNNER_BASE_URL is set in the environment

	Zephyr, err1 := agents.GetZephyrAgent("ai/qwen2.5:3B-F16")
	// IMPORTANT: the model must support the tools

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

	// BEGIN: MCP client initialization
	fmt.Println("🚀 Initializing MCP StreamableHTTP client...")
	// Create HTTP transport
	httpURL := os.Getenv("MCP_HTTP_SERVER_URL")
	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		log.Fatalf("😡 Failed to create HTTP transport: %v", err)
	}
	// Create client with the transport
	mcpClient := client.NewClient(httpTransport)
	// Start the client
	if err := mcpClient.Start(ctx); err != nil {
		log.Fatalf("😡 Failed to start client: %v", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Simple Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("😡 Failed to initialize: %v", err)
	}
	// END:  of MCP client initialization

	// TOOLS:
	var openAITools []openai.ChatCompletionToolParam
	toolsRequest := mcp.ListToolsRequest{}
	toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatalf("😡 Failed to list tools: %v", err)
	}
	// STEP 1: Convert the MCP tools to OpenAI tools
	openAITools = helpers.ConvertMCPToolsToOpenAITools(toolsResult)
	// STEP 2: Register the tools to the tools master Zepyhr agent
	toolsMaster.Params.Tools = openAITools

	// RESOURCES:

	// STEP 3: create a resource request
	rsrcReq := mcp.ReadResourceRequest{}
	rsrcReq.Params.URI = "dungeon://rooms"
	// STEP 4: read the resource from the MCP server
	rsrcResp, _ := mcpClient.ReadResource(ctx, rsrcReq)
	rsrcContent := rsrcResp.Contents[0].(mcp.TextResourceContents).Text

	// STEP 5: add the resource to the Lyralei agent
	Lyralei.Params.Messages = append(Lyralei.Params.Messages,
		openai.SystemMessage("ROOMS:\n"+rsrcContent+"\n"),
		openai.SystemMessage(`
			Si l'utilisateur te demande de lui décrire une pièce, 
			tu dois lui répondre en utilisant les ressources disponibles dans ROOMS.
		`),
	)
	// connais tu La Bibliothèque des Murmures
	// parle moi du Sanctuaire des Marées Éternelles
	// Quelles sont les pièces du donjon ?

	// PROMPTS:

	promptReq := mcp.GetPromptRequest{}
	promptReq.Params.Name = "roll_dices_prompt"
	promptReq.Params.Arguments = map[string]string{
		"numFaces": "6",
		"numDices": "3",
	}
	promptRes, err := mcpClient.GetPrompt(ctx, promptReq)
	if err != nil {
		log.Fatalf("😡 Failed to get prompt: %v", err)
	}
	promptMessages := promptRes.Messages[0]
	promptText := promptMessages.Content.(mcp.TextContent).Text

	fmt.Println("📝 Prompt response:", promptText)

	if serverInfo.Capabilities.Prompts != nil {
		fmt.Println("Fetching available prompts...")
		promptsRequest := mcp.ListPromptsRequest{}
		promptsResult, err := mcpClient.ListPrompts(ctx, promptsRequest)
		if err != nil {
			log.Printf("Failed to list prompts: %v", err)
		} else {
			fmt.Printf("Server has %d prompts available\n", len(promptsResult.Prompts))
			for i, prompt := range promptsResult.Prompts {
				fmt.Printf("  %d. %s - %s\n", i+1, prompt.Name, prompt.Description)
			}
		}
	}

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

				ui.Println(ui.Green, "💡 tool arguments:", args)

				switch toolCall.Function.Name {
				// TOOL 1:
				case "choisir_un_personnage_par_son_espece":
					//selectedCharacter = characters[tools.ChooseCharacterBySpecies(args)]

					// NOTE: Call the MCP tool with the arguments
					request := mcp.CallToolRequest{}
					request.Params.Name = toolCall.Function.Name
					request.Params.Arguments = args

					toolResponse, _ := mcpClient.CallTool(ctx, request)
					characterName := toolResponse.Content[0].(mcp.TextContent).Text

					selectedCharacter = characters[characterName]

				// TOOL 2:
				case "detecter_le_vrai_sujet_du_message_utilisateur":
					//selectedCharacter = characters[tools.ChooseCharacterFromTopic(args)]

					// NOTE: Call the MCP tool with the arguments
					request := mcp.CallToolRequest{}
					request.Params.Name = toolCall.Function.Name
					request.Params.Arguments = args

					toolResponse, _ := mcpClient.CallTool(ctx, request)
					characterName := toolResponse.Content[0].(mcp.TextContent).Text

					selectedCharacter = characters[characterName]
				// TOOL 3:
				case "lancer_des_des":
					// NOTE: Call the MCP tool with the arguments
					request := mcp.CallToolRequest{}
					request.Params.Name = toolCall.Function.Name
					request.Params.Arguments = args

					toolResponse, _ := mcpClient.CallTool(ctx, request)
					// Display the result of the dice roll
					ui.Println(ui.Green, "🎲 Dice roll result:", toolResponse.Content[0].(mcp.TextContent).Text)

					question = toolResponse.Content[0].(mcp.TextContent).Text

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

		/*
			if selectedCharacter.Name == "Aldric" {
				selectedCharacter.Params.Messages = append(selectedCharacter.Params.Messages,
					openai.UserMessage(promptText), // Use the prompt message content
				)

			} else {
				selectedCharacter.Params.Messages = append(selectedCharacter.Params.Messages,
					openai.UserMessage(question),
				)

			}
		*/

		// STEP 4: display who is speaking...
		ui.Println(ui.Magenta, "[", selectedCharacter.Avatar, "]", selectedCharacter.Name, "is speaking...")

		// STEP 5: display the answer / run the chat completion
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

	fmt.Println("Client initialized successfully. Shutting down...")
	mcpClient.Close()

}
