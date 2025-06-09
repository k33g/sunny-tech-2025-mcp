package main

import (
	"context"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Config struct {
	BaseURL        string
	EmbeddingModel string
	MaxResults     string
}

var config Config

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-dd",
		"0.0.0",
	)
	// NOTE: look at the code of getZephyrToolsCatalog()
	// Add a tool
	chooseCharacterBySpecies := mcp.NewTool("choisir_un_personnage_par_son_espece",
		mcp.WithDescription(`sélectionner une espèce parmi celles-ci: [Humain, Orc, Elfe, Nain] en disant: je veux parler à un(e) <species_name>.`),
		mcp.WithString("species_name",
			mcp.Required(),
			mcp.Description("L'espèce à détecter dans le message utilisateur. L'espèce peut être une des suivantes: [Humain, Orc, Elfe, Nain]."),
		),
	)
	s.AddTool(chooseCharacterBySpecies, chooseCharacterBySpeciesHandler)

	detectTheRealTopicInUserMessage := mcp.NewTool("detecter_le_vrai_sujet_du_message_utilisateur",
		mcp.WithDescription(`sélectionner un sujet parmi ceux-ci: [justice, guerre, combat, magie, poésie, artisanat, forge] en disant: j'ai une question sur <topic_name>.`),
		mcp.WithString("topic_name",
			mcp.Required(),
			mcp.Description("Le sujet à détecter dans le message utilisateur. Le sujet peut être un des suivant: [justice, guerre, combat, magie, poésie, artisanat, forge]."),
		),
	)
	s.AddTool(detectTheRealTopicInUserMessage, detectTheRealTopicInUserMessageHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalln("Failed to start server:", err)
		// TODO: handle error more gracefully
		return
	}

}

func chooseCharacterBySpeciesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	// Check if the species_name argument is provided
	if len(args) == 0 {
		return mcp.NewToolResultText("zephyr"), nil
	}
	var content = "zephyr" // default character
	if speciesName, ok := args["species_name"]; ok {

		switch strings.ToLower(speciesName.(string)) {
		case "humain":
			content = "aldric"
		case "orc":
			content = "grash"
		case "nain":
			content = "thorin"
		case "elfe", "elf":
			content = "lyralei"
		default:
			content = "zephyr"
		}
	}
	return mcp.NewToolResultText(content), nil

}

func detectTheRealTopicInUserMessageHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	// Check if the topic_name argument is provided
	if len(args) == 0 {
		return mcp.NewToolResultText("zephyr"), nil
	}
	var content = "zephyr" //default character
	if topicName, ok := args["topic_name"]; ok {
		switch strings.ToLower(topicName.(string)) {
		case "justice":
			content = "aldric"
		case "guerre", "combat":
			content = "grash"
		case "magie", "poésie", "poesie":
			content = "lyralei"
		case "artisanat", "forge":
			content = "thorin"
		default:
			content = "zephyr"
		}
	}
	return mcp.NewToolResultText(content), nil
}
