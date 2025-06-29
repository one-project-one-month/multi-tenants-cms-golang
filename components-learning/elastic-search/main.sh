#!/bin/bash

# Set project root directory
PROJECT_NAME="product-api"

# Define folder structure
DIRS=(
  "$PROJECT_NAME/config"
  "$PROJECT_NAME/models"
  "$PROJECT_NAME/handlers"
  "$PROJECT_NAME/services"
  "$PROJECT_NAME/database"
  "$PROJECT_NAME/elasticsearch"
  "$PROJECT_NAME/dto"
)

FILES=(
  "$PROJECT_NAME/main.go"
  "$PROJECT_NAME/config/config.go"
  "$PROJECT_NAME/models/product.go"
  "$PROJECT_NAME/handlers/product_handler.go"
  "$PROJECT_NAME/services/product_service.go"
  "$PROJECT_NAME/database/postgres.go"
  "$PROJECT_NAME/elasticsearch/client.go"
  "$PROJECT_NAME/dto/product_dto.go"
  "$PROJECT_NAME/go.mod"
  "$PROJECT_NAME/docker-compose.yml"
)

# Create directories
for dir in "${DIRS[@]}"; do
  mkdir -p "$dir"
done

# Create files
for file in "${FILES[@]}"; do
  touch "$file"
done

echo "Project structure '$PROJECT_NAME' created successfully."
