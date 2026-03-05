# TODO: End-to-End Observability Testing Plan

This document outlines the plan to create an E2E test suite that verifies our OpenTelemetry reporters are correctly sending data to our observability stack (Prometheus, Tempo, Grafana).

## 1. Project Structure
- [ ] **Create Directory**: `tests/observability` at the monorepo root.
- [ ] **Structure**:
  ```
  tests/observability/
  ├── fixtures/              # "Dummy" tests to generate known telemetry data
  │   ├── cucumber/          # Minimal Cucumber setup
  │   │   ├── feature.feature
  │   │   └── steps.ts
  │   └── playwright/        # Minimal Playwright setup
  │       └── test.spec.ts
  ├── src/                   # Helper code for the verification suite
  │   ├── runner.ts          # Logic to run the fixtures
  │   ├── prometheus-api.ts  # Client to query Prometheus
  │   └── tempo-api.ts       # Client to query Tempo
  └── verify.spec.ts         # The actual test file (using Vitest or Jest)
  ```

## 2. Implement Test Fixtures
These are lightweight tests designed specifically to produce "known" outcomes (passes, failures, tags) that we can assert on later.

- [ ] **Cucumber Fixture**:
    - [ ] Create a feature with:
        - One passing scenario.
        - One failing scenario.
        - Scenarios with tags (`@e2e-test`, `@metrics`).
    - [ ] Configure `cucumber-otel.config.mjs` for this fixture to inject unique attributes (e.g., `fixture.id`).
- [ ] **Playwright Fixture**:
    - [ ] Create a spec file with:
        - One passing test.
        - One failing test (expect(true).toBe(false)).
        - Nested steps (`test.step(...)`).
        - Annotations/Tags.
    - [ ] Configure `playwright-otel.config.mjs` for this fixture.

## 3. Implement Verification Utilities
We need tools to query the observability backend to see if the data arrived.

- [ ] **Prometheus Client**:
    - [ ] Implement `queryMetric(query: string)` fetching from `http://localhost:50081/api/v1/query`.
    - [ ] Helper to parse vector results.
- [ ] **Tempo Client**:
    - [ ] Implement `searchTraces(tags: string)` fetching from `http://localhost:50061/api/search` (or TraceQL endpoint).
    - [ ] Helper to check span hierarchy (Root -> Test -> Step).

## 4. Implement the "Meta-Test" Suite
This is the test that runs the fixtures and asserts on the observability data.

- [ ] **Test Setup (`beforeAll`)**:
    - [ ] Ensure the Observability Docker Stack is running (`docker-compose up`).
    - [ ] Run the **Cucumber Fixture** subprocess (with `ENABLE_OTEL=true`).
    - [ ] Run the **Playwright Fixture** subprocess (with `ENABLE_OTEL=true`).
    - [ ] Wait for a flush interval (e.g., 5 seconds) to ensure collectors processed the data.
- [ ] **Verify Metrics**:
    - [ ] Assert `cucumber_test_cases_total{test_status="passed"}` increases by 1.
    - [ ] Assert `cucumber_test_cases_total{test_status="error"}` increases by 1.
    - [ ] Assert `playwright_test_cases_total{test_status="passed"}` increases by 1.
    - [ ] Verify labels (e.g., `test.file` matches fixture path).
- [ ] **Verify Traces**:
    - [ ] Find traces for `service.name=cucumber-fixture`.
    - [ ] Find traces for `service.name=playwright-fixture`.
    - [ ] Verify trace structure: Check that spans have correct names, statuses, and custom attributes.

## 5. CI/CD Integration
- [ ] **GitHub Actions**:
    - [ ] Update `packages.yml` (or create `observability.yml`) to:
        - Spin up the `prototype-grafana` stack.
        - Run `pnpm test:observability`.

## Questions/Decisions
- **Runner**: Should we use Vitest or Jest for the verification suite? (Vitest is already used in some apps, might be faster/easier).
- **Isolation**: To prevent flaky tests, we should generate a unique `testRunId` for each execution of the fixtures and filter our Prometheus/Tempo queries by that ID.
