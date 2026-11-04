# Observability Report: expeditions-api

## Status: PASS

## Summary

- `GET /api/expeditions` adds no new logging. It reads the optional session via
  `resolveSessionUserId()`, filters `readExpeditions()` in memory, and returns the result —
  silent on success, matching the sibling `/api/items` and `/api/reagent-sites` endpoints'
  convention of no log line on a successful read.
- `POST /api/expeditions` logged its success with an unstructured, mixed-argument
  `console.log('Saved expedition:', newExpedition.id, 'for user:', req.userId)` — the same BAD
  pattern `reagent-sites-api`'s observability report (`docs/features/reagent-sites-api/sre-engineer.md`,
  follow-up #2) already flagged on `POST /api/orders`, but that instance predated this feature and
  was out of scope there. This instance is new code introduced by this feature's own
  `files_touched` (`apps/mock-api/server.js`), so it was fixed in place rather than recorded as a
  follow-up — see Fixes Applied.
- The only other log statements reachable through either new endpoint are the pre-existing
  `readJSON()`/`writeJSON()` error branches (`console.error(\`Failed to read/write ${filePath}\`, err)`),
  which are low-cardinality (fixed path constants, native `Error` objects) and were not changed by
  this feature.
- No PII, credentials, or session tokens are read, logged, or echoed by either handler.
  `req.userId` — an internal session-resolved user id, not an email or name — is the only
  request-derived value that reaches a log statement, and only as a bare id (architecture-guardrails.md
  #9 permits recording "an ID").
- `apps/mock-api/server.js` carries zero OTel instrumentation anywhere in the file (confirmed by
  inspection — no `@opentelemetry/*` import, no span/tracer calls on any route). Adding spans to
  only the two new expedition endpoints while every sibling endpoint (items, users, orders, auth,
  reagent-sites) stays uninstrumented would itself be the cross-endpoint inconsistency this role
  exists to prevent, and `architecture-notes.md`'s Observability Design section records no decision
  to introduce OTel for this slice. No fitness function is proposed for the same reason
  `reagent-sites-api`'s report gave: there is nothing yet to enforce.

## Service Level Indicators (SLIs)

| SLI | Definition | Target |
|---|---|---|
| Availability | % of `GET /api/expeditions` requests returning HTTP 200 while the mock-api process is up | 100% |
| Trust-boundary correctness | % of `GET /api/expeditions` responses, across all sessions, containing zero expeditions owned by a different adventurer than the caller | 100% — any non-zero rate is a security defect, not a load concern (spec AC4; enforced by `security-reviewer` and QA's trust-boundary tests) |
| Write correctness | % of `POST /api/expeditions` successes whose stored `owner` equals `req.userId` regardless of any `owner`/`userId` field the request body supplies | 100% |
| Zero-length rejection correctness | % of `POST /api/expeditions` requests with endpoints closer than 1 km that are refused with 422 and leave `expeditions.json` unchanged | 100% |
| Latency | Local JSON file read/write, no network hop (NFR: no request leaves the app/mock-api origin) | Negligible; no timeout or circuit breaker needed — this is a dev/test fixture service, not a production SLA target |

These are the fixture-service equivalent of a business SLI, matching `reagent-sites-api`'s
precedent: this is a mock API backing local development and Saturday/Sunday test suites, not
customer traffic, so "healthy" means "the trust boundary holds and the seed/save data stays
correct," not a production latency/error budget.

## Logging Cardinality Review

- **New code path, GET /api/expeditions**: no log statements. `resolveSessionUserId()` and the
  route handler are both silent on success and on a missing/invalid token — correct, no action
  needed.
- **New code path, POST /api/expeditions**: found and fixed. See Fixes Applied below.
- **Existing shared path exercised by both endpoints**: `readJSON()`'s
  `console.error(\`Failed to read ${filePath}\`, err)` and `writeJSON()`'s equivalent — `filePath`
  is always one of the module's fixed path constants (e.g. `EXPEDITIONS_PATH`), never
  request-derived, so neither adds cardinality. Unchanged by this feature; no fix needed.

## OpenTelemetry Review

- No spans exist anywhere in `apps/mock-api/server.js` today — not on `/api/items`, `/api/orders`,
  `/api/auth/*`, `/api/reagent-sites`, or the two new expedition endpoints. Instrumenting only the
  new endpoints would be premature generalization (`design-principles.md` #7) against a file no
  architecture decision chose to restructure for this slice.
- Guardrail #8 (OTel confined to adapter/interceptor layers, never domain logic) is not
  implicated — no instrumentation was added or removed.
- Judgment-only, per `architecture-guardrails.md` #7: file-wide OTel adoption is a decision for a
  dedicated cross-cutting order, not this feature. Reason: no OTel exists anywhere in
  `apps/mock-api` to extend consistently, and this slice's own `architecture-notes.md` records no
  decision to introduce it.

## PII and Secrets Review (Guardrail #9)

- Expedition payload fields (`id`, `owner`, `from`/`to` `{lat, lng}`, `savedAt`) describe a line
  between two fictional in-world locations on the invented planet Aurelis and an internal owner
  id — no PII, no credentials.
- Neither handler logs, in full, the `Authorization` header or a session token. `POST`'s fixed log
  statement now records only `expeditionId` and `owner` (an id), both permitted properties under
  guardrail #9 (id, not raw payload).
- `owner` on `GET`'s response can be another adventurer's user id only for the shop's own seeded
  expeditions (`owner: 'shop'`) — never a different adventurer's id, which is the trust boundary
  AC4 and the security-reviewer already verify at the response-body level; this report confirms
  the same id is not separately leaked through telemetry.

## Fixes Applied

- `apps/mock-api/server.js` — `POST /api/expeditions`'s success log changed from

  ```js
  console.log('Saved expedition:', newExpedition.id, 'for user:', req.userId);
  ```

  to

  ```js
  console.log({ expeditionId: newExpedition.id, owner: req.userId }, 'Expedition saved');
  ```

  The stable message text ("Expedition saved") is now separated from the variable, groupable
  context (`expeditionId`, `owner`) instead of being interleaved as alternating string/variable
  positional arguments — the same fix `reagent-sites-api`'s report recommended for the sibling
  `POST /api/orders` log line, applied here because this log line is new code within this
  feature's own `files_touched`, not a pre-existing one outside scope. No test in
  `server.test.js` asserts on this log line's shape (confirmed by inspection — no `console.log`
  or `Saved expedition` references in the test file), so no test required updating.

## Known Issues / Follow-ups

1. `POST /api/orders`' `console.log('Created order:', newOrder.id, 'for user:', req.userId)`
   (`apps/mock-api/server.js`) carries the same unstructured pattern this report just fixed on
   `POST /api/expeditions`, and was already flagged by `reagent-sites-api`'s observability report.
   Still pre-existing, still outside this feature's `files_touched` — restating the existing
   follow-up for a dedicated mock-api logging cleanup order rather than fixing it here.
2. `apps/mock-api` has no OTel instrumentation anywhere. Acceptable for a local dev/test fixture
   service; not blocking for this order, consistent with the precedent set by
   `reagent-sites-api`'s report.
3. `GET /api/expeditions` (like every other `readJSON()`-backed GET endpoint) returns HTTP 200
   with `[]` on a missing or malformed `expeditions.json` instead of signaling failure. An
   availability SLI built purely on HTTP status code would read this as 100% healthy while
   silently serving zero expeditions, including the three seeded ones. Pre-existing, shared
   behavior across the whole file — `reagent-sites-api`'s report recorded the same gap for its own
   endpoint; worth revisiting together if a later `reagent-map` UI slice comes to depend on a
   guaranteed non-empty catalog.