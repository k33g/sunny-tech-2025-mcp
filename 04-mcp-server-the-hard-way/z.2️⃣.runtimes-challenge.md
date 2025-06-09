

# MCP runtime(s) challenge

ici mettre la commande avec le code équivalent docker

l'utiliser aussi avec Claude.ai ou  Copilot?


```bash
cat /tmp/mcp_test_input.jsonl | go run main.go | jq -c '.' | jq -s '.'
cat /tmp/mcp_test_input.jsonl | ./mcp-dd | jq -c '.' | jq -s '.'
```

```dockerfile
FROM --platform=$BUILDPLATFORM golang:1.24.3-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN <<EOF
go mod tidy 
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

FROM scratch
WORKDIR /app
COPY --from=builder /app/mcp-dd .

ENTRYPOINT ["./mcp-dd"]

# docker build --platform linux/arm64 -t mcp-dd:demo .
```



```bash
cat /tmp/mcp_test_input.jsonl | docker run --rm -i mcp-dd:demo | jq -c '.' | jq -s '.'
```

## Demo(s)