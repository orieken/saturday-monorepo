# Trace & Metrics Queries for Quality Observability

This document contains useful TraceQL (for Tempo) and PromQL (for Prometheus using Span Metrics) queries to visualize test execution data.

## TraceQL (Tempo)

Use these queries in the **Explore > Tempo** view or inside Grafana dashboard panels using the Tempo data source.

### 1. Find Failed Tests
Find all spans representing a test case that failed.
```traceql
{ span.status = "error" && span.kind = "internal" } 
```
*Note: Adjust `span.kind` if your reporter uses a different kind for test cases.*

### 2. Slow Tests (> 5s)
Identify test cases taking longer than 5 seconds.
```traceql
{ duration > 5s }
```

### 3. Filter by Service (Test Suite)
Filter traces for a specific test reporter/service (e.g., `cucumber-tests` or `playwright-tests`).
```traceql
{ resource.service.name = "cucumber-tests" }
```

### 4. Failed Scenarios in a Specific Feature
Find failed spans within a specific Cucumber feature.
```traceql
{ span.status = "error" && resource.service.name = "cucumber-tests" && span.name =~ ".*Feature Name.*" }
```

---

## PromQL (Prometheus / Span Metrics)

If you are generating metrics from spans (using the OTel Collector `spanmetrics` processor), you can use these queries.
**Assumed Metric Names**: `traces_span_metrics_latency_bucket`, `traces_span_metrics_calls_total` (names may vary based on your collector config).

### 1. Test Pass Rate (Percentage)
Calculate the percentage of successful tests over time for a specific service.
```promql
sum(rate(traces_span_metrics_calls_total{service_name="cucumber-tests", status_code="OK"}[5m])) 
/ 
sum(rate(traces_span_metrics_calls_total{service_name="cucumber-tests"}[5m])) * 100
```

### 2. Failure Rate by Test Name
Top 10 failing tests by name.
```promql
topk(10, sum by (span_name) (rate(traces_span_metrics_calls_total{status_code="ERROR"}[5m])))
```

### 3. Average Test Duration (Latency)
Average duration of tests for a service.
```promql
rate(traces_span_metrics_latency_sum{service_name="cucumber-tests"}[5m]) 
/ 
rate(traces_span_metrics_latency_count{service_name="cucumber-tests"}[5m])
```

### 4. 95th Percentile Latency (P95)
The 95th percentile duration of tests (good for spotting slow outliers).
```promql
histogram_quantile(0.95, sum(rate(traces_span_metrics_latency_bucket{service_name="cucumber-tests"}[5m])) by (le))
```

### 5. Total Test Executions Count
Total number of tests executed per minute.
```promql
sum(rate(traces_span_metrics_calls_total{service_name="cucumber-tests"}[1m])) * 60
```

### 6. Flaky Tests (High Variance in Outcome)
*Conceptual*: Tests that have both Error and OK status in a short window.
This is harder to query directly in simple PromQL but you can visualize `status_code="ERROR"` vs `status_code="OK"` counts side-by-side for the same `span_name`.

---

## Dashboard Panel Ideas

| Panel Title | Visualization | Query Type | Query |
| :--- | :--- | :--- | :--- |
| **Global Pass Rate** | Gauge | PromQL | `(sum(rate(...status="OK")) / sum(rate(...))) * 100` |
| **Test Failures Over Time** | Time Series | PromQL | `sum(rate(...status="ERROR"))` |
| **Slowest Tests (Top 10)** | Bar Gauge | PromQL | `topk(10, avg_duration)` |
| **Test Duration Distribution** | Heatmap | PromQL | `rate(...latency_bucket...)` |
| **Recent Failed Traces** | Table / Trace List | TraceQL | `{ span.status = "error" }` |
