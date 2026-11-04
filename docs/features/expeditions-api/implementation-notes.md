---
feature: "expeditions-api"
bounded_context: ""
domain_terms: ["expedition", "adventurer", "Aurelis"]
files_touched: ["apps/mock-api/server.test.js"]
issue_refs: ["expeditions-api"]
linked_adrs: []
linked_kis: []
---

# Implementation Notes: expeditions-api

## Files Created

None

## Files Modified

- apps/mock-api/server.test.js

## Interface Design

None

## Named Refactoring Log

None

## Self-Review Checklist

- [x] Every public method has an intention-revealing name (no process/handle/manage)
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new functions
- [x] No primitive obsession: domain concepts wrapped in value objects, not raw strings/ints — Test literals mirror existing test style ({lat, lng} point objects already used throughout the suite).
- [x] No feature envy: methods use their own class's data more than another's
- [x] No magic numbers or strings: all literals are named constants — Added local AURELIS_RADIUS_KM and MIN_EXPEDITION_DISTANCE_KM consts mirroring server.js's unexported constants, with a comment explaining the mirroring.
- [x] Dependency direction verified: nothing in inner layers imports from outer layers — Test-only change; no production code touched.

## Simple Design Verification

- [x] Passes all tests — 18/18 pass via `npm test` (node --test) in apps/mock-api.
- [x] Reveals intention — Each new test name states the exact behavior and session/body condition under test, matching the file's existing naming convention.
- [x] No duplication — New tests reuse the existing loginAs/userIdFor/get/post helpers and the file's existing snapshot/restore (t.after) pattern rather than introducing new scaffolding.
- [x] Fewest elements — No new helpers were introduced beyond two local constants mirroring server.js; everything else reuses what the file already had.

## Key Decisions

- Construct the exact-1km boundary test by nudging the north-south offset by a 1e-9 degree epsilon rather than targeting the mathematically exact dLat.: Verified empirically (node -e) that the mathematically exact offset (MIN_EXPEDITION_DISTANCE_KM / AURELIS_RADIUS_KM radians) lands the haversine computation about 4.5e-14 km *under* 1km due to floating-point error in sin/cos, which the strict less-than check then rejects (422) — a false failure of the test, not of the feature. The epsilon pushes the computed distance to ~1.0000001 km, safely clear of that floating noise while still pinning the strict less-than boundary semantics the code reviewer asked for.
- Did not export MIN_EXPEDITION_DISTANCE_KM/AURELIS_RADIUS_KM from server.js; declared equivalent consts locally in the test file instead.: The existing geodesicDistanceKm test already hardcodes the radius (6000) rather than importing it, so mirroring that convention (with a comment noting the duplication) keeps the test file consistent with its own established style rather than widening server.js's exported surface for a test-only need.
- Added a mirrored Borin-session GET test and an invalid-token GET test as new t.test cases inside the existing 'expeditions endpoints' test, rather than a new top-level test block.: They exercise the same seeded fixtures and setup/teardown (snapshot + restore of expeditions.json and users.json) already established for that block; a new top-level test would have duplicated that scaffolding for no benefit.

## Deviations from Analysis

None

## Dependencies Added

None

## Notes for QA

- All four code-reviewer blocking/non-blocking findings addressed in apps/mock-api/server.test.js: invalid-token GET coverage, malformed-body (missing field + non-numeric) 400 coverage, mirrored Borin-session GET assertion, and an exact-1km boundary POST test.
- The pagination gap on GET /api/expeditions (architecture-guardrails.md #6) was explicitly called out by code-reviewer as pre-existing and out of scope for this feature (it mirrors GET /api/items and GET /api/reagent-sites) — not addressed here, left for a sitewide follow-up ticket per the reviewer's own note.
- Full suite verified green: 18/18 passing via `npm test` in apps/mock-api.

## Notes for DevOps

None
