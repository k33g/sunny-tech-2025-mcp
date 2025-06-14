# Application d'IA Générative
## Pattern 1: LLM augmenté
> Augmented LLM

```mermaid
flowchart LR
    In((In)) --> LLM[GenAI app + LLM]
    LLM --> Out((Out))
    
    LLM -.-> |Query/Results| Retrieval[Retrieval]
    Retrieval -.-> LLM
    
    LLM -.-> |Call/Response| Tools[Tools]
    Tools -.-> LLM
    
    LLM -.-> |Read/Write| Memory[Memory]
    Memory -.-> LLM
    
    style In fill:#FFDDDD,stroke:#DD8888
    style Out fill:#FFDDDD,stroke:#DD8888
    style LLM fill:#DDFFDD,stroke:#88DD88
    style Retrieval fill:#DDDDFF,stroke:#8888DD
    style Tools fill:#DDDDFF,stroke:#8888DD
    style Memory fill:#DDDDFF,stroke:#8888DD
```

___
[◀️ Previous](../00-intro-agenda/03-btw.md#-vous-pouvez-poser-des-questions-pendant-la-présentation) | [AI Agent ▶️](./05-ai-agent.md#application-dia-générative)
