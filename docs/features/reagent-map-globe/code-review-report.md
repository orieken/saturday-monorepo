---
feature: "reagent-map-globe"
bounded_context: ""
domain_terms: ["Reagent Site", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/api/reagentSites.ts"]
issue_refs: ["reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Code Review: reagent-map-globe

## Overall Status

**APPROVED**

## Design Narrative

This diff adds a feature-flagged /reagent-map Vue view that follows LocationsView.vue's established per-view pattern: a new view component builds an offline Cesium Viewer sized to Aurelis's 6,000 km ellipsoid, fetches reagent sites through a new adapter module mirroring items.ts's conventions, and exposes the saturday-playwright-cesium seam contract on mount/unmount. The one prior blocking finding — a WebGL context leak on partial Promise.all failure — is now fixed with explicit per-settlement handling via Promise.allSettled. Behaviour change: yes.

## Design Score

- **Clarity** 4/5
- **Cohesion** 4/5
- **Coupling** 4/5
- **Craft** 4/5

## Security Surface

- GET /api/reagent-sites — read-only, same-origin, no auth, no user-supplied input
- Dev/build-time static file server (vite.config.ts cesiumStaticAssetsPlugin) serving node_modules/cesium/Build/Cesium under /cesium/ — has a path-traversal guard via normalize() + '..' prefix check; join() semantics prevent a leading-slash escape

## Performance Surface

- Single same-origin fetch to GET /api/reagent-sites, bounded by a 5s AbortController timeout; fixed 8-site dataset needs no pagination
- Cesium Viewer/WebGL context construction runs concurrently with the sites fetch via Promise.allSettled and is now destroyed on every failure path (viewer-build rejection has nothing to clean up; sites-fetch rejection destroys the already-built viewer before rethrowing), as well as on normal unmount via onBeforeUnmount

## Test Design Review

- No tests were added or modified by this diff. apps/ye-olde-magic-shop/tests/reagent-map-feature-flag.spec.ts and tests/reagent-map-globe.spec.ts are analysis.md's QA Tasks for the next pipeline stage and do not exist in the working tree yet — qa-engineer has not run.

## Verification of Developer Self-Review

- 'Guard Clause (replace implicit happy-path destructuring with explicit rejection checks)' on ReagentMapView.vue's onMounted — verified against the code: Promise.allSettled replaces Promise.all; viewerResult rejection rethrows with nothing to clean up (correct, buildViewer() never resolved); sitesResult rejection calls newViewer.destroy() before rethrowing. Closes the leak exactly as described.
- 'Cyclomatic complexity < 7' claim — verified true by inspection: two independent, unnested rejection guards, no loops.
- 'No retry or fallback logic was added' — verified true: both rejection branches rethrow, nothing is swallowed or retried.
- 'Passes all tests' with no unit-test runner wired — honest account of a known, named gap, consistent with the prior round.
- Non-blocking DRY finding on reagentSites.ts (API-base resolution duplicated from items.ts) — still present, still correctly deferred per architecture-notes.md, unchanged from the prior round.

## Feedback for the Developer

### 1. Extract Function (shared API base resolution) (non-blocking)
- **File**: `apps/ye-olde-magic-shop/src/api/reagentSites.ts`
- **Smell**: Lines 3-9 (runtimeHost/isLocal/API_BASE) duplicate items.ts's resolution logic verbatim — a DRY violation in newly written code.
- **Instruction**: Not blocking: architecture-notes.md already made the explicit call to mirror items.ts's convention rather than introduce a shared module for two call sites. Extract a shared resolveApiBase() helper the next time a third API client module needs this same resolution, rather than copying it again.
