---
feature: "leaflet-local-tiles"
bounded_context: ""
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/tests/map.spec.ts"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Implementation Notes: leaflet-local-tiles

## Files Created

None

## Files Modified

- apps/ye-olde-magic-shop/tests/map.spec.ts

## Interface Design

None

## Named Refactoring Log

- **Replace Vacuous Check with Direct Assertion** on `apps/ye-olde-magic-shop/tests/map.spec.ts — 'locations view loads with every off-origin request blocked (NFR)'` — checkTileFailures() only inspects .leaflet-tile elements and cannot detect a failed/off-origin L.imageOverlay load, so the test could not fail against the regression it exists to prevent. Added a direct assertion on document.querySelector('.leaflet-image-layer') checking img.complete && img.naturalWidth > 0, keeping checkTileFailures() as a secondary guard.

## Self-Review Checklist

- [x] Every public method has an intention-revealing name
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new functions
- [x] No primitive obsession
- [x] No feature envy
- [x] No magic numbers or strings
- [x] Dependency direction verified

## Simple Design Verification

- [x] Passes all tests — All 10 tests in tests/map.spec.ts pass (playwright test --project=chromium tests/map.spec.ts).
- [x] Reveals intention — Assertion name imageOverlayLoaded states what it verifies; inline comment marks checkTileFailures() as the secondary guard.
- [x] No duplication
- [x] Fewest elements — Single page.evaluate() added; no new helper/abstraction introduced since this is a single-file, single-use check.

## Key Decisions

- Verify the image overlay directly via DOM query inside the off-origin-blocking test rather than extending the shared checkTileFailures() fixture helper.: code-reviewer's instruction explicitly scoped the fix to this test file, not the shared fixture (fixture.ts is outside scope.paths). Keeping checkTileFailures() unchanged avoids touching code other tests depend on.

## Deviations from Analysis

- Document known coverage gap in 'should not have broken tiles' (map.spec.ts:97): code-reviewer asked that this test's text and assertion remain untouched (per the spec's 'passes unchanged' requirement) and that the now-vacuous tile-check be recorded as a known, accepted limitation by tech-writer in delivery notes, rather than fixed by the developer. No code change made for this finding; flagging here for tech-writer to pick up.

## Dependencies Added

None

## Notes for QA

- Verified the new '.leaflet-image-layer' assertion actually catches the regression: temporarily pointed the image overlay at a nonexistent file and confirmed the off-origin-blocking test fails (Expected: true, Received: false), then restored the source and confirmed all 10 tests pass again.
- Known accepted gap (not fixed in this pass, per code-reviewer instruction): 'should not have broken tiles' (map.spec.ts:97) only queries .leaflet-tile, which no longer exists since the feature replaced the OSM tile layer with an L.imageOverlay — that assertion is now vacuous (always passes). The off-origin-blocking test's new direct image-overlay assertion is the one that actually covers this. tech-writer should record this as a documented trade-off in delivery notes.

## Notes for DevOps

None
