# @orieken/saturday-cucumber-otel-formatter

A custom CucumberJS formatter that integrates with the Saturday Framework to enhance observability. It sends test execution traces and metrics directly to OpenTelemetry (OTLP) endpoints.

## Features

- **OpenTelemetry Tracing**: Captures detailed traces for Test Runs, Scenarios, and Steps.
- **Hierarchical Context**: Steps are child spans of Scenarios, which are child spans of the Test Run.
- **Status Reporting**: Accurately reports Pass/Fail/Error statuses on spans.
- **Duration Tracking**: Precise timing info for all execution units.

## Installation

```bash
pnpm add @orieken/saturday-cucumber-otel-formatter
```

## Configuration

The formatter uses standard OpenTelemetry environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `OTEL_SERVICE_NAME` | Name of the service | `cucumber-tests` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP HTTP Endpoint | (required) |
| `OTEL_SERVICE_VERSION` | Service version | - |
| `OTEL_SERVICE_NAMESPACE` | Service namespace | - |
| `OTEL_DEFAULT_ENVIRONMENT` | Deployment environment | - |
| `OTEL_CUSTOM_CONFIG` | Path to a custom configuration file (e.g. `./otel-config.mjs`) for advanced attributes | - |
| `OTEL_DEBUG_LOGGING` | Set to `true` to enable verbose logging of active handles and initialization steps for debugging hangs. | `false` |

## Robustness Features

This formatter implements several reliability features to ensure data integrity:
- **Asynchronous Initialization**: Configuration is loaded asynchronously to support remote or complex setups.
- **Graceful Shutdown**: The formatter explicitly awaits the flushing of all OpenTelemetry spans before exiting, preventing data loss at the end of test runs.
- **Queueing**: Events are queued during initialization so no trace data is lost during startup.

## Usage

Register the formatter when running Cucumber-JS:

```bash
# Example running with environment variables
export OTEL_SERVICE_NAME="saturday-e2e-tests"
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318/v1/traces"

npx cucumber-js \
  --format @orieken/saturday-cucumber-otel-formatter \
  --require features/**/*.ts
```

## Trace Structure

*   `test-run` (Root Span)
    *   `scenario: <Scenario Name>` (Child Span)
        *   `step: <Keyword> <Step Text>` (Child Span)

