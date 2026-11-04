---
feature: "reagent-sites-api"
bounded_context: ""
domain_terms: ["Reagent Site", "Aurelis", "Adventurer", "Expedition"]
files_touched: ["apps/mock-api/server.test.js"]
issue_refs: ["reagent-sites-api"]
linked_adrs: []
linked_kis: []
---

# QA Report: reagent-sites-api

## Test Files Created

None

## Test Files Modified

- apps/mock-api/server.test.js

## Coverage Summary

- Acceptance criteria covered: 4/4
- Total new tests: 2
- Statement coverage, lowest unit: **47.9%** (`apps/mock-api/server.js (whole file, including pre-existing untested auth/orders/items endpoints)`)
  - `apps/mock-api/server.js (feature-added lines: REAGENT_SITES_PATH, readReagentSites, GET /api/reagent-sites)`: 100.0%
  - `apps/mock-api/server.js (whole file, including pre-existing untested auth/orders/items endpoints)`: 47.9%

## Test Results

- Passed: 5
- Failed: 0
- Skipped: 0

## Accessibility Check

None

## Bugs Found

- A QA-added subtest for the 'valid session' case registered a real user via POST /api/auth/register, which persists to apps/mock-api/data/users.json — leaving a non-idempotent, repo-polluting side effect on every test run (confirmed via git status after running the suite).: Snapshot users.json before the test group and restore it in t.after(), so the suite is idempotent and leaves no diff. Verified by running the suite twice and confirming git status shows no change to users.json.

## Known Gaps

- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail: The require.main === module startup block (app.listen + its console.log lines, server.js:247-263) is structurally unreachable under node:test, since tests require() the module rather than execute it as the process entrypoint. This is not feature logic (no branching, complexity 1) and was verified manually instead: the developer ran `npm start` on an ephemeral port and confirmed the server binds and GET /api/reagent-sites responds correctly outside the test harness. Automated line coverage cannot reach this block regardless of how the feature is tested, since it is mutually exclusive with being imported as a module.

## Notes for Tech Writer

- New endpoint: GET /api/reagent-sites — no auth required, returns all 8 seeded reagent sites unchanged regardless of session state (no, valid, or garbage Authorization header).
- Seed data lives at apps/mock-api/data/reagent-sites.json; each site has id/name/lat/lng and is kept in sync with docs/specs/reagent-map-spec/acceptance-criteria.json's seed.reagentSites (name/lat/lng values match exactly; ids are mock-api-local).
- apps/mock-api now runs its tests via `npm test` (node --test), added with zero new dependencies.
- Out of scope per the spec: expeditions-api and the UI/E2E scenario suites (features/reagent-map/reagent-sites.feature @AC1 @level-api and the E1/F1 suites), run by hand.
