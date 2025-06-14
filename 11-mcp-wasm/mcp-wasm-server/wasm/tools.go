package wasm

import (
	"encoding/json"
	"fmt"
	"os"
	"mcp-dd-wasm/tools"

	extism "github.com/extism/go-sdk"
)

var toolSet = make(map[string]tools.Tool)

func registerToolsOfThePlugin(pluginInst *extism.Plugin) {
	_, output, err := pluginInst.Call("tools_information", nil)
	if err != nil {
		// Handle error case
		fmt.Fprintf(os.Stderr, "Error calling tools_information: %v\n", err)
		return
	}

	var toolsList []tools.Tool
	err = json.Unmarshal(output, &toolsList)
	if err != nil {
		// Handle error case
		fmt.Fprintf(os.Stderr, "Error unmarshalling tools information: %v\n", err)
		return
	}

	// Register the tools with the server
	for _, tool := range toolsList {

		// Define the tool handler
		// This function will be called when the tool is invoked
		// It will receive the arguments as a map[string]any
		// and should return the result as a string
		tool.Handler = func(args map[string]any) (any, error) {
			// Call the plugin function with the tool name and args
			// Convert args to JSON
			argsJSON, err := json.Marshal(args)
			if err != nil {
				// Handle error case
				fmt.Fprintf(os.Stderr, "Error marshalling args: %v\n", err)
				return nil, err
			}

			_, output, err := pluginInst.Call(tool.Name, argsJSON)
			if err != nil {
				// Handle error case
				fmt.Fprintf(os.Stderr, "Error calling tool %s: %v\n", tool.Name, err)
				return nil, err
			}
			// Cast output to string
			result := string(output)
			return result, nil
		}
		// Register the tool with the server
		addTool(tool)
	}
}

func addTool(tool tools.Tool) {
	toolSet[tool.Name] = tool
}
