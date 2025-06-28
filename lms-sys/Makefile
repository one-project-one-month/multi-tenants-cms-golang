.PHONY: all build run test clean fmt vet lint tidy docker-build docker-run help

APP_NAME ?= server
DOCKER_IMAGE ?= $(APP_NAME):latest
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: fmt vet lint test build

build:
	@echo "🔨 Building application..."
	@go build -o bin/$(APP_NAME) ./cmd

run:
	@echo "🚀 Running application..."
	@go run ./cmd

test:
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

fmt:
	@echo "📝 Formatting code..."
	@gofmt -s -w $(GO_FILES)

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

tidy:
	@echo "📦 Tidying go modules..."
	@go mod tidy

docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//'