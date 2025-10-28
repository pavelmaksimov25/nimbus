# Nimbus

## Description

Nimbus is a distributed workflow engine built with Go that provides task management capabilities through a RESTful API. The application allows users to create, track, and manage tasks through different lifecycle states (NEW, IN_PROGRESS, COMPLETED, FAILED).

## Requirements

- Go 1.23.5 or higher
- No external database required (uses in-memory storage)

## Installation

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

## Todo

- [ ] Implement persistent storage (database)
- [ ] Add workflow orchestration capabilities
- [ ] Implement task queue system
- [ ] Add authentication and authorization
- [ ] Add logging and monitoring
- [ ] Implement task retry mechanisms
- [ ] Add comprehensive API documentation
- [ ] Add Docker support
- [ ] Implement distributed task execution
- [ ] Add metrics and health check endpoints
