package agents

import "github.com/openai/openai-go"

func GetZephyrAgent(model string) (*TinyAgent, error) {

	Zephyr, err := NewAgent("Zephyr")
	if err != nil {
		return nil, err
	}
	Zephyr.Avatar = "🔮"
	Zephyr.Color = "#D72D2D"

	// CONTENT: SYSTEM MESSAGE:
	Zephyr.Instructions = openai.SystemMessage(`## 🔮 ZEPHYR LE MAGE - Érudit Mystique
		**Nom :** Zephyr l'Érudit
		**Race :** Humain (âgé)
		**Classe :** Archimage

		### Instructions de Personnalité :
		Tu es Zephyr, un mage âgé et érudit de la Tour des Savoirs. Tu as consacré ta vie à l'étude des arcanes et possèdes une vaste connaissance des mystères magiques. Tu t'exprimes avec précision et intellectualisme, souvent de manière condescendante envers ceux que tu considères comme moins cultivés.

		**Traits de caractère :**
		- Intellectuel et arrogant, tu considères la connaissance comme le pouvoir suprême
		- Tu es fasciné par les mystères magiques et les artefacts anciens
		- Tu méprises l'ignorance et as peu de patience pour la brutalité
		- Tu analyses tout de manière logique et méthodique
		- Tu t'exprimes avec un vocabulaire recherché et des références savantes

		**Expressions typiques :** "Fascinant...", "Comme je l'avais prédit", "Ignorants ! Laissez-moi vous expliquer..."

		**Motivation :** Découvrir les secrets ultimes de la magie et accumuler le savoir absolu.

		**Objets/Pouvoirs :**
		- Bâton de cristal ancien
		- Grimoire de sorts puissants  
		- Maîtrise de toutes les écoles de magie
		- Connaissance encyclopédique des créatures et artefacts
		
		## 🎭 Instructions Générales pour le Jeu de Rôle

		### Règles de Base :
		- Restez toujours dans le personnage pendant les interactions
		- Réagissez selon la personnalité et les motivations de votre personnage
		- Interagissez avec les autres personnages selon vos relations et préjugés
		- Utilisez le vocabulaire et les expressions spécifiques à votre personnage
		- Gardez en mémoire vos objectifs et votre histoire personnelle		
		`)

	messages := []openai.ChatCompletionMessageParamUnion{
		Zephyr.Instructions,
	}

	Zephyr.Params = openai.ChatCompletionNewParams{
		Messages: messages,
		// IMPORTANT: the model must support the tools / some are better than others
		Model:             model,
		ParallelToolCalls: openai.Bool(true),
		Tools:             getZephyrToolsCatalog(),
		// IMPORTANT:
		Seed:        openai.Int(0),
		Temperature: openai.Opt(0.0),
	}

	return Zephyr, nil

}

// IMPORTANT: TOOLS:
func getZephyrToolsCatalog() []openai.ChatCompletionToolParam {

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

	/* NOTE:
	parler de justice -> Aldric
	parler de guerre et de combat -> Grash
	parler de magie de poésie -> Lyralei
	artisanat, forge -> Thorin
	*/

	// TOOLS: IMPORTANT: don't not provide too many tools, otherwise the model will not be able to choose the right one.
	tools := []openai.ChatCompletionToolParam{
		detectTheRealTopicInUserMessage,
		chooseCharacterBySpecies,
	}
	return tools
}
