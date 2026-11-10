---
feature: "local-fonts"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
---

# Code Review: local-fonts

## Overall Status

**APPROVED**

## Design Narrative

The change self-hosts Cinzel and Inter via @fontsource, imported as side-effect CSS in main.ts at the same site as the existing leaflet.css import — no new abstraction, no new layer. The cross-origin @import is removed from styles.css, the test's PRE_EXISTING_OFF_ORIGIN_HOSTS exception is deleted so the no-request test has no fallback, and a new Playwright test proves the fonts actually render at the exact pinned weight using document.fonts.load rather than the false-positive-prone document.fonts.check. Behaviour change: no.

## Design Score

- **Clarity** 5/5
- **Cohesion** 5/5
- **Coupling** 5/5
- **Craft** 5/5

## Security Surface

None

## Performance Surface

- Removes a cross-origin network request on every route (fonts.googleapis.com/fonts.gstatic.com) — strictly reduces external dependencies, no new surface introduced.

## Test Design Review

- tests/local-fonts.spec.ts 'headings render in self-hosted Cinzel and body text in self-hosted Inter': YES — asserts document.fonts.load matched face weight equals the requested weight per family/weight pair; developer's notesForQa record removing all four imports (fails on all four) and removing only inter/600 (fails on exactly that one) as evidence it is not vacuous and does not rely on document.fonts.check, which the analysis documented as a false positive.
- tests/reagent-map-globe.spec.ts no-request test (modified, not added): PRE_EXISTING_OFF_ORIGIN_HOSTS and its skip branch are removed, so the existing assertion now has no fallback for Google Fonts hosts — YES, consistent with AC2/AC3's requirement that it fail if the old @import is restored. Developer could not execute this test in-sandbox (needs the Docker runner's real Chrome/WebGL2 context per playwright.config.ts) and flagged that QA must capture the authoritative evidence — this is correctly deferred to qa-engineer, not a gap in this review.

## Verification of Developer Self-Review

- Function-length and complexity claims hold — the only new logic is the test file's single async callback with one ternary-equivalent conditional per font pair (complexity 2), well under the <7 budget.
- 'No duplication' claim holds — FONTS array drives all four assertions instead of four repeated blocks, matches the file as read.
- 'Dependency direction verified' claim holds — font imports sit at the existing side-effect-import site in main.ts (same layer, same pattern as leaflet.css), verified directly against main.ts:5-10.
- No discrepancies found between the self-review checklist and the actual diff.

## Feedback for the Developer

None
