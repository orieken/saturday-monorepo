---
feature: "expeditions-api"
bounded_context: ""
domain_terms: ["Adventurer", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/mock-api/server.test.js", "apps/mock-api/server.test.js", "apps/mock-api/server.js"]
issue_refs: ["expeditions-api", "reagent-map-spec"]
linked_adrs: []
linked_kis: []
---

# Code Review: expeditions-api

## Overall Status

**APPROVED**

## Design Narrative

This second pass only touches apps/mock-api/server.test.js: the developer closed the four gaps from the prior CHANGES_REQUESTED review (invalid-token GET, malformed-body POST, mirrored Borin GET, and a 1 km boundary POST) without touching server.js. The endpoints and their trust-boundary allow-list construction are unchanged from the previous round. Behaviour change: no.

## Design Score

- **Clarity** 4/5
- **Cohesion** 4/5
- **Coupling** 4/5
- **Craft** 4/5

## Security Surface

- New API endpoints (GET/POST /api/expeditions) — unchanged this round
- Authorization: per-adventurer owner filtering — now covered symmetrically (Elara and Borin) plus an invalid-token case proving a bad token behaves like no session, never like a valid one
- Trust boundary: POST ignores a forged owner field — covered by an existing test, unchanged this round

## Performance Surface

- GET /api/expeditions remains unbounded with no pagination — pre-existing, matches GET /api/items and GET /api/reagent-sites, not a regression, carried forward as a non-blocking note

## Test Design Review

- GET with an invalid token returns exactly the three seeded expeditions — YES: asserts length===3 and every owner==='shop'; would fail if an invalid/expired token were mistakenly treated as a valid session (leak) or triggered a 401 instead of the required optional-auth 200.
- GET with Borin's session returns seeded plus his own, never Elara's — YES: mirrors the already-approved Elara case; would fail if the owner filter were asymmetric or Elara's expedition leaked into Borin's view.
- POST with endpoints just at the 1 km boundary is accepted (strict less-than) — YES, with a caveat: the 1e-9-degree epsilon (documented and justified by a measured floating-point underflow in the exact-boundary calculation) places the tested distance a hair above 1.0 km rather than exactly at it. It would fail if the accept threshold shifted upward (e.g. a unit or radius error), but it cannot distinguish `distance < 1` from `distance <= 1` at the literal boundary, so it doesn't fully deliver the 'pins the strict less-than semantics' claim in its own name. Non-blocking: this was already a non-blocking ask, and an exactly-1.0-km input is not constructible without the same floating-point noise the implementation notes measured.
- POST with a missing to field is refused with 400 and nothing is stored — YES: asserts 400, an /from and to/i error, and a byte-identical data file; would fail if isValidPoint's guard were removed or weakened.
- POST with non-numeric coordinates is refused with 400 and nothing is stored — YES: same shape as above with a string lat; would fail if isValidPoint stopped checking typeof.

## Verification of Developer Self-Review

- Developer's claim of 18/18 passing verified independently: ran node --test twice from a clean working tree, 18/18 both times.
- Developer's claim that running the suite twice leaves git status on apps/mock-api/data/ unchanged verified independently — confirmed no diff before vs. after two consecutive runs.
- Developer's claim that this round touched only apps/mock-api/server.test.js verified against the actual diff — no production code changed this round.
- Magic-number mirroring (AURELIS_RADIUS_KM, MIN_EXPEDITION_DISTANCE_KM declared locally in the test file with a comment) matches what the diff shows — accurate.

## Feedback for the Developer

### 1. None — all four prior findings addressed (non-blocking)
- **File**: `apps/mock-api/server.test.js`
- **Smell**: N/A
- **Instruction**: No action required.
### 2. Clarify test name or add exact-boundary comment (non-blocking)
- **File**: `apps/mock-api/server.test.js`
- **Smell**: The boundary test's name says it pins 'strict less-than' semantics, but the 1e-9 epsilon needed to dodge floating-point noise means the tested input is measurably above 1.0 km, not at it — so the test cannot distinguish a `<` implementation from a `<=` one at the literal boundary.
- **Instruction**: Non-blocking: either soften the test name to 'just above the 1 km boundary' to match what it actually proves, or add a one-line comment (next to the existing one) stating explicitly that this test verifies 'at-or-above is accepted' rather than the exact `<` vs `<=` distinction, which is untestable here without reintroducing the floating-point false-failure the epsilon was added to avoid.
### 3. Carry forward pre-existing gap (non-blocking)
- **File**: `apps/mock-api/server.js`
- **Smell**: GET /api/expeditions remains unbounded with no pagination (architecture-guardrails.md #6), matching GET /api/items and GET /api/reagent-sites.
- **Instruction**: Still out of scope for this feature — a sitewide follow-up ticket across apps/mock-api's collection endpoints remains the right fix, not a rework here.
