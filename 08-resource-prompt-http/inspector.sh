#!/bin/bash
: <<'COMMENT'
Inspector project:
https://modelcontextprotocol.io/docs/tools/inspector
COMMENT

npx @modelcontextprotocol/inspector@0.13.0


# List available tools
#npx @modelcontextprotocol/inspector@0.13.0 --cli ./mcp-dd --method tools/list
#npx @modelcontextprotocol/inspector@0.14.0 --cli ./mcp-dd --method tools/list | jq -c '.' | jq -s '.'


# Call a specific tool
#npx @modelcontextprotocol/inspector --cli node build/index.js --method tools/call --tool-name mytool --tool-arg key=value --tool-arg another=value2
