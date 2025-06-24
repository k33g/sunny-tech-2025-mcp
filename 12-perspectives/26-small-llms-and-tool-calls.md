# Et si on améliorait la détection des "tool calls" avec des plus petits modèles ?

## C'était mieux avant

```mermaid
flowchart TD
    A[📦 Catalogue d'outils] --> B[🔄 Marshal vers JSON]
    B --> C["📝 Création toolsContent<br/>[AVAILABLE_TOOLS]....[/AVAILABLE_TOOLS]"]
    
    C --> D[📋 Instructions système]
    D --> E["💬 Message système complet:<br/>Introduction + toolsContent + Instructions"]
    
    E --> F[🤖 1ère Complétion API]
    F --> G[📄 Résultat: JSON brut des tool calls]
    
    G --> H["🔧 Construction nouveaux messages:<br/>System: 'Return all function calls...'<br/>User: résultat précédent"]
    
    H --> I[⚙️ Configuration format JSON forcé]
    I --> J[🤖 2ème Complétion API]
    J --> K[📊 Résultat: JSON structuré avec 'function_calls']
    
    K --> L[🔍 Parse JSON vers struct FunctionCalls]
    L --> M[🔄 Boucle: Pour chaque command]
    
    M --> N[🔧 Marshal arguments vers JSON]
    N --> O[🆔 Génération UUID unique]
    O --> P[🏗️ Création ChatCompletionMessageToolCall]
    
    P --> Q{Fin de boucle?}
    Q -->|Non| M
    Q -->|Oui| R[✅ Retour: Liste des Tool Calls]
    
    subgraph "Étape 1: Préparation du catalogue"
        B
        C
    end
    
    subgraph "Étape 2: Instructions et 1ère completion"
        D
        E
        F
        G
    end
    
    subgraph "Étape 3: 2ème completion avec format JSON"
        H
        I
        J
        K
    end
    
    subgraph "Étape 4: Transformation en Tool Calls"
        L
        M
        N
        O
        P
    end
    

    
    style A fill:#e3f2fd
    style R fill:#c8e6c9
    style F fill:#fff3e0
    style J fill:#fff3e0
    style B fill:#f3e5f5
    style L fill:#f3e5f5
```



1. Exporter le catalogue d'outils dans une string JSON: 
    ```go
    toolsJson, err := json.Marshal(catalog)
    ```
2. Créer avec cette string JSON, une nouvelle string: 
    ```go
    toolsContent := "[AVAILABLE_TOOLS]" + string(toolsJson) + "[/AVAILABLE_TOOLS]"
    ```
3. Créer des instructions pour expliquer au modèle comment utiliser ce catalogue:
    ```go
    systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
    To call a tool, respond with a JSON object with the following structure: 
    [
        {
            "name": <name of the called tool>,
            "arguments": {
                <name of the argument>: <value of the argument>
            }
        },
    ]

    search the name of the tool in the list of tools with the Name field
    `
    ```
4. Faire une 1ère complétion avec ces instructions système et la question de l'utilisateur
5. Ensuite construire de nouvelles instructions à partir du résultat:
    ```go
    agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
        openai.SystemMessage("Return all function calls wrapped in a container object with a 'function_calls' key."),
        openai.UserMessage(result),
    }
    ```
6. Faire une 2ème complétion en forçant le format JSON:
    ```go
    agent.Params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
        OfJSONObject: &openai.ResponseFormatJSONObjectParam{
            Type: "json_object",
        },
    }

    completionNext, err := agent.clientEngine.Chat.Completions.New(agent.ctx, agent.Params)
    ```
7. Transformer chaque item du résultat en des messages de type "Tool Call"
    ```go
    for i, command := range commands.FunctionCalls {
        argumentsJson, err := json.Marshal(command.Arguments)
        //...

        toolCalls[i] = openai.ChatCompletionMessageToolCall{
            ID: uuid.New().String(), // Generate a unique ID for the tool call
            Function: openai.ChatCompletionMessageToolCallFunction{
                Name:      command.Name,
                Arguments: string(argumentsJson),
            },
            Type: "function",
        }
    }
    ```

## Et ça fonctionne 🎉