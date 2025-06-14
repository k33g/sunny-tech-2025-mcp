# 🛠️ Model Context Protocol

## Comment ça marche ?

```mermaid
sequenceDiagram
    participant MCP Server
    participant Host App
    participant LLM

    MCP Server->>Host App: Expose tools (add_numbers, subtract_numbers...)
    Host App->>Host App: Format tools for LLM
    Host App->>LLM: Send prompt + formatted tools
    Note over LLM: Process request
    LLM->>Host App: Return tool_calls JSON
    Host App->>MCP Server: Convert and send tool request
    MCP Server->>Host App: Return operation result
```
> Le serveur MCP ne se contente pas d'exposer des tools (on en parlera plus tard)

___
[◀️ Previous](./13-mcp.md#️-model-context-protocol) | [📺 STDIO ▶️](./15-stdio.md#fonctionnement-dun-programme-utilisant-stdio)


<!-- TODO: explain

-->