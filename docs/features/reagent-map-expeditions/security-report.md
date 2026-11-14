---
feature: "reagent-map-expeditions"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: ["reagent-map-expeditions", "reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Security Review: reagent-map-expeditions

## Threat Model Summary

This slice adds a new GET-only client (expeditions.ts) and renders its results as Cesium polyline entities plus clickable site-list buttons. There is no authentication requirement, no mutation, no user-supplied input accepted by the client, and no new dependency. The one trust boundary this feature newly crosses — GET /api/expeditions — is pre-existing (expeditions-api dependency, out of this feature's scope.paths) and already resolves ownership server-side from the session token rather than any client-supplied identifier, so no IDOR path opens from this change. No Critical or High findings. Two Low/Info observations noted for awareness; no code changes were required.

## Dependency Audit

- No new npm/pnpm dependencies introduced by this feature — git diff shows no changes to apps/ye-olde-magic-shop/package.json or pnpm-lock.yaml, confirmed via `git diff --stat HEAD -- apps/ye-olde-magic-shop/package.json pnpm-lock.yaml` (empty).
- pnpm audit not re-run: no new dependency surface exists for this slice to introduce a new CVE.

## STRIDE Analysis

None

### Spoofing

fetchExpeditions() attaches an optional Bearer token read from localStorage, never required, mirroring orders.ts's existing pattern. The server (apps/mock-api/server.js, pre-existing, out of scope) resolves the caller's identity from its own session map via resolveSessionUserId(), not from any client-asserted field — the frontend change introduces no new spoofing surface.

### Tampering

Expedition coordinates feed directly into Cesium.Cartesian3.fromDegreesArray() as numbers; no string is parsed as code or markup. fetchExpeditions() and fetchReagentSites() both use a bare `as Expedition[]`/`as ReagentSite[]` type assertion instead of runtime schema validation (Zod) — pre-existing pattern this feature mirrors rather than introduces, see finding below.

### Repudiation

No state-changing action exists in this slice (GET only); audit-trail concerns don't apply. Expedition creation (POST /api/expeditions) is out of scope for this order per the spec's Out of Scope section.

### Information Disclosure

Expedition/site `name` renders via Vue template interpolation and `:aria-label` attribute binding, both auto-escaping — no innerHTML/v-html, no DOM XSS path. The `owner` field on Expedition travels over the wire into the client but is never read or rendered by ReagentMapView.vue — see finding below for the egress angle. Server already filters the response to the shop's seeded expeditions plus the caller's own, so no cross-user disclosure originates from this change.

### Denial of Service

fetchExpeditions() carries the same explicit 5s AbortController timeout as fetchReagentSites(), satisfying architecture-guardrails.md #5. addExpeditionEntities() is a single bounded forEach over a small, server-filtered dataset (three seeded expeditions plus at most one signed-in adventurer's own) — no unbounded loop over attacker-controlled volume.

### Elevation of Privilege

No role or permission check is added or bypassed by this slice; it is a read-only view with no admin or privileged action.

## Findings

### [LOW] — fetchExpeditions() trusts the API response shape via type assertion, not runtime validation
**Location**: `apps/ye-olde-magic-shop/src/api/expeditions.ts:31`
**Threat category**: Tampering
**Description**: `return (await res.json()) as Expedition[]` asserts the shape at compile time only; a malformed or unexpected response (e.g. a missing `from`/`to`) would reach Cesium.Cartesian3.fromDegreesArray() and fail there rather than at the trust boundary. This is not a new risk introduced by this feature — fetchReagentSites() (apps/ye-olde-magic-shop/src/api/reagentSites.ts:19) has the identical pattern already on `main`, and this change is an exact mirror of it, consistent with the developer's and code-reviewer's explicit note that the duplicated API_BASE/auth pattern across clients is a disclosed, deferred refactor. Not blocking: the mock API is same-origin, developer-controlled, and not attacker-reachable in this slice's threat model.
**Fix applied**: Recommendation only
**Verification**: Not required for this feature; tracked as pre-existing tech debt. If/when a shared typed API client or Zod validation layer is introduced across apps/ye-olde-magic-shop/src/api/*, extend it to expeditions.ts at the same time as reagentSites.ts and orders.ts rather than singling this file out.
### [INFO] — Expedition.owner travels to the client unused by the UI
**Location**: `apps/ye-olde-magic-shop/src/api/expeditions.ts:2`
**Threat category**: Information Disclosure
**Description**: The `Expedition` type includes `owner: string`, and the server already scopes the response to the shop's seeded expeditions (`owner === 'shop'`) plus the signed-in caller's own — so the only owner value a caller ever sees for their own expeditions is their own userId, not another user's. ReagentMapView.vue never reads or renders `owner`. No exploit path exists today; flagged only because the field is on the wire and in the TS surface for a UI that has no use for it, so a future change that does render it should treat it as already a per-user identifier, not a safe-to-display label.
**Fix applied**: Recommendation only
**Verification**: No action required now. If a future slice renders `owner` in the UI, confirm it never displays another adventurer's raw userId/email without a documented decision in architecture-notes.md.

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, tokens, or credentials introduced
- [x] Auth token attached only when present, never required, on a bodyless GET — matches orders.ts pattern
- [x] Server-side authorization derives identity from session token, not client-supplied userId — no IDOR surface
- [x] No innerHTML/v-html — all dynamic text goes through Vue's auto-escaping interpolation/attribute binding
- [x] Explicit timeout (AbortController, 5s) on the new network call
- [x] No unbounded loop over attacker-controlled data
- [x] No new dependencies introduced — dependency audit not applicable
- [ ] Runtime schema validation (Zod) on fetched response — Mirrors pre-existing fetchReagentSites() pattern; tracked as LOW finding, not blocking for this slice.

## Notes for QA

- No security-driven test additions are required before QA proceeds — no Critical/High findings exist.
- The three selector-path test repairs (reagent-map-globe.spec.ts, reagent-map-pan-zoom.spec.ts) and the new resilience test for fetchExpeditions() rejection were read in full: none weaken an assertion, none touch the test-repair-contract.md boundary — confirmed independently of code-reviewer's Test Design Review.

## Notes for Tech Writer

None
