# Performance Report: saturday-playwright-cesium

---
feature: "saturday-playwright-cesium"
bounded_context: "Test Automation"
files_touched: ["packages/saturday-playwright-cesium/fixtures/serve.mjs", "packages/saturday-playwright-cesium/src/fixture.ts", "packages/saturday-playwright-cesium/src/geometry.ts", "packages/saturday-playwright-cesium/playwright.config.ts"]
issue_refs: ["saturday-playwright-cesium"]
linked_adrs: []
linked_kis: []
---

## Summary

This feature is a Playwright test-fixture package with no production API surface, no database, and no external network calls — it serves a single local fixture page over `127.0.0.1` and drives an in-page Cesium viewer. The four governing risks (idempotency, timeouts, N+1 queries, caching) were checked against the salvage implementation at commit `926ad56` (`src/fixture.ts`, `src/geometry.ts`, `fixtures/serve.mjs`, `playwright.config.ts`), read directly rather than assumed from the analyst/architect summaries. **No critical or blocking risk was found.** One caching gap is flagged as a non-blocking recommendation for the developer to consider, since it affects local suite run-time rather than correctness or a guardrail.

## Idempotency on Mutations

Not applicable. The package exposes no POST/PUT/DELETE-shaped operations — `clickAt` and `flyTo` drive a mouse click and a camera move against a locally rendered page, the same category of action a human tester performs, not a state-mutating API call a network retry could duplicate. No idempotency key or dedupe mechanism is required here.

## Strict Timeouts

**Sound.** Every wait in `src/fixture.ts` goes through Playwright's `page.waitForFunction(..., { timeout: timeoutMs })` with an explicit `timeoutMs` (default `DEFAULT_SETTLE_TIMEOUT_MS = 15000`), and each failure path throws a message naming the specific stuck condition (`seam`, or the joined list from `describeUnsettledConditions()`: `tiles`, `camera`, `render`) rather than hanging indefinitely. `flyTo` reuses `waitForSettled()`'s bounded wait rather than adding its own unbounded wait. This matches `architecture-guardrails.md` #5's explicit-timeout requirement using the framework-provided bounded-wait primitive (`page.waitForFunction`), not a hand-rolled `for`/`while` + `sleep` loop — confirming the architect's ratified judgment call in `architecture-notes.md`. No gap found.

## Prevent N+1 Queries

**No risk.** `extractEntities()`'s `for (const entity of viewer.entities.values)` loop runs entirely inside the page via `page.evaluate()` — it is in-memory iteration over Cesium's own entity collection, not a loop issuing a database query or network round-trip per iteration. There is exactly one `page.evaluate()` call per fixture method (`getEntities`, `toScreen`, `clickAt`'s underlying `toScreen`, `flyTo`, `metersPerPixelAt`, `globeRadiusMeters`), each a single serialized round-trip into the page, never one per entity or per point. No batching or eager-loading pattern is needed because there is nothing to batch.

## Caching Strategies

**One non-blocking recommendation.** `fixtures/serve.mjs`'s `handleRequest()` does an unconditional `stat` + `readFile` from disk and returns `200` with no `Cache-Control`, `ETag`, or `Last-Modified` header on every request — including requests for Cesium's own `Build/Cesium` assets and `fixtures/assets/aurelis.jpg`, which never change within a test run. This is a legitimate "hot path" by this package's own shape: `playwright.config.ts` has no per-test fixture reuse (`fullyParallel: true`, default one-page-per-test), so all 11 Playwright tests independently reload `index.html` and re-request the same static asset set, and `retries: process.env.CI ? 2 : 0` multiplies that further on any CI retry. None of these repeated reads are served from cache; each is a fresh disk read and full response body.

This does not block AC9 (the suite already passes offline in both local and `--network none` docker runs per the spec) and is not a correctness issue — Cesium's asset set here is scoped to core `Build/Cesium` plus one local texture, not terrain/3D-tiles (explicitly out of scope), so the absolute cost is bounded. It is flagged because it is the one genuine "hot path, no caching" shape this feature has, and conditional-GET support (`ETag`/`If-None-Match` or a simple `Cache-Control: no-cache` + 304 on unchanged `mtime`) in `serve.mjs` would shorten the local dev test loop across repeated runs without any risk to the `--network none` isolation requirement, since caching here is same-origin and local-disk, not a third-party CDN dependency.

**Recommendation for developer, non-gating:** if `serve.mjs`'s request latency is observed to matter during QA's fail-proof verification pass (AC10) or local iteration, add `Cache-Control` / conditional-GET support to the two static routes. Do not block the pipeline on this — it is a suggestion, not a required change, and the architect's existing fitness functions for this file (127.0.0.1 bind, `--network none` pass) remain the binding checks.

## Verdict

**No immediate risks identified that block the developer from proceeding.** Idempotency is not applicable, timeouts are correctly bounded end-to-end, no N+1 pattern exists because there is no network/database loop anywhere in the fixture, and the one caching gap found is a local-dev-loop efficiency note rather than a macro-architectural bottleneck. The developer may proceed to implementation (applying the salvage diff) without addressing the caching note as a precondition.