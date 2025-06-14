package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

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

//go:export tools_information
func ToolsInformation() {

	orcGreetings := Tool{
		Name:        "orc_greetings",
		Description: "make greetings as an Orc",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"name"},
			Properties: map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the person to greet",
				},
			},
		},
	}

	vulcanGreetings := Tool{
		Name:        "vulcan_greetings",
		Description: "make greetings as Vulcan",
		InputSchema: InputSchema{
			Type:     "object",
			Required: []string{"name"},
			Properties: map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the person to greet",
				},
			},
		},
	}


	tools := []Tool{orcGreetings, vulcanGreetings}

	jsonData, _ := json.Marshal(tools)
	pdk.OutputString(string(jsonData))
}

//go:export orc_greetings
func OrcGreetings() {
	type Arguments struct {
		Name string `json:"name"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	pdk.OutputString("Throm-ka " + args.Name)
}

//go:export vulcan_greetings
func VulcanGreetings() {
	type Arguments struct {
		Name string `json:"name"`
	}
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	pdk.OutputString("🖖 Peace and Prosper " + args.Name)
}

func main() {}
