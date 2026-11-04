---
feature: "reagent-sites-api"
bounded_context: ""
domain_terms: ["Reagent Site", "Aurelis", "Adventurer"]
files_touched: []
issue_refs: ["reagent-sites-api"]
linked_adrs: []
linked_kis: []
---

# Security Review: reagent-sites-api

## Threat Model Summary

GET /api/reagent-sites is a new, intentionally unauthenticated, read-only endpoint that returns a static 8-entry array (id, name, lat, lng) of fictional game-world coordinates from a local JSON file. There is no user-supplied input (no body, no query params, no path params), no write path, no PII, no secrets, and no outbound network call — the handler only reads a bundled static file via the existing readJSON() helper shared with /api/items. Trust boundary: public internet -> mock API, same boundary /api/items already crosses with the identical pattern. Reviewed implementation files directly (server.js, reagent-sites.json, server.test.js, package.json) rather than relying on notes. No Critical, High, or Medium findings. Dependency audit: no new dependencies were introduced (cors and express are pre-existing, unchanged); npm audit could not run in apps/mock-api because this directory has no local lockfile (dependencies are resolved by the monorepo's root package manager) — not a blocker since nothing new was added to the dependency tree.

## Dependency Audit

- No new dependencies introduced by this feature — package.json only gained a "test": "node --test" script, which uses Node's built-in test runner.
- npm audit could not execute in apps/mock-api (ENOLOCK — no local lockfile; this app resolves deps via the monorepo root). Not a blocker: existing cors/express versions are unchanged by this diff.

## STRIDE Analysis

None

### Spoofing

Endpoint deliberately bypasses requireAuth and returns an identical response regardless of a missing, valid, or garbage Authorization header (verified in server.js:84-86 and asserted by server.test.js's garbage-header test). There is no identity to spoof because the endpoint grants no identity-based access or capability — it is public by design per the spec's AC1. No spoofing surface introduced.

### Tampering

No request body, query string, or path parameter is accepted or parsed by this handler (app.get('/api/reagent-sites', ...) takes no :id or input), so there is no attacker-controlled input to tamper the response with. The response data comes only from a static file (reagent-sites.json) committed to the repo, mirroring readItems()'s existing pattern. No tampering surface.

### Repudiation

A GET of public, non-sensitive catalog data does not require an audit trail — there is no action a user could later deny performing, and no other endpoint in this file logs GET /api/items either (consistent existing convention). Not applicable.

### Information Disclosure

Read side: the response contains only fictional game-world site names and lat/lng coordinates for an invented planet (Aurelis) — no PII, no credentials, no internal paths, no stack traces, and the data is explicitly meant to be public per the spec. Egress side: the handler makes no outbound call of any kind (confirmed by direct code read of the 1-line handler and readJSON()) so there is no channel for this data to be exfiltrated anywhere beyond the HTTP response itself. CORS is wide-open (app.use(cors()) at server.js:14) but that is pre-existing, applies identically to /api/items already, and is not part of this feature's diff — not a new exposure, and out of scope to fix here since it would change behavior for every other endpoint in the file.

### Denial of Service

readReagentSites() re-reads and re-parses the 8-entry (51-line) JSON file from disk on every request, same as readItems()/readUsers()/readOrders() — already flagged by code-reviewer as an existing, unbounded-but-small-scale convention rather than a new smell. No rate limiting exists on this or any other GET endpoint in this mock API, which is consistent with its purpose as a local test double, not a production service. No new DoS surface relative to the rest of the file.

### Elevation of Privilege

No role or permission model applies to this endpoint — it is uniformly public for every caller, authenticated or not, by explicit design (AC1). There is no higher-privilege state to elevate into. Not applicable.

## Findings

None

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, API keys, or credentials introduced
- [x] No user-supplied input is parsed, interpolated, or trusted by the new handler
- [x] No PII, credentials, or sensitive data present in the response payload or data file
- [x] No outbound network call introduced (verified by direct read of the handler and readJSON())
- [x] Response does not vary by Authorization header, eliminating any auth-bypass ambiguity (verified by test)
- [x] No new dependency added
- [x] Error responses (malformed data file) degrade to an empty array via existing readJSON() catch, not a stack trace or 500

## Notes for QA

- No security-driven test additions are required — the existing server.test.js already covers the two properties that matter for this trust boundary: identical 200 response with no auth header and with a garbage Bearer token.
- If QA adds further coverage, confirm the endpoint still returns 200 (never 401/403) for any Authorization header value, since that is the one behavior a future accidental requireAuth addition would silently break.

## Notes for Tech Writer

- Worth documenting explicitly (not just implicitly via code) that GET /api/reagent-sites is intentionally public and unauthenticated — this is a deliberate API contract decision (AC1), not an oversight, so a future reader shouldn't 'fix' it by adding auth.
