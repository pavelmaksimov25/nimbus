# Nimbus Constitution

## Core Principles

### I. Clean Architecture Layers
Four-layer structure per domain module: domain, application, adapters, presenter. Dependencies flow inward only. The domain layer has zero external dependencies. Each layer has a single, well-defined responsibility: domain defines entities and interfaces, application implements business logic, adapters provide infrastructure, and presenter handles HTTP concerns.

### II. Interface-Based Dependency Injection
All cross-layer dependencies defined as interfaces in `domain/`. Wiring performed in `cmd/app/main.go` via constructors. No service locator, no global state, no init() side effects. Every dependency is explicit and passed through constructor parameters: `NewXxxService(repo XxxRepository)`.

### III. Test-First with Co-located Tests
Tests live alongside source files (`service_test.go` next to `service.go`). Use `testify` for assertions and `uber/mock` for generated mocks. Naming convention: `TestAction_Scenario` (e.g., `TestCreateTask_WhenValidInput`). Mock repositories via interfaces defined in `domain/repository/`.

### IV. Custom Error Types with Unwrap
Domain errors defined in `domain/types/` implementing `Error()` and `Unwrap()` methods. Sentinel errors enable `errors.Is()` matching across layers (e.g., `types.ErrNotFound`). Application layer wraps domain errors; presenter layer maps them to HTTP status codes. No naked string errors in business logic.

### V. Convention over Configuration
`snake_case.go` file names. Constructor naming: `NewXxxRepository(db)`, `NewXxxService(repo)`. Goose-annotated SQL migrations with sequential numbering in `migrations/` directory. Configuration via `.env` file loaded with `godotenv`. Environment variables: `DB_DSN`, `CORS_ALLOWED_ORIGINS`, `RATE_LIMIT`.

## Technology Stack Constraints

- **Language**: Go 1.25+
- **HTTP Framework**: Gin
- **ORM**: GORM with PostgreSQL 16
- **Migrations**: goose (SQL-based, sequential numbering)
- **Configuration**: godotenv (`.env` file)
- **Mocking**: uber/mock (generated mocks in `domain/repository/mocks/`)
- **Assertions**: testify
- **Security Middleware**: gin-contrib/requestid, gin-contrib/secure, gin-contrib/cors, ulule/limiter

## Development Workflow

- Run `make lint` before every commit (gofmt + go vet + golangci-lint)
- Run `make test` with race detection (`-race` flag) and coverage
- Run `make security` for vulnerability scanning (govulncheck + gosec + staticcheck)
- CI pipeline: 4 parallel jobs (Lint, Test, Security, Build)
- GitHub Actions triggers on push to `main` and PRs to `main`, `feature/*`, `bugfix/*`
- Run full CI locally with `make ci`

## Governance

This constitution supersedes all other development practices and conventions. Amendments require a pull request with clear rationale documenting what changed and why. All code reviews must verify compliance with these principles. Complexity beyond what the constitution prescribes must be justified in writing.

**Version**: 1.0.0 | **Ratified**: 2026-02-15 | **Last Amended**: 2026-02-15
