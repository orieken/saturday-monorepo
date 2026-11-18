# Performance & Reliability Report: reagent-map-save

## Status: PASS

## Summary

This order adds one write operation — `saveExpedition(from, to)` POSTing to the already-implemented
`POST /api/expeditions` — plus the client-side gating, in-flight guarding, and rendering around it,
and a cross-file test-isolation lock confined to Playwright spec files. No new loop touches a
database or network call, no new collection endpoint is introduced, and no hot-path read is added.
The double-click/non-idempotent-POST gap flagged in this report's prior pass is now fully specified
in `analysis.md` (developer tasks 4, 5, 8: `isSaving` ref, early-return guard, `:disabled` binding)
and acknowledged in `architecture-notes.md`'s Stability Design section — no gate remains open on that
finding. No new critical or high risk was found in this pass.

## Idempotency on Mutations

- `POST /api/expeditions` remains **not idempotent** by design: each accepted save creates a new
  `expedition-${Date.now()}` record (`apps/mock-api/server.js`, out of scope). The in-scope mitigation
  is client-side prevention, not a server idempotency key.
- The double-click gap identified previously is now closed at the spec level: `analysis.md` task 4
  adds `isSaving = ref(false)`, task 5 makes `saveLine()` return early when `isSaving.value` is true
  and sets/clears it around the request in `try/finally`, and task 8 binds the save button's
  `:disabled` to `!canSaveLine || isSaving`. `architecture-notes.md`'s Stability Design section
  names this the same way ("in-flight guard … since the server defines no idempotency key for this
  endpoint"). A proposed fitness function in both `analysis.md` and the Playwright task list
  (rapid-click → exactly one POST, exactly one stored expedition) gives this a concrete, checkable
  verification — confirmed as owned by `qa-engineer`.
- On failure (timeout, network, 4xx), task 7 leaves the draft line and its entities untouched with no
  automatic resend — the only path to a duplicate POST is now a second *deliberate* click after the
  first request has already resolved, which is an accepted, in-scope outcome, not a retry hazard.
- No immediate risk identified here.

## Strict Timeouts

- `saveExpedition()` is specified (`analysis.md` developer task 1) to reuse `fetchExpeditions()`'s
  existing `AbortController` / `REQUEST_TIMEOUT_MS = 5_000` shape verbatim, satisfying
  architecture-guardrails.md #5 — no unbounded network call is introduced.
- No second timeout constant is introduced; `architecture-notes.md`'s own fitness functions and
  Stability Design section reconfirm the single shared `REQUEST_TIMEOUT_MS` budget.
- The new test-infra lock (`acquireExpeditionsLock()`) is itself a bounded wait: it polls via
  Playwright's `expect(fn).toPass({ timeout: 30_000 })`, not an unbounded or hand-rolled
  `while`/`sleep` loop (architecture-guardrails.md #5 is satisfied in test tooling the same way it is
  in production code). A `STALE_LOCK_MS = 30_000` recovery path prevents a crashed holder from
  deadlocking the suite permanently. This is test infrastructure, not a production network call, but
  it is reviewed here because it gates every spec run and a hang there blocks delivery just as
  effectively as a hung request would.
- No immediate risk identified here.

## Prevent N+1 Queries

- The one loop-shaped operation in this feature — `addExpeditionEntities(cesium, viewer, ellipsoid,
  [returned expedition])` — is called once per successful save with a single-element array
  (`analysis.md` task 6), not inside a per-item network loop.
- `architecture-notes.md`'s decision to append the server's returned `Expedition` directly to
  `expeditions.value` rather than re-calling `fetchExpeditions()` after a successful POST avoids an
  unnecessary extra round trip against data already in hand, and is backed by a fitness function
  (one network request per successful save, verified by the Playwright request-count assertion in
  `analysis.md`'s Proposed Fitness Functions).
- `GET /api/expeditions` is unaffected by this order and was already reviewed for N+1 risk in a prior
  order; nothing here changes that endpoint's shape.
- No immediate risk identified here.

## Caching Strategies

- No hot path is introduced. `saveLine()` is a low-frequency, user-initiated write (at most once per
  completed draft line), not a high-read path.
- The existing expedition list (`expeditions.value`) is already held in memory for the view's
  lifetime and is extended in place on save rather than re-fetched — the correct in-memory "cache"
  for this access pattern. No additional caching layer is warranted for a component-local,
  single-session list of this size.
- No immediate risk identified here.

## Observation (Out of Scope, Not a Gate)

- `apps/mock-api/server.js`'s `POST /api/expeditions` handler does a synchronous
  `readExpeditions()` → push → `writeExpeditions()` against a shared JSON file with no locking.
  Concurrent saves (from different adventurers, or a legitimate second deliberate click once the
  `isSaving` window has closed) read-modify-write the same file without server-side serialization,
  so a lost update remains theoretically possible under concurrent load. `apps/mock-api/server.js` is
  explicitly out of scope (`files_touched` excludes it; `analysis.md`'s Out of Scope section forbids
  re-implementing it client-side). Noted for the record, not gated — unchanged pre-existing behavior
  at the demo mock API's current single-process scale.
- The new Playwright-side mutex (`acquireExpeditionsLock`/`releaseExpeditionsLock`) is the in-scope
  answer to the *test* version of this same read-modify-write hazard, and is reviewed above under
  Strict Timeouts rather than repeated here.

## Gate Decision

- **No critical or high-risk N+1, caching, timeout, or idempotency vulnerability remains open.** The
  double-click/non-idempotent-POST finding raised in the prior pass of this review is now fully
  specified in `analysis.md` (tasks 4, 5, 8) and `architecture-notes.md`'s Stability Design section,
  with an owned fitness function for verification. The developer may proceed to implementation
  without further performance-engineering changes to the plan.