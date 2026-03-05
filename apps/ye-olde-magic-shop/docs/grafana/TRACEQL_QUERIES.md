# Useful TraceQL Queries for Cucumber OpenTelemetry Analysis

## Quick Reference Guide

This document contains ready-to-use TraceQL queries for analyzing your Cucumber test data in Grafana/Tempo.

---

## Basic Queries

### Find All Test Runs
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="test-run"}
```

### Find All Scenarios
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"}
```

### Find All Steps
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"}
```

---

## Filtering by Status

### All Failed Scenarios
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.status="error"}
```

### All Passed Scenarios
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.status="ok"}
```

### Failed Steps Only
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.status="failed"}
```

### Skipped Steps
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.status="skipped"}
```

---

## Feature-Based Queries

### Scenarios from Specific Feature
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.custom.feature.name="Login"}
```

### All Failures in a Feature
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.custom.feature.name="Checkout"} && {span.test.status="error"}
```

### Feature Performance Comparison
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| quantile_over_time(duration, 0.95) by (span.custom.feature.name)
```

---

## Tag-Based Queries

### Smoke Tests Only
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.custom.scenario.tags=~".*@smoke.*"}
```

### Regression Tests
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.custom.scenario.tags=~".*@regression.*"}
```

### Happy Path Tests
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.custom.scenario.tags=~".*@happy_path.*"}
```

### Tests with Multiple Tags
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.custom.scenario.tags=~".*@smoke.*"} && {span.custom.scenario.tags=~".*@critical.*"}
```

---

## Environment & Browser Queries

### Tests Run in Development
```traceql
{resource.service.name="ye-olde-magic-shop"} && {resource.test.environment="development"} && {span.test.type="scenario"}
```

### Chrome-Specific Tests
```traceql
{resource.service.name="ye-olde-magic-shop"} && {resource.test.browser="chrome"} && {span.test.type="scenario"}
```

### Firefox Failures
```traceql
{resource.service.name="ye-olde-magic-shop"} && {resource.test.browser="firefox"} && {span.test.type="scenario"} && {span.test.status="error"}
```

### Cross-Browser Comparison
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (resource.test.browser, span.test.status)
```

### Production Environment Tests
```traceql
{resource.service.name="ye-olde-magic-shop"} && {resource.test.environment="production"} && {span.test.type="scenario"}
```

---

## Performance Analysis

### Slowest Scenarios (Top 10)
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && duration > 5s
```

### Find Scenarios Taking Over 10 Seconds
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && duration > 10s
```

### P95 Duration by Environment
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| quantile_over_time(duration, 0.95) by (resource.test.environment)
```

### Average Step Duration by Keyword
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} 
| avg(duration) by (span.test.step.keyword)
```

### Fastest vs Slowest Features
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| avg(duration) by (span.custom.feature.name)
```

---

## Error Analysis

### Most Common Error Messages
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.status="error"} 
| rate() by (span.test.error)
```

### Scenarios with Specific Error
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.error=~".*timeout.*"}
```

### Assertion Failures
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.error=~".*AssertionError.*"}
```

### Network-Related Failures
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.error=~".*(network|connection|timeout).*"}
```

---

## Step-Level Analysis

### Failed "When" Steps
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.step.keyword="When"} && {span.test.status="failed"}
```

### All "Given" Steps Performance
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.step.keyword="Given"} 
| quantile_over_time(duration, 0.95)
```

### Steps Matching Pattern
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.step.text=~".*login.*"}
```

### Most Frequently Failing Step
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="step"} && {span.test.status="failed"} 
| rate() by (span.test.step.name)
```

---

## Time-Based Analysis

### Tests Run Today
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="test-run"}
```
*Note: Use Grafana's time range selector to set to "Today"*

### Tests Run in Last Hour
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"}
```
*Note: Use Grafana's time range selector to set to "Last 1 hour"*

### Success Rate Trend
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (span.test.status)
```

---

## Flaky Test Detection

### Scenarios That Sometimes Fail
Find scenarios with both passed and failed executions in the time range:
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (span.test.scenario.name, span.test.status)
```
*Look for scenarios appearing in both "ok" and "error" status*

### Scenarios with Recent Failures
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.status="error"}
```
*Group results by scenario name to find repeat offenders*

---

## Aggregation Queries

### Total Scenario Count
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} | count()
```

### Success Rate Percentage
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (span.test.status)
```

### Average Scenarios per Test Run
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} | count()
```
*Divide by number of test runs in the time range*

### Test Execution Frequency
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="test-run"} | rate()
```

---

## Multi-Dimensional Grouping

### Status by Feature and Environment
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (span.custom.feature.name, resource.test.environment, span.test.status)
```

### Browser and Platform Matrix
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (resource.test.browser, resource.test.platform, span.test.status)
```

### Tag and Status Breakdown
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| rate() by (span.custom.scenario.tags, span.test.status)
```

---

## Advanced Queries

### Scenarios with No Steps (Potential Issues)
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {name!~".*step.*"}
```

### Long-Running Test Runs
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="test-run"} && duration > 5m
```

### Scenarios from Specific File
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.scenario.file=~".*login.feature"}
```

### Failed Scenarios Grouped by File
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.status="error"} 
| rate() by (span.test.scenario.file)
```

### Distribution of Test Status
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| histogram(span.test.status)
```

---

## Custom Attribute Queries

### Query by Service Instance
```traceql
{resource.service.instance.id="dev-instance-1"} && {span.test.type="scenario"}
```

### Query by Platform Type
```traceql
{resource.service.name="ye-olde-magic-shop"} && {resource.test.platform="darwin"} && {span.test.type="scenario"}
```

### Query by OS Type
```traceql
{resource.service.name="ye-olde-magic-shop"} && {resource.test.os="Darwin"} && {span.test.type="scenario"}
```

---

## Combining Filters

### Smoke Tests in Chrome that Failed
```traceql
{resource.service.name="ye-olde-magic-shop"} 
&& {resource.test.browser="chrome"} 
&& {span.test.type="scenario"} 
&& {span.custom.scenario.tags=~".*@smoke.*"} 
&& {span.test.status="error"}
```

### Production Environment, Firefox, Failed "When" Steps
```traceql
{resource.service.name="ye-olde-magic-shop"} 
&& {resource.test.environment="production"} 
&& {resource.test.browser="firefox"} 
&& {span.test.type="step"} 
&& {span.test.step.keyword="When"} 
&& {span.test.status="failed"}
```

### Specific Feature, Specific Tag, Performance Analysis
```traceql
{resource.service.name="ye-olde-magic-shop"} 
&& {span.custom.feature.name="Checkout"} 
&& {span.custom.scenario.tags=~".*@regression.*"} 
&& {span.test.type="scenario"} 
| quantile_over_time(duration, 0.95)
```

---

## Export Queries

### Export All Failed Scenario Details
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} && {span.test.status="error"}
```
*Use Grafana's table view and export to CSV*

### Export Performance Metrics by Feature
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
| avg(duration) by (span.custom.feature.name)
```
*Export to analyze in spreadsheet*

---

## Tips for Writing TraceQL Queries

1. **Always filter by service name first** to improve query performance
2. **Use regex carefully** - they're powerful but can be slow
3. **Combine specific filters** before broader aggregations
4. **Use proper span types** - test-run vs scenario vs step
5. **Remember attribute namespaces**:
   - `resource.*` for resource attributes
   - `span.*` for span attributes
6. **Test incrementally** - build complex queries step by step
7. **Use appropriate aggregations**:
   - `rate()` for counts over time
   - `quantile_over_time()` for percentiles
   - `avg()`, `min()`, `max()` for simple stats
8. **Consider time ranges** - shorter ranges = faster queries

---

## Common Patterns

### Pattern: Find Outliers
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"} 
&& duration > 3 * avg(duration)
```

### Pattern: Compare Before/After
Run same query with different time ranges:
- Time Range 1: Last week
- Time Range 2: This week
- Compare results manually or in spreadsheet

### Pattern: Drill Down
1. Start broad: All scenarios
2. Filter by status: Failed only
3. Group by feature: Which feature?
4. Filter by feature: Focus on one
5. View steps: What failed?

### Pattern: Identify Regressions
```traceql
{resource.service.name="ye-olde-magic-shop"} && {span.test.type="scenario"}
| quantile_over_time(duration, 0.95)
```
Compare P95 week-over-week

---

## Query Templates

Replace placeholders with your values:

```traceql
{resource.service.name="YOUR-SERVICE"} 
&& {resource.test.environment="YOUR-ENV"} 
&& {span.test.type="TYPE"} 
&& {span.custom.feature.name="FEATURE"} 
&& {span.test.status="STATUS"}
```

Common replacements:
- `YOUR-SERVICE`: ye-olde-magic-shop
- `YOUR-ENV`: development, staging, production
- `TYPE`: test-run, scenario, step
- `FEATURE`: Your feature name
- `STATUS`: ok, error, passed, failed, skipped

---

**Last Updated**: December 2025
**Version**: 1.0