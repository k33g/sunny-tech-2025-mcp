# Fonctionnement d’un serveur MCP utilisant stdio

Un **serveur** MCP utilisant le transport **stdio** s’appuie sur les flux d’entrée/sortie standards (`stdin` et `stdout`) pour communiquer avec un **client** MCP, comme un IDE (Visual Studio, Cursor, Claude Desktop, etc.) ou un assistant IA.

## Communication

1. Le client envoie des messages structurés (en JSON-RPC 2.0) sur le flux `stdin` du serveur.
2. Le serveur lit ces messages depuis `stdin`, traite la demande (exécution d’un outil, accès à une ressource, etc.), puis écrit la réponse structurée sur `stdout`.

## Cycle de vie typique

1. Initialisation : Le serveur MCP s’initialise, attend les requêtes sur `stdin`.
2. Découverte des outils : Le client interroge le serveur pour découvrir les outils disponibles.
3. Appel d’outil : Lorsqu’une action est demandée (ex : "lire un fichier", "fetch une page web"), le client envoie la requête sur `stdin`.
4. Traitement : Le serveur exécute la logique associée.
5. Réponse : Le serveur écrit la réponse (succès, données, erreur) sur `stdout`, que le client lit et transmet à l’IA ou à l’utilisateur final.

## ✋ Cela se passe en local

___
[◀️ Previous](./15-stdio.md#fonctionnement-dun-programme-utilisant-stdio) | [1st MCP STDIO server ▶️](./main.go) **+ tests**
