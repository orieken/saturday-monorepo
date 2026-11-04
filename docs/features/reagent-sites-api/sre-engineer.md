# Observability Report: reagent-sites-api

## Status: PASS

## Summary

- `GET /api/reagent-sites` adds no new logging, metrics, or spans. It returns the eight seeded sites straight from `readReagentSites()` with no `console.log`/`console.error` call on its success path, matching the sibling `/api/items` endpoint's silent-on-success convention.
- The only log statement reachable through this endpoint is the pre-existing `readJSON()` error branch (`console.error(\`Failed to read ${filePath}\`, err)`), which is low-cardinality (fixed path constant, native `Error` object) and was not changed by this feature.
- `architecture-notes.md`'s Observability Design section is explicitly "None," and this review confirms that judgment: `apps/mock-api` carries zero OTel instrumentation anywhere in the file, so adding spans to one endpoint while every sibling endpoint (items, users, orders, auth) stays uninstrumented would itself be the inconsistency this agent exists to prevent.
- No PII, credentials, or session data is read, logged, or echoed by the handler — it ignores `req.headers` entirely, per AC1's identical-response-regardless-of-session requirement.
- No fixes were required. Two pre-existing, out-of-scope observability gaps elsewhere in `server.js` are recorded below as follow-ups rather than fixed here, since they fall outside `reagent-sites-api`'s `files_touched` and fixing them would violate this order's scope.paths discipline (analysis.md Task List; `approval-gates.md` #6).

## Service Level Indicators (SLIs)

| SLI | Definition | Target |
|---|---|---|
| Availability | % of `GET /api/reagent-sites` requests returning HTTP 200 while the mock-api process is up | 100% |
| Response completeness | Returned array length and contents equal `apps/mock-api/data/reagent-sites.json` (8 sites) | 100% equality — any drift is a data defect, not a load concern |
| Latency | Local static-file read, no network hop (NFR: no request leaves the app/mock-api origin) | Negligible; no timeout or circuit breaker needed — this is a dev/test fixture service, not a production SLA target |

These are the fixture-service equivalent of a business SLI: this is a mock API backing local development and Saturday/Sunday test suites, not customer traffic, so "healthy" means "the seed catalog it serves is complete and correct," not a production latency/error budget.

## Logging Cardinality Review

- **New code path**: no log statements added. `function readReagentSites() { return readJSON(REAGENT_SITES_PATH); }` and `app.get('/api/reagent-sites', (req, res) => { res.json(readReagentSites()); })` are both silent on success — correct, no action needed.
- **Existing shared path exercised by this endpoint**: `readJSON()`'s `console.error(\`Failed to read ${filePath}\`, err)` — `filePath` is a fixed constant (`REAGENT_SITES_PATH`), not request-derived, so this does not add cardinality. Unchanged by this feature; no fix needed.
- **Found, not fixed (out of scope)**: `app.post('/api/orders', ...)` at `apps/mock-api/server.js:239` logs `console.log('Created order:', newOrder.id, 'for user:', req.userId)` — this is exactly the BAD pattern this role exists to catch (unstructured string concatenation mixing a stable message with variable identifiers inline, ungroupable by any log aggregator). It predates `reagent-sites-api`, is not in this feature's `files_touched`, and correcting it here would expand the diff beyond `scope.paths` for no reason tied to this change. Recorded below as a follow-up rather than edited.

## OpenTelemetry Review

- No spans exist anywhere in `apps/mock-api/server.js` today — not on `/api/items`, `/api/orders`, `/api/auth/*`, or the new `/api/reagent-sites`. Adding OTel to only the new endpoint would be premature generalization (`design-principles.md` #7) against a file the architect explicitly decided not to restructure for this slice.
- Guardrail #8 (OTel confined to adapter/interceptor layers, never domain logic) is not implicated — no instrumentation was added or removed.
- No fitness function is proposed here, consistent with `architecture-notes.md`'s own "Observability Design: None" and "judgment-only" fitness function for this slice — there is nothing yet to enforce.

## PII and Secrets Review (Guardrail #9)

- Response payload fields (`id`, `name`, `lat`, `lng`) describe fictional in-world locations on the invented planet Aurelis — no PII, no credentials, no user-identifying data.
- The handler does not inspect, log, or echo the `Authorization` header or any session token, by design (AC1).
- No new telemetry properties are recorded by this change, so Guardrail #9's opt-in-per-field declaration requirement does not yet apply here — not applicable.

## Fixes Applied

None required. No high-cardinality or unstructured logging was introduced by this change.

## Known Issues / Follow-ups

1. `GET /api/reagent-sites` (and every other `readJSON()`-backed GET endpoint) returns HTTP 200 with `[]` on a missing or malformed data file instead of signaling failure. An availability SLI built purely on HTTP status code would read this as 100% healthy while silently serving zero sites. This is pre-existing, shared behavior across the whole file, not a regression introduced here — worth revisiting if a later reagent-map slice (`reagent-map-globe`) comes to depend on a guaranteed non-empty catalog.
2. `POST /api/orders`' `console.log('Created order:', newOrder.id, 'for user:', req.userId)` (`apps/mock-api/server.js:239`) should be converted to structured logging, e.g. `console.log({ orderId: newOrder.id, userId: req.userId }, 'Order created')`. Pre-existing and outside `reagent-sites-api`'s scope — flag for a dedicated mock-api logging cleanup order.
3. `apps/mock-api` has no OTel instrumentation anywhere. Acceptable for a local dev/test fixture service; not blocking for this order, and consistent with the architect's explicit decision to defer this.