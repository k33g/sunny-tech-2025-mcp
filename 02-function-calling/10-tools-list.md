# Catalogue de tools

#### Pour pouvoir **`identifier`** les fonctions à appeler et **`fournir les paramètres nécessaires`**, le LLM a besoin d'une liste d'outils:

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
> ✋ La description de l'outil est super importante (ainsi que les paramètres)
___
[◀️ Previous](./09-function-calling.md#function-calling) | [📝 some code ▶️](./main.go)

