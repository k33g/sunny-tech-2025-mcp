package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {
	chooseCharacterBySpecies := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "choisir_un_personnage_par_son_espece",
			Description: openai.String(`sélectionner une espèce parmi celles-ci: [Humain, Orc, Elfe, Nain] en disant: je veux parler à un(e) <species_name>.`),
			// NOTE: the species list is more to give the model a hint about the species, it can be any string.
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"species_name": map[string]string{
						"type":        "string",
						"description": "L'espèce à détecter dans le message utilisateur. L'espèce peut être une des suivantes: [Humain, Orc, Elfe, Nain].",
					},
				},
				"required": []string{"species_name"},
			},
		},
	}

	detectTheRealTopicInUserMessage := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "detecter_le_vrai_sujet_du_message_utilisateur",
			Description: openai.String(`sélectionner un sujet parmi ceux-ci: [justice, guerre, combat, magie, poésie, artisanat, forge] en disant: j'ai une question sur <topic_name>.`),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"topic_name": map[string]string{
						"type":        "string",
						"description": "Le sujet à détecter dans le message utilisateur. Le sujet peut être un des suivant: [justice, guerre, combat, magie, poésie, artisanat, forge].",
					},
				},
				"required": []string{"message"},
			},
		},
	}

	/*
		Some small model are able to detect several tool calls in a single request.
		So let ParallelToolCalls set to openai.Bool(true), then the model will detect all the tool calls in a single request.

		Models that are able to detect several tool calls in a single request:
		- ignaciolopezluna020/llama-xlam:8B-Q4_K_M
		- k33g/llama-xlam-2:8b-fc-r-q2_k
		  - https://huggingface.co/Salesforce/Llama-xLAM-2-8b-fc-r-gguf
		  - Llama-xLAM-2-8B-fc-r-Q2_K.gguf

	*/
	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				//Model: "k33g/llama-xlam-2:8b-fc-r-q2_k",
				//Model: "ai/qwen2.5:latest",
				Model: "ai/qwen2.5:1.5B-F16",
				Temperature: openai.Opt(0.0), // IMPORTANT: set temperature to 0.0 to ensure the agent uses the tool
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						je souhaite parler avec un nain

						Je voudrais aborder le sujet de la justice

						je veux discuter avec une elfe

						j'ai une question sur la magie					
					`),
				},
			},
		),
		agents.WithTools([]openai.ChatCompletionToolParam{chooseCharacterBySpecies, detectTheRealTopicInUserMessage}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("🤖 Bob is ready to assist!", bob.Params.Tools)

	// Generate the tools detection completion
	detectedToolCalls, err := bob.AltenativeToolsCompletion() // TODO: test is with Ollama
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Number of Tool Calls:\n", len(detectedToolCalls))

	detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
	if err != nil {
		fmt.Println("Error converting tool calls to JSON string:", err)
		return
	}
	fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)


}
