---
feature: "reagent-map-line-tool"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: ["reagent-map-line-tool"]
linked_adrs: []
linked_kis: []
---

# Code Review: reagent-map-line-tool

## Overall Status

**APPROVED**

## Design Narrative

The implementation adds line-drawing to ReagentMapView.vue by extending the file's existing per-entity Cesium pattern: a single LEFT_CLICK handler registered once in onMounted resolves a screen click to a geographic point via camera.pickEllipsoid followed by an EllipsoidalOccluder visibility check (two independent guard clauses implementing the off-globe and far-side rejection paths), accumulates up to two draft points as component-local refs, and tags every entity -- including the two new draft kinds -- through one shared typed tagEntity() helper that replaces the previous per-call-site untyped cast. All new state and logic stay in the same file and layer as the code they extend, with no new module, store, or persisted state, exactly as architecture-notes.md decided. Behaviour change: yes.

## Design Score

- **Clarity** 5/5
- **Cohesion** 5/5
- **Coupling** 5/5
- **Craft** 4/5

## Security Surface

None

## Performance Surface

None

## Test Design Review

- No test file was added or modified in this diff -- implementation-notes.md's Files Modified lists only ReagentMapView.vue, and analysis.md's Task List assigns reagent-map-line-tool.spec.ts to qa-engineer, not developer. There is nothing to judge for 'would it fail' in this stage's diff.

## Verification of Developer Self-Review

- Claim 'No function exceeds 30 LOC' -- confirmed; longest new function (pickLandPoint) is 12 lines.
- Claim 'Cyclomatic complexity < 7 on all new functions' -- confirmed; pickLandPoint and handleMapClick each have 3 branches, the rest fewer.
- Claim 'No duplication -- extracted tagEntity()' -- confirmed; the same cast-and-assign pattern previously inlined at addSiteEntities and addExpeditionEntities is now one helper used at all four call sites.
- Claim 'No magic numbers or strings -- MAX_LINE_POINTS' -- confirmed; the literal 2 is replaced by the named constant at both its use sites.
- The checklist does not surface the non-null-assertion pattern in addDraftPointEntity/addDraftLineEntity (see findings) -- not a false claim, since those items it does check are accurate, but it is the one place a closer look finds something the checklist didn't flag.

## Feedback for the Developer

### 1. Introduce Guard Clause (non-blocking)
- **File**: `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue`
- **Smell**: addDraftPointEntity (line 213) and addDraftLineEntity (line 222) reach straight for viewer!, cesium!, ellipsoid! non-null assertions, while flyToSite (line 194) and pickLandPoint (line 242) -- the file's own established convention -- both open with an explicit `if (!viewer || !cesium || !ellipsoid) return` guard clause. The two new functions are only safe today because handleMapClick only calls placePoint (which calls them) after pickLandPoint has already confirmed all three are non-null -- an invariant that holds by the current call chain but isn't visible locally in either function.
- **Instruction**: Either pass viewer/cesium/ellipsoid into addDraftPointEntity/addDraftLineEntity as already-narrowed parameters from placePoint, or add the same `if (!viewer || !cesium || !ellipsoid) return null` guard used elsewhere in this file and widen the return type to `CesiumModule.Entity | null`. Non-blocking: the current invariant is sound and the functions are small enough that the risk is low, but matching the house convention removes the implicit cross-function dependency.
