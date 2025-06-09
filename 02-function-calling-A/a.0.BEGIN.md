# Function Calling

## Function calling is the ability for AI models to **`identify`** when a function **`should be executed`** and provide the **`necessary parameters`** for that execution


# ✋ les "petits" LLMs ne sont pas très bons pour le function calling


Pour faire du function calling, il faut que le modèle puisse **`identifier`** les fonctions à appeler et **`fournir les paramètres nécessaires`**.

Il faut donc lui fournir un catalogue de fonctions/tools.

> exemple de fonction/tool : "vulcan salute"
```golang
vulcanSaluteTool := openai.ChatCompletionToolParam{
    Function: openai.FunctionDefinitionParam{
        Name:        "vulcan_salute",
        Description: openai.String("Give a vulcan salute to the given person name"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]interface{}{
                "name": map[string]string{
                    "type": "string",
                },
            },
            "required": []string{"name"},
        },
    },
}
```