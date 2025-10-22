# Nimbus

Resilient load-leveling facade for asynchronous, high-latency energy data computations.

Nimbus decomposes portfolio requests into per‑property Jobs, manages distributed asynchronous processing via SQS queues and the Energy Light service, applies controlled retries for transient failures, and exposes aggregated task status and results through a polling API.

---

## Documentation Index

| Topic | File |
|-------|------|
| Entity Model (Task & Job) | `docs/entity-schema.md` |
| New Task Workflow | `docs/workflow-new-task.md` |
| Processing Workflow | `docs/workflow-processing.md` |
| Retrieval Workflow | `docs/workflow-retrieval.md` |
| System Design Overview | `docs/system-design.md` |
| State Machines (Task & Job) | `docs/state-machines.md` |

---

## Overview

### Core Responsibilities
1. Accept and validate batch portfolio requests.
2. Persist Task and decomposed per-property Jobs.
3. Queue Jobs for downstream energy computations.
4. Ingest asynchronous results and update Job + Task aggregates.
5. Provide a simple polling contract for clients (Location header pattern).

### Key Concepts
- **Task**: Aggregates overall progress & status for a portfolio request.
- **Job**: Atomic unit: energy calculation for a single property.
- **Retry Count (`retryCount`)**: Incremented on each transient/system retry attempt.
- **DLQ**: Ensures isolation of poison messages and forensic analysis.

### Technology Constraints
- Requires highly concurrent durable storage (e.g. PostgreSQL / DynamoDB).
- SQLite is unsuitable due to single-writer limitations.
- Downstream Energy Light processing must be idempotent (at-least-once delivery).

---

## API Summary

| Endpoint | Method | Purpose | Response |
|----------|--------|---------|----------|
| `/api/tasks` | POST | Submit new portfolio task | `202 Accepted` + `Location` or `400 Bad Request` |
| `/api/tasks/{taskID}` | GET | Poll task progress/results | `200 OK`, final results or `303 See Other` |

---

## Resilience & Retry Strategy

Jobs failing due to transient/system errors may be retried (bounded attempts). Permanent data quality issues (incomplete property data) do not retry and contribute to partial or failed aggregate states. DLQs configured on both processing and results queues prevent snowball failure loops.

---

## Future Improvements

- External workflow engine (Temporal / Step Functions) for declarative retries & visibility.
- Structured error taxonomy for analytics & SLA monitoring.
- Enhanced metrics: retry rate, permanent failure segmentation, time-to-first-result.

---

## License

Add license information here.

<!-- Diagrams moved to /docs. See Documentation Index above. -->