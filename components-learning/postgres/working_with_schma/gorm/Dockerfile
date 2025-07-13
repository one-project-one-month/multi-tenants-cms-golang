FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:latest
WORKDIR /root
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server .

CMD ["./server"]