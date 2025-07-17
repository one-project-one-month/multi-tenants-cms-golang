# ==============================================================================
# LMS Application Makefile
# ==============================================================================

# ------------------------------------------------------------------------------
# Configuration
# ------------------------------------------------------------------------------

# Application
PROTO_DIR       = proto
PROTOGEN_DIR    = protogen
APP_DIR         = app
BINARY_NAME     = server
MAIN_FILE       = cmd/main.go
BUILD_DIR       = bin

# Database
DB_URL          = postgresql://neondb_owner:npg_MEB4CYJS7TKh@ep-withered-salad-a27oo7gn-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require
MIGRATION_DIR   = database/migrations

# Colors
RED             = \033[0;31m
GREEN           = \033[0;32m
YELLOW          = \033[0;33m
BLUE            = \033[0;34m
NC              = \033[0m

# ------------------------------------------------------------------------------
# Targets
# ------------------------------------------------------------------------------

.PHONY: all
all: deps lint test build ## Build everything (default target)

# ------------------------------------------------------------------------------
# Database Migration Targets
# ------------------------------------------------------------------------------

.PHONY: migrate-create
migrate-create: ## Create a new migration file
	@echo "${YELLOW}Creating new migration...${NC}"
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $${name}
	@echo "${GREEN}Migration file created${NC}"

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@echo "${YELLOW}Applying migrations...${NC}"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up
	@echo "${GREEN}Migrations applied successfully${NC}"

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@echo "${YELLOW}Rolling back last migration...${NC}"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1
	@echo "${GREEN}Migration rolled back successfully${NC}"

.PHONY: migrate-force
migrate-force: ## Force a specific migration version
	@echo "${YELLOW}Forcing migration version...${NC}"
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $${version}
	@echo "${GREEN}Migration version forced${NC}"

.PHONY: migrate-drop
migrate-drop: ## Drop all database objects (DESTRUCTIVE)
	@echo "${RED}WARNING: Dropping all database objects${NC}"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" drop
	@echo "${GREEN}Database objects dropped${NC}"

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

.PHONY: migrate-reset
migrate-reset: migrate-drop migrate-up ## Reset database (drop all then migrate up)

.PHONY: fix-migration-perms
fix-migration-perms: ## Fix permissions on migration files
	@echo "${YELLOW}Fixing migration file permissions...${NC}"
	@chmod -R 755 $(MIGRATION_DIR)
	@echo "${GREEN}Permissions fixed${NC}"

# ------------------------------------------------------------------------------
# Build Targets
# ------------------------------------------------------------------------------

.PHONY: build
build: ## Build the application
	@echo "${YELLOW}Building application...${NC}"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "${GREEN}Build successful${NC}"

.PHONY: run
run: ## Run the application
	@echo "${YELLOW}Starting application...${NC}"
	@go run $(MAIN_FILE)

.PHONY: clean
clean: ## Clean build artifacts
	@echo "${YELLOW}Cleaning build artifacts...${NC}"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(PROTOGEN_DIR)/*
	@echo "${GREEN}Clean complete${NC}"

# ------------------------------------------------------------------------------
# Code Generation Targets
# ------------------------------------------------------------------------------

.PHONY: protoc
protoc: ## Generate protobuf files
	@echo "${YELLOW}Generating protobuf files...${NC}"
	@cd $(PROTO_DIR) && protoc \
		--go_out=../$(PROTOGEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=../$(PROTOGEN_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=../$(PROTOGEN_DIR) --grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out=../doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api \
		--experimental_allow_proto3_optional \
		./**/*.proto
	@echo "${YELLOW}Generating swagger docs...${NC}"
	@statik -src=doc/swagger -dest=doc -f
	@echo "${GREEN}Protobuf generation complete${NC}"

.PHONY: sqlc
sqlc: ## Generate database code with sqlc
	@echo "${YELLOW}Generating database code...${NC}"
	@sqlc generate
	@echo "${GREEN}Database code generation complete${NC}"

# ------------------------------------------------------------------------------
# Development Targets
# ------------------------------------------------------------------------------

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "${YELLOW}Installing dependencies...${NC}"
	@go mod download
	@go mod tidy
	@echo "${GREEN}Dependencies installed${NC}"

.PHONY: test
test: ## Run tests
	@echo "${YELLOW}Running tests...${NC}"
	@go test -v -cover ./...
	@echo "${GREEN}Tests completed${NC}"

.PHONY: fmt
fmt: ## Format code
	@echo "${YELLOW}Formatting code...${NC}"
	@go fmt ./...
	@echo "${GREEN}Formatting complete${NC}"

.PHONY: lint
lint: ## Lint code
	@echo "${YELLOW}Linting code...${NC}"
	@golangci-lint run
	@echo "${GREEN}Linting complete${NC}"

.PHONY: dev
dev: fmt lint test build ## Run standard development checks

# ------------------------------------------------------------------------------
# Docker Targets
# ------------------------------------------------------------------------------

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "${YELLOW}Building Docker image...${NC}"
	@docker build -t $(BINARY_NAME) .
	@echo "${GREEN}Docker build complete${NC}"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "${YELLOW}Starting Docker container...${NC}"
	@docker run -p 8082:8082 -p 8083:8083 $(BINARY_NAME)

# ------------------------------------------------------------------------------
# Setup Targets
# ------------------------------------------------------------------------------

.PHONY: setup
setup: deps ## Setup development environment
	@echo "${YELLOW}Setting up development tools...${NC}"
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/rakyll/statik@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "${GREEN}Development environment setup complete${NC}"

# ------------------------------------------------------------------------------
# Help
# ------------------------------------------------------------------------------

.PHONY: grp-sqlc
grp-sqlc:
	@grpcframe sqlc generate
.PHONY: help
help: ## Show this help message
	@echo "${BLUE}Available targets:${NC}"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-20s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
	@echo ""