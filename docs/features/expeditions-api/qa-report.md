---
feature: "expeditions-api"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/mock-api/server.test.js"]
issue_refs: ["expeditions-api"]
linked_adrs: []
linked_kis: []
---

# QA Report: expeditions-api

## Test Files Created

None

## Test Files Modified

- apps/mock-api/server.test.js

## Coverage Summary

- Acceptance criteria covered: 7/8
- Total new tests: 15
- Statement coverage, lowest unit: **69.4%** (`apps/mock-api/server.js (whole file, pre-existing + new routes)`)
  - `apps/mock-api/server.js (whole file, pre-existing + new routes)`: 69.4%
  - `apps/mock-api/server.js (expeditions-api changed lines only)`: 100.0%

## Test Results

- Passed: 21
- Failed: 0
- Skipped: 0

## Accessibility Check

None

## Bugs Found

- isValidPoint (apps/mock-api/server.js:99) has cyclomatic complexity of 7, measured with ESLint's complexity rule (max allowed 6, i.e. the framework's '< 7' budget per shared/rules/testing-conventions.md and design-principles.md). It is a single boolean return with six && operators (point truthiness, two Number.isFinite calls, and four range comparisons added by the security-reviewer's fix for the NaN/Infinity bypass). Verified independently: `eslint --rule complexity:[error,6]` flags exactly this one function in the whole file; every other new/changed function (geodesicDistanceKm, resolveSessionUserId, the GET/POST /api/expeditions handlers) measures well under the cap.: Not fixed by QA — this is a structural/complexity violation, not a behavioral defect discovered through a failing test, so it falls outside QA's mandate to modify source only for bugs found while testing. Flagged here for the developer/code-reviewer to extract the range check (e.g. a separate isValidLatitude/isValidLongitude helper) into a second function, which would bring both under the cap without changing behavior. No test needs to change for that refactor — all 21 tests, including the three new 400-rejection tests below, pin isValidPoint's current behavior and would catch a regression.

## Known Gaps

- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change.: Coverage and file-scope portions of this criterion are fully met (verified: all changed lines in server.js fall outside ESLint's reported uncovered-line ranges; git status shows changes confined to apps/mock-api/** and docs/features/expeditions-api/**). The complexity portion is not met for one function, isValidPoint, which measures 7 against the <7 budget — see bugsFound. Marking this criterion as not-yet-fully-covered rather than silently rounding it up.

## Notes for Tech Writer

- New endpoints to document: GET /api/expeditions (optional-auth: seeded 'shop' expeditions always visible, plus the caller's own when a valid session is present) and POST /api/expeditions (requires a session; body is { from: {lat,lng}, to: {lat,lng} }; owner is always taken from the session, any owner/userId field in the body is ignored; 401 with no session, 400 for a malformed/non-finite/out-of-range point, 422 for a zero-length line under 1 km).
- New seed file apps/mock-api/data/expeditions.json holds the three shop-owned expeditions exactly as specified in docs/specs/reagent-map-spec/acceptance-criteria.json's seed.seededExpeditions, resolved to concrete lat/lng via reagent-sites.json. apps/mock-api/data/users.json gained two new seed adventurers, Elara (user-4) and Borin (user-5), each with one pre-saved expedition matching seed.adventurers.
- Security hardening applied by security-reviewer and verified by QA: isValidPoint now rejects NaN/Infinity (via Number.isFinite, not typeof) and out-of-range lat/lng; the GET owner filter requires userId !== null before matching, closing a latent null-coincidence path. Both fixes are covered by tests that fail against the pre-fix code (verified directly by temporarily reverting the fix and confirming the new tests catch it).
- Known pre-existing gap, out of scope for this feature: GET /api/expeditions returns an unbounded, unpaginated array, same as GET /api/items and GET /api/reagent-sites — flagged by code-reviewer as a sitewide follow-up, not addressed here.
- A minor complexity cleanup is outstanding on isValidPoint (complexity 7, cap is 6) — see bugsFound/knownGaps. Does not block shipping on its own but should be swept up by the next touch to that function.
