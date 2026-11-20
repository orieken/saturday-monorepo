---
feature: "leaflet-local-tiles"
bounded_context: "ye-olde-magic-shop frontend — shop-locations map (Leaflet)"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/src/views/LocationsView.vue", "apps/ye-olde-magic-shop/public/", "apps/ye-olde-magic-shop/tests/map.spec.ts"]
issue_refs: ["leaflet-local-tiles"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: leaflet-local-tiles

## Summary

LocationsView.vue currently loads OpenStreetMap tiles from tile.openstreetmap.org, so the shop-locations map (and its Playwright suite) requires network access. This swap replaces the L.tileLayer call with L.imageOverlay over an in-repo image, keeping Leaflet's default EPSG:3857 CRS so map.spec.ts's real lat/lng, bounds, geolocation and Turf-distance assertions stay valid. Markers, the Sword Coast territory polygon, popups and the layer control are untouched; a new Playwright test proves the view loads with every non-origin request blocked, and the new image's licence is recorded next to it.

### Acceptance Criteria

- The locations view's map imagery MUST come only from the app's own origin; no tile or image request MUST leave the page.
  - Example: Only same-origin asset requests occur when /locations loads, e.g. a request to /sword-coast-map.png served by the app's own build output
- The existing markers, territory polygon and popups MUST keep the same position relative to each other after the tile layer is replaced.
  - Example: Waterdeep Branch marker stays at [45.0, -120.0]
  - Example: The Sword Coast polygon still spans [-122.0, 48.0] to [-118.0, 41.0]
- Every test in tests/map.spec.ts MUST pass unchanged under the default Leaflet CRS.
  - Example: Map centre stays close to 45.0, -120.0
  - Example: Zoom-10 bounds include 45.0, -120.0 but exclude 42.0 latitude
  - Example: Geolocation pan moves the centre to 47.6, -121.1
  - Example: Turf distance between Waterdeep and Baldur's Gate stays between 300 and 400 km
- A Playwright test that aborts every request not addressed to the app's own origin MUST show the locations view loading, leaflet.waitForMapLoad() resolving, and leaflet.checkTileFailures() returning false.
- The new map image's licence MUST be recorded next to it, either as an original in-repo creation or under a licence that permits redistribution.
- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except the run's own delivery record under docs/features/leaflet-local-tiles/.

### Non-Functional Requirements

- **security**: The locations view MUST make zero requests to any origin other than the app's own while loading and rendering the map.

## Proposed Fitness Functions

- **No tile or image request leaves the app's own origin when the locations view loads**
  - **Verification**: Playwright: new test in tests/map.spec.ts that page.route()-aborts any request whose URL is not same-origin, then asserts leaflet.waitForMapLoad() resolves and leaflet.checkTileFailures() is false
  - **Owner**: qa-engineer
- **Real-geography assertions (centre, zoom-10 bounds, geolocation pan, 300-400km Turf distance) keep passing under the default CRS**
  - **Verification**: tests/map.spec.ts unchanged, run via scripts/ui-suite-in-chrome.sh
  - **Owner**: qa-engineer
- **Every new or changed function stays below cyclomatic complexity 7**
  - **Verification**: analyze-complexity skill / lint run on changed files in apps/ye-olde-magic-shop/src/views/LocationsView.vue
  - **Owner**: devops-engineer
- **The in-repo map image's licence is recorded next to the asset**
  - **Verification**: code-reviewer / tech-writer confirms a licence note sits beside the image in apps/ye-olde-magic-shop/public/; no automated check can judge licence sufficiency
  - **Owner**: tech-writer

## Out of Scope

- Changing the Leaflet package version or swapping it for another mapping library.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: ye-olde-magic-shop frontend — shop-locations map (Leaflet)
- **Context Crossings**: packages/saturday-playwright-leaflet's test fixture contract (waitForMapLoad/checkTileFailures) — relied on, not modified, tests/reagent-map-globe.spec.ts's off-origin-request-blocking pattern — reused for the new no-network test

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/src/views/LocationsView.vue` — Remove the L.tileLayer OSM call in onMounted and replace it with L.imageOverlay over an in-repo image and a lat/lng bounds box; keep the default CRS, markers, territory layer, popups, layer control and window.__SATURDAY_LEAFLET_MAP__ exposure exactly as they are
- `apps/ye-olde-magic-shop/public/` — Add the new in-repo map image asset (e.g. sword-coast-map.svg or .png) and a licence note file placed next to it
- `apps/ye-olde-magic-shop/tests/map.spec.ts` — Add the new 'blocks every off-origin request' test alongside the existing tests, which must pass unchanged

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Create an in-repo map image asset under apps/ye-olde-magic-shop/public/ (e.g. sword-coast-map.svg or .png) depicting the Sword Coast region; author it originally in-repo or use an image whose licence explicitly permits redistribution
2. Add a licence/attribution note file directly beside the new image asset (e.g. sword-coast-map.LICENSE.txt), recording authorship or the permitting licence
3. In LocationsView.vue's onMounted, remove the L.tileLayer('https://{s}.tile.openstreetmap.org/...') call and add L.imageOverlay(imageUrl, bounds, { className: 'map-tiles-parchment' }).addTo(map.value)
4. Size the imageOverlay's lat/lng bounds to cover all three shop coordinates (45.0,-120.0 / 42.0,-119.5 / 47.5,-121.0), the Sword Coast polygon (-122.0,48.0 to -118.0,41.0), and the zoom-6 geolocation flyTo target (47.6,-121.1), with margin so no marker or the flyTo view sits outside the image
5. Keep Leaflet's default CRS (EPSG:3857) — do not pass crs: L.CRS.Simple or any other CRS to L.map()
6. Confirm the existing .map-tiles-parchment CSS filter (sepia/hue-rotate/saturate/contrast/brightness) still reads acceptably against the new image; adjust the image's own colours rather than removing the filter if it does not

### QA Tasks

1. While developing, run only tests/map.spec.ts and the new test; confirm every existing test in tests/map.spec.ts passes unchanged
2. Add a new Playwright test that uses page.route() to abort any request whose URL origin is not the app's own, then loads /locations, calls leaflet.waitForMapLoad(), and asserts leaflet.checkTileFailures() is false
3. Run the full suite via scripts/ui-suite-in-chrome.sh one image at a time, waiting for each run to finish before starting the next or ending the stage
4. Confirm cyclomatic complexity stays below 7 for every changed function, and every changed line is covered by a test that can fail
5. Confirm no file outside scope.paths changed, other than the delivery record under docs/features/leaflet-local-tiles/

### Tech Writer Tasks

1. Confirm the new map image's licence note sits next to the asset and is accurate (original work vs. a named redistribution-permitting licence)
2. Record the tile-layer-to-image-overlay change and the licence decision in the feature's delivery record under docs/features/leaflet-local-tiles/

### DevOps Tasks

None

## Edge Cases and Risks

- Image overlay bounds are too narrow and crop the geolocation flyTo view (47.6,-121.1) or a shop marker out of frame: Size bounds generously around all three shop coordinates, the territory polygon, and the zoom-6 flyTo target; verify visually and via the existing bounds-reading assertions in map.spec.ts
- L.ImageOverlay causes waitForMapLoad() to hang because some layer unexpectedly exposes isLoading(): packages/saturday-playwright-leaflet/src/fixture.ts only waits on layers that define isLoading(); L.ImageOverlay does not define it, so waitForMapLoad() resolves immediately — confirm this holds for the Leaflet version pinned in this repo
- checkTileFailures() looks for `.leaflet-tile` elements with naturalWidth===0, but an image overlay renders under the `.leaflet-image-layer` class, not `.leaflet-tile`: With the tile layer removed, no `.leaflet-tile` elements exist at all, so the existing check's `.some()` over an empty array returns false automatically — no fixture change is needed
- The chosen image licence turns out not to permit redistribution: Prefer an original in-repo creation (hand-authored SVG or PNG) over any third-party map image, per the spec's settled preference, to avoid licence risk entirely
- A developer runs two ui-suite-in-chrome.sh images concurrently, exceeding the 30-minute stage budget: Run one image at a time to completion, per the spec's Test Approach, before starting the next

## Definition of Done

- All 6 acceptance criteria verified; the locations view loads with zero off-origin requests
- L.imageOverlay over an in-repo image replaces the OSM L.tileLayer, with Leaflet's default EPSG:3857 CRS kept
- Markers, the territory polygon, popups and the layer control render in the same relative positions as before
- tests/map.spec.ts passes unchanged; the new off-origin-blocking test passes
- The map image's licence is recorded next to the asset
- Every new or changed function has cyclomatic complexity below 7; every changed line is covered by a test that can fail
- No file outside scope.paths changed, except docs/features/leaflet-local-tiles/
- scripts/ui-suite-in-chrome.sh passes, run one image at a time to completion
