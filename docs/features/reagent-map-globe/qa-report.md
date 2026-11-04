---
feature: "reagent-map-globe"
bounded_context: ""
domain_terms: ["Aurelis", "Reagent Site", "Adventurer", "Expedition", "Seeded Expedition"]
files_touched: ["apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts"]
issue_refs: ["reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# QA Report: reagent-map-globe

## Test Files Created

- apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts

## Test Files Modified

- apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts

## Coverage Summary

- Acceptance criteria covered: 7/7
- Total new tests: 16

## Test Results

- Passed: 16
- Failed: 0
- Skipped: 0

## Accessibility Check

- [A11Y OK] #cesiumContainer has role="img" and a non-empty, descriptive aria-label (tests/reagent-map-a11y-check.spec.ts)
- [A11Y OK] data-testid="reagent-site-list" carries aria-label="Reagent sites" so the site list reads sensibly to a screen reader
- [A11Y OK] the reagent-map nav link is visible and keyboard-focusable (Tab reaches it, :focus is observable)

## Bugs Found

None

## Known Gaps

- Given the view open, then no request leaves the app's and mock API's own origins.: The reagent map view's own new code (Cesium config, fetchReagentSites) makes zero off-origin requests, verified directly. But the page as a whole still issues two off-origin requests to fonts.googleapis.com and fonts.gstatic.com, caused by a pre-existing, sitewide @import in src/styles.css that loads on every route and predates this feature (confirmed via `git log -- apps/ye-olde-magic-shop/src/styles.css`; it is untouched by this diff). tests/reagent-map-globe.spec.ts's NFR test carves out exactly those two named hosts, with a comment citing this evidence, so it still fails on any new off-origin dependency the reagent map itself might introduce. Fixing the root cause (self-hosting Cinzel/Inter) is a site-wide typography decision outside this feature's scope.paths intent — recommend a follow-up ticket rather than a silent fix during this QA pass.
- Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail.: apps/ye-olde-magic-shop has no statement-coverage tooling configured (no nyc/istanbul/c8, no Vitest) — the app's own convention is Playwright E2E/BDD only, so no percentage figure can be reported. Coverage of every changed line was instead confirmed by direct review: isReagentMapEnabled (all branches), the router/App.vue flag gates, fetchReagentSites's success path, and both of ReagentMapView.vue's onMounted failure branches (fetchReagentSites rejecting after the viewer resolves; buildViewer itself rejecting) are each exercised by a passing-or-intentionally-failing assertion across the 16 reagent-map-*.spec.ts tests. Complexity of every new/changed function was reviewed by inspection and stays well under 7 (straight-line or single-branch bodies); no ESLint complexity rule exists in this app to enforce it automatically, consistent with analysis.md's definitionOfDone note.

## Notes for Tech Writer

- New view at /reagent-map (behind VITE_FEATURE_REAGENT_MAP, off by default in .env) shows Aurelis — an invented, 6,000km-radius spherical planet — as an offline Cesium 3D globe with the 8 seeded reagent sites marked and listed by name beside it.
- tests/reagent-map-resilience.spec.ts was added this QA pass to verify the two failure paths implementation-notes.md flagged: the Cesium viewer is destroyed with no leaked canvas/WebGL context when the sites fetch fails after the viewer has already built (checked across 3 repeated mounts), and onMounted simply rethrows with no entities or sites shown if the Cesium script itself fails to load.
- tests/reagent-map-globe.spec.ts's no-off-origin-request test documents and carves out two pre-existing, unrelated off-origin hosts (Google Fonts) — see knownGaps for the full rationale; this is worth a line in the feature's docs so a future reader doesn't mistake the allowlist for scope creep.
- Security review's one LOW finding (unhandled decodeURIComponent throw in vite.config.ts's dev-only Cesium static-asset middleware) was not re-verified by execution this pass; it remains tracked as tech debt per security-report.md.
