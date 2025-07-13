# Content Management System (CMS) - Project Architecture Documentation

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║                         🏗️  CMS ARCHITECTURE OVERVIEW                         ║
║                                                                               ║
║   A modern, scalable Content Management System built with Go                  ║
║   following clean architecture principles and domain-driven design            ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

## Overview

This Content Management System (CMS) is architected as a high-performance, scalable web application built with Go. The system employs a layered architecture pattern that ensures maintainability, testability, and clear separation of concerns across all components.

## Project Structure

```
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd/
│   └── main.go
├── go.mod
├── go.sum
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   └── requestHandler.go
│   ├── repository/
│   │   └── auth_repo.go
│   ├── routes/
│   │   └── auth_route.go
│   ├── service/
│   │   └── auth_service.go
│   └── types/
│       ├── model_types.go
│       ├── request.go
│       └── response.go
├── pkg/
│   └── utils/
│       ├── context.go
│       ├── database.go
│       ├── env.go
│       ├── health.go
│       ├── jwt_utils.go
│       ├── logger.go
│       ├── password.go
│       └── response.go
└── test/
```

## Architecture Layers

### 🚀 Application Entry Point (`cmd/`)

**Purpose**: Application bootstrap and initialization
- `main.go` - Entry point that orchestrates server startup and component initialization

### 🏛️ Core Business Logic (`internal/`)

The heart of the CMS application, implementing the layered architecture pattern:

#### 🎯 Handler Layer (`internal/handler/`)
**Responsibility**: HTTP request/response handling and API endpoint controllers
- `auth_handler.go` - Authentication endpoint handlers (login, registration, logout)
- `requestHandler.go` - General request processing and routing logic

#### ⚙️ Service Layer (`internal/service/`)
**Responsibility**: Business logic implementation and orchestration
- `auth_service.go` - Authentication business logic, validation, and processing workflows

#### 🗄️ Repository Layer (`internal/repository/`)
**Responsibility**: Data persistence and database operations
- `auth_repo.go` - Authentication-related database queries and data access patterns

#### 🛣️ Routing Layer (`internal/routes/`)
**Responsibility**: API route definitions and middleware configuration
- `auth_route.go` - Authentication route registration and endpoint mapping

#### 📋 Types Layer (`internal/types/`)
**Responsibility**: Data structure definitions and contracts
- `model_types.go` - Domain models and database entities
- `request.go` - HTTP request payload structures
- `response.go` - HTTP response payload structures

### 🔧 Shared Utilities (`pkg/`)

**Purpose**: Reusable components and cross-cutting concerns

#### 🛠️ Utility Package (`pkg/utils/`)
- `context.go` - Context management and request tracing utilities
- `database.go` - Database connection management and configuration
- `env.go` - Environment variable handling and configuration loading
- `health.go` - Application health check implementations
- `jwt_utils.go` - JWT token generation, validation, and management
- `logger.go` - Structured logging configuration and utilities
- `password.go` - Password hashing, verification, and security utilities
- `response.go` - Standardized HTTP response formatting and error handling

### 🚢 Development & Deployment

- **Dockerfile** - Container configuration for deployment environments
- **Makefile** - Build automation and development workflow commands
- **README.md** - Project documentation and setup instructions
- **go.mod/go.sum** - Go module dependency management
- **test/** - Test suites, fixtures, and testing utilities

## Design Principles

This CMS architecture implements several key architectural principles:

- **🎯 Separation of Concerns**: Each layer has a distinct responsibility
- **🔄 Dependency Inversion**: Higher-level modules don't depend on lower-level modules
- **🏗️ Clean Architecture**: Business logic is independent of external frameworks
- **🧩 Modularity**: Components are loosely coupled and highly cohesive
- **🧪 Testability**: Clear boundaries enable comprehensive testing strategies

---

