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
| `make install` | Install dependencies, security tools, and linters |
| `make run` | Run the application locally |
| `make build` | Build the application binary |
| `make ci` | Run full CI pipeline locally (lint + test + security + build) |
| `make migrate` | Run database migrations |
| `make migrate-create name=<name>` | Create a new migration |
| `make docker-up` | Start all services with docker-compose |
| `make docker-down` | Stop all services |
| `make docker-logs` | View docker-compose logs |
| `make db-start` | Start only the database |
| `make db-shell` | Open PostgreSQL shell |
| `make test` | Run tests with race detection |
| `make lint` | Run go vet and golangci-lint |
| `make security` | Run all security scanners (govulncheck, gosec, staticcheck) |
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

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_DSN` | Yes | - | PostgreSQL connection string |
| `CORS_ALLOWED_ORIGINS` | No | `*` | Comma-separated list of allowed CORS origins |
| `RATE_LIMIT` | No | `100-M` | Rate limit per IP (format: `count-period`, e.g. `100-M` = 100/min) |

**Example `.env` file:**
```
DB_DSN="postgres://nimbus:nimbus_password@localhost:5432/nimbus_db?sslmode=disable&charset=utf8mb4&parseTime=True&loc=Local"
CORS_ALLOWED_ORIGINS=http://localhost:3000
RATE_LIMIT=100-M
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
- Security middleware (CORS, secure headers, rate limiting, request ID tracing)
- Extensible architecture for future adapters (queue, database)

## Docker Support

The project includes Docker support with multi-stage builds for optimal image size and security:

- **Dockerfile**: Multi-stage build with Go 1.24 and Alpine Linux
- **docker-compose.yml**: Orchestrates app and PostgreSQL 16 database
- Non-root user execution for security
- Health checks for database readiness
- Persistent volume for database data

## Database Migrations

The project uses [goose](https://github.com/pressly/goose) for database schema management. Migrations run automatically when using Docker Compose or Make commands.

**Quick Migration Commands:**
```bash
# Run all pending migrations
make migrate

# Create a new migration
make migrate-create name=add_user_table

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status
```

For detailed information on creating and managing migrations, see [Migration Documentation](docs/migration.md).

## Security

The application includes security middleware applied to all routes:

- **Request ID** — `X-Request-Id` header on every response for request tracing
- **Secure Headers** — `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, XSS filter, `Content-Security-Policy`, `Referrer-Policy`
- **CORS** — configurable allowed origins via `CORS_ALLOWED_ORIGINS` env var
- **Rate Limiting** — per-IP rate limiting configurable via `RATE_LIMIT` env var (default: 100 requests/minute)
- **Body Size Limit** — 1 MB max request body

### Security Scanning

Run security scanners locally:

```bash
make security                # Run all scanners
make security-govulncheck    # Check known CVEs in dependencies
make security-gosec          # Static application security testing
make security-staticcheck    # Advanced static analysis
```

Security scanning also runs automatically in CI on every push and pull request.

## CI/CD

GitHub Actions runs 4 parallel jobs on every push to `main` and PRs:

| Job | Tools | Purpose |
|-----|-------|---------|
| Lint | gofmt, go vet, golangci-lint | Code quality and formatting |
| Test | go test -race, Codecov | Correctness and coverage |
| Security | govulncheck, gosec | Vulnerability and SAST scanning |
| Build | go build | Binary compilation check |

Run the full pipeline locally with `make ci`.

## Todo
- [x] Implement persistent storage (database)
- [x] Add security middleware (CORS, secure headers, rate limiting)
- [x] Add security scanning to CI pipeline
- [ ] Add workflow orchestration capabilities
- [ ] Implement task queue system
- [ ] Add authentication and authorization
- [ ] Add logging and monitoring
- [ ] Implement task retry mechanisms
- [ ] Add comprehensive API documentation
- [ ] Implement distributed task execution
- [ ] Add metrics and health check endpoints
