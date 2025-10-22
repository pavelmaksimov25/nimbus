# State Machine Diagrams

## Task State Machine
```mermaid
stateDiagram-v2
    [*] --> New
    New --> In_Progress: Jobs dispatched
    In_Progress --> Done: All jobs succeed
    In_Progress --> Failed: All jobs permanent failures
    In_Progress --> Partially_Done: Mixed success & permanent failures
    Partially_Done --> Done: Remaining retryable succeed
    Partially_Done --> Failed: Remaining retryable exhaust
    Done --> [*]
    Failed --> [*]
```

## Job State Machine
```mermaid
stateDiagram-v2
    [*] --> New
    New --> In_Progress: Picked for processing
    In_Progress --> Done: Success
    In_Progress --> Failed: Error captured
    Failed --> Retry_Pending: Transient error AND attempts < maxAttempts
    Retry_Pending --> In_Progress: Re-queued
    Failed --> Exhausted: Non-retryable OR attempts >= maxAttempts
    Done --> [*]
    Exhausted --> [*]
```

**Notes**
- `Retry_Pending` enables controlled reprocessing.
- `Exhausted` differentiates permanent vs transient failures for aggregation logic.
