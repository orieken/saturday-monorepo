---
feature: "reagent-map-globe"
bounded_context: ""
domain_terms: ["Reagent Site", "Aurelis"]
files_touched: []
issue_refs: ["reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Security Review: reagent-map-globe

## Threat Model Summary

Read-only, unauthenticated feature: an offline Cesium globe view at /reagent-map (behind a build-time-only feature flag, default off) renders site data fetched once from GET /api/reagent-sites. No forms, no auth, no write paths, no PII. The only new attack surface is (1) a dev-only Vite middleware serving Cesium's static assets from node_modules, and (2) a new pinned dependency (cesium@1.146.0). Both were inspected directly; no Critical or High findings. The one prior Critical-class issue in this feature — a WebGL context leak when fetchReagentSites() rejected after buildViewer() resolved — was already found by code-reviewer and fixed by the developer (Promise.allSettled + explicit destroy()-on-failure); verified correct by reading the current code, not relied on secondhand.

## Dependency Audit

- pnpm audit --prod run against the full workspace: 33 findings (2 low, 15 moderate, 14 high, 2 critical), all under apps__mock-api>express (qs, body-parser) — pre-existing transitive advisories unrelated to this feature's files_touched scope (mock-api's package.json was not modified by reagent-map-globe).
- cesium@1.146.0 (new, pinned exact version, no caret range): zero advisories reported by pnpm audit for this package.

## STRIDE Analysis

None

### Spoofing

No authentication or identity on this route; it is a public, read-only view with no session or token handling introduced. No spoofing surface — nothing to impersonate.

### Tampering

Site data (id/name/lat/lng) comes from the same-origin mock API, not user input, and is rendered via Vue's auto-escaping `{{ site.name }}` interpolation (no v-html) — no injection path. The feature flag (VITE_FEATURE_REAGENT_MAP) is a build-time constant baked into the bundle via import.meta.env, not a runtime-mutable value a client request could flip.

### Repudiation

No state-changing action exists in this feature (view-only), so there is nothing a user could perform and later deny. No audit trail gap applicable.

### Information Disclosure

Checked what can be read: site names/coordinates are non-sensitive seed data, already served by the pre-existing reagent-sites-api with no auth gate (out of scope here). Checked what can carry it out: the new dev-only cesiumStaticAssetsPlugin in vite.config.ts is a path-traversal candidate — verified safe (see findings). window.__SATURDAY_CESIUM__ exposes only the Viewer/Cesium handle for the Playwright test seam, no credentials or secrets attached.

### Denial of Service

fetchReagentSites() has an explicit 5s AbortController timeout (no unbounded network call). The previously-found WebGL-context-leak-on-partial-failure (repeated failed mounts exhausting the browser's per-page WebGL context limit) was fixed in this round via Promise.allSettled + newViewer.destroy() on a sites-fetch rejection — verified by reading the current onMounted implementation at ReagentMapView.vue:119-138; both rejection branches are handled and nothing is left orphaned.

### Elevation of Privilege

No roles, permissions, or per-user data in this feature — every visitor sees the same fixed site list. No privilege boundary exists for this view to cross.

## Findings

### [LOW] — cesiumStaticAssetsPlugin's resolveRequestedFile can throw on malformed percent-encoding
**Location**: `apps/ye-olde-magic-shop/vite.config.ts:44-48`
**Threat category**: Denial of Service
**Description**: resolveRequestedFile calls decodeURIComponent(urlPath.slice(...)) with no try/catch. A request like `/cesium/%` is invalid percent-encoding and decodeURIComponent throws a URIError synchronously inside the middleware, uncaught. The path-traversal guard itself is correct — decode-then-normalize ordering is right, the '..'-prefix check after normalize correctly rejects traversal attempts (e.g. '../../etc/passwd', '..%2f..%2fetc%2fpasswd'), and path.join's semantics mean an absolute-looking segment like '/etc/passwd' gets appended under CESIUM_SOURCE_DIR rather than escaping it, as path.resolve would allow. The only gap is the unhandled throw on bad encoding. Blast radius is minimal: this middleware only runs in `configureServer` (vite dev server, not shipped to production or dist/), and the dev server binds to localhost by default with no `host` change in this diff — so exploitation requires a request already on the loopback interface, which already has full filesystem access.
**Fix applied**: Recommendation only
**Verification**: Not applicable — Low severity per the severity table (track as tech debt), no fix required before QA. If QA wants to exercise it: request `GET /cesium/%` while the dev server is running and confirm whether the dev server process surfaces a clean 500/next() vs crashing; expected is Vite/connect's middleware dispatcher catching the throw, but this was not executed, only read.

## Files Modified

None

## Security Checklist

- [x] Authentication/session handling reviewed — None present — feature has no auth surface.
- [x] Authorization/permission checks reviewed — None present — no per-user data or roles.
- [x] User-supplied input validated/encoded — No user input in this feature; site names rendered via Vue's auto-escaping interpolation.
- [x] API endpoints reviewed for exposure — GET /api/reagent-sites is read-only, same-origin (via proxy/mock API), no auth; pre-existing endpoint, not modified by this diff.
- [x] Secrets/credentials/tokens handling reviewed — None introduced. No ion token, no API keys — Cesium runs fully offline per spec AC.
- [x] Rate limiting / lockout reviewed — Not applicable — no auth endpoints in this feature.
- [x] Trust-boundary crossings identified — Single boundary: browser fetch to mock API's own origin, bounded by a 5s timeout; no third-party origin contacted (verifies the spec's 'no request leaves the app's and mock API's own origins' NFR for this file set).
- [x] New dependency (cesium@1.146.0) checked for known vulnerabilities — pnpm audit: zero advisories for cesium.
- [x] Dev-only static file server (cesiumStaticAssetsPlugin) checked for path traversal — Guard is correct: decode-before-normalize, '..'-prefix rejection, path.join cannot be escaped by an absolute-looking segment. One Low-severity unhandled-throw note recorded, not blocking.
- [x] WebGL/resource cleanup on failure paths verified — Promise.allSettled + destroy()-on-sites-rejection fix from code review verified correct by reading the current implementation.

## Notes for QA

- Security review found no Critical or High issues; nothing blocks QA from proceeding.
- One Low finding is recorded as tech debt (unhandled decodeURIComponent throw in the dev-only Cesium static-asset middleware) — no fix required, informational only.
- Confirm during testing that no network request in the reagent-map view's traffic log targets any origin other than the app's own origin and the mock API's origin (the spec's explicit NFR) — this review verified it by reading the code's fetch/Cesium-config paths, not by running the app.

## Notes for Tech Writer

None
