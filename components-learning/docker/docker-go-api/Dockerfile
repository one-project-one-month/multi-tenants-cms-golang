#AUTHOR: Swan Htet Aung Phyo

# Building stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o server .

# Runtime Stage
FROM alpine:latest
WORKDIR /root
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server .
EXPOSE 5000

# We should check the service internally
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:5000/health || exit 1

CMD ["./server"]