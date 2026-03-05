---
description: Understanding the flow of Custom Metrics and Span Metrics from code to Grafana
---
# Workflow: OTEL Metrics Pipeline Debugging

This workflow helps you trace where metrics might be dropping in the OpenTelemetry pipeline.

## 1. Verify Code Instrumentation

Ensure your Go/TS code is actually emitting metrics.
```bash
# Check OTel debug logs in the application container
cat otel-debug.log | grep "Measurement recorded"
```

## 2. Check OTel Collector (Sidecar/Agent)

The application sends metrics to the Collector usually on `localhost:4317` (gRPC) or `4318` (HTTP).
Check if the Collector is receiving them.
```bash
# Check Collector logs for ZPages or debug exporters
docker logs saturday-otel-collector
```

## 3. Prioritize Tests based on Usage (New)

If you have production metrics available (e.g. from Prometheus/Datadog exported as JSON), use various saturday tools to find what to test.

1.  Export metrics to `metrics.json`.
2.  Run `prioritize_tests`:
    ```json
    {
      "metricsFile": "/path/to/metrics.json"
    }
    ```
3.  The agent will recommend creating tests for high-traffic/high-error pages.
