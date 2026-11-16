---
feature: "reagent-map-line-tool"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: ["reagent-map-line-tool"]
linked_adrs: []
linked_kis: []
---

# Implementation Notes: reagent-map-line-tool

## Files Created

None

## Files Modified

- apps/ye-olde-magic-shop/src/views/ReagentMapView.vue

## Interface Design

- type SaturdayMetadata = { id: string; name: string; kind: 'site' | 'expedition' | 'draft-point' | 'draft-line' }
- function startLine(): void — resets draft state and arms the tool
- function handleMapClick(click: CesiumModule.ScreenSpaceEventHandler.PositionedEvent): void — registered once on viewer.screenSpaceEventHandler
- function pickLandPoint(position: CesiumModule.Cartesian2): GeoPoint | null — AC9/AC10 guard clauses
- function placePoint(point: GeoPoint): void
- function addDraftPointEntity(point: GeoPoint): CesiumModule.Entity
- function addDraftLineEntity(points: GeoPoint[]): CesiumModule.Entity
- function clearDraftEntities(): void
- function tagEntity(entity: CesiumModule.Entity, metadata: SaturdayMetadata): void — shared by all four saturdayMetadata assignment sites
- <button data-testid="start-line-tool" aria-label="Draw a new expedition"> — new template control, sibling of (not inside) header/#cesiumContainer/reagent-site-list/reagent-map-error

## Named Refactoring Log

- **Extract Function** on `ReagentMapView.vue: tagEntity(entity, metadata)` — addSiteEntities and addExpeditionEntities each inlined an identical `(entity as unknown as { saturdayMetadata: unknown }).saturdayMetadata = {...}` cast with an untyped literal. Extracting one typed helper removes the duplication (design-principles.md Rule 3) and gives all four call sites (including the two new draft helpers) a single checked SaturdayMetadata shape instead of four independent literals — satisfies architecture-notes.md's typing fitness function.

## Self-Review Checklist

- [x] Every public method has an intention-revealing name (no process/handle/manage) — handleMapClick is the one Cesium-convention exception (matches setInputAction's callback-naming idiom used nowhere else in this file, so it stays explicit about being an event handler).
- [x] No function exceeds 30 LOC — Longest new function (pickLandPoint) is 12 lines.
- [x] Cyclomatic complexity < 7 on all new functions — Highest is pickLandPoint/handleMapClick at 3 branches each.
- [x] No primitive obsession: domain concepts wrapped in value objects, not raw strings/ints — GeoPoint (imported, not redeclared) and SaturdayMetadata are the two domain shapes; MAX_LINE_POINTS replaces the literal 2.
- [x] No feature envy
- [x] No magic numbers or strings — MAX_LINE_POINTS constant added for the 2-point line-completion check.
- [x] Dependency direction verified: nothing in inner layers imports from outer layers — No layered architecture exists for this slice (confirmed in architecture-notes.md); all new code stays in the same Frameworks & Drivers (Vue view) layer as the code it extends.

## Simple Design Verification

- [x] Passes the tests — No unit-test runner exists in this app (Playwright/Cucumber only, per architecture-notes.md); `npm run build` (vite build) succeeds cleanly against the changed file. The Playwright spec proving AC6-AC10/AC14 is qa-engineer's deliverable per analysis.md's Task List split (Developer Tasks vs QA Tasks) and testing-conventions.md's E2E/UI row.
- [x] Reveals intention — startLine/clearDraftEntities/pickLandPoint/placePoint/addDraftPointEntity/addDraftLineEntity each name exactly the one thing they do.
- [x] No duplication — Extracted tagEntity to remove the saturdayMetadata-cast duplication across all four tagging call sites.
- [x] Fewest elements — No new file, module, store, or class — everything lives beside the existing addSiteEntities/addExpeditionEntities/flyToSite pattern in ReagentMapView.vue, per architecture-notes.md's structural decision against premature extraction.

## Key Decisions

- Kept all line-tool state, the click handler, and the draft-entity helpers inside ReagentMapView.vue with no new file.: Matches architecture-notes.md's structural decision 1 and developer handoff notes — no second consumer of this logic exists yet, and extracting it now would be premature generalization.
- Introduced a shared SaturdayMetadata type and a single tagEntity() helper used at all four saturdayMetadata assignment sites, replacing the two existing untyped inline-literal casts.: architecture-guardrails.md #4 forbids raw `any`; the existing per-call-site literal duplicated the cast and had no cross-checking of `kind` values across call sites. One typed helper removes that duplication per design-principles.md's No Duplication rule and architecture-notes.md's typing fitness function.
- Implemented pickLandPoint as two independent guard clauses — camera.pickEllipsoid returning undefined (AC9), then EllipsoidalOccluder.isPointVisible returning false (AC10) — using Cesium's own picking/occlusion APIs with no custom geometry math.: architecture-notes.md's structural decision 3: AC9/AC10 are two independent rejection checks, not branches of a shared conditional, so guard clauses are the correct shape per ARCHITECTURE_RULES.md and Simple Design Rule #4.
- Registered the LEFT_CLICK handler exactly once, in onMounted's happy path, with no new onBeforeUnmount cleanup.: viewer.destroy() already tears down its own ScreenSpaceEventHandler (architecture-notes.md structural decision 5); adding redundant cleanup would be dead code.
- Imported GeoPoint from ../api/expeditions.ts instead of declaring a second { lat, lng } type for draft points.: Keeps the draft-point shape identical to Expedition.from/to so reagent-map-save can consume it without a translation layer (architecture-notes.md structural decision 4).
- Placed the new 'Draw a new expedition' button in its own div, as a sibling between the error banner and the flex row — outside <header>, #cesiumContainer, reagent-site-list, and reagent-map-error.: reagent-map-visual.spec.ts's four @visual baselines each screenshot one specific locator's own element, not the full page, so a sibling addition elsewhere on the page cannot change any of their pixels — confirmed by reading that spec directly.

## Deviations from Analysis

- QA Tasks 1-7 (writing apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts and running/confirming the visual baselines): Not a developer task. analysis.md's Task List explicitly splits this into 'Developer Tasks' (11 items, all completed — ReagentMapView.vue only) and 'QA Tasks' (writing and running the Playwright spec, confirming baselines). testing-conventions.md's test-level table and this agent's own Rules section ('Acceptance and integration suites are the QA engineer's') assign E2E/UI-level Playwright specs to qa-engineer, not developer — this app has no unit-test runner (Vitest) for the developer to use instead (confirmed in architecture-notes.md). The implementation is written so every acceptance criterion (AC6, AC7, AC8, AC9, AC10, AC14) is directly verifiable by that spec once qa-engineer writes it.

## Dependencies Added

None

## Notes for QA

- pickLandPoint's AC9 path: viewer.camera.pickEllipsoid returns undefined for a click outside the rendered globe disc — use a canvas-corner pixel via #cesiumContainer's boundingBox(), as the context-manifest.md and architecture-notes.md both note cesium.clickAt() cannot express this (it throws when toScreen() is null).
- pickLandPoint's AC10 path: camera.pickEllipsoid can intersect the ellipsoid at a grazing/limb pixel where EllipsoidalOccluder.isPointVisible then reports false; per test-repair-contract.md, prove this guard actually fires by showing the test fails if the isPointVisible check is deleted.
- Every placed point/line is readable via cesium.getEntities() with saturdayMetadata.kind === 'draft-point' or 'draft-line' — id and name are both '' on draft entities (never persisted), unlike 'site'/'expedition' entities which carry real id/name.
- AC14: clicking #start-line-tool again calls startLine(), which removes every existing draft-point/draft-line entity before resetting draftPoints to empty — assert getEntities() returns zero draft-* entities after a second click with one point already placed.
- The button's data-testid is 'start-line-tool' and its aria-label is 'Draw a new expedition', matching the fly-to-site button's convention (data-testid + aria-label, no visual-style distinction).
- No new baseline should change: the button was placed outside all four regions reagent-map-visual.spec.ts screenshots (header, reagent-site-list, reagent-map-error, #cesiumContainer), each captured as that locator's own element screenshot, not a full-page screenshot.

## Notes for DevOps

None
