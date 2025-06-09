package helpers

import (
	"encoding/json"

	"github.com/openai/openai-go"
)

// ToolCallsToJSONString converts a slice of openai.ChatCompletionMessageToolCall to a JSON string.
// It extracts the tool call ID and function arguments, converting them to a generic interface
// for JSON marshaling. The resulting JSON string is formatted with indentation for readability.
// If the tool calls are empty, it returns an empty JSON array.
// If any error occurs during the conversion, it returns an error.
// The function is useful for logging or storing tool calls in a structured format.
// It returns a JSON string representation of the tool calls.
func ToolCallsToJSONString(tools []openai.ChatCompletionMessageToolCall) (string, error) {
	var jsonData []any

	// Convert tools to generic interface
	for _, tool := range tools {
		var args any
		if err := json.Unmarshal([]byte(tool.Function.Arguments), &args); err != nil {
			return "", err
		}

		jsonData = append(jsonData, map[string]any{
			"id": tool.ID,
			"function": map[string]any{
				"name":      tool.Function.Name,
				"arguments": args,
			},
		})
	}

	// Marshal back to JSON with indentation
	jsonString, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}


func ExtractArgsFromJSONString(jsonString string) (map[string]string) {
	var args map[string]string
	err := json.Unmarshal([]byte(jsonString), &args)
	if err != nil {
		return nil
	}
	return args
}