# Nimbus

## Description

Nimbus is a distributed workflow engine built with Go that provides task management capabilities through a RESTful API. The application allows users to create, track, and manage tasks through different lifecycle states (NEW, IN_PROGRESS, COMPLETED, FAILED).

## Requirements

- Go 1.24.0 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL 16 (when using database, or use Docker Compose)

## Installation

### Option 1: Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone git@github.com:pavelmaksimov25/nimbus.git
cd nimbus
```

2. Start all services (app + PostgreSQL database):
```bash
docker-compose up -d
```

3. View logs:
```bash
docker-compose logs -f
```

4. Stop services:
```bash
docker-compose down
```

The server will be available at `http://localhost:8080`

**Environment Variables (configured in docker-compose.yml):**
- `DB_HOST=database`
- `DB_PORT=5432`
- `DB_USER=nimbus`
- `DB_PASSWORD=nimbus_password`
- `DB_NAME=nimbus_db`

### Option 2: Using Docker Only

Build and run the application container:
```bash
docker build -t nimbus:latest .
docker run -p 8080:8080 nimbus:latest
```

### Option 3: Local Development

1. Clone the repository:
```bash
git clone git@github.com:pavelmaksimov25/nimbus.git
cd nimbus
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run cmd/app/main.go
```

The server will start on `http://localhost:8080`

## Architecture

The project follows **Clean Architecture** principles with clear separation of concerns:

- **Domain Layer** (`internal/workflow/domain/`)
  - Defines core business entities (Task, Workflow)
  - Task statuses and domain types
  - Business rules and validation logic

- **Application Layer** (`internal/workflow/application/`)
  - Task service implementing business logic
  - Orchestrates task lifecycle operations (create, start, complete, fail)

- **Adapters Layer** (`internal/workflow/adapters/`)
  - Storage implementations (in-memory task storage)
  - Repository and queue interfaces for future extensions

- **Presenter Layer** (`internal/workflow/presenter/`)
  - REST API implementation using Gin framework
  - HTTP handlers for task operations
  - Endpoints: `/api/v1/tasks/*`

**Key Features:**
- Task lifecycle management (NEW → IN_PROGRESS → COMPLETED/FAILED)
- Thread-safe in-memory storage
- RESTful API with proper error handling
- Extensible architecture for future adapters (queue, database)

## Docker Support

The project includes Docker support with multi-stage builds for optimal image size and security:

- **Dockerfile**: Multi-stage build with Go 1.24 and Alpine Linux
- **docker-compose.yml**: Orchestrates app and PostgreSQL 16 database
- Non-root user execution for security
- Health checks for database readiness
- Persistent volume for database data

## Todo
- [ ] Implement persistent storage (database)
- [ ] Add workflow orchestration capabilities
- [ ] Implement task queue system
- [ ] Add authentication and authorization
- [ ] Add logging and monitoring
- [ ] Implement task retry mechanisms
- [ ] Add comprehensive API documentation
- [ ] Implement distributed task execution
- [ ] Add metrics and health check endpoints
