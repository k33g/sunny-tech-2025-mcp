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