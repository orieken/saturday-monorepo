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
| `OTEL_CUSTOM_CONFIG` | Path to a custom configuration file. If not specified, looks for `cucumber-otel.config.mjs` or `cucumber-otel.config.js` in the project root. | - |
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


## Custom Configuration

You can define a `cucumber-otel.config.mjs` file in your project root to customize resource and scenario attributes:

```javascript
export default {
  resourceAttributes: {
    'service.version': '1.0.0',
    'service.runner': 'cucumber'
  },
  // Custom logic for adding attributes to scenarios
  scenarioAttributes: (pickle, gherkinDocument) => {
    return {
      'custom.scenario.tags': pickle.tags.map(t => t.name).join(',')
    };
  },
  // Custom logic for adding attributes to steps
  stepAttributes: (pickleStep, gherkinStep) => {
    return {
      'custom.step.keyword': gherkinStep.keyword
    };
  }
};
```

## Metrics

In addition to traces, this formatter automatically collects the following metrics:

- `cucumber.test.cases` (Counter): Counts the number of executed test cases.
  - Labels: `test.status` ('passed' | 'error'), `test.file`

## Sample Queries (PromQL)

Here are some example queries you can use in Grafana (Prometheus datasource) to visualize your test data:

**Total Pass/Fail Count per Feature:**
```promql
sum by (test_status, test_file) (cucumber_test_cases_total)
```

**Success Rate (%):**
```promql
sum(cucumber_test_cases_total{test_status="passed"}) / sum(cucumber_test_cases_total) * 100
```


## Debugging & Diagnostics

| Variable | Description | Default |
|----------|-------------|---------|
| `OTEL_SAVE_PAYLOADS` | If `true`, saves the exact JSON payload sent to the OTel collector to `./reports/otel-cucumber-spans-<TIMESTAMP>.json` | `false` |
| `OTEL_DEBUG_LOGGING` | Set to `true` to enable verbose logging to console and `otel-debug-init.log` | `false` |

## Trace Structure

*   `test-run` (Root Span)
    *   `scenario: <Scenario Name>` (Child Span)
        *   `step: <Keyword> <Step Text>` (Child Span)

