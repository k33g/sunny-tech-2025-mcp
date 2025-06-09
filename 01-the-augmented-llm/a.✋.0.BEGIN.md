# Application d'IA Générative
## Augmented LLM

```mermaid
flowchart LR
    In((In)) --> LLM[LLM]
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

## AI Agents - "Micro" Agent pattern


```mermaid
flowchart LR
    Human((Human)) <-.-> LLM[LLM Call]
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
> donc ce à quoi l'on souhaite arriver.


## What is a prompt?


## Docker Model Runner

TODO: explain (mac, linux, windows)


## The Heroic Agents
> qui ne sont pas encore des agents...

Imaginons je veux écrire un jeu de rôle


