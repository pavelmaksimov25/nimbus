# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Nimbus is a distributed workflow engine in Go providing task management through a REST API. Tasks follow a state machine: `NEW -> IN_PROGRESS -> COMPLETED|FAILED`.

## Common Commands

```bash
make dev              # Start dev environment (db + migrations + app on :8080)
make test             # Run tests with race detection and coverage
make test-coverage    # Generate HTML coverage report (coverage.html)
make lint             # Run golangci-lint
make build            # Build Linux binary to bin/nimbus
make run              # Run app locally (requires DB running)

# Database
make db-start         # Start PostgreSQL via docker-compose
make migrate          # Run pending migrations
make migrate-create name=<name>  # Create new migration file
make migrate-down     # Rollback last migration
make migrate-status   # Show migration status

# Docker
make docker-up        # Start all services
make docker-down      # Stop all services
```

Run a single test:
```bash
go test -v -run TestFunctionName ./internal/workflow/application/task/
```

## Architecture

Clean Architecture with four layers per domain module. Currently one domain module exists (`internal/workflow/`), with `internal/task_runnner/` scaffolded.

```
internal/workflow/
├── domain/           # Entities, interfaces, custom error types
│   ├── entity/       # Workflow and Task structs (GORM-annotated)
│   ├── repository/   # Repository interfaces + generated mocks (uber/mock)
│   ├── service/      # Service interfaces (TaskService, WorkflowService)
│   └── types/        # Domain error types (RecordNotFoundError, UnprocessableEntityError)
├── application/      # Service implementations with business logic
│   ├── task/         # TaskService - task lifecycle operations
│   └── workflow/     # WorkflowService - workflow CRUD
├── adapters/         # Infrastructure implementations
│   └── repository/   # GORM-backed repository implementations
└── presenter/        # HTTP layer
    └── rest_api/     # Gin handlers, router setup, route registration
```

**Dependency flow:** presenter -> service interfaces <- application -> repository interfaces <- adapters

**Key patterns:**
- Interfaces defined in `domain/`, implementations in `application/` and `adapters/`
- Constructor-based dependency injection wired in `cmd/app/main.go`
- Custom error types with `Unwrap()` for `errors.Is()` matching (e.g., `types.ErrNotFound`)
- Mocks generated with `go.uber.org/mock` in `domain/repository/mocks/`
- Tests use `testify` assertions with mock repositories

## Tech Stack

- **Go 1.24**, **Gin** (HTTP), **GORM** (ORM), **PostgreSQL 16**
- **goose** for SQL migrations (`migrations/` directory)
- **uber/mock** for mock generation, **testify** for assertions
- Config via `.env` file loaded with `godotenv` (key var: `DB_DSN`)

## CI

GitHub Actions (`.github/workflows/ci.yaml`): `gofmt` check -> tests with race detection -> Codecov upload. Triggers on push to `main` and PRs to `main`, `feature/*`, `bugfix/*`.

## API

All endpoints under `/api/v1`. REST handlers in `presenter/rest_api/`. OpenAPI spec at `docs/apidoc/apidoc.yaml`.
