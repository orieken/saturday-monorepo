# Context Manifest: leaflet-local-tiles

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `.claude/factory/claims/leaflet-local-tiles.md` — The feature spec. Settles the approach (keep default CRS, replace the OSM tile layer with an L.imageOverlay) and lists every acceptance criterion the analyst must turn into tasks.
- `apps/ye-olde-magic-shop/src/views/LocationsView.vue` — The single file to change. Contains the OSM L.tileLayer call to remove, the markers/territory GeoJSON/popups/layer control that must keep their relative positions, and the window.__SATURDAY_LEAFLET_MAP__ exposure the fix must not disturb.
- `apps/ye-olde-magic-shop/tests/map.spec.ts` — The suite that must pass unchanged. Fixes the real lat/lng values (center 45.0/-120.0, zoom-10 bounds, geolocation pan, 300-400km Turf distance) the imageOverlay's geographic bounds must still satisfy under the default CRS.
- `packages/saturday-playwright-leaflet/src/fixture.ts` — Defines waitForMapLoad (resolves immediately unless a layer exposes isLoading()) and checkTileFailures (queries .leaflet-tile elements for naturalWidth===0). An L.imageOverlay has neither an isLoading method nor a .leaflet-tile class, so the analyst needs this to confirm both fixture calls stay satisfied with no special-casing.
- `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` — Holds the app's only existing 'no request leaves the app's own origin' Playwright pattern (collectOffOriginRequests, APP_ORIGIN_HOSTNAME/ALLOWED_PORTS). AC4's new blocking test should reuse or extend this mechanism rather than invent a second one.
- `docs/features/local-fonts/analysis.md` — Most recent prior delivery in this bounded context (ye-olde-magic-shop frontend) that removed an external-origin dependency. Its Out of Scope section explicitly names this feature ('leaflet-local-tiles') as the deferred map-tile half of the same problem, and its AC/task shape (exact-pin new assets, record failing-test evidence for negative cases, no-request test with no host exceptions) is the template this feature should follow. No retrospective.md exists for it, so no What-Went-Poorly/What-To-Improve content could be pulled.
- `apps/ye-olde-magic-shop/src/main.ts` — App entry point; currently imports leaflet/dist/leaflet.css globally. Relevant if the fix needs any app-wide leaflet setup beyond LocationsView.vue itself.
- `packages/saturday-playwright-leaflet/README.md` — Documents the public contract of the `leaflet` Playwright fixture (waitForMapLoad, checkTileFailures, getBounds, distanceBetween, etc.) that both map.spec.ts and the new AC4 test depend on — needed to know what the harness guarantees versus what the view must provide.

## Token Budget

- **Status**: OK — ~12159 tokens, 10.1% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `.claude/factory/claims/leaflet-local-tiles.md` | 2902 | 725 |
| `apps/ye-olde-magic-shop/src/views/LocationsView.vue` | 4780 | 1195 |
| `apps/ye-olde-magic-shop/tests/map.spec.ts` | 5064 | 1266 |
| `packages/saturday-playwright-leaflet/src/fixture.ts` | 11148 | 2787 |
| `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` | 8276 | 2069 |
| `docs/features/local-fonts/analysis.md` | 11146 | 2786 |
| `apps/ye-olde-magic-shop/src/main.ts` | 727 | 181 |
| `packages/saturday-playwright-leaflet/README.md` | 4600 | 1150 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue -- different bounded context (Cesium-based reagent map, not Leaflet locations); uses a different mapping library and route, no structural overlap with this fix`
- `docs/testing-cesium.md -- documents the Cesium app's no-network testing conventions for a different feature/route; not Leaflet, no docs/testing-leaflet.md equivalent exists yet for this feature to extend`
- `apps/ye-olde-magic-shop/public/aurelis.jpg -- binary imagery asset for the unrelated Cesium reagent-map feature; same 'local image replaces network imagery' idea but different bounded context and no recorded licence precedent to borrow`
- `Repo-graph blast-radius tools (mcp__loom__impact/dependents_of/find_files) were not available in this session (not registered/deferred) -- blastRadius is omitted rather than hand-computed; pinning above was done by Grep over imports/usages instead`
- `No Knowledge Items or ADRs exist anywhere in this repository (shared/knowledge/, .claude/knowledge/, docs/adrs/ are all empty or absent) -- search-ki found no corpus to search against for 'leaflet tiles offline map image overlay no-network testing'; this is an empty-corpus result, not a missed match`
- `Infrastructure/migration/CI files -- excluded by the change-surface rule: this is a UI-only fix (no data model or API changes per the spec's Trust Boundaries and NFR sections)`
