# Content Management System (CMS) - Complete Repository Architecture

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║                    🏗️   CMS PLATFORM ARCHITECTURE                             ║
║                                                                               ║
║   A comprehensive, microservices-based Content Management System              ║
║   with Learning Management capabilities, built for enterprise scale           ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

## 🌟 Repository Overview

This  CMS platform implements a sophisticated microservices architecture with multiple backend services, comprehensive infrastructure management, and extensive learning components for educational content delivery.

## 📁 Complete Project Structure

```
.
├── README.md
├── backend/
│   ├── cms-sys/                          # Core CMS Service (Go)
│   ├── gateway/                          # API Gateway Service (Spring Boot)
│   └── lms-sys/                          # Learning Management System (Go)
├── components-learning/                   # Educational Components & Examples
│   ├── certificate_signing/              # Digital Certificate Generation
│   ├── docker/                           # Docker Learning Examples
│   ├── elastic-search/                   # Search Implementation
│   ├── mfa/                              # Multi-Factor Authentication
│   └── postgres/                         # Database Examples
├── docker-compose.yml                     # Main Orchestration
├── infra/                                # Infrastructure as Code
│   ├── migration/                        # Database Migrations
│   ├── github/                           # CI/CD Configuration
│   ├── schema/                           # Database Schemas
│   └── terraform/                        # Infrastructure Provisioning
├── keys/                                 # Cryptographic Keys
├── scripts/                              # Automation Scripts
└── vault-init.sh                         # Vault Initialization
```

## 🏛️ Architecture Overview

### 🚀 Backend Services (`backend/`)

#### 1. **CMS Core System** (`cms-sys/`)
**Technology**: Go | **Purpose**: Primary content management functionality

```
cms-sys/
├── cmd/main.go                           # Application entry point
├── internal/                             # Core business logic
│   ├── handler/                          # HTTP handlers
│   ├── repository/                       # Data access layer
│   ├── routes/                           # API route definitions
│   ├── service/                          # Business logic layer
│   └── types/                            # Data structures
├── pkg/utils/                            # Shared utilities
└── test/                                 # Test suites
```

**Key Features**:
- Authentication & authorization
- Content CRUD operations
- JWT token management
- Database abstraction
- Health monitoring

#### 2. **API Gateway** (`gateway/`)
**Technology**: Spring Boot (Java) | **Purpose**: Service orchestration and routing

```
gateway/
├── src/main/java/com/content_management_system/gateway/
│   └── GatewayApplication.java           # Spring Boot application
├── src/main/resources/
│   └── application.yml                   # Configuration
└── src/test/                             # Integration tests
```

**Key Features**:
- Request routing and load balancing
- Cross-cutting concerns (logging, monitoring)
- API versioning and documentation
- Security policies enforcement

#### 3. **Learning Management System** (`lms-sys/`)
**Technology**: Go | **Purpose**: Educational content delivery

```
lms-sys/
├── cmd/main.go                           # LMS application entry
├── internal/                             # LMS-specific logic
│   ├── config/                           # Configuration management
│   ├── handler/                          # LMS API handlers
│   ├── model/                            # Domain models
│   ├── repository/                       # Data persistence
│   └── service/                          # Business services
└── pkg/utils/                            # LMS utilities
```

**Key Features**:
- Course management
- Student progress tracking
- Assessment and grading
- Content delivery optimization

### 🧩 Learning Components (`components-learning/`)

#### **Certificate Signing** (`certificate_signing/`)
- Digital certificate generation and validation
- PDF certificate creation with cryptographic signatures
- Private key management and security

#### **Docker Orchestration** (`docker/`)
- Multi-service containerization examples
- Service discovery and communication
- Container networking and scaling

#### **Elasticsearch Integration** (`elastic-search/`)
- Full-text search implementation
- Product catalog search optimization
- Real-time indexing and querying

#### **Multi-Factor Authentication** (`mfa/`)
- TOTP (Time-based One-Time Password) implementation
- Security token generation and validation
- Authentication flow optimization

#### **PostgreSQL Advanced Features** (`postgres/`)
- **Schema Management**: Multi-tenant architecture patterns
- **Geometric Data**: PostGIS spatial data handling
- **Functions & Domains**: Advanced SQL programming
- **GORM Integration**: Go ORM implementation

### 🏗️ Infrastructure (`infra/`)

#### **Database Management** (`migration/` & `schema/`)
- **CMS Schema**: Content management database structure
- **LMS Schema**: Learning management database design
- **Migration Tools**: Database version control and updates
- **Diagrams**: Visual database architecture documentation

#### **Infrastructure as Code** (`terraform/`)
- Cloud resource provisioning
- Environment-specific configurations
- Scalability and high availability setup

#### **CI/CD Pipeline** (`github/`)
- Automated testing and deployment
- Code quality enforcement
- Security scanning integration

## 🔧 Development & Operations

### **Security & Cryptography** (`keys/`)
- RSA key pair management
- Encryption/decryption utilities
- Secure communication setup

### **Automation** (`scripts/`)
- Deployment automation
- Environment setup scripts
- Maintenance and monitoring tools

### **Service Orchestration** (`docker-compose.yml`)
- Multi-service local development
- Service dependency management
- Environment variable configuration

## 🎯 Design Principles & Patterns

### **Microservices Architecture**
- **Service Independence**: Each service can be developed, deployed, and scaled independently
- **Technology Diversity**: Go for performance-critical services, Spring Boot for gateway complexity
- **Data Isolation**: Each service maintains its own database and data models

### **Clean Architecture Implementation**
- **Separation of Concerns**: Clear boundaries between presentation, business, and data layers
- **Dependency Inversion**: Business logic independent of external frameworks
- **Testability**: Comprehensive testing at unit, integration, and system levels

### **Infrastructure as Code**
- **Reproducible Environments**: Consistent deployment across development, staging, and production
- **Version Control**: Infrastructure changes tracked and reviewed
- **Automation**: Minimal manual intervention in deployment processes

### **Security-First Design**
- **Authentication**: Multi-factor authentication and JWT token management
- **Authorization**: Role-based access control across services
- **Encryption**: End-to-end encryption for sensitive data
- **Certificate Management**: Digital certificates for document authenticity

## Deployment & Scalability

This architecture supports multiple deployment strategies:

- **Local Development**: Docker Compose for rapid development cycles
- **Cloud Deployment**: Terraform for S3 cloud infrastructure
