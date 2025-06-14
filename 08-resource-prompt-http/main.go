package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-dd-http",
		"0.0.0",
	)

	// =================================================
	// TOOLS:
	// =================================================
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

	rollDices := mcp.NewTool("lancer_des_des",
		mcp.WithDescription("Lancez des dés pour obtenir un résultat aléatoire."),
		mcp.WithNumber("nb_dices",
			mcp.Required(),
			mcp.Description("Le nombre de dés à lancer. Par défaut, 1 dé est lancé."),
		),
		mcp.WithNumber("sides",
			mcp.Required(),
			mcp.Description("Le nombre de faces du dé. Par défaut, un dé à 6 faces est lancé."),
		),
	)

	s.AddTool(rollDices, rollDicesHandler)

	// =================================================
	// RESOURCES:
	// =================================================

	// Static resource example - exposing a README file
	resource := mcp.NewResource(
		"dungeon://rooms",
		"dungeon project",
		mcp.WithResourceDescription("Therooms of the dungeon project"),
		mcp.WithMIMEType("text/markdown"),
	)

	// Add resource with its handler
	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {

		content, err := os.ReadFile("dungeon.md")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "dungeon://rooms",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})

	// =================================================
	// PROMPTS:
	// =================================================


	prompt := mcp.NewPrompt(
		"roll_dices_prompt",
		mcp.WithPromptDescription("A roll dices prompt example"),
		mcp.WithArgument("numFaces",
			mcp.ArgumentDescription("nombre de faces du dé"),
		),
		mcp.WithArgument("numDices",
			mcp.ArgumentDescription("nombre de dés à lancer"),
		),
	)

	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		numFaces := request.Params.Arguments["numFaces"]
		// Default to 6 faces if not provided
		if numFaces == "" {
			numFaces = "6"
		}
		numDices := request.Params.Arguments["numDices"]
		// Default to 1 dice if not provided
		if numDices == "" {
			numDices = "1"
		}

		// Create the prompt content
		return mcp.NewGetPromptResult(
			"Roll Dices Prompt",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf("Lance un jet de %s dé(s) avec %s faces", numDices, numFaces)),
				),
			},
		), nil
		
	})


	// Start the stdio server
	// if err := server.ServeStdio(s); err != nil {
	// 	log.Fatalln("Failed to start server:", err)
	// 	return
	// }

	// Start the HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	).Start(":" + httpPort)
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

func rollDicesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nbDices := request.GetInt("nb_dices", 1)
	sides := request.GetInt("sides", 6)

	roll := func(n, x int) int {
		if n <= 0 || x <= 0 {
			return 0
		}

		results := make([]int, n)
		sum := 0

		for i := range n {
			roll := rand.Intn(x) + 1 // +1 because rand.Intn(x) donne 0 à x-1
			results[i] = roll
			sum += roll
		}

		return sum
	}

	// Simulate rolling dice
	result := roll(nbDices, sides)

	return mcp.NewToolResultText("Le résultat du lancer de " + strconv.Itoa(nbDices) + " dés à " + strconv.Itoa(sides) + " faces est: " + strconv.Itoa(result)), nil

}
