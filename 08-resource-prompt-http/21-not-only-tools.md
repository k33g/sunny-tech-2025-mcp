# Ressources, Prompts et Streamable HTTP

## Ressources & Prompts

<!-- TODO: 🚧 -->

## Avantages de Streamable HTTP:

- Serveur MCP remote
- Accessible par plusieurs Applications d'IA Générative (et Agents)
- Partage d'informations
- Authentification, autorisations, 
- ...

## Demo

```bash
docker build --platform linux/arm64 -t mcp-dd-http:demo .
docker run --rm -p 9090:9090 mcp-dd-http:demo
```

- Initialize: `initialize.sh`
- Resources: `resources.list.sh`, `resources.read.sh` 
- Prompts: `prompts.sh`
- Tools: `tools.list.sh`, `use.tool.roll.dice.sh`

