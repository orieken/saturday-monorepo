---
feature: "reagent-map-line-tool"
bounded_context: "Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7"
domain_terms: ["Expedition", "Reagent Site", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts"]
issue_refs: ["reagent-map-line-tool"]
linked_adrs: []
linked_kis: []
---

# Architecture Notes: reagent-map-line-tool

## Structural Decisions

### Keep all new line-tool code (toolStarted/draftPoints state, the LEFT_CLICK handler, pickLandPoint, addDraftPointEntity, addDraftLineEntity, clearDraftEntities, and the SaturdayMetadata type) inside ReagentMapView.vue. No new file, module, composable, or package is created.
- **Rationale**: Extends the exact pattern reagent-map-expeditions already established in this file (addSiteEntities/addExpeditionEntities/flyToSite as plain functions over module-scope viewer/cesium/ellipsoid). This app has no Vitest/unit-test runner — only Playwright and Cucumber — so extracting pickLandPoint into a separately-testable module would buy no real testability without adding new test infrastructure, which this order does not request. Analyst flagged uncertainty about where click-to-geo-point logic should live; resolving that now: it lives beside the other per-entity helpers.
- **Trade-offs**:
  - pickLandPoint's two rejection paths (AC9/AC10) are only exercised through Playwright + real WebGL, not fast unit tests
  - ReagentMapView.vue grows from ~220 lines to roughly 320-350 — still well short of anything requiring a split, but a third interactive tool should trigger re-evaluation (see Refactoring Opportunities)
- **Fitness Function**: No new src/ file or module is introduced for line-tool logic outside ReagentMapView.vue and its spec file — verified by code-reviewer checks the diff's file list for this feature; any new file requires a stated reason before approval
### Introduce one shared type, SaturdayMetadata = { id: string; name: string; kind: 'site' | 'expedition' | 'draft-point' | 'draft-line' }, in ReagentMapView.vue and use it at all four tagging call sites (addSiteEntities, addExpeditionEntities, addDraftPointEntity, addDraftLineEntity) in place of the existing `(entity as unknown as { saturdayMetadata: unknown })` cast's untyped literal.
- **Rationale**: architecture-guardrails.md #4 forbids raw `any`; the existing cast already correctly routes through `unknown`, but the `kind` string literals are not checked against each other today. reagent-map-expeditions already accepted duplicating this cast once (its architecture-notes.md deferred the extraction as 'not justified by this order alone'); this order adds a third and fourth call site with two more kind values, which is the point design-principles.md's 'No Duplication' rule intends to catch.
- **Trade-offs**:
  - One more type to extend the day a fifth entity kind (e.g. a future polygon tool) arrives
  - Still routes through a cast at the assignment point, since Cesium's own Entity type has no saturdayMetadata field — this narrows what can be assigned, it doesn't remove the cast
- **Fitness Function**: Every saturdayMetadata assignment in ReagentMapView.vue is typed via the shared SaturdayMetadata type, not an inline object literal — verified by code-reviewer confirms by diff review — no ESLint config exists anywhere in this repo yet, so this is manual judgment, not an automated gate
### Resolve a screen click to a geographic point with a single pickLandPoint(position) function using Cesium's native camera.pickEllipsoid followed by `new EllipsoidalOccluder(ellipsoid, camera.positionWC).isPointVisible(...)`, written as two independent guard clauses (off-globe, far-side) returning GeoPoint | null. No Strategy class, no Adapter, no custom geometry math.
- **Rationale**: ARCHITECTURE_RULES.md's guard-clause convention and Simple Design Rule #4 (fewest elements) apply directly: AC9 and AC10 are two sequential, independent rejection checks, not branches of one conditional on a type code, so Decompose Conditional / Replace Conditional with Polymorphism don't apply here. Using Cesium's own occluder avoids reinventing visibility math.
- **Fitness Function**: judgment-only — ARCHITECTURE_RULES.md's guard-clause convention and Simple Design Rule #4 (fewest elements) apply directly: AC9 and AC10 are two sequential, independent rejection checks, not branches of one conditional on a type code, so Decompose Conditional / Replace Conditional with Polymorphism don't apply here. Using Cesium's own occluder avoids reinventing visibility math.
### Draft state (toolStarted, draftPoints, draftEntities) stays as component-local Vue refs/plain arrays in ReagentMapView.vue. It is not promoted to a Pinia store.
- **Rationale**: Identical lifecycle to the existing sites/expeditions/viewer state: ephemeral, scoped to one view instance, with exactly one consumer. The spec explicitly defers persistence to reagent-map-save; promoting to a store now would be premature generalization (design-principles.md Anti-Pattern Radar) for a requirement this order doesn't have.
- **Trade-offs**:
  - reagent-map-save will need to re-decide this once a Save action needs to read the draft from a different place in the tree — this decision doesn't block that, it just doesn't do it preemptively
- **Fitness Function**: judgment-only — Identical lifecycle to the existing sites/expeditions/viewer state: ephemeral, scoped to one view instance, with exactly one consumer. The spec explicitly defers persistence to reagent-map-save; promoting to a store now would be premature generalization (design-principles.md Anti-Pattern Radar) for a requirement this order doesn't have.
### Import GeoPoint from ../api/expeditions.ts in ReagentMapView.vue rather than declaring a second { lat, lng } type for draft points.
- **Rationale**: design-principles.md's No Duplication rule, and it keeps the draft-point shape identical to Expedition.from/to — the same shape reagent-map-save's SaveExpedition call will need next, so no translation layer is required between the draft and the eventual save payload; both live in the same bounded context (Expedition Planning).
- **Fitness Function**: No second GeoPoint-shaped type literal is declared in ReagentMapView.vue — verified by code-reviewer greps the diff for a second `{ lat: number; lng: number }`-shaped declaration
### Register the LEFT_CLICK handler exactly once, via viewer.screenSpaceEventHandler.setInputAction(...), inside onMounted's happy path alongside the rest of viewer setup. startLine() never registers or re-registers a handler. No new teardown code is added to onBeforeUnmount.
- **Rationale**: Cesium's ScreenSpaceEventHandler is itself Observer/PubSub (ARCHITECTURE_RULES.md section III); using it directly avoids a custom event system. The Viewer owns and destroys its own ScreenSpaceEventHandler when viewer.destroy() runs, so the existing onBeforeUnmount cleanup already covers it — no new resource to release.
- **Fitness Function**: Exactly one setInputAction(..., LEFT_CLICK) call exists in ReagentMapView.vue, registered in onMounted; no second registration appears in startLine() or elsewhere — verified by code-reviewer confirms by diff review

## Component Placement

- `toolStarted/draftPoints refs, draftEntities array, startLine(), clearDraftEntities(), handleMapClick(), pickLandPoint(), placePoint(), addDraftPointEntity(), addDraftLineEntity(), SaturdayMetadata type` → Frameworks & Drivers (Vue UI layer, ARCHITECTURE_RULES.md section I.4) — not Entities/UseCases, since draft points/lines are not validated against business rules or persisted in this order layer, `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Extends the file's existing per-entity helper pattern (addSiteEntities, addExpeditionEntities, flyToSite); same layer, same file, same calling convention.
- `New start-line button` → Frameworks & Drivers (Vue template) layer, `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Must sit outside <header>, #cesiumContainer, reagent-site-list and reagent-map-error so none of reagent-map-visual.spec.ts's four captured regions change, per analyst's developerTasks.
- `reagent-map-line-tool.spec.ts` → Test (E2E/UI, Saturday/Playwright conventions) layer, `apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts` — Follows the house style of reagent-map-fly-to-site.spec.ts and reagent-map-expeditions.spec.ts; must compute raw screen-pixel click targets itself since @orieken/saturday-playwright-cesium's clickAt() throws on toScreen() === null and can't express AC9/AC10.

## Bounded Context

- **Owning Context**: Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7
- **Context Crossings**: None

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] Domain/UseCases layers must not import Adapters or Frameworks (ARCHITECTURE_RULES.md section I) — This app has no separate domain/use-case layer for this slice — ReagentMapView.vue and api/expeditions.ts are both adapter/framework-level already. Confirmed by reading both files directly; mcp__loom__dependencies_of/dependents_of were not invoked for this two-file change, consistent with the fallback the architect prompt describes for a change this small.
- [PASS] New code imports only from the Interface Adapters layer it already depends on (api/expeditions.ts), introducing no new cross-package dependency — The only new import is the type-only GeoPoint from the already-imported ../api/expeditions.ts module; no new package dependency is added.

## Anti-Pattern Check

- Premature generalization: not found — Explicitly rejected extracting pickLandPoint/draft state into a new module or store ahead of a second consumer — see structural decisions 1 and 4.
- God object: not found — ReagentMapView.vue grows by roughly 9 small functions and 3 refs but keeps one responsibility (render the Aurelis map and its controls). Flagging as a watch-item for the next feature touching this file — see Refactoring Opportunities.
- Shotgun surgery: not found — Single file touched (plus its new spec); no other view, store, or API client changes.
- Leaky abstraction: not found — pickLandPoint's GeoPoint | null return hides Cesium's Cartesian3/Cartographic types from callers; no caller outside this file needs to know Cesium's picking internals.
- Anemic domain model: not found — Not applicable — no domain entity is introduced; draft-point/draft-line are UI-level saturdayMetadata tags, same status as the existing 'site'/'expedition' tags.

## Fitness Functions

- **No new src/ file or module is introduced for line-tool logic outside ReagentMapView.vue and its spec file**
  - **Verification**: code-reviewer checks the diff's file list for this feature; any new file requires a stated reason before approval
  - **Owner**: code-reviewer
- **Every saturdayMetadata assignment in ReagentMapView.vue is typed via the shared SaturdayMetadata type, not an inline object literal**
  - **Verification**: code-reviewer confirms by diff review — no ESLint config exists anywhere in this repo yet, so this is manual judgment, not an automated gate
  - **Owner**: code-reviewer
- **No second GeoPoint-shaped type literal is declared in ReagentMapView.vue**
  - **Verification**: code-reviewer greps the diff for a second `{ lat: number; lng: number }`-shaped declaration
  - **Owner**: code-reviewer
- **Exactly one setInputAction(..., LEFT_CLICK) call exists in ReagentMapView.vue, registered in onMounted; no second registration appears in startLine() or elsewhere**
  - **Verification**: code-reviewer confirms by diff review
  - **Owner**: code-reviewer

## Refactoring Opportunities

- Extract Function/Module (deferred) on `pickLandPoint() and the draft-entity helpers in ReagentMapView.vue` — If a second interactive drawing tool (beyond the line tool and the eventual save flow) is added later and needs the same click-to-geo-point logic, extract it into its own module then — not before. Flag for whichever order introduces that second consumer.

## Developer Handoff Notes

- Do not create any new file or module for this feature besides the Playwright spec — all state, the click handler, pickLandPoint, and the draft-entity helpers belong in ReagentMapView.vue beside addSiteEntities/addExpeditionEntities/flyToSite.
- Add the SaturdayMetadata type ({ id, name, kind: 'site'|'expedition'|'draft-point'|'draft-line' }) and use it at all four saturdayMetadata assignment sites, replacing the existing untyped `unknown` cast literal — this touches addSiteEntities and addExpeditionEntities too, not just the two new functions.
- Register the LEFT_CLICK handler exactly once in onMounted's happy path, after the existing viewer/sites/expeditions setup succeeds. Do not add any new cleanup in onBeforeUnmount — viewer.destroy() already tears down its own ScreenSpaceEventHandler.
- Import GeoPoint from ../api/expeditions.ts; do not redeclare it.
- Keep pickLandPoint, placePoint, addDraftPointEntity, and addDraftLineEntity as small guard-clause functions (ARCHITECTURE_RULES.md's 30-LOC ceiling, design-principles.md's 5-10 line Sandi Metz guidance) — do not combine them into one larger handleMapClick body.
- Confirm the new button's placement never falls inside <header>, #cesiumContainer, reagent-site-list, or reagent-map-error — re-check against reagent-map-visual.spec.ts's four locators before assuming zero baseline diff, per analyst's architecturalFlags.

## Open Architectural Questions

- No ESLint config exists anywhere in this repo today (confirmed by search), so the complexity and no-raw-any fitness functions analyst proposed, and the SaturdayMetadata typing fitness above, rely on code-reviewer's manual judgment rather than a CI gate. This is a pre-existing gap, not something this order should fix — flagging for devops-engineer/human awareness only.
