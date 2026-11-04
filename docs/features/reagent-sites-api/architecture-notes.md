---
feature: "reagent-sites-api"
bounded_context: "Expedition Planning — Reagent Site catalog (apps/mock-api backend)"
domain_terms: ["Reagent Site", "ReagentSite", "ReagentSiteCatalog"]
files_touched: ["apps/mock-api/data/reagent-sites.json", "apps/mock-api/server.js", "apps/mock-api/server.js"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: reagent-sites-api

## Structural Decisions

### Skip deep architectural redesign: GET /api/reagent-sites follows the exact existing pattern already used by GET /api/items — a readJSON-backed helper, an unauthenticated app.get handler, a static data file under apps/mock-api/data/, listed in the startup console.log block.
- **Rationale**: Verified apps/mock-api/server.js directly: readItems()/ITEMS_PATH and the GET /api/items handler are structurally identical to what analysis.md proposes for readReagentSites()/REAGENT_SITES_PATH. No new base class, interface, layer, or cross-cutting concern is introduced. This meets CLAUDE.md's architect-skip criteria ('purely additive within an established pattern') exactly, so no RFC and no new component-placement decisions beyond confirming the existing shape.
- **Trade-offs**:
  - None beyond those already inherent in the existing items.json/readItems() pattern (unbounded in-memory JSON reload per request, acceptable at this app's scale).
- **Fitness Function**: judgment-only — Verified apps/mock-api/server.js directly: readItems()/ITEMS_PATH and the GET /api/items handler are structurally identical to what analysis.md proposes for readReagentSites()/REAGENT_SITES_PATH. No new base class, interface, layer, or cross-cutting concern is introduced. This meets CLAUDE.md's architect-skip criteria ('purely additive within an established pattern') exactly, so no RFC and no new component-placement decisions beyond confirming the existing shape.
### server.js must export the Express app (module.exports = app) and guard app.listen(...) behind if (require.main === module) so node:test can require() the app without binding a port.
- **Rationale**: server.js currently calls app.listen() unconditionally at module load with no module boundary for testability. This is the one structural gap that blocks every acceptance criterion requiring node:test coverage (AC2: unauthenticated-vs-session parity; AC4: complexity/coverage). The change is additive and backward compatible: `npm start` still runs `node server.js`, which still satisfies require.main === module and listens exactly as before.
- **Trade-offs**:
  - None — this is the standard Express testability pattern (export app, guard listen); it changes no runtime behavior for the existing start script.
- **Fitness Function**: apps/mock-api/server.js can be required by a test file without binding PORT, and still binds PORT when run directly via `node server.js` / `npm start`. — verified by node --test apps/mock-api — the new server.test.js imports server.js via require() and makes requests against the in-memory app (e.g. via supertest-free direct invocation or an ephemeral listen() in the test itself); a passing run with no EADDRINUSE and no port conflict confirms the guard works.
### New reagent-sites data and handler stay inside apps/mock-api — no new package, no shared library extraction — because this is the first of six slices splitting reagent-map-spec's D1 work, and the other five slices (reagent-map-globe, reagent-map-save, etc.) are separate, not-yet-started orders that may change the final shape of the Reagent Site domain.
- **Rationale**: Extracting a ReagentSite abstraction, a repository interface, or a shared types package now would be premature generalization (design-principles.md #7) against work that doesn't exist yet. apps/mock-api has no Clean Architecture layering today (flat server.js) — introducing one for a single read-only endpoint would be inconsistent with every other endpoint in the file and outside this order's scope.paths.
- **Trade-offs**:
  - If a later slice needs a richer ReagentSite domain model (e.g. expeditions referencing sites), that slice's analyst/architect will need to revisit server.js's flat structure — deferred cost, not eliminated cost.
- **Fitness Function**: judgment-only — Extracting a ReagentSite abstraction, a repository interface, or a shared types package now would be premature generalization (design-principles.md #7) against work that doesn't exist yet. apps/mock-api has no Clean Architecture layering today (flat server.js) — introducing one for a single read-only endpoint would be inconsistent with every other endpoint in the file and outside this order's scope.paths.

## Component Placement

- `reagent-sites.json` → data layer, `apps/mock-api/data/reagent-sites.json` — Matches the existing items.json/users.json/orders.json convention — static seed data read via readJSON().
- `readReagentSites()` → framework layer, `apps/mock-api/server.js` — Mirrors readItems()/readUsers()/readOrders() — a thin readJSON(REAGENT_SITES_PATH) wrapper, same file, same section style.
- `GET /api/reagent-sites handler` → framework layer, `apps/mock-api/server.js` — Mirrors the unauthenticated GET /api/items handler — no requireAuth, placed in its own '============ REAGENT SITES ENDPOINTS ============' block near /api/items.

## Bounded Context

- **Owning Context**: Expedition Planning — Reagent Site catalog (apps/mock-api backend)
- **Context Crossings**: features/reagent-map UI slices (reagent-map-globe, reagent-map-save) — future consumers of this endpoint, not yet built, expeditions-api — sibling endpoint for the same bounded context, out of scope here

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] No existing Clean Architecture layer boundaries apply — apps/mock-api/server.js is a flat, single-file Express app with no domain/application/adapter separation today. — Confirmed by reading server.js directly (loom MCP graph tools were not used — not needed for a single-file app with no directory structure to query). No dependency-direction violation is possible to introduce here because there is no inner/outer layer split to violate.

## Anti-Pattern Check

- Premature generalization: not found — Considered and rejected extracting a ReagentSite domain abstraction now; deferred until a slice actually needs it (see structuralDecisions #3).
- God object: not found — server.js is already a single large file housing all mock-api endpoints; adding one more unauthenticated GET in the established style does not materially worsen this. Not in scope to refactor here — scope.paths restricts this order to additive changes.

## Fitness Functions

- **apps/mock-api/server.js can be required by a test file without binding PORT, and still binds PORT when run directly.**
  - **Verification**: node --test apps/mock-api — server.test.js requires server.js and exercises the app in-process; a clean test run with no port-binding error confirms the require.main === module guard works.
  - **Owner**: qa-engineer
- **Every changed line in server.js and the new reagent-sites.json is covered by a test that can fail**
  - **Verification**: node --test apps/mock-api — assertions deep-equal the full 8-site array against docs/specs/reagent-map-spec/acceptance-criteria.json's seed.reagentSites (per test-repair-contract.md, no weakened assertions).
  - **Owner**: qa-engineer
- **Cyclomatic complexity stays below 7 for readReagentSites() and the new handler**
  - **Verification**: Manual review by code-reviewer — apps/mock-api has no ESLint config today, so this cannot run as an automated CI check for this file yet; both new functions are single-line pass-throughs so the risk is low.
  - **Owner**: code-reviewer

## Refactoring Opportunities

None

## Developer Handoff Notes

- Confirmed by direct read of apps/mock-api/server.js: follow the readItems()/GET /api/items shape exactly for readReagentSites()/GET /api/reagent-sites — no deviation needed.
- Add module.exports = app and wrap the existing app.listen(...) in if (require.main === module) { ... } — this is required for qa-engineer's node:test suite to import the app without binding the port, and must not change the behavior of `npm start`.
- Do not introduce a ReagentSite class, repository interface, or shared package in this slice — stay inside server.js's existing flat style per the architect's premature-generalization check above.

## Open Architectural Questions

None
