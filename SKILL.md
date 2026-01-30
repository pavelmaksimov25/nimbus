---
name: nimbus-patterns
description: Coding patterns extracted from the nimbus workflow engine repository
version: 1.0.0
source: local-git-analysis
analyzed_commits: 26
---

# Nimbus Patterns

## Commit Conventions

The project uses **conventional commit prefixes** (lowercase, colon-separated). Based on 26 commits:

| Prefix | Count | Usage |
|--------|-------|-------|
| `feat:` | 9 (35%) | New features and endpoints |
| `chore:` | 9 (35%) | Internal work, tests, maintenance |
| `doc:` | 3 (12%) | Documentation updates |
| `refactor:` | 2 (8%) | Renaming, restructuring |
| `fix:` | 1 (4%) | Bug fixes |

Remaining commits used unstructured messages (`init`, `Database integration (#1)`).

**Convention:** `<type>: <short description in lowercase>`. No scope parentheses. No capitalization after prefix. Examples:
- `feat: introduced startTask action`
- `chore: task service tests`
- `doc: workflows apidoc`
- `refactor: renamed store to storage`

## Code Architecture

Clean Architecture with domain modules under `internal/`:

```
internal/{module}/
├── domain/
│   ├── entity/       # Business entities (GORM-annotated structs)
│   ├── repository/   # Repository interfaces + mock implementations
│   ├── service/      # Service interfaces
│   └── types/        # Custom error types
├── application/
│   └── {resource}/   # Service implementations (one per resource)
├── adapters/
│   └── repository/   # GORM-backed repository implementations
└── presenter/
    └── rest_api/     # Gin HTTP handlers and router
```

**Naming conventions:**
- Entities: singular nouns (`Task`, `Workflow`)
- Repositories: `{Resource}Repository` interface + `New{Resource}Repository` constructor
- Services: `{Resource}Service` interface + `New{Resource}Service` constructor
- Handlers: `{resource}_handler.go` with methods on a handler struct
- Files: `snake_case.go`

## File Co-Change Patterns

Files that consistently change together (high coupling):

### Feature Addition (all 4 layers touched)
When adding a new resource operation, these files change together:
1. `domain/entity/entity.go` — Add/modify entity fields or status types
2. `application/{resource}/service.go` — Implement business logic
3. `application/{resource}/service_test.go` — Add test cases
4. `presenter/rest_api/{resource}_handler.go` — Add HTTP handler

This pattern appeared in **8 out of 26 commits** (31%).

### New Resource Workflow
When introducing a new resource type (e.g., workflows):
1. `domain/entity/entity.go` — Define entity struct
2. `adapters/storage/` or `adapters/repository/` — Storage implementation
3. `application/{resource}/service.go` — Service implementation
4. `presenter/rest_api/{resource}_handler.go` — Handler
5. `presenter/rest_api/rest_webserver.go` — Register routes
6. `cmd/app/main.go` — Wire dependencies

### Test Changes
`service_test.go` changed in **9 out of 26 commits** (35%), always alongside `service.go`. Tests are co-located with implementation in the same package.

## Workflows

### Adding a New Endpoint to an Existing Resource
1. Define or update service interface in `domain/service/interfaces.go`
2. If new repo method needed, update `domain/repository/interfaces.go`
3. Implement repo method in `adapters/repository/{resource}/`
4. Regenerate mocks if repo interface changed (`go.uber.org/mock`)
5. Implement service method in `application/{resource}/service.go`
6. Write tests in `application/{resource}/service_test.go`
7. Add HTTP handler in `presenter/rest_api/{resource}_handler.go`
8. Register route in `presenter/rest_api/rest_webserver.go`

### Adding a New Domain Module
1. Create directory structure under `internal/{module}/` following the 4-layer pattern
2. Define entities in `domain/entity/`
3. Define repository and service interfaces in `domain/`
4. Implement adapters and application logic
5. Create presenter handlers
6. Wire everything in `cmd/app/main.go`
7. Create SQL migration files via `make migrate-create name=<name>`

### Database Migration
1. `make migrate-create name=<descriptive_name>` — Creates up/down SQL files
2. Edit `migrations/NNNNNN_<name>.up.sql` and `.down.sql`
3. `make migrate` — Apply migration
4. Down migration must fully reverse the up migration

## Testing Patterns

- **Framework:** `testify/assert` + `go.uber.org/mock`
- **Location:** Tests co-located with implementation (`service_test.go` next to `service.go`)
- **Style:** Table-driven subtests are not used; each test case is a separate `Test*` function
- **Naming:** `Test{Action}_{Scenario}` (e.g., `TestStartTask_NotFound`, `TestCompleteTask_Success`)
- **Mocking:** Mock interfaces generated via uber/mock in `domain/repository/mocks/`
- **Coverage:** Covers happy path + error cases for each service method
- **Run:** `make test` (race detection + coverage) or `go test -v -run TestName ./path/`

## Dependency Injection Pattern

Constructor-based injection wired manually in `cmd/app/main.go`:

```
DB connection → Repository(db) → Service(repo) → Handler registration(service)
```

No DI framework. Each layer depends only on interfaces from the domain layer.

## Error Handling Pattern

Custom error types in `domain/types/` with sentinel errors:
- `ErrNotFound` + `RecordNotFoundError{Resource, ID, Msg}`
- `ErrUnpocessable` + `UnprocessableEntityError{Resource, ID, Msg}`

Error types implement `Unwrap()` for `errors.Is()` matching. HTTP handlers use `errors.Is()` to map domain errors to HTTP status codes (404, 422).

## Key Technologies

| Component | Technology |
|-----------|-----------|
| Language | Go 1.24 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Migrations | golang-migrate (SQL files) |
| Mocking | go.uber.org/mock |
| Assertions | testify |
| Config | godotenv (.env files) |
| IDs | google/uuid |
