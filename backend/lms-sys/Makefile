# Makefile
.PHONY: help install-tools create-env db-url migrate-up migrate-down migrate-create migrate-force migrate-version migrate-drop sqlc-generate clean build run test docker-up docker-down dev-setup prod-setup

# Load environment variables
include .env
export

# Database URL construction
DB_URL = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-tools: ## Install required tools
	@echo "Installing required tools..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Tools installed successfully!"

create-env: ## Create .env file from .env.example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created! Please update it with your actual values."; \
	else \
		echo ".env file already exists!"; \
	fi

db-url: ## Show database URL (for debugging)
	@echo "Database URL: $(DB_URL)"

# Migration targets
migrate-up: ## Apply all pending migrations
	@echo "Applying migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
	@echo "Migrations applied successfully!"

migrate-down: ## Rollback all migrations
	@echo "Rolling back all migrations..."
	@read -p "Are you sure you want to rollback ALL migrations? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down; \
		echo "All migrations rolled back!"; \
	else \
		echo "Migration rollback cancelled."; \
	fi

migrate-down-1: ## Rollback last migration
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1
	@echo "Last migration rolled back!"

migrate-up-1: ## Apply next migration
	@echo "Applying next migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up 1
	@echo "Next migration applied!"

migrate-create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo "Migration name cannot be empty!"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name; \
	echo "Migration files created for: $$name"

migrate-force: ## Force migration to specific version
	@read -p "Enter version to force: " version; \
	if [ -z "$$version" ]; then \
		echo "Version cannot be empty!"; \
		exit 1; \
	fi; \
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $$version; \
	echo "Migration forced to version: $$version"

migrate-version: ## Show current migration version
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-drop: ## Drop everything in database
	@echo "This will DROP ALL DATA in the database!"
	@read -p "Are you absolutely sure? Type 'DROP' to confirm: " confirm; \
	if [ "$$confirm" = "DROP" ]; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop; \
		echo "Database dropped!"; \
	else \
		echo "Database drop cancelled."; \
	fi

# SQLC targets
sqlc-generate: ## Generate Go code from SQL
	@echo "Generating Go code from SQL..."
	sqlc generate
	@echo "Go code generated successfully!"

sqlc-verify: ## Verify SQL queries
	@echo "Verifying SQL queries..."
	sqlc verify
	@echo "SQL queries verified successfully!"

sqlc-vet: ## Lint SQL queries
	@echo "Linting SQL queries..."
	sqlc vet
	@echo "SQL queries linted successfully!"

# Database management
db-create: ## Create database
	@echo "Creating database: $(DB_NAME)"
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "Database created successfully!"

db-drop: ## Drop database
	@echo "Dropping database: $(DB_NAME)"
	@read -p "Are you sure you want to drop the database? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME); \
		echo "Database dropped!"; \
	else \
		echo "Database drop cancelled."; \
	fi

db-reset: ## Drop and recreate database, then run migrations
	@echo "Resetting database..."
	make db-drop
	make db-create
	make migrate-up
	@echo "Database reset complete!"

# Application targets
build: ## Build the application
	@echo "Building application..."
	go build -o bin/lms-app ./cmd/server
	@echo "Application built successfully!"

run: ## Run the application
	@echo "Running application..."
	go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker targets
docker-up: ## Start docker containers
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Docker containers started!"

docker-down: ## Stop docker containers
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "Docker containers stopped!"

docker-logs: ## Show docker logs
	docker-compose logs -f

# Setup targets
dev-setup: create-env install-tools migrate-up sqlc-generate ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Update .env file with your actual database credentials"
	@echo "2. Run 'make db-create' to create the database"
	@echo "3. Run 'make migrate-up' to apply migrations"
	@echo "4. Run 'make sqlc-generate' to generate Go code"
	@echo "5. Run 'make run' to start the application"

prod-setup: ## Setup production environment
	@echo "Setting up production environment..."
	@if [ "$(APP_ENV)" != "production" ]; then \
		echo "Warning: APP_ENV is not set to 'production'"; \
	fi
	make migrate-up
	make sqlc-generate
	make build
	@echo "Production setup complete!"

# Cleanup targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

clean-all: clean ## Clean everything including generated code
	@echo "Cleaning all generated files..."
	rm -rf $(GENERATED_PATH)
	@echo "All clean complete!"

# Seed data (example)
seed: ## Seed database with sample data
	@echo "Seeding database with sample data..."
	@if [ -f ./scripts/seed.sql ]; then \
		psql "$(DB_URL)" -f ./scripts/seed.sql; \
		echo "Database seeded successfully!"; \
	else \
		echo "Seed file not found: ./scripts/seed.sql"; \
	fi

# Backup and restore
backup: ## Backup database
	@echo "Creating database backup..."
	pg_dump "$(DB_URL)" > backups/$(DB_NAME)_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup created in backups/ directory"

restore: ## Restore database from backup
	@read -p "Enter backup filename: " filename; \
	if [ -f "backups/$$filename" ]; then \
		psql "$(DB_URL)" < "backups/$$filename"; \
		echo "Database restored from: $$filename"; \
	else \
		echo "Backup file not found: backups/$$filename"; \
	fi

# Proto targets
proto-setup: ## Create proto folder structure
	@echo "Creating proto folder structure..."
	chmod +x scripts/generate_proto_structure.sh
	./scripts/generate_proto_structure.sh

proto-clean: ## Clean proto directories
	@echo "Cleaning proto directories..."
	rm -rf proto/
	@echo "Proto directories cleaned!"

proto-tree: ## Show proto directory structure
	@echo "Proto directory structure:"
	@if command -v tree >/dev/null 2>&1; then \
		tree proto/; \
	else \
		find proto/ -type d | sed 's|[^/]*/|  |g'; \
	fi
protoc:
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

# Show current configuration
show-config: ## Show current configuration
	@echo "Current Configuration:"
	@echo "====================="
	@echo "Database Host: $(DB_HOST)"
	@echo "Database Port: $(DB_PORT)"
	@echo "Database Name: $(DB_NAME)"
	@echo "Database User: $(DB_USER)"
	@echo "Migrations Path: $(MIGRATIONS_PATH)"
	@echo "Schema Path: $(SCHEMA_PATH)"
	@echo "Queries Path: $(QUERIES_PATH)"
	@echo "Generated Path: $(GENERATED_PATH)"
	@echo "App Port: $(APP_PORT)"
	@echo "App Environment: $(APP_ENV)"