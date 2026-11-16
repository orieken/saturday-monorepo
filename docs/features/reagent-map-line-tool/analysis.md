---
feature: "reagent-map-line-tool"
bounded_context: "Expedition Planning (Ye Olde Magic Shop demo app) -- DOMAIN_DICTIONARY.md section 7"
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts", "apps/ye-olde-magic-shop/.playwright-baselines/debian/**", "apps/ye-olde-magic-shop/.playwright-baselines/ubi9/**", "docs/features/reagent-map-line-tool/"]
issue_refs: ["reagent-map-line-tool"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: reagent-map-line-tool

## Summary

Adds a line-drawing tool to the existing Aurelis reagent map (ReagentMapView.vue), still gated behind the VITE_FEATURE_REAGENT_MAP build flag reagent-map-globe introduced (off by default until reagent-map-save turns it on; no new flag needed). A new 'start the line tool' control resets any in-progress draft and arms a click handler on the globe; placing a point on the near side of the globe (land or sea, per the spec's land/sea correction) adds a saturdayMetadata kind:'draft-point' entity, and a second point adds a geodesic kind:'draft-line' entity joining them -- the only observable sense in which a line 'can be saved' this order, since the actual save ships in reagent-map-save. A click that misses the globe entirely (Cesium's camera.pickEllipsoid returns nothing) or lands on the limb where the surface faces away (EllipsoidalOccluder.isPointVisible reports false) places nothing and leaves the tool waiting for its first point. Starting the tool again discards any unsaved points, so a new line starts empty. Purely additive within the view's established per-entity saturdayMetadata / window.__SATURDAY_CESIUM__ pattern -- no new package, endpoint, or persisted state.

### Acceptance Criteria

- The line tool MUST join two points placed on the land with a completed geodesic line that is ready to save.
  - Example: Placing a point at Moonwell Spring then a point at Ember Crag joins them with a completed line (spec AC6).
  - Example: The completed line is exposed only as a saturdayMetadata kind:'draft-line' entity readable through getEntities(), since no save action exists in this order.
- After only one point is placed, the line tool MUST NOT show a completed line and that line MUST NOT be saveable.
  - Example: Placing a single point at Moonwell Spring leaves only a kind:'draft-point' entity; no kind:'draft-line' entity exists (spec AC7).
- A completed line MUST follow the shortest path over Aurelis's surface between its two points, not a straight line on a flat projection.
  - Example: A line from Thornveil Glade to Glimmerdeep Caves is drawn with arcType GEODESIC, the same convention addExpeditionEntities already uses (spec AC8).
- The line tool MUST place a point anywhere on the near side of the globe, sea included, since 'the land' means the planet's surface for now.
  - Example: docs/specs/reagent-map-spec/corrections.md: only space beside the globe and the far side refuse a point; no separate land/water check is implemented.
- A click on space beside the globe MUST place no point, show nothing, and leave the line tool waiting for its first point.
  - Example: camera.pickEllipsoid(windowPosition, ellipsoid) returns undefined for a click outside the rendered globe (spec AC9).
- A click on the globe's edge, where the surface faces away from the camera, MUST place no point.
  - Example: pickEllipsoid intersects the sphere but EllipsoidalOccluder(ellipsoid, camera.positionWC).isPointVisible(cartesian) returns false at the grazing limb (spec AC10).
- Starting the line tool again MUST discard any points from the line already in progress, so the new line has no points.
  - Example: Given one point placed, starting the line tool again removes its draft-point entity and clears draft state (spec AC14).
- Every point and line the tool draws MUST carry saturdayMetadata with kind 'draft-point' or 'draft-line' so tests can read them through getEntities().
- Only a @visual baseline that changes for a stated reason MUST be regenerated; no other baseline MUST change.
  - Example: This order's design places the new start-line control outside the four regions reagent-map-visual.spec.ts captures (header, site list, globe, error), so no baseline is expected to change; the full suite is still run in both runner images to confirm.
- Every new or changed function MUST keep cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except this order's own delivery record.

### Non-Functional Requirements

- **security**: No request the line tool triggers MUST leave the app's own origin or the mock API's own origin. This order adds no new network call, so the existing off-origin NFR test continues to cover it.
- **accessibility**: The control that starts the line tool MUST be operable by keyboard (Enter/Space activates it, it is a real <button>) and MUST expose an accessible name announcing that it starts drawing a new expedition line, matching the fly-to-site button convention.

## Proposed Fitness Functions

- **Every new or changed function in this order keeps cyclomatic complexity below 7**
  - **Verification**: ESLint complexity rule capped at 6 per typescript-conventions.md, run in CI lint
  - **Owner**: devops-engineer
- **Every changed line is covered by a Playwright spec that can fail**
  - **Verification**: No coverage tool is wired for this Playwright-only app. qa-engineer proves each new assertion fails against the pre-change view, per test-repair-contract.md, and records this in qa-report.md
  - **Owner**: qa-engineer
- **AC9/AC10's rejection paths (pickEllipsoid returning nothing, or EllipsoidalOccluder.isPointVisible returning false) are each exercised by a test that fails if the corresponding check is removed**
  - **Verification**: reagent-map-line-tool.spec.ts's off-globe and limb-click tests; each is shown to fail against a build where the corresponding guard is deleted, per test-repair-contract.md
  - **Owner**: qa-engineer
- **No @visual baseline changes except one this order's acceptance criteria name a reason for**
  - **Verification**: scripts/ui-suite-in-chrome.sh and --ubi9 run with no --update-baselines flag before ship; zero diff is required unless the delivery record states which baseline changed and why
  - **Owner**: qa-engineer

## Out of Scope

- Saving a completed line: reagent-map-save.
- The full scenario suites E1 and F1: run by hand from the spec package, not by this pipeline.
- Rejecting a completed line whose two points are closer than 1 km (the 'same place' rule): that rejection happens at save time (spec AC13), not at placement, so this order does not implement it.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Expedition Planning (Ye Olde Magic Shop demo app) -- DOMAIN_DICTIONARY.md section 7
- **Context Crossings**: None

### Domain Events (Event Storming Lite)

- **Commands**: StartLine, PlacePoint
- **Events Produced**: LineStarted, PointPlaced, PointRejected
- **Owning Aggregates**: Expedition (draft, in-memory only -- not persisted until reagent-map-save)
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Add a 'start the line tool' button, draft-point/draft-line state (toolStarted, draftPoints refs), a ScreenSpaceEventHandler LEFT_CLICK handler that resolves a click to a geographic point via pickEllipsoid + EllipsoidalOccluder, and helpers to add/clear the draft-point and draft-line entities.
- `apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts` — New Playwright spec covering draw-line.feature's five @level-ui scenarios (AC6, AC7, AC8, AC9, AC10) plus AC14, following the house style of reagent-map-expeditions.spec.ts and reagent-map-fly-to-site.spec.ts.
- `apps/ye-olde-magic-shop/.playwright-baselines/debian/**` — Re-run (not necessarily regenerate) to confirm the new control, placed outside all four captured regions, produces zero diff; regenerate only if a diff appears and record why.
- `apps/ye-olde-magic-shop/.playwright-baselines/ubi9/**` — Same check as the debian baselines, ubi9 runner image.
- `docs/features/reagent-map-line-tool/` — Delivery record stating the visual-suite outcome (changed or unchanged, and why) per this order's acceptance criteria.

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Add toolStarted (ref<boolean>) and draftPoints (ref<GeoPoint[]>) state to ReagentMapView.vue, plus a component-scope draftEntities: CesiumModule.Entity[] array (same lifetime as viewer).
2. Add a <button type="button" data-testid="start-line-tool" aria-label="Draw a new expedition"> in a new element that is NOT inside <header>, #cesiumContainer, or the reagent-site-list/reagent-map-error elements, so no existing @visual baseline is affected; wire it to startLine().
3. Add startLine(): calls clearDraftEntities(), resets draftPoints.value to [], sets toolStarted.value to true.
4. Add clearDraftEntities(): removes every entity in draftEntities from viewer.entities, then empties the array.
5. Register viewer.screenSpaceEventHandler.setInputAction(handleMapClick, Cesium.ScreenSpaceEventType.LEFT_CLICK) once the viewer exists in onMounted's happy path.
6. Add handleMapClick(click): returns early when !toolStarted.value, draftPoints.value.length >= 2, or viewer/cesium/ellipsoid is null; otherwise calls pickLandPoint(click.position) and, if it returns a point, calls placePoint(point).
7. Add pickLandPoint(position): cartesian = viewer.camera.pickEllipsoid(position, ellipsoid); return null if undefined (AC9). Otherwise build new Cesium.EllipsoidalOccluder(ellipsoid, viewer.camera.positionWC) and return null if !occluder.isPointVisible(cartesian) (AC10). Otherwise convert to Cartographic and return { lat, lng } in degrees.
8. Add placePoint(point): pushes point onto draftPoints.value, pushes addDraftPointEntity(point) onto draftEntities; if draftPoints.value.length === 2, also pushes addDraftLineEntity(draftPoints.value) onto draftEntities.
9. Add addDraftPointEntity(point): adds a point entity at Cesium.Cartesian3.fromDegrees(point.lng, point.lat, 0, ellipsoid) with point: { pixelSize: SITE_MARKER_PIXEL_SIZE }, tagged saturdayMetadata { id: '', name: '', kind: 'draft-point' }; returns the entity.
10. Add addDraftLineEntity(points): adds a polyline entity from Cesium.Cartesian3.fromDegreesArray([points[0].lng, points[0].lat, points[1].lng, points[1].lat], ellipsoid) with arcType: Cesium.ArcType.GEODESIC, tagged saturdayMetadata { id: '', name: '', kind: 'draft-line' }; returns the entity.
11. Import GeoPoint from ../api/expeditions.ts rather than redeclaring it, keeping the draft points in the same shape reagent-map-save will later hand to SaveExpedition.

### QA Tasks

1. Write apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts covering draw-line.feature: start the tool, place Moonwell Spring then Ember Crag via cesium.clickAt(), assert getEntities() has one kind:'draft-point' per placed point plus one kind:'draft-line' joining them (AC6); after only the first point, assert no kind:'draft-line' entity exists (AC7); assert the drawn line's geometry is a 2-point geodesic arc for Thornveil Glade to Glimmerdeep Caves, consistent with arcType GEODESIC (AC8).
2. In the same spec, for AC9: compute the canvas bounding box via page.locator('#cesiumContainer').boundingBox() and page.mouse.click at a corner pixel outside the rendered globe disc (space, at the default zoomed-out view); assert getEntities() still has zero draft-point/draft-line entities.
3. For AC10: compute a near-limb screen pixel inline in the spec file (using the camera height and globe radius already exposed via cesium.globeRadiusMeters(), or an inline page.evaluate reading window.__SATURDAY_CESIUM__ the same way reagent-map-fly-to-site.spec.ts's readCamera() does) and click it; assert no point is placed. Do not extend packages/saturday-playwright-cesium -- it is outside scope.paths.
4. Write the AC14 case: place one point, click start-line-tool again, assert getEntities() has zero draft-point/draft-line entities.
5. For every new assertion, confirm it fails against the pre-change view (or cite why the behavior is new) per test-repair-contract.md, and record this in qa-report.md.
6. Run the reagent-map-*.spec.ts set under --project=chrome locally while developing; before ship, run scripts/ui-suite-in-chrome.sh and --ubi9 (no flags) and confirm all four existing baselines are unchanged.
7. Only if a baseline actually differs (e.g. a different control placement than this analysis assumed), regenerate it with --update-baselines in both images, view the diff, and record exactly what changed and why in the delivery record. Do not raise the globe's maxDiffPixelRatio tolerance.

### Tech Writer Tasks

1. Write the docs/features/reagent-map-line-tool/ delivery record per CLAUDE.md's Delivery Summary Format, stating whether any @visual baseline changed and, if so, what and why; otherwise state plainly that all four baselines stayed byte-identical.

### DevOps Tasks

None

## Edge Cases and Risks

- A third click occurs after a line is already complete (two points placed).: handleMapClick's draftPoints.value.length >= 2 guard makes this a no-op; the only way to draw a new line is pressing start-line-tool again (AC14).
- Both points are placed less than 1 km apart (the spec's 'same place' rule).: Placement still completes the line and adds the draft-line entity; the zero-length rejection is enforced only at save time (spec AC13, reagent-map-save's scope), not here.
- A click lands on space beside the globe.: camera.pickEllipsoid returns undefined; pickLandPoint returns null; no entity is added and draftPoints is unchanged (AC9).
- A click lands on the globe's limb, where the surface faces away from the camera.: pickEllipsoid intersects, but EllipsoidalOccluder.isPointVisible returns false; pickLandPoint returns null; no entity is added (AC10).
- start-line-tool is pressed while a line is already complete, or mid-draw with one point placed.: clearDraftEntities() removes any existing draft-point/draft-line entities and draftPoints is reset to [] before the tool re-arms (AC14).
- The adventurer leaves the reagent map view with an unsaved line in progress (domain-model.md hotspot H2's other half).: onBeforeUnmount already calls viewer.destroy(), which tears down every entity including drafts; no additional handling is needed.
- Placing a point is a mouse-only gesture with no keyboard equivalent.: Out of scope for this order -- mirrors the globe's existing lack of keyboard interaction for anything except the site-list buttons; flagged for accessibility-engineer, not fixed here.

## Definition of Done

- Two points placed on the near side of the globe join into a completed geodesic line exposed as a kind:'draft-line' entity, verified by reagent-map-line-tool.spec.ts (AC6, AC8).
- One placed point shows no completed line -- no kind:'draft-line' entity exists -- verified by the same spec (AC7).
- A click on space beside the globe and a click on the globe's limb where the surface faces away each place nothing and leave the tool waiting for its first point, verified by the same spec (AC9, AC10).
- Starting the line tool again after one point is placed clears it, verified by the same spec (AC14).
- Every placed point and completed line carries saturdayMetadata kind:'draft-point'/'draft-line', readable through getEntities().
- scripts/ui-suite-in-chrome.sh and --ubi9 are run; any baseline that differs is regenerated deliberately, reviewed by eye, and the reason is recorded in the delivery record; baselines that should not change do not.
- No file outside scope.paths changes except this order's own delivery record under docs/features/reagent-map-line-tool/.
- Every new or changed function keeps cyclomatic complexity below 7.
- code-reviewer has approved the implementation.
- qa-engineer has proven every new assertion can fail, per test-repair-contract.md.
