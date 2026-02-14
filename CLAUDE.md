# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Nimbus is a distributed workflow engine in Go providing task management through a REST API. Tasks follow a state machine: `NEW -> IN_PROGRESS -> COMPLETED|FAILED`.

## Common Commands

```bash
make dev              # Start dev environment (db + migrations + app on :8080)
make test             # Run tests with race detection and coverage
make test-coverage    # Generate HTML coverage report (coverage.html)
make lint             # Run go vet + golangci-lint
make build            # Build Linux binary to bin/nimbus
make run              # Run app locally (requires DB running)
make ci               # Run full CI pipeline locally (lint + test + security + build)

# Security
make security                # Run all security scanners
make security-govulncheck    # Check known CVEs in dependencies
make security-gosec          # SAST scanner for Go
make security-staticcheck    # Advanced static analysis

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
- Config via `.env` file loaded with `godotenv` (key vars: `DB_DSN`, `CORS_ALLOWED_ORIGINS`, `RATE_LIMIT`)

## Security Middleware

HTTP security middleware is configured in `internal/workflow/presenter/rest_api/middleware.go` and wired in `rest_webserver.go`:

- **Request ID** (`gin-contrib/requestid`) — `X-Request-Id` header for tracing
- **Secure Headers** (`gin-contrib/secure`) — X-Frame-Options, X-Content-Type-Options, XSS filter, CSP, Referrer-Policy
- **CORS** (`gin-contrib/cors`) — configurable origins via `CORS_ALLOWED_ORIGINS` env var
- **Rate Limiting** (`ulule/limiter/v3`) — per-IP, configurable via `RATE_LIMIT` env var (default: `100-M`)

## CI

GitHub Actions (`.github/workflows/ci.yaml`) with 4 parallel jobs:

- **Lint** — gofmt + go vet + golangci-lint
- **Test** — tests with race detection + Codecov upload
- **Security** — govulncheck + gosec
- **Build** — compile Linux binary

Triggers on push to `main` and PRs to `main`, `feature/*`, `bugfix/*`. Run locally with `make ci`.

## API

All endpoints under `/api/v1`. REST handlers in `presenter/rest_api/`. OpenAPI spec at `docs/apidoc/apidoc.yaml`.
