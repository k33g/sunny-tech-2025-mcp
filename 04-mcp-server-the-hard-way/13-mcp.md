# 🛠️ Model Context Protocol
> le function calling a un niveau supérieur

## MCP (par Anthropic) ?

```mermaid
flowchart TD
  subgraph Host["MCP Host (AI Application)"]
      App["Application LLM<br>(e.g., Claude Desktop)"]
  end

  subgraph Client["MCP Client"]
      Protocol["Protocol Client<br>Transport"]
  end

  subgraph Servers["MCP Servers"]
      Server1["MCP Server 1<br>(<b>👋 Local</b> Data)"]
      Server2["MCP Server 2<br>(Remote Data)"]
      Server3["MCPServer 3<br>(Remote Data)"]
  end

  App --> Protocol
  Protocol <--> |STDIO<br>JSON-RPC 2.0| Server1
  Protocol <--> |SSE<br>JSON-RPC 2.0| Server2
  Protocol <--> |Streamable HTTP<br>JSON-RPC 2.0| Server3

style Server1 fill:#FFDDDD,stroke:#DD8888
style App fill:#DDFFDD,stroke:#88DD88
style Server2 fill:#DDDDFF,stroke:#8888DD
style Server3 fill:#7ECF7E,stroke:#195919
```

### Pour le moment le plus **utilisé** est le serveur **`MCP STDIO`** (et ce n'est pas vraiment un serveur)
### **`MCP SSE`** est déprécié
### **`MCP Streamable HTTP`** est arrivé
> on en parlera un peu plus tard


___
[◀️ Previous](./12-we-have-a-problem.md#ok-on-sait-faire-exécuter-du-code-à-un-llm-) | [🎉 MCP ▶️](./14-mcp.md#️-model-context-protocol)