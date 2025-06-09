package agents

import (
	"context"
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type TinyAgent struct {
	ctx    context.Context
	client openai.Client
	Params openai.ChatCompletionNewParams
	Name   string
	Avatar string
	Color string // used for UI display
	Instructions openai.ChatCompletionMessageParamUnion
}

func NewAgent(name string) (*TinyAgent, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("failed to load .env file: " + err.Error())
	}
	// Ensure MODEL_RUNNER_BASE_URL is set in the environment
	if os.Getenv("MODEL_RUNNER_BASE_URL") == "" {
		return nil, errors.New("MODEL_RUNNER_BASE_URL environment variable is not set")
	}

	llmURL := os.Getenv("MODEL_RUNNER_BASE_URL") + "/engines/llama.cpp/v1/"

	agent := &TinyAgent{}
	agent.Name = name
	agent.ctx = context.Background()

	agent.client = openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)

	return agent, nil
}

// ChatCompletion handles the chat completion request using the DMR client.
// It sends the parameters set in the Agent and returns the response content or an error.
// It is a synchronous operation that waits for the completion to finish.
func (agent *TinyAgent) ChatCompletion() (string, error) {
	completion, err := agent.client.Chat.Completions.New(agent.ctx, agent.Params)

	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		return completion.Choices[0].Message.Content, nil
	} else {
		return "", errors.New("no choices found")

	}
}

// ChatCompletionStream handles the chat completion request using the DMR client in a streaming manner.
// It takes a callback function that is called for each chunk of content received.
// The callback function receives the Agent instance, the content of the chunk, and any error that occurred.
// It returns the accumulated response content and any error that occurred during the streaming process.
// The callback function should return an error if it wants to stop the streaming process.
func (agent *TinyAgent) ChatCompletionStream(callBack func(self *TinyAgent, content string, err error) error) (string, error) {
	response := ""
	stream := agent.client.Chat.Completions.NewStreaming(agent.ctx, agent.Params)
	var cbkRes error

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			cbkRes = callBack(agent, chunk.Choices[0].Delta.Content, nil)
			response += chunk.Choices[0].Delta.Content
		}

		if cbkRes != nil {
			break
		}
	}
	if cbkRes != nil {
		return response, cbkRes
	}
	if err := stream.Err(); err != nil {
		return response, err
	}
	if err := stream.Close(); err != nil {
		return response, err
	}

	return response, nil
}

// ToolsCompletion handles the tool calls completion request using the DMR client.
// It sends the parameters set in the Agent and returns the detected tool calls or an error.
// It is a synchronous operation that waits for the completion to finish.
func (agent *TinyAgent) ToolsCompletion() ([]openai.ChatCompletionMessageToolCall, error) {

	completion, err := agent.client.Chat.Completions.New(agent.ctx, agent.Params)
	if err != nil {
		return nil, err
	}
	detectedToolCalls := completion.Choices[0].Message.ToolCalls
	if len(detectedToolCalls) == 0 {
		return nil, errors.New("no tool calls detected")
	}
	return detectedToolCalls, nil
}
