package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/extism/go-pdk"
)

// Available races
const (
	Human = "human"
	Elf   = "elf"
	Dwarf = "dwarf"
	Orc   = "orc"
)

// Prefixes and suffixes by race
var transformations = map[string]map[string][]string{
	Human: {
		"prefixes": {"Sir", "Lady", "Lord", "Dame", "Baron", "Count"},
		"suffixes": {"the Brave", "the Noble", "of Valenhall", "the Just", "of Ironhold", "the Wise"},
	},
	Elf: {
		"prefixes": {"Aer", "Sil", "Gal", "Fin", "Elen", "Thal", "Cel"},
		"suffixes": {"wen", "dil", "rion", "ithil", "las", "dir", "thel", "andreth", "ion"},
	},
	Dwarf: {
		"prefixes": {"Thorin", "Balin", "Dain", "Nain", "Borin", "Gloin", "Oin"},
		"suffixes": {"beard", "forge", "hammer", "stone", "iron", "axe", "mountain", "shield"},
	},
	Orc: {
		"prefixes": {"Grok", "Uruk", "Gash", "Skar", "Thok", "Grish", "Ugluk"},
		"suffixes": {"grond", "dum", "hai", "gul", "agh", "burz", "tark", "snaga"},
	},
}

// TransformNameDnD transforms a first name according to the chosen D&D race
func TransformNameDnD(firstName, race string) (string, error) {
	// Normalize race to lowercase
	race = strings.ToLower(strings.TrimSpace(race))
	firstName = strings.TrimSpace(firstName)

	if firstName == "" {
		return "", fmt.Errorf("first name cannot be empty")
	}

	// Check if race exists
	transformation, exists := transformations[race]
	if !exists {
		return "", fmt.Errorf("unsupported race: %s. Available races: human, elf, dwarf, orc", race)
	}

	// Initialize random generator
	rand.Seed(time.Now().UnixNano())

	var newName string

	switch race {
	case Human:
		// For humans, add a title and noble family name
		prefixes := transformation["prefixes"]
		suffixes := transformation["suffixes"]

		title := prefixes[rand.Intn(len(prefixes))]
		familyName := suffixes[rand.Intn(len(suffixes))]

		newName = fmt.Sprintf("%s %s %s", title, firstName, familyName)

	case Elf:
		// For elves, modify the name with elvish sounds
		prefixes := transformation["prefixes"]
		suffixes := transformation["suffixes"]

		// Combine the first name with elvish elements
		if len(firstName) > 3 {
			base := firstName[:len(firstName)/2]
			suffix := suffixes[rand.Intn(len(suffixes))]
			newName = base + suffix
		} else {
			prefix := prefixes[rand.Intn(len(prefixes))]
			suffix := suffixes[rand.Intn(len(suffixes))]
			newName = prefix + suffix
		}

	case Dwarf:
		// For dwarves, add a clan name related to craftsmanship
		suffixes := transformation["suffixes"]
		suffix := suffixes[rand.Intn(len(suffixes))]

		newName = fmt.Sprintf("%s %s%s", firstName, strings.Title(firstName), strings.Title(suffix))

	case Orc:
		// For orcs, create a brutal and guttural name
		prefixes := transformation["prefixes"]
		suffixes := transformation["suffixes"]

		// Replace part of the first name with orcish sounds
		if len(firstName) > 2 {
			base := strings.ToUpper(firstName[:2])
			suffix := suffixes[rand.Intn(len(suffixes))]
			newName = base + suffix
		} else {
			prefix := prefixes[rand.Intn(len(prefixes))]
			suffix := suffixes[rand.Intn(len(suffixes))]
			newName = prefix + suffix
		}
		newName = strings.ToUpper(newName)
	}

	return newName, nil
}

// Utility function to list available races
func AvailableRaces() []string {
	races := make([]string, 0, len(transformations))
	for race := range transformations {
		races = append(races, race)
	}
	return races
}


// -------------------------------------------------
//  Tools
// -------------------------------------------------

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	//InputSchema map[string]any `json:"inputSchema"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string         `json:"type"`
	Required   []string       `json:"required"`
	Properties map[string]any `json:"properties"`
}

// Transform a name in a heroic fantasy style
func TransformName(name string) string {
	if name == "" {
		return "Hero"
	}
	return name + " the Brave"
}

//go:export tools_information
func ToolsInformation() {

	generateName := Tool{
		Name:        "generate_name",
		Description: "transforms a first name into D&D style based on the chosen race",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"name"},
			Properties: map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the person to transform",
				},
				"race": map[string]any{
					"type":        "string",
					"description": "the chose race for the transformation (human, elf, dwarf, orc)",
				},
			},
		},
	}



	tools := []Tool{generateName}

	jsonData, _ := json.Marshal(tools)
	pdk.OutputString(string(jsonData))
}

//go:export generate_name
func GenerqteName() {
	type Arguments struct {
		Name string `json:"name"`
		Race string `json:"race"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	newName, err := TransformNameDnD(args.Name, args.Race)
	if err != nil {
		pdk.OutputString("Error: " + err.Error())
		return
	}

	pdk.OutputString(newName)
}



func main() {}



