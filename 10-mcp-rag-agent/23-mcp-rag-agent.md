# Un usage pratique: MCP + RAG
> Utiliser un modèle dans un serveur MCP

## RAG: retrieval augmented generation

- Découper des données en morceaux
- Sélectionner uniquement les morceaux qui correspondent à une question (proches de la question)
- Fournir les morceaux avec le prompt
- Laisser le LLM se débrouiller avec tout ça

### Embedding ? 🤨

- Un embedding = transformer des données (ex: texte) en vecteurs de nombres. 
- Donner des "coordonnées" à des données pour les placer dans un espace mathématique où :
  - Les éléments **similaires** sont **proches**
  - Les éléments **différents** sont **éloignés**
- **==> Recherche de similarités**

## Serveur MCP (StreamableHTTP)

- Un serveur MCP qui utilise un modèle d'embeddings pour calculer les vecteurs des morceaux de textes d'un document (ici les règles du jeu d'un JDR)
- Expose un "tool": **`question_about_something`** pour de la recherche de similarité


> Démarrage:
```bash
cd mcp-rag-server
go run main.go
```

## Agent Zephyr (Google ADK Python) pour utiliser `question_about_something`

Zephyr va utiliser le serveur MCP pour chercher des informations à propos des règles du jeu.

> Démarrage
```bash
cd agents
adk web
```

### Questions

- trouve dans ta base la liste des monstres
- trouve dans ta base qui est keegorg


