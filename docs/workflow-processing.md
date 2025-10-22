# Task Processing Workflow

```mermaid
sequenceDiagram
    participant NimbusAPI
    participant DB
    participant SQS_Out
    participant EnergyLight
    participant SQS_In

    NimbusAPI->>DB: Decompose Task into Jobs
    NimbusAPI->>SQS_Out: Publish Job messages
    EnergyLight->>SQS_Out: Poll for Job messages
    EnergyLight->>EnergyLight: Process Job
    EnergyLight->>SQS_In: Publish results/errors
    NimbusAPI->>SQS_In: Poll for results
    NimbusAPI->>DB: Update Job status/results
    NimbusAPI->>DB: Aggregate Task status
```

**Description**
- Outbound queue decouples ingestion from processing latency.
- Results queue enables durable asynchronous completion reporting.
- Idempotency required in downstream calculations.
