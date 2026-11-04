---
feature: "reagent-sites-api"
bounded_context: ""
domain_terms: ["Reagent Site", "Adventurer", "Aurelis"]
files_touched: ["apps/mock-api/data/reagent-sites.json", "apps/mock-api/server.test.js", "apps/mock-api/server.js", "apps/mock-api/package.json"]
issue_refs: ["reagent-sites-api"]
linked_adrs: []
linked_kis: []
---

# Implementation Notes: reagent-sites-api

## Files Created

- apps/mock-api/data/reagent-sites.json
- apps/mock-api/server.test.js

## Files Modified

- apps/mock-api/server.js
- apps/mock-api/package.json

## Interface Design

- GET /api/reagent-sites -> 200, JSON array of 8 { id, name, lat, lng } objects, no auth required, identical response regardless of Authorization header
- function readReagentSites(): object[] — thin readJSON(REAGENT_SITES_PATH) wrapper, mirrors readItems()/readUsers()/readOrders()
- module.exports = app — server.js now exports the Express app so node:test can require() it without binding PORT; app.listen(...) is guarded by if (require.main === module)

## Named Refactoring Log

- **Extract Function** on `apps/mock-api/server.js — app.listen(...) call wrapped in if (require.main === module) { ... }` — server.js previously bound the port unconditionally at module load, which is not a testability smell introduced by this feature but a pre-existing gap that blocked every acceptance criterion requiring node:test coverage (architecture-notes.md structural decision #2). Guarding listen() and exporting app is the standard Express testability pattern; npm start still runs node server.js, which satisfies require.main === module and listens exactly as before.

## Self-Review Checklist

- [x] Every public method has an intention-revealing name (no process/handle/manage)
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new functions — readReagentSites() and the new GET handler are both single-line, complexity 1
- [x] No primitive obsession — Mirrors the existing flat-object convention used by every other endpoint in server.js (items, users, orders) — introducing a value object here would be inconsistent with the file's established style and out of scope per architecture-notes.md
- [x] No feature envy
- [x] No magic numbers or strings
- [x] Dependency direction verified: nothing in inner layers imports from outer layers — apps/mock-api has no Clean Architecture layering (flat server.js); not applicable per architecture-notes.md Layer Boundary Checks

## Simple Design Verification

- [x] Passes all tests — node --test apps/mock-api: 4/4 passing
- [x] Reveals intention — readReagentSites(), REAGENT_SITES_PATH, and the '============ REAGENT SITES ENDPOINTS ============' block name match the existing items/users/orders naming exactly
- [x] No duplication — Reuses the existing readJSON() helper rather than writing a new read path
- [x] Fewest elements — No new abstraction, class, or package added — matches architecture-notes.md's explicit instruction not to introduce a ReagentSite domain model in this slice

## Key Decisions

- Assigned each seeded site a kebab-case id (e.g. moonwell-spring) derived from its name, since acceptance-criteria.json's seed.reagentSites has no id field.: AC says each site needs an id, name, lat, lng; analysis.md task 1 specifies kebab-case slugging matching items.json's id convention (e.g. wand-embers).
- No requireAuth on GET /api/reagent-sites.: AC1/analysis.md explicitly require identical responses with or without a session — this is a public, read-only catalog endpoint by design.
- Wrote apps/mock-api/server.test.js myself as the developer rather than leaving it fully to qa-engineer.: The developer contract (ADR-008/009, 'Unit Tests Are Part of the Change') requires the developer to cover the lines they add/change with a test that can fail; the new endpoint and the listen-guard change are both covered. qa-engineer may still extend this suite with broader acceptance/integration coverage per the pipeline's normal division of labor.

## Deviations from Analysis

None

## Dependencies Added

None

## Notes for QA

- node --test apps/mock-api passes (4/4). Verified a mutant (lat 10 -> 11 on Moonwell Spring) turns 2 of the 4 subtests red, confirming the assertions are not vacuous.
- server.test.js reads the seeded expectation straight from docs/specs/reagent-map-spec/acceptance-criteria.json (seed.reagentSites) rather than hardcoding the 8 values, so it stays in sync with the spec package automatically.
- Verified with a live `npm start` run (PORT=8099) that the require.main === module guard still binds the port and serves /api/reagent-sites correctly outside the test harness, and that a garbage Authorization header returns the identical 200 payload.
- features/reagent-map/reagent-sites.feature's @AC1 @level-api scenario and the E1/F1 full scenario suites are explicitly out of scope for this order (run by hand per the spec package) — not covered by server.test.js.

## Notes for DevOps

None
