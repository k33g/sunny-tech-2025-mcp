FROM --platform=$BUILDPLATFORM golang:1.24.3-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN <<EOF
go mod tidy 
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/dungeon.md .
COPY --from=builder /app/mcp-dd-http .

EXPOSE 9090

ENTRYPOINT ["./mcp-dd-http"]

# docker build --platform linux/arm64 -t mcp-dd-http:demo .
# docker run --rm -p 9090:9090 mcp-dd-http:demo

# http://host.docker.internal:9090/mcp
# curl http://0.0.0.0:5050/mcp
# curl http://0.0.0.0:9090/mcp