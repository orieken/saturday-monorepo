# Playwright OTel Reporter Implementation Plan

## Phase 1: Test Coverage & Baseline
- [ ] Replicate Cucumber `cart.feature` scenarios as native Playwright tests in `apps/ye-olde-magic-shop/tests/cart.spec.ts`.
  - [x] Add item to cart
  - [x] Increase quantity
  - [ ] Remove item
  - [x] Verify totals
- [x] Verify tests pass with `pnpm --filter @orieken/ye-olde-magic-shop test:pw`

## Phase 2: Create Reporter Package
- [x] Initialize new package `packages/saturday-playwright-otel-reporter`.
- [x] Configure `tsconfig.json`, `package.json`, and `tsup` build.
- [x] Install dependencies:
  - `@playwright/test` (as peer/dev)
  - `@opentelemetry/api`
  - `@opentelemetry/sdk-trace-node`
  - `@opentelemetry/resources`
  - `@opentelemetry/semantic-conventions`

## Phase 3: Implement Reporter Logic
- [x] Create `OtelReporter` class implementing `Reporter` interface.
- [x] Implement Lifecycle Methods:
  - `onBegin(config, suite)`: Initialize OTel provider. Start "Test Run" span.
  - `onTestBegin(test)`: Start "Test Case" span. Link to parent.
  - `onStepBegin(test, result, step)`: Start "Step" span.
  - `onStepEnd(test, result, step)`: End "Step" span. Record status/error.
  - `onTestEnd(test, result)`: End "Test Case" span. Record status/error.
  - `onEnd(result)`: End "Test Run" span. Flush/Shutdown OTel.
- [x] Reuse `TracerSetup` logic from cucumber-otel-formatter (or refactor to shared lib if appropriate, otherwise copy for independence).
- [x] Implement `SpanManager` to handle the span stack.

## Phase 4: Integration & Verification
- [x] Configure `apps/ye-olde-magic-shop/playwright.config.ts` to use the new reporter.
- [x] Run tests and verify traces in Tempo/Grafana.
- [x] Ensure "Step" spans appear correctly nested under "Test" spans.
- [x] Verify metadata (browser, platform, tags) is attached to spans.

## Phase 5: Polish
- [x] Add configuration support (env vars `OTEL_XXX`, config options).
- [x] Add robust error handling and logging (debug mode).
- [x] Write README.md.
- [x] Publish package.
