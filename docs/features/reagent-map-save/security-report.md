---
feature: "reagent-map-save"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: ["reagent-map-save"]
linked_adrs: []
linked_kis: []
---

# Security Review: reagent-map-save

## Threat Model Summary

Small, additive diff: a new saveExpedition() adapter function in expeditions.ts and a save handler in ReagentMapView.vue, plus flipping VITE_FEATURE_REAGENT_MAP on. No new dependencies, no new server endpoint (POST /api/expeditions and its requireAuth/ownership/zero-length checks already existed, unchanged, in apps/mock-api/server.js from the prior reagent-map-expeditions feature). The client-side authStore.isAuthenticated check in saveLine() is UX convenience only; the server's requireAuth middleware remains the actual trust boundary, deriving ownership from req.userId (the session token) and never from the request body, which the POST body (JSON.stringify({ from, to })) confirms — no owner/adventurer identifier is sent by the client. The server's own error text (including the AC13 zero-length message) is surfaced verbatim via a thin Error wrapper, which is a correct leaky-abstraction avoidance, not a disclosure of anything the server didn't already decide to say. No auth/session/token logic was added or changed; no secrets, headers, or PII appear in any new code. Net new attack surface is minimal and no Critical or High findings were found.

## Dependency Audit

- No package.json, package-lock.json, or pnpm-lock.yaml changed in this diff (confirmed via git diff) — no new dependencies introduced, so no new CVE surface to audit for this stage.

## STRIDE Analysis

None

### Spoofing

saveExpedition() reuses getAuthHeaders() (same localStorage 'dnd_auth_token' key authStore.ts also reads/writes), attaching Authorization: Bearer <token> exactly as fetchExpeditions() already does. The server's requireAuth middleware (unchanged) is the actual identity check; the client's authStore.isAuthenticated gate in saveLine() is a UX shortcut only and cannot be used to impersonate another adventurer since the server never trusts a client-supplied identity.

### Tampering

POST body is JSON.stringify({ from, to }) only — no owner/id field a client could tamper with to claim another adventurer's identity. Server-side isValidPoint()/geodesicDistanceKm() validation (unchanged) still bounds and checks the two points before persisting; this diff adds no new server-side input handling.

### Repudiation

No new logging was added by this diff. The server's existing console.log({expeditionId, owner}, 'Expedition saved') (unchanged) already records who saved what; not a regression introduced here.

### Information Disclosure

Egress check: saveExpedition() targets the same API_BASE resolution fetchExpeditions() already uses (build-time VITE_API_BASE, else localhost:8001/mock-api:8001) — no new outbound destination, satisfying the spec's NFR that no request leaves the app's and mock API's own origins.

### Denial of Service

saveExpedition() sets the same REQUEST_TIMEOUT_MS=5000ms AbortController timeout as fetchExpeditions() — no unbounded wait. isSaving guards against a double-click producing two in-flight POSTs from the UI. Server-side rate limiting/idempotency for POST /api/expeditions is unchanged and remains a client-side-only prevention, already accepted as a judgment-only trade-off in architecture-notes.md; noted below as pre-existing, not a new gap.

### Elevation of Privilege

Ownership (req.userId) is still derived entirely server-side from the session token; the GET /api/expeditions filter (owner === 'shop' || owner === userId, unchanged) and the POST handler's newExpedition.owner = req.userId (unchanged) mean this diff introduces no path for an adventurer to read or write another adventurer's expedition.

## Findings

### [LOW] — Double-submit and save-volume limits are enforced client-side only
**Location**: `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue:132,275 (isSaving); apps/mock-api/server.js:141 (POST /api/expeditions, unchanged)`
**Threat category**: Denial of Service
**Description**: isSaving prevents a double-click from the UI producing two in-flight POSTs, but a caller bypassing the UI (devtools, curl, a buggy retry loop) with a valid token can still issue unlimited POST /api/expeditions requests — there is no server-side rate limit or idempotency key. Blast radius is limited to the caller's own stored expeditions (ownership is still server-derived, so this cannot affect other adventurers' data), and the endpoint itself is unchanged by this diff — architecture-notes.md already accepted the client-side-only guard as a judgment-only trade-off. Flagged for awareness; no fix applied since server.js is out of this stage's changed files and the trade-off was already reviewed.
**Fix applied**: Recommendation only
**Verification**: Not applicable — tracked as existing tech debt, not a regression from this diff.
### [INFO] — Non-JSON error body would surface a raw parse error instead of a friendly message
**Location**: `apps/ye-olde-magic-shop/src/api/expeditions.ts:52`
**Threat category**: Information Disclosure
**Description**: On a non-ok response, `await res.json()` can throw if the server ever returns a non-JSON error body (e.g. an HTML 500 page from an intermediary). That SyntaxError's message (e.g. 'Unexpected token < in JSON at position 0') would then be shown to the adventurer via saveMessage instead of a friendly fallback. This discloses only that parsing failed, not any data, path, or secret — already identified by code-reviewer as non-blocking because apps/mock-api/server.js's documented contract always returns a JSON error body for non-ok /api/expeditions responses today. No fix applied; confirming the prior call was correct from a security lens as well.
**Fix applied**: Recommendation only
**Verification**: Not applicable — no exploitable disclosure; documented for awareness only.

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, API keys, or tokens in the diff
- [x] No new dependencies introduced (verified via git diff on lockfiles)
- [x] Server derives ownership from the session token only, never from client-supplied request body fields — Unchanged server.js behavior; POST body carries only { from, to }
- [x] Auth enforcement happens server-side (requireAuth), not only in client UI logic — ReagentMapView.vue's isAuthenticated check is UX convenience only
- [x] Every new network call has an explicit timeout (architecture-guardrails.md #5) — saveExpedition() uses the existing REQUEST_TIMEOUT_MS=5000ms AbortController pattern
- [x] Error text shown to the user is rendered via auto-escaping template interpolation, not v-html
- [x] No PII, tokens, or secrets appear in any new log statement — No new logging added by this diff
- [x] No new outbound destination introduced (egress stays within app/mock-api origins)

## Notes for QA

- Verify the AC12 signed-out path shows 'Sign in to save this expedition.' and triggers zero POST /api/expeditions calls (page.route interception), per implementation-notes.md.
- Also worth a quick manual/exploratory check (not required to block): simulate a mid-session token expiry (clear the server-side session but leave authStore hydrated) and confirm the resulting 401 surfaces the server's 'Unauthorized' text in save-line-message without crashing the view — this is a UX edge, not a security gap, since the server still correctly rejects the write.
- No new dependencies and no auth/session code changed, so no additional security-specific test scenarios are required beyond what analysis.md's QA Tasks already scope (the expeditions.json cross-test lock and reagent-map-save.spec.ts).

## Notes for Tech Writer

- No new environment secrets were introduced — .env's only change is flipping the existing VITE_FEATURE_REAGENT_MAP build-time flag to true; nothing new to document in secrets/config handling.
- Saving an expedition requires sign-in, enforced server-side via the existing Bearer-token session check; this is worth one line in any user-facing or API doc update QA/tech-writer produce for this feature.
