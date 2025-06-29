FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o service-b .

FROM alpine:latest
WORKDIR /root
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server .
EXPOSE 8001

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:8001/health || exist 1

CMD ["./server"]