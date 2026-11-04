---
feature: "reagent-map-globe"
bounded_context: ""
domain_terms: ["ReagentSite", "Viewer"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Implementation Notes: reagent-map-globe

## Files Created

None

## Files Modified

- apps/ye-olde-magic-shop/src/views/ReagentMapView.vue

## Interface Design

- onMounted's viewer/sites bootstrap now resolves via Promise.allSettled with explicit per-settlement status narrowing instead of Promise.all — no new public surface, internal control flow only.

## Named Refactoring Log

- **Guard Clause (replace implicit happy-path destructuring with explicit rejection checks)** on `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue:onMounted` — code-reviewer finding: Promise.all discarded a fully-constructed Cesium Viewer (WebGL context, canvas, render loop) with no reference to destroy it whenever fetchReagentSites() rejected after buildViewer() had already resolved, leaking a WebGL context per failed visit. Switched to Promise.allSettled, checked viewerResult.status and sitesResult.status explicitly, and call newViewer.destroy() on a sites-fetch rejection before rethrowing, so a fetch failure never orphans a live viewer. No retry or fallback logic was added.

## Self-Review Checklist

- [x] Every public method has an intention-revealing name
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new functions — onMounted gained two guard-clause branches; still well under 7.
- [x] No primitive obsession
- [x] No feature envy
- [x] No magic numbers or strings
- [x] Dependency direction verified — No import changes.

## Simple Design Verification

- [x] Passes all tests — No unit-test framework is wired for this app (E2E/BDD only, per package.json); tsc/eslint binaries aren't installed in this environment, so the fix was verified by manual type-narrowing inspection rather than a tool run.
- [x] Reveals intention — viewerResult/sitesResult names and explicit 'rejected' checks make the cleanup-on-partial-failure path explicit instead of implicit in a destructuring Promise.all.
- [x] No duplication
- [x] Fewest elements — No new helpers added — reused buildViewer/fetchReagentSites, changed only the combinator and added two guard clauses.

## Key Decisions

- Used Promise.allSettled with per-result status checks rather than wrapping the original Promise.all in try/catch.: allSettled gives independent settlement for buildViewer() and fetchReagentSites() without needing a catch block that would have to re-inspect which of the two promises actually failed; a rejected viewerResult has nothing to clean up and can just rethrow, while a rejected sitesResult can call newViewer.destroy() on the already-resolved viewer before rethrowing — exactly the cleanup the code-reviewer asked for, with no retry/fallback behavior added.

## Deviations from Analysis

None

## Dependencies Added

None

## Notes for QA

- Verify: if fetchReagentSites() rejects (e.g. mock API down or past its 5s timeout) after buildViewer() has resolved, the Cesium viewer's destroy() is called and no WebGL context or #cesiumContainer canvas is left behind — repeat the failing-fetch scenario across multiple mounts to confirm no accumulation.
- Verify: if buildViewer() itself rejects, onMounted simply rethrows (no destroy call attempted on a viewer that never existed) and the component shows no entities/sites.
- The non-blocking DRY finding on apps/ye-olde-magic-shop/src/api/reagentSites.ts (duplicated API-base resolution with items.ts) was explicitly deferred per architecture-notes.md and intentionally left unchanged.

## Notes for DevOps

None
