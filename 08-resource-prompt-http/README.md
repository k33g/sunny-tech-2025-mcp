curl -i -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "id": "init-uuid",
    "params": {
      "protocolVersion": "2024-11-05"
    }
  }' \
  http://localhost:9090/mcp

---
Mcp-Session-Id: mcp-session-806e191c-6358-410d-8835-8520c6111e20

curl -X POST \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: mcp-session-806e191c-6358-410d-8835-8520c6111e20" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": "tools-list-uuid",
    "params": {}
  }' \
  http://localhost:9090/mcp