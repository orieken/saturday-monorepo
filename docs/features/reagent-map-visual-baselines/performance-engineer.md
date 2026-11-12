# Performance Report: reagent-map-visual-baselines

## Summary

This feature adds Playwright `@visual` screenshot tests, one error-state DOM element to `ReagentMapView.vue`, a `RUNNER_VARIANT`-keyed `snapshotPathTemplate` in `playwright.config.ts`, and a `--update-baselines` flag to `scripts/ui-suite-in-chrome.sh`. There is no new API endpoint, no database access, no request/response payload, and no production hot path — the only "runtime" surface is a CI/local test-runner shell script and a static-file comparison step. **No immediate risks identified** against the four governing principles below; one existing pattern is flagged for the developer to preserve rather than regress.

## Idempotency on Mutations

- No state-mutating API endpoint is introduced or modified. The sole "mutation" this feature adds is writing baseline PNG files to `.playwright-baselines/{debian,ubi9}/`, and the spec already requires that write to be the explicit, opt-in effect of `--update-baselines` — never an implicit side effect of a normal test run. Re-running `--update-baselines` against an unchanged UI reproduces the same bytes (or a visually-equivalent re-render), so the operation is naturally idempotent; it does not need a request-id/dedupe-key pattern the way a payment or order-creation endpoint would.
- Gate condition: `run_chrome_suite()` MUST NOT append `--update-snapshots` unless the flag was passed (already called out in architecture-notes.md's fitness functions). This is the one place a developer could accidentally make a read-only test run silently mutate committed baselines — verify the parameter threading, not an env-var read inside the function, as architecture-notes.md already specifies.

## Strict Timeouts

- `wait_for_app_healthy()` already bounds its poll loop to `APP_HEALTH_TIMEOUT_SECONDS` (180s) with a distinct `EXIT_APP_UNHEALTHY` exit code, and each `curl` probe inside it carries `--max-time 5`. This feature's NFR explicitly requires this timeout stay unchanged — confirm the developer does not need to raise it to accommodate five-run repeat checks (each run restarts the container fresh, so this should not compound).
- The new `#cesiumContainer` `@visual` test waits on `cesium.waitForSettled()` before capturing. Confirm this helper has its own bounded timeout (inherits Playwright's default test timeout, or an explicit one) rather than an unbounded poll — an indefinitely-hanging "wait for settled" on a flaky globe render would stall the whole suite, not just fail one assertion. If no explicit bound exists today, the developer should not add one as new scope here (it's inherited behavior, not introduced by this feature) — flag it as an observation for the developer to confirm, not a gate-blocking finding.
- `scripts/ui-suite-in-chrome.sh` makes no new network calls; the only new flag-parsing path (`--update-baselines`) is pure in-process argument handling with no I/O of its own.

## Prevent N+1 Queries

- Not applicable. There is no database, no API client, and no loop issuing per-item network or query calls anywhere in this feature's scope. The route-interception pattern (`page.route('**/api/reagent-sites', ...)`) reused for the error-state `@visual` test is a single intercepted call, matching the existing `reagent-map-resilience.spec.ts` pattern — no new fan-out.

## Caching Strategies

- Not applicable in the hot-path sense this principle targets (no high-read production data). The "cache" here is the committed baseline PNG set itself, which is already correctly scoped per runner image (`debian`/`ubi9`) so a debian-rendered baseline is never compared against an ubi9 run — this is the right design, not a gap. No CDN/Redis/in-memory caching decision is needed for a local test-runner artifact.

## Missing Reliability Structures

- None required. `parse_args`'s restructuring from a single-positional `case` to a loop over `"$@"` is a correctness change for flag composition, not a reliability primitive — no `CircuitBreaker` or `ExponentialBackoffStrategy` applies to in-process argument parsing.
- `wait_for_app_healthy()`'s poll loop is a bounded, one-shot readiness check of a process the script itself started, already documented in-code as `PRECONDITION`/`ENFORCED-BY` judgment-only per architecture-notes.md's Stability Design section — this is the correct exemption shape under `architecture-guardrails.md` #5 (that guardrail targets retries against unreliable *external* dependencies, not a bounded startup probe of a child process), and this feature does not touch or extend that loop.

## Gate Decision

**No critical or high-severity performance/reliability risk found. The developer may proceed.** The only item carried forward for the developer and QA to confirm during implementation is `cesium.waitForSettled()`'s own timeout bound, noted above as an observation, not a blocker.