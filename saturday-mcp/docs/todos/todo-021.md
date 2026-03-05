# TODO-021: Observability Integration

## Status: 📋 Planning

**Objective**: Connect test generation and prioritization to production usage data. This allows the agent to answer questions like "What are the most critical flows?" or "Which tests should I write first?".

## Goals
1.  **Metrics Interface**: Define a generic interface for retrieving usage metrics (visits, error rates).
2.  **Usage Analyzer**: Implement an analyzer that ingests metrics and maps them to existing Page Objects.
3.  **Prioritization**: Add logic to rank Page Objects by importance.
4.  **Mock Implementation**: Create a predefined `metrics.json` loader to simulate a production environment (Google Analytics / Datadog).

## Components

### 1. `internal/observability/metrics.go`
- `MetricsProvider` interface.
- `PageMetric` struct (Path, Visits, ErrorRate).

### 2. `internal/analyzers/usage_analyzer.go`
- Logic to match specific URL paths (e.g., `/products/123`) to generalized Saturday Framework Paths (e.g., `/products/[id]`).

### 3. Tool: `prioritize_tests`
- **Input**: `projectPath`, `metricsFile`.
- **Output**: List of Pages/Flows ordered by risk/traffic.

## MVP Scope
We will not connect to a real live API like Datadog for this MVP. We will build the *adapter* pattern and a file-based implementation (`json`) to demonstrate the capability.
