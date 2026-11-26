---
feature: "golden-test-ids"
bounded_context: ""
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/tests/locator-conventions.spec.ts"]
issue_refs: ["golden-test-ids"]
linked_adrs: []
linked_kis: []
---

# QA Report: golden-test-ids

## Test Files Created

- apps/ye-olde-magic-shop/tests/locator-conventions.spec.ts

## Test Files Modified

None

## Coverage Summary

- Acceptance criteria covered: 5/5
- Total new tests: 3

## Test Results

- Passed: 118
- Failed: 0
- Skipped: 0

## Accessibility Check

- No accessibility impact: the only app-markup change is a data-testid attribute rename on CartView.vue's existing complete-transaction button — no change to markup, ARIA attributes, visible text, or tab order. HomePage.ts/CartPage.ts are test-only page objects, not rendered UI.

## Bugs Found

None

## Known Gaps

- No file outside scope.paths changes, other than this run's delivery record (AC5): analysis.md's qaTasks explicitly required adding a regression test scanning lib/site/**/*.ts for has-text(/:text(/:nth-/nested selectors, which this stage added at apps/ye-olde-magic-shop/tests/locator-conventions.spec.ts. That path is not literally inside the claim's declared scope.paths (apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/lib/**, docs/**) — Playwright's testDir is fixed to ./tests in playwright.config.ts, so there is no in-scope path where this test would actually be discovered and run. Surfacing the contradiction between the explicit qaTask and the literal scope.paths wording rather than resolving it silently; the test itself passed on both runner images.
- Every scenario that passed before this change passes after it (AC4): tests/reagent-map-pan-zoom.spec.ts's 'panning and zooming move the camera, never the sites' timed out once (cesium.waitForSettled() camera timeout) during the first full debian-runner-image Playwright run, under 3-worker contention. It passed in isolated re-run and passed cleanly in the full ubi9 run (102/102). This file is untouched by this change (no reagent-map/Cesium code in diff) and the failure mode matches known software-GL timing sensitivity under Docker emulation on this host, not a locator regression. Noting it rather than silently treating it as a clean pass.

## Notes for Tech Writer

- This is an internal test-hygiene chore (locator stability) with no user-facing behavior change — likely no README/API doc updates needed.
- A new standing fitness-function test exists at apps/ye-olde-magic-shop/tests/locator-conventions.spec.ts: it fails the build if any lib/site/ locator string regresses to has-text(/:text(/a CSS :nth- pseudo-class/a nested descendant-combinator selector, and pins CartPage's and HomePage's specific testid locators. Worth a one-line mention in any Saturday testing-conventions doc for this app if one exists.
- Flag for delivery sign-off: see knownGaps entry about the new test file sitting outside the claim's literal scope.paths.
