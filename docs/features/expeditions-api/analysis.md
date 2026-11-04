---
feature: "expeditions-api"
bounded_context: "Expedition Planning — Expedition storage (apps/mock-api backend)"
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/mock-api/server.js", "apps/mock-api/data/expeditions.json", "apps/mock-api/data/users.json", "apps/mock-api/server.test.js"]
issue_refs: ["expeditions-api", "reagent-map-spec"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: expeditions-api

## Summary

Adds GET /api/expeditions and POST /api/expeditions to the mock API so an adventurer can save an expedition (a line between two points on the planet Aurelis) and find it again, the shop's three seeded expeditions are visible to everyone, and no adventurer ever sees another adventurer's expedition. This is the API-only slice of reagent-map-spec's D1 order (sibling to the already-delivered reagent-sites-api); it follows that delivery's established flat-server.js patterns exactly — readJSON/writeJSON helpers, requireAuth for the authenticated write, and owner-filtering already proven by GET/POST /api/orders — with one new wrinkle: GET /api/expeditions must read an optional session (unlike any existing endpoint, which is either fully public or fully gated) to decide whether to add the caller's own expeditions to the seeded set. No feature flag is needed: this is backend-only and additive, inert until a UI slice (reagent-map-globe/reagent-map-save, gated separately by VITE_FEATURE_REAGENT_MAP) calls it.

### Acceptance Criteria

- The seed data MUST contain exactly the three shop-owned expeditions in seed.seededExpeditions, and adventurers Elara and Borin MUST exist as users with the expeditions seed.adventurers assigns them.
  - Example: shop: Spring to Crag (Moonwell Spring → Ember Crag), Glade to Caves (Thornveil Glade → Glimmerdeep Caves), Cairn to Shoals (Saltmarsh Cairn → Tidewrack Shoals)
  - Example: Elara: Moonwell Spring → Frostcap Hollow
  - Example: Borin: Ember Crag → Saltmarsh Cairn
- GET /api/expeditions with no session MUST return exactly the three seeded expeditions (spec AC3).
- GET /api/expeditions with Elara's session MUST return the seeded expeditions plus Elara's own, and MUST NOT ever return Borin's (spec AC4). This is the trust boundary the security gate checks.
- POST /api/expeditions with a session and a line whose two endpoints are at least 1 km apart on Aurelis (sphere, radius 6,000 km) MUST store the expedition for the session's adventurer, never an owner the request body supplies, and the stored expedition MUST then be returned by GET (spec AC11).
- POST /api/expeditions with no session MUST be refused with 401 and MUST NOT store anything (spec AC12).
- POST /api/expeditions with endpoints less than 1 km apart MUST be refused with 422 naming a zero-length line, and MUST NOT store anything (spec AC13).
- Tests MUST run through node:test via the mock-api's existing `test` script (node --test, added by reagent-sites-api) and MUST leave the repository unchanged: running the suite twice MUST show no git status change, including under apps/mock-api/data/. POST writes to a data file, as registering a user does, so tests MUST isolate or restore what they write.
- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths (apps/mock-api/**, docs/features/expeditions-api/**) MUST change.

### Non-Functional Requirements

- **security**: No request this feature makes MUST leave the app's and mock API's own origins — the geodesic distance check MUST be computed locally (plain arithmetic), never via an external geocoding or mapping service.
- **other**: Every new or changed function's cyclomatic complexity MUST stay below the framework-wide budget. (cyclomatic complexity per function 7McCabe complexity (exclusive upper bound))

## Proposed Fitness Functions

- **GET /api/expeditions never returns another adventurer's expedition**
  - **Verification**: node --test apps/mock-api — assert Borin's expedition is absent from Elara's session response and vice versa
  - **Owner**: qa-engineer
- **POST /api/expeditions always stores req.userId as owner, ignoring any owner/userId field the request body supplies**
  - **Verification**: node --test apps/mock-api — POST with a forged owner field in the body, assert the stored and returned owner equals the session's adventurer id
  - **Owner**: security-reviewer
- **A line whose endpoints are closer than 1 km apart on Aurelis (radius 6000 km) is rejected with 422 naming a zero-length line, and nothing is persisted**
  - **Verification**: node --test apps/mock-api — assert 422 and an unchanged expeditions.json for near-identical coordinates
  - **Owner**: qa-engineer
- **Running the mock-api test suite twice leaves git status unchanged, including apps/mock-api/data/**
  - **Verification**: manual/CI check: run `node --test apps/mock-api` twice, confirm `git status --short` shows no diff
  - **Owner**: qa-engineer
- **Every new or changed function in server.js related to expeditions has cyclomatic complexity below 7**
  - **Verification**: Manual review by code-reviewer — apps/mock-api has no ESLint complexity rule configured yet (same gap reagent-sites-api's architecture-notes.md recorded)
  - **Owner**: code-reviewer

## Out of Scope

- Reagent sites (reagent-sites-api) — already delivered, untouched here.
- Any UI — this slice is API only.
- Editing or deleting an expedition.
- The scenario suites E1 and F1 — run by hand from the spec package, not by this slice.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Expedition Planning — Expedition storage (apps/mock-api backend)
- **Context Crossings**: Identity (mock API session store) — reads req.userId via the existing requireAuth middleware and sessions map; does not change sign-in., reagent-sites-api — reads apps/mock-api/data/reagent-sites.json's site ids/positions to resolve seeded expeditions' from/to names into coordinates; does not modify that file., features/reagent-map UI slices (reagent-map-globe, reagent-map-save) — future, not-yet-built consumers of these two endpoints, gated separately by VITE_FEATURE_REAGENT_MAP.

### Domain Events (Event Storming Lite)

- **Commands**: SaveExpedition
- **Events Produced**: ExpeditionSaved, ExpeditionRejected
- **Owning Aggregates**: Expedition
- **Read Models / Projections**: VisibleExpeditions

### Affected Components

- `apps/mock-api/server.js` — Add readExpeditions()/writeExpeditions() helpers (mirroring readItems/readOrders), a geodesic-distance helper for the 1 km zero-length check, GET /api/expeditions (new optional-auth pattern: read the session if present, don't require it), and POST /api/expeditions behind the existing requireAuth middleware, mirroring POST /api/orders' req.userId-as-owner pattern.
- `apps/mock-api/data/expeditions.json` — New seed data file holding the three shop-owned expeditions (resolved from seed.seededExpeditions via reagent-sites.json's ids/lat/lng) plus Elara's and Borin's saved expeditions, following the same flat JSON-array convention as items.json/orders.json/reagent-sites.json.
- `apps/mock-api/data/users.json` — Add Elara and Borin as users (matching the existing id/email/password/name/createdAt/address shape) so a real session can be created for them in tests and so AC4's trust-boundary check has two distinct real adventurers to compare.
- `apps/mock-api/server.test.js` — Add node:test cases for both new endpoints (AC3, AC4, AC11, AC12, AC13) following the existing startServer/stopServer helpers, the t.after() snapshot-and-restore pattern for data-file isolation, and the @issue/@ac annotation convention already used for reagent-sites-api's tests.

### Data Model Changes

- New Expedition read model/data file (apps/mock-api/data/expeditions.json): array of { id, owner ('shop' or an adventurer's user id), from: {lat, lng}, to: {lat, lng}, savedAt }. No name field is required by the API contract (domain-model.md H1: 'a line and its owner' only) — any name present in seed data is a human-readable label only. (**expand** phase)
- apps/mock-api/data/users.json gains two new user records (Elara, Borin) — purely additive, no existing record changes shape. (**expand** phase)

### API Changes

- `GET /api/expeditions` — New endpoint. No requireAuth gate, but reads the Authorization header optionally: with no valid session, returns only the shop's (owner: 'shop') expeditions; with a valid session, also includes that session's own expeditions. Never includes another adventurer's.
- `POST /api/expeditions` — New endpoint behind requireAuth. Body carries from/to {lat, lng}. Computes the geodesic distance between from and to on Aurelis's sphere (radius 6000 km); rejects with 422 (naming a zero-length line) when the two points are closer than 1 km apart. On success, stores a new expedition with owner = req.userId (ignoring any owner/userId the body supplies) and savedAt = now, and returns it. With no session, requireAuth already returns 401 before any body is read.

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Add EXPEDITIONS_PATH, readExpeditions(), and writeExpeditions() to apps/mock-api/server.js following the existing readJSON/writeJSON helper pattern.
2. Create apps/mock-api/data/expeditions.json seeded with the three shop-owned expeditions from docs/specs/reagent-map-spec/acceptance-criteria.json's seed.seededExpeditions, resolving from/to site names to lat/lng via apps/mock-api/data/reagent-sites.json.
3. Add Elara and Borin to apps/mock-api/data/users.json (matching the existing user shape) and seed their saved expeditions into expeditions.json with owner set to each adventurer's user id.
4. Implement a pure geodesic-distance helper (haversine or spherical law of cosines), parameterized by Aurelis's radius (6000 km), used only to decide whether two endpoints are closer than 1 km apart.
5. Implement GET /api/expeditions: no requireAuth; look up the Authorization header's token in the sessions map without rejecting a missing or invalid one; return seeded (owner: 'shop') expeditions plus, when a valid session exists, that session's own expeditions.
6. Implement POST /api/expeditions behind requireAuth: validate body.from/body.to as {lat, lng}; reject with 422 naming a zero-length line when the geodesic distance is below 1 km; otherwise store a new expedition with owner = req.userId (never from the request body) and savedAt = new Date().toISOString(), and return it.
7. Add both new endpoints to the startup console.log endpoint list, matching the existing style.

### QA Tasks

1. Add tests to apps/mock-api/server.test.js for GET /api/expeditions: no session returns exactly the three seeded expeditions (AC3); Elara's session returns seeded + Elara's, never Borin's (AC4). Annotate with @issue expeditions-api and the matching @ac.
2. Add tests for POST /api/expeditions: a valid line is stored under the session's adventurer and visible via a subsequent GET, and a forged owner field in the request body is ignored (AC11).
3. Add tests for POST with no session → 401, nothing persisted (AC12); POST with endpoints less than 1 km apart → 422 naming a zero-length line, nothing persisted (AC13).
4. Snapshot and restore apps/mock-api/data/expeditions.json (and users.json, since Elara/Borin are added) via t.after(), mirroring server.test.js's existing USERS_PATH snapshot-restore pattern; verify the full suite run twice leaves no git diff.
5. Review every new/changed function for cyclomatic complexity below 7 (no ESLint complexity rule exists in apps/mock-api yet, so this is a manual check).

### Tech Writer Tasks

None

### DevOps Tasks

None

## Edge Cases and Risks

- Request body includes an owner or userId field on POST /api/expeditions: Ignored entirely — the stored and returned owner is always req.userId from the session, never trusted from the client body, matching POST /api/orders' existing pattern.
- Two points exactly 1 km apart (the tolerance boundary): The settled rule is 'closer than the tolerance' (strict less-than), so exactly 1.0 km is accepted as a valid, non-zero-length line.
- GET /api/expeditions with a garbage or expired Authorization header: Since GET uses an optional-auth lookup rather than requireAuth, an invalid token behaves identically to no session: seeded expeditions only, no 401.
- Distance calculation across the antimeridian or near a pole (e.g. Tidewrack Shoals at lng 180, Frostcap Hollow at lat 85): A standard haversine or spherical law-of-cosines great-circle formula is correct at any latitude/longitude, including poles and the antimeridian — no special-casing needed.
- Repeated test runs appending to the same expeditions.json file: Tests snapshot apps/mock-api/data/expeditions.json (and users.json) before writing and restore them in t.after(), so the suite is repeatable and leaves no git diff.

## Definition of Done

- All 8 acceptance criteria above are verified by passing node:test cases.
- GET /api/expeditions and POST /api/expeditions are implemented in apps/mock-api/server.js following the existing readJSON/writeJSON, requireAuth, and owner-filtering patterns.
- apps/mock-api/data/expeditions.json is seeded with exactly the three shop-owned expeditions from seed.seededExpeditions, resolved to concrete lat/lng via reagent-sites.json.
- apps/mock-api/data/users.json includes Elara and Borin with the saved expeditions seed.adventurers assigns them.
- `node --test apps/mock-api` passes twice in a row with no git status change.
- Every new or changed function has cyclomatic complexity below 7 (code-reviewer verified).
- No file outside apps/mock-api/** or docs/features/expeditions-api/** is modified.
- DOMAIN_DICTIONARY.md already defines Adventurer/Reagent Site/Expedition/Seeded Expedition/Aurelis (section 7) — no new terms needed.
