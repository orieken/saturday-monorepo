# Context Manifest: reagent-map-globe

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `docs/specs/reagent-map-spec/acceptance-criteria.md` — Primary spec source for the two criteria this order covers (AC1: all eight sites shown within 1 km, listed by name; AC2: pole/180th-meridian sites shown once, correctly). Also carries the seed data table (eight sites with lat/lng) and the settled gate answers (1 km tolerance, Aurelis is a 6,000 km sphere) the analyst needs to restate.
- `features/reagent-map/reagent-sites.feature` — The Gherkin scenarios for AC1/AC2 in business language -- the AC-to-scenario traceability the analyst must preserve. The other four .feature files in this directory belong to the next three orders and are excluded per the spec's Out of Scope section.
- `packages/saturday-playwright-cesium/README.md` — Defines the seam contract this order must implement verbatim: window.__SATURDAY_CESIUM__ = { viewer, Cesium }, and saturdayMetadata (id, name, kind) on every site entity. This is the context crossing into the Test Automation bounded context the spec explicitly names -- only the public-interface doc is pinned, not the package's internal source.
- `packages/saturday-playwright-cesium/fixtures/assets/fixture-globe.js` — Concrete, already-reviewed reference implementation of the exact Cesium Viewer configuration AC2 demands: an Ellipsoid sized to a 6,000 km radius, baseLayerPicker/geocoder/ion off, a local SingleTileImageryProvider image, and the saturdayMetadata-tagged entity pattern. Saves the analyst from re-deriving the Cesium option list from scratch.
- `apps/mock-api/server.js` — The existing GET /api/reagent-sites handler (around lines 127-129) and its data loading (around lines 11-26, 342-352) -- the API contract the view consumes. No backend change is in scope for this order, so the analyst only needs to document it.
- `apps/mock-api/data/reagent-sites.json` — The actual seed data, including the two edge-case sites named in AC2 by id: frostcap-hollow (lat 85, near the north pole) and tidewrack-shoals (lng 180, on the meridian).
- `apps/ye-olde-magic-shop/src/views/LocationsView.vue` — The only existing map view in this app (Leaflet-based) -- the established pattern for lifecycle, exposing a window seam for Playwright, and attaching metadata to markers. The new reagent map view is a sibling of this file, adapted to Cesium's seam contract.
- `apps/ye-olde-magic-shop/src/router/index.ts` — The file the new /reagent-map route (and its feature-flag gate) must be added to; shows the existing route list so the analyst can describe the change as additive.
- `apps/ye-olde-magic-shop/playwright.config.ts` — Already has the 'chrome' project (real Chrome, software WebGL) that the spec requires reagent-map-*.spec.ts to run under, and shows that the default chromium project must stay unaffected -- needed for the test-approach/DoD criteria.
- `apps/ye-olde-magic-shop/.env` — Confirms no VITE_FEATURE_REAGENT_MAP entry exists yet; AC1 requires this file to set it to false.

## Token Budget

- **Status**: OK — ~10160 tokens, 8.5% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `docs/specs/reagent-map-spec/acceptance-criteria.md` | 7969 | 1992 |
| `features/reagent-map/reagent-sites.feature` | 1335 | 333 |
| `packages/saturday-playwright-cesium/README.md` | 6486 | 1621 |
| `packages/saturday-playwright-cesium/fixtures/assets/fixture-globe.js` | 4700 | 1175 |
| `apps/mock-api/server.js` | 11361 | 2840 |
| `apps/mock-api/data/reagent-sites.json` | 781 | 195 |
| `apps/ye-olde-magic-shop/src/views/LocationsView.vue` | 4780 | 1195 |
| `apps/ye-olde-magic-shop/src/router/index.ts` | 731 | 182 |
| `apps/ye-olde-magic-shop/playwright.config.ts` | 2378 | 594 |
| `apps/ye-olde-magic-shop/.env` | 133 | 33 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `DOMAIN_DICTIONARY.md`
- `features/reagent-map/draw-line.feature`
- `features/reagent-map/fly-to-site.feature`
- `features/reagent-map/expeditions-visibility.feature`
- `features/reagent-map/save-expedition.feature`
- `apps/ye-olde-magic-shop/src/api/items.ts`
- `apps/ye-olde-magic-shop/vite.config.ts`
- `apps/ye-olde-magic-shop/package.json`
- `apps/ye-olde-magic-shop/tests/map.spec.ts`
- `apps/ye-olde-magic-shop/tests/webgl.spec.ts`
- `packages/saturday-playwright-cesium/fixtures/index.html`
- `packages/saturday-playwright-cesium/fixtures/serve.mjs`
- `packages/saturday-playwright-cesium/package.json`
- `docs/features/reagent-sites-api/analysis.md`
- `docs/features/expeditions-api/analysis.md`
- `docs/features/saturday-playwright-cesium/analysis.md`
