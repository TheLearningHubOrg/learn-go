# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server cmd/server/main.go

# Minimal Production Runner
FROM scratch
WORKDIR /
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
