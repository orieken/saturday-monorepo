---
feature: "expeditions-api"
bounded_context: "Expedition Planning — Expedition storage (apps/mock-api backend)"
domain_terms: ["Expedition", "Seeded Expedition", "Adventurer", "Reagent Site", "Aurelis"]
files_touched: ["apps/mock-api/data/expeditions.json", "apps/mock-api/data/users.json", "apps/mock-api/server.js", "apps/mock-api/server.js", "apps/mock-api/server.js", "apps/mock-api/server.js", "apps/mock-api/server.js"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: expeditions-api

## Structural Decisions

### Implement both endpoints inside apps/mock-api/server.js following the exact readJSON/writeJSON + flat-handler pattern already used by items, orders, and reagent-sites. No new package, module, class, or Clean Architecture layering is introduced.
- **Rationale**: Verified server.js directly: readItems/readOrders/readReagentSites and their handlers are structurally identical to what this feature needs for readExpeditions/writeExpeditions. apps/mock-api has no domain/application/adapter split today — introducing one for two endpoints would be inconsistent with every other endpoint in the file and is explicitly outside scope.paths. This matches CLAUDE.md's architect-skip criterion for purely additive work within an established pattern, and repeats the precedent architecture-notes.md already set for reagent-sites-api.
- **Trade-offs**:
  - Carries forward the existing unbounded in-memory JSON reload-per-request cost — acceptable at this app's demo scale, same as every other endpoint.
  - Defers (again) the cost reagent-sites-api's architecture notes already flagged: a later slice needing a richer Expedition domain model will have to revisit server.js's flat structure.
- **Fitness Function**: judgment-only — Verified server.js directly: readItems/readOrders/readReagentSites and their handlers are structurally identical to what this feature needs for readExpeditions/writeExpeditions. apps/mock-api has no domain/application/adapter split today — introducing one for two endpoints would be inconsistent with every other endpoint in the file and is explicitly outside scope.paths. This matches CLAUDE.md's architect-skip criterion for purely additive work within an established pattern, and repeats the precedent architecture-notes.md already set for reagent-sites-api.
### Introduce a new, explicitly named optional-session helper — resolveSessionUserId(req) — that returns the session's user id or null, and never rejects the request. GET /api/expeditions uses this helper; it MUST NOT use requireAuth, and requireAuth MUST NOT be changed to support an optional mode.
- **Rationale**: This is the one genuinely new structural element in this feature: every existing endpoint is either fully public or fully gated by requireAuth's all-or-nothing contract (missing/invalid token -> 401). GET /api/expeditions needs a third shape — proceed regardless, personalize if possible. Solving this by adding an 'optional' flag to requireAuth would corrupt a contract every other call site already depends on being all-or-nothing. A separate, named function keeps the two contracts impossible to confuse: requireAuth guards; resolveSessionUserId informs.
- **Trade-offs**:
  - Two session-lookup code paths now exist (requireAuth's and resolveSessionUserId's) instead of one — acceptable because they have genuinely different contracts, and the alternative (overloading requireAuth) is worse.
  - Any future endpoint needing optional auth must reuse resolveSessionUserId rather than reinventing it — noted in developer handoff.
- **Fitness Function**: GET /api/expeditions returns 200 for every request regardless of Authorization header content (absent, malformed, expired, or valid), and a valid session's response set equals the shop's seeded expeditions plus exactly that session's own expeditions — never another adventurer's. — verified by node --test apps/mock-api — qa-engineer asserts: (1) no Authorization header -> 200 with exactly the 3 seeded expeditions; (2) garbage/expired token -> 200 with exactly the 3 seeded expeditions (never 401); (3) Elara's valid session -> seeded + Elara's own, Borin's absent; (4) Borin's valid session -> seeded + Borin's own, Elara's absent.
### The endpoint-to-endpoint distance check (haversine geodesic distance) is a pure function, geodesicDistanceKm(from, to, radiusKm), colocated in server.js near the other helpers, with AURELIS_RADIUS_KM = 6000 and MIN_EXPEDITION_DISTANCE_KM = 1 as named constants at the top of the file — not inline magic numbers, and not extracted into a separate module or package.
- **Rationale**: ARCHITECTURE_RULES.md IV forbids magic numbers; naming the constants satisfies that without introducing new structure. Extracting a separate geo-math module or an Expedition domain package now would be premature generalization (design-principles.md #7) against work that doesn't exist yet — apps/mock-api has no domain layer to receive it, and no second caller exists yet. Keeping it a pure function (no I/O, no req/res) means it is trivially unit-testable in isolation from the Express handler, which is what keeps the handler itself under the complexity budget despite combining JSON I/O, auth, and validation.
- **Trade-offs**:
  - If a later slice (e.g. multi-point expeditions, routing) needs richer geo-math, this function will need to move into whatever domain layer that slice introduces — deferred, not eliminated.
  - No test currently proves the constant is actually used (vs. a hardcoded duplicate) short of code-reviewer's manual check, since apps/mock-api has no magic-number lint rule configured.
- **Fitness Function**: geodesicDistanceKm is a standalone, side-effect-free function with no dependency on req/res/fs; AURELIS_RADIUS_KM and MIN_EXPEDITION_DISTANCE_KM are declared as named constants, used by both the handler and directly by tests, never duplicated as literals elsewhere in the expedition-related code. — verified by Manual review by code-reviewer (apps/mock-api has no magic-number/no-restricted-syntax ESLint rule configured today — same tooling gap reagent-sites-api's architecture-notes.md already recorded); qa-engineer additionally unit-tests geodesicDistanceKm directly with known coordinate pairs to prove it can fail independently of the HTTP layer.
### POST /api/expeditions MUST construct the persisted record from an explicit allow-list — { id, owner: req.userId, from, to, savedAt } — and MUST NOT spread or copy req.body wholesale into the stored object. owner and savedAt are never read from the request body.
- **Rationale**: This is the feature's core trust-boundary invariant (AC11, and the security gate's named property: 'POST always stores req.userId as owner, ignoring any owner/userId field the request body supplies'). The existing POST /api/orders handler already destructures named fields rather than spreading req.body, so this follows the codebase's own established safe pattern — it is not a new convention, just the one that must not be skipped here under time pressure.
- **Trade-offs**:
  - Slightly more verbose than `{ ...req.body, id, owner: req.userId }` — that shorthand is explicitly rejected because object key order in a spread does not protect against a later field reordering reintroducing the vulnerability; the allow-list form has no such failure mode.
- **Fitness Function**: Submitting a POST body that includes a forged owner or userId field never changes the stored or returned owner from req.userId. — verified by node --test apps/mock-api — security-reviewer/qa-engineer POST with { from, to, owner: 'shop' } (or another user's id) using Elara's session, assert the stored and returned expedition's owner equals Elara's user id, not the forged value.
### apps/mock-api/data/expeditions.json stores pre-resolved { from: {lat,lng}, to: {lat,lng} } coordinate pairs for every seeded expedition, copied at seed-authoring time from reagent-sites.json's existing site records. Neither GET nor POST /api/expeditions reads reagent-sites.json at request time; the two data files are not runtime-coupled.
- **Rationale**: The API contract (per domain-model.md H1, carried into analysis.md) is 'a line and its owner' — two coordinate pairs, not two site-id references. Resolving site-name-to-coordinate once, by hand, when authoring the seed file keeps GET /api/expeditions a single readJSON call (matching every other GET handler's shape) and avoids introducing a foreign-key-style lookup between two otherwise-independent data files for a read-only seed. This also means expeditions-api introduces zero runtime dependency on reagent-sites-api's data file, keeping the two orders' sibling relationship (noted in both analyses) read-only and one-directional at authoring time only.
- **Trade-offs**:
  - If reagent-sites.json's coordinates for Moonwell Spring, Ember Crag, etc. ever change, expeditions.json's seeded coordinates will silently go stale unless someone re-derives them by hand — acceptable for a demo seed file, called out here so it isn't mistaken for a bug later.
- **Fitness Function**: judgment-only — The API contract (per domain-model.md H1, carried into analysis.md) is 'a line and its owner' — two coordinate pairs, not two site-id references. Resolving site-name-to-coordinate once, by hand, when authoring the seed file keeps GET /api/expeditions a single readJSON call (matching every other GET handler's shape) and avoids introducing a foreign-key-style lookup between two otherwise-independent data files for a read-only seed. This also means expeditions-api introduces zero runtime dependency on reagent-sites-api's data file, keeping the two orders' sibling relationship (noted in both analyses) read-only and one-directional at authoring time only.
### No CircuitBreaker, ExponentialBackoffStrategy, or network timeout is introduced for this feature, and none is needed: the distance check is synchronous in-memory arithmetic with zero network calls, consistent with the NFR that no request may leave the app's/mock API's own origins.
- **Rationale**: architecture-guardrails.md #5 only applies to network calls; there are none here. This decision exists to make explicit and testable the NFR that a future contributor must not 'upgrade' the distance check to call an external geocoding/mapping API.
- **Fitness Function**: No new outbound HTTP/network client (fetch, axios, a geocoding SDK, etc.) is imported or used anywhere in the expedition-related code added by this feature. — verified by code-reviewer greps the diff for new require()/import of an HTTP client library in the expedition handlers and helper; absence confirms the check stayed local.

## Component Placement

- `expeditions.json (seeded shop expeditions + Elara's and Borin's saved expeditions)` → Data (flat JSON seed file, no ORM/domain layer) layer, `apps/mock-api/data/expeditions.json` — Matches the existing items.json/orders.json/reagent-sites.json convention — a static array read via readJSON().
- `users.json additions (Elara, Borin)` → Data layer, `apps/mock-api/data/users.json` — Purely additive records matching the existing id/email/password/name/createdAt/address shape — no schema change, consistent with guardrail #2 (no destructive migrations).
- `readExpeditions() / writeExpeditions()` → Framework (flat server.js, no domain layer exists) layer, `apps/mock-api/server.js` — Mirrors readItems/readOrders/writeOrders — thin readJSON/writeJSON wrappers, same file, same section style.
- `geodesicDistanceKm(from, to, radiusKm) + AURELIS_RADIUS_KM / MIN_EXPEDITION_DISTANCE_KM constants` → Domain logic, colocated (no domain layer exists in apps/mock-api to receive it yet) layer, `apps/mock-api/server.js` — Pure, side-effect-free function kept next to the other helpers so it stays independently unit-testable without pulling in Express req/res.
- `resolveSessionUserId(req) — optional-session helper` → Framework layer, `apps/mock-api/server.js` — New auth-adjacent helper, deliberately separate from requireAuth so the all-or-nothing and optional contracts can never be confused; placed next to requireAuth/generateToken.
- `GET /api/expeditions handler` → Framework layer, `apps/mock-api/server.js` — Uses resolveSessionUserId, never requireAuth; placed in its own '============ EXPEDITIONS ENDPOINTS ============' block near /api/reagent-sites.
- `POST /api/expeditions handler` → Framework layer, `apps/mock-api/server.js` — Gated by requireAuth, mirrors POST /api/orders' req.userId-as-owner pattern with an explicit field allow-list instead of a body spread.

## Bounded Context

- **Owning Context**: Expedition Planning — Expedition storage (apps/mock-api backend)
- **Context Crossings**: Identity (mock API session store) — reads req.userId via the existing requireAuth middleware and the new resolveSessionUserId helper; does not change sign-in or session creation., reagent-sites-api — expeditions.json's seeded coordinates are copied from reagent-sites.json's site records at seed-authoring time only; no runtime read of that file by either new endpoint., features/reagent-map UI slices (reagent-map-globe, reagent-map-save) — future, not-yet-built consumers of these two endpoints, gated separately by VITE_FEATURE_REAGENT_MAP; out of scope here.

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] No Clean Architecture layer boundary applies — apps/mock-api/server.js is a flat, single-file Express app with no domain/application/adapter separation. — Confirmed by reading server.js directly (loom MCP graph tools were not used — not needed for a single-file app with no directory structure to query). No dependency-direction violation is possible to introduce because there is no inner/outer layer split to violate.
- [PASS] The new optional-auth helper (resolveSessionUserId) must not be reachable from, or merged into, requireAuth's gating contract. — Enforced by naming and placement decision above (structuralDecisions #2) rather than by a graph query; verified by reading requireAuth and confirming no shared code path is proposed between it and the new helper.

## Anti-Pattern Check

- Premature Generalization: not found — Considered and rejected extracting an Expedition domain class, a geo-math module, or a shared sites/expeditions linkage now; deferred until a slice actually needs it (structuralDecisions #1, #3, #5).
- God Object: not found — server.js continues to grow (now six endpoint groups), same as reagent-sites-api's note recorded. Not worsened materially by two more handlers in the established style; refactoring server.js is outside scope.paths for this order.
- Leaky Abstraction: not found — resolveSessionUserId returns a plain user id or null — callers don't need to know how sessions are stored, same abstraction level requireAuth already exposes.
- Shotgun Surgery: not found — All changes are additive and localized to server.js plus two new/extended data files; no existing endpoint's behavior changes.

## Fitness Functions

- **GET /api/expeditions never returns 401 and never returns another adventurer's expedition for any session.**
  - **Verification**: node --test apps/mock-api — qa-engineer asserts response sets for no-session, invalid-session, Elara's session, and Borin's session per structuralDecisions #2.
  - **Owner**: qa-engineer
- **POST /api/expeditions always stores req.userId as owner, ignoring any owner/userId field the request body supplies.**
  - **Verification**: node --test apps/mock-api — POST with a forged owner field in the body, assert the stored and returned owner equals the session's adventurer id.
  - **Owner**: security-reviewer
- **A line whose endpoints are closer than 1 km apart on Aurelis (radius 6000 km) is rejected with 422 naming a zero-length line, and nothing is persisted.**
  - **Verification**: node --test apps/mock-api — assert 422 and an unchanged expeditions.json for near-identical coordinates.
  - **Owner**: qa-engineer
- **Running the mock-api test suite twice leaves git status unchanged, including apps/mock-api/data/.**
  - **Verification**: Run `node --test apps/mock-api` twice, confirm `git status --short` shows no diff — tests must snapshot and restore any data file they write to, per the t.after() pattern already established in server.test.js.
  - **Owner**: qa-engineer
- **Every new or changed function related to expeditions has cyclomatic complexity below 7.**
  - **Verification**: Manual review by code-reviewer — apps/mock-api has no ESLint complexity rule configured yet (same gap reagent-sites-api's architecture-notes.md recorded).
  - **Owner**: code-reviewer
- **No new outbound HTTP/network client is introduced for the distance check.**
  - **Verification**: code-reviewer greps the diff for a new HTTP client import in the expedition-related code; absence confirms the NFR holds.
  - **Owner**: code-reviewer

## Refactoring Opportunities

None

## Developer Handoff Notes

- Add EXPEDITIONS_PATH, readExpeditions(), writeExpeditions() to server.js following readJSON/writeJSON exactly as readItems/readOrders do.
- Add AURELIS_RADIUS_KM = 6000 and MIN_EXPEDITION_DISTANCE_KM = 1 as named constants near the top of server.js (not inline literals), and implement geodesicDistanceKm(from, to, radiusKm) as a pure function with no req/res/fs access, so qa-engineer can unit-test it directly against known coordinate pairs.
- Add resolveSessionUserId(req) next to requireAuth: reads the Authorization header, looks up sessions.get(token), returns the user id or null — it must never call res.status() or next() itself; it is a plain lookup, not middleware.
- GET /api/expeditions: no requireAuth. Call resolveSessionUserId(req); always return the shop's (owner: 'shop') expeditions, plus — only when resolveSessionUserId returned a non-null id — that session's own expeditions filtered by owner === that id.
- POST /api/expeditions: behind requireAuth. Validate body.from/body.to as {lat, lng}; call geodesicDistanceKm; reject with 422 naming a zero-length line when distance < MIN_EXPEDITION_DISTANCE_KM. On success, build the stored record as { id, owner: req.userId, from, to, savedAt: new Date().toISOString() } — an explicit allow-list, never a spread of req.body — and return it.
- Build apps/mock-api/data/expeditions.json by hand: copy the three seeded expeditions' from/to coordinates out of the current apps/mock-api/data/reagent-sites.json records (Moonwell Spring, Ember Crag, Thornveil Glade, Glimmerdeep Caves, Saltmarsh Cairn, Tidewrack Shoals), plus Elara's (Moonwell Spring → Frostcap Hollow) and Borin's (Ember Crag → Saltmarsh Cairn) expeditions with owner set to their respective user ids. Do not add a runtime read of reagent-sites.json into either handler.
- Add Elara and Borin to apps/mock-api/data/users.json matching the existing user record shape exactly (id, email, password, name, createdAt, address).
- Add both new endpoints to the startup console.log endpoint list, matching the existing style.
- Do not introduce an Expedition class, repository interface, or shared package in this slice — stay inside server.js's existing flat style, per structuralDecisions #1 and #5.

## Open Architectural Questions

None
