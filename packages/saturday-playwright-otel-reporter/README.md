# @orieken/saturday-playwright-otel-reporter

An OpenTelemetry reporter for Playwright, integrated with the Saturday Framework. It instruments your Playwright tests with distributed tracing and metrics.

## Features

- **Automatic Tracing**: Captures Test, Step, and Hook execution as OTel spans.
- **Hierarchical Trace Context**: Preserves the parent-child relationship of steps and fixtures within a test.
- **Status Reporting**: Accurately maps Playwright test results to OTel Span Status codes (`OK`, `ERROR`).
- **Debugging Tools**: Built-in support for debugging logs and payload inspection.

## Installation

```bash
pnpm add -D @orieken/saturday-playwright-otel-reporter
```

## Configuration

Add the reporter to your `playwright.config.ts`:

```typescript
import { defineConfig } from '@playwright/test';

export default defineConfig({
  reporter: process.env.ENABLE_OTEL === 'true' 
    ? [['list'], ['@orieken/saturday-playwright-otel-reporter']]
    : [['list']],
  // ...
});
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENABLE_OTEL` | Set to `true` to activate the reporter | `false` |
| `OTEL_SERVICE_NAME` | Name of the service | `playwright-tests` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP HTTP Endpoint (Traces) | (required) |
| `OTEL_CUSTOM_CONFIG` | Path to a custom configuration file. If not specified, looks for `playwright-otel.config.mjs` or `playwright-otel.config.js` in the project root. | - |
| `OTEL_SAVE_PAYLOADS` | If `true`, saves the generated OTel spans to `./reports/otel-playwright-spans-<TIMESTAMP>.json` | `false` |
| `OTEL_DEBUG_LOGGING` | Set to `true` for verbose console output during test execution | `false` |


## Custom Configuration

You can define a `playwright-otel.config.mjs` file in your project root to customize attributes:

```javascript
export default {
  resourceAttributes: {
    'service.version': '1.0.0',
    'service.runner': 'playwright'
  },
  // Custom logic for adding attributes to tests
  testAttributes: (test) => {
    return {
      'custom.test.tags': test.tags ? test.tags.map(t => t.tag).join(',') : ''
    };
  }
};
```

## Trace Attributes

The reporter automatically adds the following Resource Attributes:
- `test.reporter`: `playwright`
- `test.browser`: The browser name (if available)

## Metrics

This reporter automatically collects the following metrics:

- `playwright.test.cases` (Counter): Counts the number of executed test cases.
  - Labels: `test.status`, `test.browser`, `test.file`

## Sample Queries (PromQL)

**Total Pass/Fail Count per Browser:**
```promql
sum by (test_status, test_browser) (playwright_test_cases_total)
```

**Success Rate (%):**
```promql
sum(playwright_test_cases_total{test_status="passed"}) / sum(playwright_test_cases_total) * 100
```

