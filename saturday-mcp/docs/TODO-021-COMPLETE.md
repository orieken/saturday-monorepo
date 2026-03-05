# TODO-021: Observability Integration

## Status: ✅ Complete

We have implemented the ability to ingest production usage metrics and prioritize test efforts accordingly.

## Implemented Components

### 1. Observability Package (`internal/observability`)
- `MetricsProvider` interface for fetching data.
- `FileMetricsProvider` for reading local JSON dumps (MVP).
- `PageMetric` model.

### 2. Usage Analyzer (`internal/analyzers/usage_analyzer.go`)
- **Scoring**: Ranks pages by `Visits * (1 + ErrorRate)`.
- **Logic**: Identifies high-traffic and high-risk areas.

### 3. Tool: `prioritize_tests`
- **Input**: `metricsFile` (path to json).
- **Output**: Ranked list of pages with explanations (e.g., "High traffic, High error rate").

## Next Steps

To make this fully production-ready, we would:
1.  Add a `DatadogMetricsProvider` or `PrometheusMetricsProvider`.
2.  Connect `UsageAnalyzer` to `GraphAnalyzer` to map generic URLs (like `/products/123`) to `Page` matches in the Knowledge Graph.
