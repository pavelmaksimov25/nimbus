# Nimbus

## Description

Nimbus is a distributed workflow engine built with Go that provides task management capabilities through a RESTful API. The application allows users to create, track, and manage tasks through different lifecycle states (NEW, IN_PROGRESS, COMPLETED, FAILED).

## Requirements

- Go 1.24.0 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL 16 (when using database, or use Docker Compose)
- Make (optional, but recommended for easier project management)

## Installation

### Quick Start with Make (Recommended)

The project includes a Makefile for easy setup and management. To see all available commands:

```bash
make help
```

#### Complete Setup

```bash
# Clone the repository
git clone git@github.com:pavelmaksimov25/nimbus.git
cd nimbus

# Run complete setup (installs dependencies and starts services)
make setup

# Start development environment (database + migrations + application)
make dev
```

#### Common Make Commands

| Command | Description |
|---------|-------------|
| `make help` | Display all available commands |
| `make install` | Install dependencies and setup environment |
| `make run` | Run the application locally |
| `make build` | Build the application binary |
| `make migrate` | Run database migrations |
| `make migrate-create name=<name>` | Create a new migration |
| `make docker-up` | Start all services with docker-compose |
| `make docker-down` | Stop all services |
| `make docker-logs` | View docker-compose logs |
| `make db-start` | Start only the database |
| `make db-shell` | Open PostgreSQL shell |
| `make test` | Run tests |
| `make docs` | Show documentation information |
| `make clean` | Clean build artifacts |

For a complete list of commands and their descriptions, run `make help`.

### Option 1: Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone git@github.com:pavelmaksimov25/nimbus.git
cd nimbus
```

2. Set up environment variables:
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env if needed (default values work with docker-compose)
```

The default `.env` file contains:
```
DB_DSN="postgres://nimbus:nimbus_password@database:5432/nimbus_db?sslmode=disable&charset=utf8mb4&parseTime=True&loc=Local"
```

3. Start all services (database + migrations):
```bash
docker-compose up -d
```

The migrations will run automatically on startup.

4. View logs:
```bash
docker-compose logs -f
```

5. Stop services:
```bash
docker-compose down
```

**Note:** The current docker-compose setup runs the database and migrations. To run the application, use Option 2 or Option 3.

### Option 2: Using Docker Only

1. Ensure the database is running:
```bash
docker-compose up -d database
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env to use localhost instead of 'database' hostname
```

Update `.env` for local Docker:
```
DB_DSN="postgres://nimbus:nimbus_password@localhost:5432/nimbus_db?sslmode=disable&charset=utf8mb4&parseTime=True&loc=Local"
```

3. Build and run the application container:
```bash
docker build -t nimbus:latest .
docker run -p 8080:8080 --env-file .env --network host nimbus:latest
```

The server will be available at `http://localhost:8080`

### Option 3: Local Development

1. Clone the repository:
```bash
git clone git@github.com:pavelmaksimov25/nimbus.git
cd nimbus
```

2. Set up environment variables:
```bash
cp .env.example .env
```

Update `.env` for local development:
```
DB_DSN="postgres://nimbus:nimbus_password@localhost:5432/nimbus_db?sslmode=disable&charset=utf8mb4&parseTime=True&loc=Local"
```

3. Start the database:
```bash
docker-compose up -d database
```

4. Run migrations:
```bash
docker-compose up migration
```

5. Install dependencies:
```bash
go mod download
```

6. Run the application:
```bash
go run cmd/app/main.go
```

The server will start on `http://localhost:8080`

**Environment Variables:**

The application uses the [godotenv](https://github.com/joho/godotenv) library to load environment variables from a `.env` file. The following variables are supported:

- `DB_DSN` - PostgreSQL connection string (required)

**Example `.env` file:**
```
DB_DSN="postgres://nimbus:nimbus_password@localhost:5432/nimbus_db?sslmode=disable&charset=utf8mb4&parseTime=True&loc=Local"
```

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

## Database Migrations

The project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. Migrations run automatically when using Docker Compose or Make commands.

**Quick Migration Commands:**
```bash
# Run all pending migrations
make migrate

# Create a new migration
make migrate-create name=add_user_table

# Rollback last migration
make migrate-down

# Check current migration version
make migrate-version
```

For detailed information on creating and managing migrations, see [Migration Documentation](docs/migration.md).

## Todo
- [x] Implement persistent storage (database)
- [ ] Add workflow orchestration capabilities
- [ ] Implement task queue system
- [ ] Add authentication and authorization
- [ ] Add logging and monitoring
- [ ] Implement task retry mechanisms
- [ ] Add comprehensive API documentation
- [ ] Implement distributed task execution
- [ ] Add metrics and health check endpoints
