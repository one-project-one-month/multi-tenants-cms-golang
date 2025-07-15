# Variables
PROTO_DIR = proto
PROTOGEN_DIR = protogen
APP_DIR = app
BINARY_NAME = server
MAIN_FILE = cmd/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

# Run the application
.PHONY: run
run:
	go run $(MAIN_FILE)

protoc: ## Generate protobuf files
	@echo "Generating proto..."
	cd proto && protoc --go_out=../protogen --go_opt=paths=source_relative \
	--go-grpc_out=../protogen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=../protogen --grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out=../doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api \
	--experimental_allow_proto3_optional \
	./**/*.proto
	@echo "Generating swagger statik..."
	statik -src=doc/swagger -dest=doc -f
# Generate protobuf files
.PHONY: proto
proto:
	buf generate

# Generate database code with sqlc
.PHONY: sqlc
sqlc:
	sqlc generate

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	rm -rf $(PROTOGEN_DIR)/*

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Test
.PHONY: test
test:
	go test ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME) .

# Docker run
.PHONY: docker-run
docker-run:
	docker run -p 8082:8082 -p 8083:8083 $(BINARY_NAME)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  proto       - Generate protobuf files"
	@echo "  sqlc        - Generate database code"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  help        - Show this help message"