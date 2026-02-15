# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.25
**Primary Dependencies**: Gin v1.11, GORM v1.31, goose, godotenv
**Storage**: PostgreSQL 16
**Testing**: go test + testify + uber/mock (co-located *_test.go)
**Target Platform**: Linux server (Docker), macOS dev
**Project Type**: Single Go module with Clean Architecture
**Performance Goals**: [NEEDS CLARIFICATION]
**Constraints**: [NEEDS CLARIFICATION]
**Scale/Scope**: [NEEDS CLARIFICATION]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[Gates determined based on constitution file]

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
cmd/app/main.go                                    # DI wiring
internal/{module}/
├── domain/entity/entity.go                        # GORM entities
├── domain/repository/interfaces.go                # Repo interfaces
├── domain/service/interfaces.go                   # Service interfaces
├── domain/types/errors.go                         # Custom error types
├── application/{resource}/service.go              # Service impl
├── application/{resource}/service_test.go         # Co-located tests
├── adapters/repository/{resource}/                # GORM repos
├── adapters/mocks/                                # Generated mocks
├── presenter/rest_api/{resource}_handler.go       # Gin handlers
└── presenter/rest_api/rest_webserver.go           # Route registration
migrations/00NNN_description.sql                   # goose migrations
```

**Structure Decision**: Clean Architecture with domain modules under `internal/`.
Currently one domain module exists (`internal/workflow/`), with `internal/task_runner/`
scaffolded. New features add to existing modules or create new ones following the
same four-layer pattern.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 3rd domain module] | [current need] | [why existing modules insufficient] |
| [e.g., additional adapter type] | [specific problem] | [why GORM repo insufficient] |
