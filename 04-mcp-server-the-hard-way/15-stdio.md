# Fonctionnement d’un programme utilisant stdio

> `name.txt`
```text
Bob Morane
```

> `my-serveur.sh`
```bash
#!/bin/bash
echo -n "👋 Hello "
cat -
```
- `echo "👋 Hello "` écrit "👋 Hello " sur **stdout**.
- `cat -` lit tout ce qui arrive sur **stdin** (par ex le contenu de `name.txt`) et l’affiche à la suite.

> utilisation
```bash
cat name.txt | ./my-server.sh 
```

## ✋ Cela se passe en local

___
[◀️ Previous](./14-mcp.md#️-model-context-protocol) | [MCP Serveur STDIO ▶️](./16-mcp-stdio.md#fonctionnement-dun-serveur-mcp-utilisant-stdio)
