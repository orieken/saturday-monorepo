# Trace & Metrics Queries for Quality Observability

This document contains useful TraceQL (for Tempo) and PromQL (for Prometheus using Span Metrics) queries to visualize test execution data.

## TraceQL (Tempo)

Use these queries in the **Explore > Tempo** view or inside Grafana dashboard panels using the Tempo data source.

### 1. Find Failed Tests
Find all spans representing a test case that failed.
```traceql
{ span.status = "error" && span.kind = "internal" } 
```

### 2. Slow Tests (> 5s)
Identify test cases taking longer than 5 seconds.
```traceql
{ duration > 5s }
```

### 3. Filter by Service (Test Suite)
Filter traces for a specific test reporter/service.
```traceql
{ resource.service.name = "cucumber-tests" }
# OR
{ resource.service.name = "playwright-tests" }
```

### 4. Failed Scenarios by Feature Name (Cucumber)
Find failed spans within a specific Cucumber feature using the newly added `test.feature.name` attribute.
```traceql
{ span.status = "error" && test.feature.name = "User Login" }
```

### 5. Find Playwright Tests by Suite (File)
Filter for tests belonging to a specific test file/suite using `test.suite`.
```traceql
{ test.suite = "tests/auth/login.spec.ts" }
```

---

## PromQL (Prometheus / Span Metrics)

We use the OTel Collector `spanmetrics` processor to generate metrics from traces.
**Dimensions Available**: `test_status`, `test_browser`, `test_platform`, `test_type`, `test_scenario_name`, `test_feature_name`, `test_case_title` (Playwright), `test_suite` (Playwright), `test_case_file`.

### 1. Test Pass Rate (Percentage) by Feature
Calculate the percentage of successful tests per Feature (Cucumber).
```promql
sum(rate(traces_span_metrics_calls_total{job="otel-collector", status_code="OK"}[5m])) by (test_feature_name)
/ 
sum(rate(traces_span_metrics_calls_total{job="otel-collector"}[5m])) by (test_feature_name) * 100
```

### 2. Failure Rate by Scenario / Case Name
Top 10 failing tests by individual test name. Matches both Cucumber Scenarios and Playwright Cases if you query by the specific label, or use `span_name` as a fallback.
```promql
# By Scenario Name (Cucumber)
topk(10, sum by (test_scenario_name) (rate(traces_span_metrics_calls_total{status_code="ERROR"}[5m])))

# By Case Title (Playwright)
topk(10, sum by (test_case_title) (rate(traces_span_metrics_calls_total{status_code="ERROR"}[5m])))
```

### 3. Average Test Duration (Latency) by Browser
Compare performance across browsers (e.g., Chromium vs Firefox).
```promql
rate(traces_span_metrics_latency_sum{}[5m]) by (test_browser)
/ 
rate(traces_span_metrics_latency_count{}[5m]) by (test_browser)
```

### 4. Total Execution Count by Suite
See which test suites are running the most.
```promql
sum(rate(traces_span_metrics_calls_total{service_name="playwright-tests"}[1m])) by (test_suite) * 60
```

### 5. Flaky Tests (High Variance)
Identify tests that fail frequently but also pass (Conceptual).
```promql
sum by (test_scenario_name) (traces_span_metrics_calls_total{status_code="ERROR"}) 
> 0 
and 
sum by (test_scenario_name) (traces_span_metrics_calls_total{status_code="OK"}) > 0
```

---

## Dashboard Panel Ideas

| Panel Title | Visualization | Query Type | Query |
| :--- | :--- | :--- | :--- |
| **Global Pass Rate** | Gauge | PromQL | `(sum(rate(...status="OK")) / sum(rate(...))) * 100` |
| **Pass Rate by Feature** | Bar Chart | PromQL | `... by (test_feature_name)` |
| **Failures by Browser** | Pie Chart | PromQL | `sum(traces_span_metrics_calls_total{status_code="ERROR"}) by (test_browser)` |
| **Slowest Scenarios** | Bar Gauge | PromQL | `topk(10, avg_duration by (test_scenario_name))` |
| **Test Throughput (Tests/min)** | Stat | PromQL | `sum(rate(traces_span_metrics_calls_total[1m])) * 60` |
| **Recent Failed Traces** | Table | TraceQL | `{ span.status = "error" }` |
