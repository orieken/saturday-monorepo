# Context Manifest: reagent-map-line-tool

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — The file this order extends. Already shows the per-entity saturdayMetadata tagging pattern (kind:'site', kind:'expedition'), the window.__SATURDAY_CESIUM__ seam, geodesic polylines via Cesium.ArcType.GEODESIC, and the existing fly-to-site <button> control convention (data-testid, aria-label, no visual-style change) that a new 'start the line tool' control should match. A new click handler on the cesiumContainer div (or a Cesium ScreenSpaceEventHandler) is needed to turn a canvas click into PlacePoint/PointRejected; none exists yet.
- `features/reagent-map/draw-line.feature` — The exact Gherkin for AC6-AC10 (this order's five @level-ui scenarios): two points complete a line that can be saved; one point cannot; the line follows the shortest surface path; a click beside the globe or on the far-side edge places nothing. AC14's scenario text lives in save-expedition.feature instead, but the claim's own fourth AC bullet already restates it verbatim.
- `packages/saturday-playwright-cesium/src/fixture.ts` — Defines clickAt(lat,lng) (drives a real page.mouse click, but throws if toScreen() returns null -- i.e. it cannot express AC9's 'click on space beside the globe' or AC10's 'click on the far-side edge', since those are exactly the null cases) and getEntities() (reads saturdayMetadata + geometry, needed to assert draft-point/draft-line entities exist or don't). The analyst needs to flag that the new draw-line test will need a raw-coordinate click gesture for AC9/AC10, not clickAt, since no lat/lng maps to 'off the globe'.
- `packages/saturday-playwright-cesium/src/geometry.ts` — distance() and pointInCircle(), radius-parameterized (never Earth's) -- what the new test asserts the 1 km tolerance with, both for 'line follows the shortest surface path between two known sites' (AC8) and for the spec's settled 'two points closer than 1 km are the same place' rule this order's line-drawing must respect when placing a second point.
- `apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts` — House style for a getEntities()-based assertion test in this app: cesium fixture import, waiting for settle, and checking saturdayMetadata.kind plus geometry shape. The new draw-line spec will assert kind:'draft-point'/'draft-line' the same way this file asserts kind:'expedition'.
- `apps/ye-olde-magic-shop/tests/reagent-map-fly-to-site.spec.ts` — House style for a clickAt()/toScreen()-driven gesture test against the production view -- the closest existing sibling to a click-to-place-a-point test, including how it waits for settle after a camera-affecting action before asserting.
- `apps/mock-api/data/reagent-sites.json` — Seed coordinates for every site draw-line.feature's scenarios name by name (Moonwell Spring, Ember Crag, Thornveil Glade, Glimmerdeep Caves) -- the new test needs these lat/lng pairs to call cesium.clickAt() at a known on-land point.
- `docs/specs/reagent-map-spec/corrections.md` — Eleven lines, directly cited by the claim: the post-ship correction that 'the land' means any point on the near side of the globe, sea included -- only space beside the globe and the far side refuse a point. This is the one settled gate a developer could get wrong by assuming an actual land/water check is needed.
- `docs/specs/reagent-map-spec/domain-model.md` — Names the commands/events this order implements (StartLine -> LineStarted, PlacePoint -> PointPlaced/PointRejected) and hotspot H2: choosing the line tool again, or leaving the view, discards the unsaved line -- directly backs AC14 and AC15 (AC15 is out of scope here, but the discard-on-restart behavior is not).
- `apps/ye-olde-magic-shop/src/api/expeditions.ts` — The existing Expedition/GeoPoint shape ({id, name?, owner, from, to}) the drawn-but-unsaved line's in-memory data should mirror (minus id/owner/savedAt, since nothing is saved yet), so the next order (reagent-map-save) can hand this order's draft state to SaveExpedition without reshaping it.

## Token Budget

- **Status**: OK — ~9771 tokens, 8.1% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` | 7804 | 1951 |
| `features/reagent-map/draw-line.feature` | 1758 | 439 |
| `packages/saturday-playwright-cesium/src/fixture.ts` | 10891 | 2722 |
| `packages/saturday-playwright-cesium/src/geometry.ts` | 2061 | 515 |
| `apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts` | 5269 | 1317 |
| `apps/ye-olde-magic-shop/tests/reagent-map-fly-to-site.spec.ts` | 3930 | 982 |
| `apps/mock-api/data/reagent-sites.json` | 781 | 195 |
| `docs/specs/reagent-map-spec/corrections.md` | 1297 | 324 |
| `docs/specs/reagent-map-spec/domain-model.md` | 3700 | 925 |
| `apps/ye-olde-magic-shop/src/api/expeditions.ts` | 1605 | 401 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `features/reagent-map/save-expedition.feature -- carries AC14's actual Gherkin scenario ('Starting a new line abandons the unfinished one'), but the claim's fourth acceptance-criteria bullet already restates it verbatim; AC11-AC13 and AC15 in the same file are the save order, out of scope here.`
- `docs/specs/reagent-map-spec/acceptance-criteria.md, acceptance-criteria.json, domain-model.json, discovery.md, discovery.json, example-map.md, example-map.json, README.md -- the full spec package; the claim already pins the digest and restates every AC and settled gate this order needs (1 km tolerance, geodesic, same-place rule, land/sea correction) in plain language, and draw-line.feature carries the exact Gherkin for AC6-AC10.`
- `docs/features/reagent-map-globe/analysis.md, implementation-notes.md; docs/features/reagent-map-expeditions/analysis.md, implementation-notes.md -- the two prior deliveries in this bounded context. Their established pattern (per-entity saturdayMetadata tagging, window.__SATURDAY_CESIUM__ lifecycle, geodesic arcType, button-as-list-item control convention) is already carried forward via the pinned ReagentMapView.vue itself. No retrospective.md exists anywhere under docs/features/ (checked reagent-map-globe, reagent-map-expeditions, reagent-sites-api, expeditions-api, reagent-map-visual-baselines) to check for What-Went-Poorly/What-To-Improve notes.`
- `apps/ye-olde-magic-shop/src/api/reagentSites.ts, src/api/orders.ts -- this order adds no new API call (no saving yet, per the claim's Out of Scope); the existing fetch clients are unaffected and nothing in them bears on drawing an in-memory line.`
- `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts, reagent-map-pan-zoom.spec.ts, reagent-map-resilience.spec.ts, reagent-map-a11y-check.spec.ts, reagent-map-feature-flag.spec.ts -- no draw-line coverage today; resilience.spec.ts specifically covers only the two existing fetch-rejection paths, and this order introduces no new fetch call for it to extend.`
- `apps/ye-olde-magic-shop/.playwright-baselines/**/*.png -- binary baseline images, not readable as text context; regenerated per the claim's visual-baseline acceptance criterion only if a control to start the tool changes a baseline's visible chrome.`
- `apps/ye-olde-magic-shop/src/router/index.ts, src/App.vue, src/featureFlags.ts -- the VITE_FEATURE_REAGENT_MAP route/nav-link plumbing shipped by reagent-map-globe; this order is additive within that existing flag and route, touching neither.`
- `DOMAIN_DICTIONARY.md section 7 -- already restated in the claim's own Domain Language section; no separate read needed.`
- `apps/ye-olde-magic-shop/src/views/LocationsView.vue, src/components/ItemCard.vue, src/stores/cartStore.ts -- different bounded context (shop/cart, Leaflet locations view); excluded by bounded-context exclusion.`
- `mcp__loom__impact / mcp__loom__dependents_of / mcp__loom__find_files -- not available as registered tools in this session, so the blast radius below was not graph-derived; pinned files were chosen by reading the claim, the spec package summaries, the two prior deliveries' artifacts, and Grep/manual exploration of the saturday-playwright-cesium fixture instead.`
