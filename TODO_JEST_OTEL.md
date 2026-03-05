# TODO: Jest OpenTelemetry Reporter Implementation

This document outlines the tasks required to implement the `@orieken/saturday-jest-otel-reporter`, ensuring consistency with the existing Cucumber and Playwright OTel reporters.

## 1. Project Setup
- [ ] **Initialize Package**: Create `packages/saturday-jest-otel-reporter`.
- [ ] **Dependencies**: Install necessary `@opentelemetry/*` packages (api, core, sdk-trace-base, resources, semantic-conventions, etc.) and `jest` types.
- [ ] **Build Configuration**: Configure `tsup` to build CJS and ESM formats with type definitions.
- [ ] **CI/CD**: Add the new package to `.github/workflows/packages.yml`.

## 2. Core Architecture
- [ ] **Config Loader (`src/config-loader.ts`)**:
    - [ ] Implement `loadConfig` to support `jest-otel.config.mjs`, `jest-otel.config.js`, or `OTEL_CUSTOM_CONFIG`.
    - [ ] Define `OtelJestConfig` interface (resource attributes, test attributes).
- [ ] **Tracer Setup (`src/tracer-setup.ts`)**:
    - [ ] Initialize `NodeTracerProvider` with `OTEL_EXPORTER_OTLP_ENDPOINT`.
    - [ ] Implement `getResource()` to handle service name and other resource attributes.
- [ ] **Metrics Setup (`src/metrics-setup.ts`)**:
    - [ ] Initialize `MeterProvider` with `PeriodicExportingMetricReader`.
    - [ ] Create `jest.test.cases` counter metric.
- [ ] **Span Manager (`src/span-manager.ts`)**:
    - [ ] Manage span lifecycle (Run -> Suite -> Test).
    - [ ] Handle parent-child context propagation.

## 3. Reporter Implementation (`src/index.ts`)
- [ ] **Class Structure**: Implement Jest's `Reporter` interface (`onRunStart`, `onTestStart`, `onTestResult`, `onRunComplete`).
- [ ] **Initialization**: Include `OTEL_DEBUG_LOGGING` and `OTEL_SAVE_PAYLOADS` support.
- [ ] **Event Handling**:
    - [ ] `onRunStart`: Start root span.
    - [ ] `onTestStart`: Start test span.
    - [ ] `onTestResult`: End test span, record status, record exception if failed.
    - [ ] `onRunComplete`: End root span, flush traces and metrics.
- [ ] **Metrics Integration**:
    - [ ] Increment `jest.test.cases` counter on test completion.
    - [ ] Add labels: `test.status`, `test.file`.

## 4. Features & Consistency
- [ ] **Status Mapping**: Map Jest statuses (passed, failed, skipped, pending, todo) to OTel `SpanStatusCode`.
- [ ] **Attribute Mapping**:
    - [ ] `resourceAttributes`: Custom service info.
    - [ ] `testAttributes`: Map Jest test names and file paths to custom attributes.
- [ ] **Error Handling**: Graceful degradation if OTel endpoint is unreachable.

## 5. Documentation
- [ ] **README.md**:
    - [ ] Installation instructions.
    - [ ] Configuration guide (Env vars + Config file).
    - [ ] Sample `jest-otel.config.mjs`.
    - [ ] Sample PromQL queries (`jest_test_cases_total`).

## 6. Testing
- [ ] **Unit Tests**: Test the reporter logic using Jest (mocking OTel components).
- [ ] **Integration Test**: Verify traces appear in Tempo and metrics in Prometheus.
