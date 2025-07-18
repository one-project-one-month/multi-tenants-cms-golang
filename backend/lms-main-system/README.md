# github.com/multi-tenants-cms-golang/lms-sys

A gRPC-based microservice built with Go, featuring both gRPC API and HTTP gateway support.

## Project Structure

```
.
├── Dockerfile
├── README.md
├── app
│   ├── gateway
│   │   └── gateway.go
│   └── rpc-api
│       └── server.go
├── buf.yaml
├── cmd
│   └── main.go
├── database
│   ├── migration
│   ├── queries
│   └── schema
├── go.mod
├── internal
│   ├── repo
│   └── services
├── makefile
├── pkg
│   └── utils
├── proto
├── protogen
└── sqlc.yaml
```

## Directory Structure

- **app/**: Application layer containing API implementations
    - **gateway/**: HTTP gateway server implementation
    - **rpc-api/**: gRPC server implementation
- **cmd/**: Application entry point
- **database/**: Database-related files
    - **migration/**: Database migration files
    - **queries/**: SQL query files for sqlc
    - **schema/**: Database schema definitions
- **internal/**: Private application code
    - **repo/**: Repository layer for data access
    - **services/**: Business logic layer
- **pkg/**: Public library code and utilities
- **proto/**: Protocol buffer definitions
- **protogen/**: Generated protobuf code
- **buf.yaml**: Buf configuration for protobuf generation
- **sqlc.yaml**: SQLC configuration for database code generation

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker (optional)
- Buf CLI (for protobuf generation)
- SQLC (for database code generation)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   make deps
   ```

### Development

#### Generate Code

Generate protobuf files:

```bash
make proto
```

Generate database code:

```bash
make sqlc
```

#### Build and Run

Build the application:

```bash
make build
```

Run the application:

```bash
make run
```

#### Testing

Run tests:

```bash
make test
```

Format code:

```bash
make fmt
```

Lint code:

```bash
make lint
```

### Docker

Build Docker image:

```bash
make docker-build
```

Run Docker container:

```bash
make docker-run
```

### Available Make Targets

Run ```make help``` to see all available targets.

## API Endpoints

- gRPC Server: :8082
- HTTP Gateway: :8083

## Configuration

Configuration can be set through environment variables or configuration files.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

# Current routes

# User Registration

```zsh
  curl -X POST --location "http://localhost:8082/lms/v1/auth/register" \
    -H "accept: application/json" \
    -H "X-Organisation: <tenants-namespace>" \
    -H "Content-Type: application/json" \
    -d '{
          "username": "SwanhtetAungphyo",
          "email": "",
          "password": "",
          "address": "Swan Krakow",
          "phoneNumber": "",
          "mfaVerificationOption": true
        }'
```

