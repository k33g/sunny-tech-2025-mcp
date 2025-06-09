package tools

import "strings"

func ChooseCharacterBySpecies(args map[string]string) string {
	if speciesName, ok := args["species_name"]; ok {
		switch strings.ToLower(speciesName) {
		case "humain":
			return "aldric"
		case "orc":
			return "grash"
		case "nain":
			return "thorin"
		case "elfe", "elf":
			return "lyralei"

		default:
			return "zephyr"
		}
	} else {
		return "zephyr"
	}

}

func ChooseCharacterFromTopic(args map[string]string) string {
	if topicName, ok := args["topic_name"]; ok {
		switch strings.ToLower(topicName) {
		case "justice":
			return "aldric"
		case "guerre", "combat":
			return "grash"
		case "magie", "poésie", "poesie":
			return "lyralei"
		case "artisanat", "forge":
			return "thorin"
		default:
			return "zephyr"
		}
	} else {
		return "zephyr"
	}
}
