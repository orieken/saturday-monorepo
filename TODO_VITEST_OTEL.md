# TODO: Vitest OpenTelemetry Reporter Implementation

This document outlines the tasks required to implement the `@orieken/saturday-vitest-otel-reporter`, ensuring consistency with the existing Cucumber and Playwright OTel reporters.

## 1. Project Setup
- [ ] **Initialize Package**: Create `packages/saturday-vitest-otel-reporter`.
- [ ] **Dependencies**: Install necessary `@opentelemetry/*` packages and `vitest` types.
- [ ] **Build Configuration**: Configure `tsup` to build CJS and ESM formats with type definitions.
- [ ] **CI/CD**: Add the new package to `.github/workflows/packages.yml`.

## 2. Core Architecture
- [ ] **Config Loader (`src/config-loader.ts`)**:
    - [ ] Implement `loadConfig` to support `vitest-otel.config.mjs`, `vitest-otel.config.js`, or `OTEL_CUSTOM_CONFIG`.
    - [ ] Define `OtelVitestConfig` interface.
- [ ] **Tracer Setup (`src/tracer-setup.ts`)**:
    - [ ] Initialize `NodeTracerProvider` with `OTEL_EXPORTER_OTLP_ENDPOINT`.
    - [ ] Implement `getResource()` logic.
- [ ] **Metrics Setup (`src/metrics-setup.ts`)**:
    - [ ] Initialize `MeterProvider`.
    - [ ] Create `vitest.test.cases` counter metric.
- [ ] **Span Manager (`src/span-manager.ts`)**:
    - [ ] Manage span lifecycle (Run -> Suite -> Test).
    - [ ] Handle parent-child context propagation (Suites can be nested).

## 3. Reporter Implementation (`src/index.ts`)
- [ ] **Class Structure**: Implement Vitest's `Reporter` interface (`onInit`, `onFinished`, `onTaskUpdate` or specific test hooks).
    - *Note: Check if Vitest's `onTaskUpdate` is the best place for granular updates or if there are explicit start/end hooks for tests.*
- [ ] **Initialization**: Include `OTEL_DEBUG_LOGGING` and `OTEL_SAVE_PAYLOADS`.
- [ ] **Event Handling**:
    - [ ] `onInit`: Start root span.
    - [ ] Task Start: Start span for File/Suite/Test.
    - [ ] Task End: End span, record status, record errors.
    - [ ] `onFinished`: Flush traces and metrics.
- [ ] **Metrics Integration**:
    - [ ] Increment `vitest.test.cases` counter on test completion.
    - [ ] Add labels: `test.status`, `test.file`.

## 4. Features & Consistency
- [ ] **Status Mapping**: Map Vitest statuses (pass, fail, skip, todo) to OTel `SpanStatusCode`.
- [ ] **Attribute Mapping**:
    - [ ] `resourceAttributes`: Custom service info.
    - [ ] `testAttributes`: Map tags/meta to OTel attributes.
- [ ] **Error Handling**: Graceful degradation.

## 5. Documentation
- [ ] **README.md**:
    - [ ] Installation instructions.
    - [ ] Configuration guide.
    - [ ] Sample `vitest-otel.config.mjs`.
    - [ ] Sample PromQL queries (`vitest_test_cases_total`).

## 6. Testing
- [ ] **Unit Tests**: Test the reporter logic using Vitest (dogfooding where possible or using Jest if preferred for isolation).
- [ ] **Integration Test**: Verify traces appear in Tempo and metrics in Prometheus.
