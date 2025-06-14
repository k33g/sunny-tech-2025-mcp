# Application d'IA Générative

## Pattern 2: AI Agent(s) - "Micro" Agent pattern
> "Micro" Agent pattern >> LLM Locaux

```mermaid
flowchart LR
    Human((Human)) <-.-> LLM[GenAI App + LLM + CALL]
    LLM -->|Action| Environment((Environment))
    Environment -->|Feedback| LLM
    LLM -.-> Stop[Stop]
    
    classDef human fill:#FFDDDD,stroke:#DD8888
    classDef llm fill:#DDFFDD,stroke:#88DD88
    classDef stop fill:#DDDDFF,stroke:#8888DD
    classDef env fill:#FFDDDD,stroke:#DD8888
    
    class Human human
    class LLM llm
    class Stop stop
    class Environment env
```

## C'est le `"CALL"` qui est important
> ce qui va nous amener à MCP

___
[◀️ Previous](./04-augmented-llm.md#application-dia-générative) | [Prompt ▶️](./06-prompt.md#le-prompt-le-nerf-de-la-guerre)



