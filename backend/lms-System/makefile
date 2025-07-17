# Makefile for LMS Application with Database Migrations

# Application Variables
PROTO_DIR = proto
PROTOGEN_DIR = protogen
APP_DIR = app
BINARY_NAME = server
MAIN_FILE = cmd/main.go

# Database Variables
DB_URL = postgresql://neondb_owner:npg_MEB4CYJS7TKh@ep-withered-salad-a27oo7gn-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require
MIGRATION_DIR = database/migrations

# Colors
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
NC = \033[0m

# Default target
.PHONY: all
all: build

#-------------------------------------------------------------------------------
# Database Migration Targets
#-------------------------------------------------------------------------------

## Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $${name}

## Apply all up migrations
migrate-up:
	@echo "$(YELLOW)Applying migrations...$(NC)"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up
	@echo "$(GREEN)Migrations applied successfully$(NC)"

## Apply one down migration
migrate-down:
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1
	@echo "$(GREEN)Migration rolled back successfully$(NC)"

## Force a specific migration version (to fix dirty state)
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $${version}

## Drop everything in the database
migrate-drop:
	@echo "$(RED)Dropping all database objects...$(NC)"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" drop
	@echo "$(GREEN)Database objects dropped successfully$(NC)"

## Show current migration version
migrate-version:
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

## Fix permissions on migration files
fix-migration-perms:
	@chmod -R 755 $(MIGRATION_DIR)

## Reset database (drop, create, migrate)
migrate-reset: migrate-drop migrate-up
	@echo "$(GREEN)Database reset complete$(NC)"

#-------------------------------------------------------------------------------
# Application Build Targets
#-------------------------------------------------------------------------------

## Build the application
build:
	@echo "$(YELLOW)Building application...$(NC)"
	@go build -o bin/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)Build successful$(NC)"

## Run the application
run:
	@echo "$(YELLOW)Starting application...$(NC)"
	@go run $(MAIN_FILE)

## Generate protobuf files
proto:
	@echo "$(YELLOW)Generating protobuf files...$(NC)"
	@cd proto && protoc --go_out=../protogen --go_opt=paths=source_relative \
	--go-grpc_out=../protogen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=../protogen --grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out=../doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api \
	--experimental_allow_proto3_optional \
	./**/*.proto
	@echo "$(YELLOW)Generating swagger statik...$(NC)"
	@statik -src=doc/swagger -dest=doc -f
	@echo "$(GREEN)Protobuf generation complete$(NC)"

## Generate database code with sqlc
sqlc:
	@echo "$(YELLOW)Generating database code...$(NC)"
	@sqlc generate
	@echo "$(GREEN)Database code generation complete$(NC)"

## Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -rf $(PROTOGEN_DIR)/*
	@echo "$(GREEN)Clean complete$(NC)"

## Install dependencies
deps:
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies installed$(NC)"

## Run tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test ./...
	@echo "$(GREEN)Tests completed$(NC)"

## Format code
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Formatting complete$(NC)"

## Lint code
lint:
	@echo "$(YELLOW)Linting code...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)Linting complete$(NC)"

## Docker build
docker-build:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(BINARY_NAME) .
	@echo "$(GREEN)Docker build complete$(NC)"

## Docker run
docker-run:
	@echo "$(YELLOW)Starting Docker container...$(NC)"
	@docker run -p 8082:8082 -p 8083:8083 $(BINARY_NAME)

## Setup development environment (run once)
setup: deps
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/rakyll/statik@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)Development environment setup complete$(NC)"

## Show help message
help:
	@echo "Available targets:"
	@echo ""
	@echo "Database Migration:"
	@echo "  migrate-create      - Create a new migration file"
	@echo "  migrate-up         - Apply all up migrations"
	@echo "  migrate-down       - Rollback the last migration"
	@echo "  migrate-force      - Force a specific migration version"
	@echo "  migrate-drop       - Drop all database objects"
	@echo "  migrate-version    - Show current migration version"
	@echo "  migrate-reset      - Reset database (drop, create, migrate)"
	@echo "  fix-migration-perms - Fix migration file permissions"
	@echo ""
	@echo "Application:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  proto              - Generate protobuf files"
	@echo "  sqlc               - Generate database code"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Install dependencies"
	@echo "  test               - Run tests"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  setup              - Setup development environment"
	@echo "  help               - Show this help message"