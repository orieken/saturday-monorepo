---
feature: "leaflet-local-tiles"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: ["leaflet-local-tiles"]
linked_adrs: []
linked_kis: []
---

# QA Report: leaflet-local-tiles

## Test Files Created

None

## Test Files Modified

None

## Coverage Summary

- Acceptance criteria covered: 6/6
- Total new tests: 1
- Statement coverage, lowest unit: **100.0%** (`apps/ye-olde-magic-shop/src/views/LocationsView.vue onMounted (no coverage-instrumentation tool wired for this app's Playwright/Cucumber-only suite; measured by line-execution inspection — every added statement runs unconditionally on every /locations page load, exercised by all 10 tests in map.spec.ts)`)
  - `apps/ye-olde-magic-shop/src/views/LocationsView.vue onMounted (no coverage-instrumentation tool wired for this app's Playwright/Cucumber-only suite; measured by line-execution inspection — every added statement runs unconditionally on every /locations page load, exercised by all 10 tests in map.spec.ts)`: 100.0%

## Test Results

- Passed: 83
- Failed: 0
- Skipped: 0

## Accessibility Check

- No new a11y review needed from QA: accessibility-engineer.md (pre-existing artifact) already reviewed this exact diff and returned APPROVED.
- Verified against Leaflet 1.9.4 source: L.imageOverlay defaults to alt="" and interactive=false, so the new <img class="leaflet-image-layer"> stays out of the tab order and is correctly marked decorative (shop info lives in marker popups, not the base image).
- No template markup changed; the 'Find Nearest Shop' button and existing keyboard-pannable map/controls are unaffected.

## Bugs Found

None

## Known Gaps

- Every test in tests/map.spec.ts MUST pass unchanged under the default Leaflet CRS.: 'should not have broken tiles' (map.spec.ts:97) only queries '.leaflet-tile' elements, which no longer exist since the OSM tile layer was replaced by L.imageOverlay — the assertion is now vacuously true (checkTileFailures() always returns false on an empty NodeList). The spec requires this test's text stay unchanged, and code-reviewer confirmed the real regression coverage now lives in the new NFR test (map.spec.ts:140), which directly asserts the image overlay's own <img> loaded (complete && naturalWidth > 0) under full off-origin blocking. Verified independently: a test run against the full suite confirms both tests still pass, and the developer documented breaking the image path to confirm the new test (not the old one) catches the regression.

## Notes for Tech Writer

- Document the known, accepted limitation on tests/map.spec.ts:97 ('should not have broken tiles'): it's now a vacuous check against the removed tile layer, kept unchanged per the spec's 'passes unchanged' requirement; the new off-origin-blocking test (map.spec.ts:140) is the one that actually verifies local-image loading under a zero-egress network.
- The map image's licence record sits at apps/ye-olde-magic-shop/public/sword-coast-map.LICENSE.txt: an original in-repo SVG creation by Oscar Rieken, CC BY 4.0 — no third-party redistribution risk to call out.
- Full pre-ship suite (scripts/ui-suite-in-chrome.sh, debian image) passed 83/83, confirming zero regressions in the magic-shop, cart, heatmap, and reagent-map suites from this change.
