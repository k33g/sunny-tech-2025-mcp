# Application d'IA Générative

## AI Agents - "Micro" Agent pattern
> Micro, c'est de moi ...

```mermaid
flowchart LR
    Human((Human)) <-.-> LLM[GenAI App + LLM + Call]
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


