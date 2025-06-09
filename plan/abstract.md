Conference Sunnytech
##### Abstract

Le Model Context Protocol (MCP) d'Anthropic redéfinit l'interaction entre les LLMs et nos systèmes. D'un côté vous aurez un serveur MCP capable d'exécuter des commandes, programmes, fonctions, ... et de l'autre une application hôte utilisant à la fois un LLM et un client MCP qui donne l'opportunité au LLM d'envoyer des commandes au serveur MCP: "ajoute une issue au projet GitLab demo dans le groupe Zeira Corp avec le titre hello world", "dis bonjour à Bob", ...

Ce qu'il faut s'avoir, c'est que vous n'êtes pas obligés d'utiliser les modèles d'Anthropic et que vous pouvez faire ça avec des LLMs locaux (et comme ça on évite le soulèvement des machines et on consomme moins d'énergie).

Avec cette présentation, je vous expliquerai:

- MCP et ses grands principes
- Comment utiliser des serveurs MCP existants avec Docker (nous verrons comment demander à un LLM de créer des issues dans GitLab ou "s'amuser" avec SQLite)
- Comment faire votre propre serveur MCP et une application hôte pour l'utiliser avec Ollama et un LLM local (j'utiliserai le langage Go pour cela et nous ferons exécuter des fonctions WebAssembly à notre LLM)

L'objectif de cette présentation est aussi de démontrer comment nous pouvons conserver le contrôle de nos IA, utiliser des ressources locales et pas forcément avoir besoin de ressources importantes (par exemple, nous n'avons pas toujours besoin de GPU)

---
idees
la suite de Sarah Connor
expliquer les tools avec les LLMs
lui faire utiliser SQLIte pour sa mémoire (MCP local) ... ou GitLab
faire son serveur MCP (en go ...) mais pas le SSE, un truc simple
	wasm: un simulateur de combat
re-utiliser le module wasm avec wasimancer
et le faire utiliser par Sarah Connor et un autre bot Terminator
-> comment les faire discuter tous les deux ?

Est-ce que j'utilise Parakeet?
Est-ce que j'utilise l'API Go d'Ollama (ce que j'avais fait pour le devfest Toulouse)
Ou directement l'API Go d'OpenAI (mais fonctionne-t-elle avec Ollama)?
https://github.com/k33g/sarah-connor-at-devfest-toulouse-2024

Zeira corp
https://terminator.fandom.com/wiki/Zeira_Corporation

_[Terminator: The Sarah Connor Chronicles](https://terminator.fandom.com/wiki/Terminator:_The_Sarah_Connor_Chronicles "Terminator: The Sarah Connor Chronicles")_



the-sarah-connor-chronicles-at-sunny-tech-2025

https://github.com/search?q=repo%3Aopenai%2Fopenai-go%20EmbeddingNewParamsInputUnion&type=code


Pour présenter le fonctionnement de MCP:
https://deadprogrammersociety.com/2025/03/calling-mcp-servers-the-hard-way



### Attention
Penser à parler de la partie sécurité

✋ Faire une partie sur:
mais comment on intègre tout ça dans les messages de prompt


### MCP est un truc tout con
