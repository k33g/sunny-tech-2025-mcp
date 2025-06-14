# Création d'un agent

## Structure

```mermaid
classDiagram
    class TinyAgent {
        -ctx context.Context
        -client openai.Client
        +Params openai.ChatCompletionNewParams
        +Name string
        +Avatar string
        +Color string
        +Instructions openai.ChatCompletionMessageParamUnion
        +ChatCompletion() (string, error)
        +ChatCompletionStream(callBack func(*TinyAgent, string, error) error) (string, error)
        +ToolsCompletion() ([]openai.ChatCompletionMessageToolCall, error)
    }

 
    class NewAgent {
        <<function>>
        +NewAgent(name string) (*TinyAgent, error)
    }


    NewAgent --> TinyAgent : creates

    note for TinyAgent "Main agent struct that wraps OpenAI client"
    
    note for NewAgent "Constructor function"
```

## Avec :

### ⦿ OpenAI Golang SDK
### ⦿ 🐳 Docker Model Runner

## 👋 J'ai un bout de code avec 2 agents (futurs agents)...
## ... Alors en fait, ce sont des personnages de D&D 🧝‍♀️
___
[◀️ Previous](./06-prompt.md#le-prompt-le-nerf-de-la-guerre) | [📝 some code ▶️](./main.go)

