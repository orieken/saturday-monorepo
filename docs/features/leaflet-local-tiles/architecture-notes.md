---
feature: "leaflet-local-tiles"
bounded_context: "ye-olde-magic-shop frontend — shop-locations map (Leaflet)"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/src/views/LocationsView.vue (call site), apps/ye-olde-magic-shop/public/ (image + licence note)"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: leaflet-local-tiles

## Structural Decisions

### Replace the L.tileLayer(OSM) call in LocationsView.vue's onMounted with a single L.imageOverlay() call over an in-repo image, keeping Leaflet's default EPSG:3857 CRS on the existing L.map('map') instance. No new component, class, or package is introduced.
- **Rationale**: This is a one-line substitution of a tile source inside an already-established pattern (a single Vue view owning its own Leaflet map instance). CLAUDE.md's architect-skip criterion — 'purely additive within an established pattern' — applies directly: no new base class, no layer-boundary change, no cross-cutting concern. The analyst's own architecturalFlags entry already reached this conclusion; this review confirms it against the actual file rather than taking it on faith.
- **Trade-offs**:
  - A static image overlay cannot pan/zoom-stream new detail the way real tiles can, but the feature's own purpose (deterministic, network-free test fixture) makes that a non-issue, not a cost.
  - The image's lat/lng bounds box must be hand-sized to cover all three shop coordinates, the Sword Coast polygon, and the zoom-6 flyTo target — a magic-number risk that lands on the developer, not on structure.
- **Fitness Function**: judgment-only — This is a one-line substitution of a tile source inside an already-established pattern (a single Vue view owning its own Leaflet map instance). CLAUDE.md's architect-skip criterion — 'purely additive within an established pattern' — applies directly: no new base class, no layer-boundary change, no cross-cutting concern. The analyst's own architecturalFlags entry already reached this conclusion; this review confirms it against the actual file rather than taking it on faith.
### Keep the new map image and its licence note physically beside each other under apps/ye-olde-magic-shop/public/ (e.g. sword-coast-map.png + sword-coast-map.LICENSE.txt), matching the existing flat-file convention already used there (aurelis.jpg, logo.png, favicon.svg).
- **Rationale**: apps/ye-olde-magic-shop/public/ is Vite's static-asset root; files here are served byte-for-byte from the app's own origin with no build-time processing, which is exactly the same-origin guarantee acceptance criterion #1 depends on. No new directory or asset pipeline is justified for one image.
- **Trade-offs**:
  - Flat public/ placement works at the current scale (four static assets); it would need revisiting only if this app accumulates many map/image assets, which is not this feature's concern.
- **Fitness Function**: The map image is served from the app's own origin with zero cross-origin requests during /locations load — verified by Playwright: new test in tests/map.spec.ts that page.route()-aborts any request whose URL is not same-origin, asserting the route handler never fires for a non-origin request and leaflet.checkTileFailures() is false
### Do not introduce L.CRS.Simple or any non-default CRS; L.map('map') keeps its current no-crs-argument call, so Leaflet's default EPSG:3857 (Web Mercator) continues to back every coordinate in the view.
- **Rationale**: map.spec.ts's assertions (map centre ≈ [45.0, -120.0], zoom-10 bounds, geolocation flyTo to [47.6, -121.1], Turf.js distance 300–400km between Waterdeep and Baldur's Gate) are all computed against real-world lat/lng under EPSG:3857. Switching to L.CRS.Simple — the usual pairing for a static image overlay — would silently invalidate every one of those assertions by changing what a 'degree' means on the map. This is the one decision in this feature where getting it wrong breaks the pipeline's own stated acceptance criteria, which is why it's called out explicitly rather than left as an implementation detail.
- **Trade-offs**:
  - The image's bounds box has to be sized in real lat/lng degrees (a trapezoid-shaped projection distortion at this latitude) rather than simple pixel-to-unit image coordinates, which is slightly less intuitive to size by hand than L.CRS.Simple would be.
- **Fitness Function**: map.spec.ts's real-geography assertions (centre, zoom-10 bounds, geolocation pan, Turf distance) keep passing unchanged — verified by tests/map.spec.ts run via scripts/ui-suite-in-chrome.sh in CI; any CRS regression fails these tests directly, no separate lint rule needed since the existing test suite already encodes this invariant

## Component Placement

- `Sword Coast map image overlay (L.imageOverlay call + its source image)` → Frameworks & Drivers / UI layer, `apps/ye-olde-magic-shop/src/views/LocationsView.vue (call site), apps/ye-olde-magic-shop/public/ (image + licence note)` — LocationsView.vue already owns its Leaflet map instance end-to-end (markers, polygon, popups, layer control); the tile-to-image swap is a change inside that existing ownership boundary, not a new one. public/ is the app's existing static-asset root.

## Bounded Context

- **Owning Context**: ye-olde-magic-shop frontend — shop-locations map (Leaflet)
- **Context Crossings**: packages/saturday-playwright-leaflet's test fixture contract (waitForMapLoad/checkTileFailures) — relied on read-only, not modified by this feature, tests/reagent-map-globe.spec.ts's off-origin-request-blocking pattern — reused, not owned, by the new no-network test

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] LocationsView.vue (adapter/UI layer) may depend on Leaflet (framework) directly — no domain or use-case layer exists in this app for map rendering to violate — This app has no Clean-Architecture domain/use-case split around the map feature; LocationsView.vue is itself the framework-facing component. loom MCP graph tools were not used — the repo graph/index availability wasn't verified, and the change is small enough that a direct read of the single file (apps/ye-olde-magic-shop/src/views/LocationsView.vue) was sufficient to confirm this, stated here per the fallback-disclosure requirement.

## Anti-Pattern Check

- Premature generalization: not found — A single L.imageOverlay() call was chosen over introducing an asset-loading abstraction or configurable tile-source strategy; nothing in this feature's scope needs more than one static image, so no abstraction is justified.
- Leaky abstraction: not found — No new abstraction is introduced, so there is nothing to leak.

## Fitness Functions

- **No tile or image request leaves the app's own origin when the locations view loads**
  - **Verification**: Playwright: new test in tests/map.spec.ts that page.route()-aborts any non-same-origin request, then asserts leaflet.waitForMapLoad() resolves and leaflet.checkTileFailures() is false
  - **Owner**: qa-engineer
- **Real-geography assertions (centre, zoom-10 bounds, geolocation pan, 300-400km Turf distance) keep passing under the default EPSG:3857 CRS**
  - **Verification**: tests/map.spec.ts unchanged, run via scripts/ui-suite-in-chrome.sh
  - **Owner**: qa-engineer
- **Every new or changed function in LocationsView.vue stays below cyclomatic complexity 7**
  - **Verification**: analyze-complexity skill / lint run on changed files in apps/ye-olde-magic-shop/src/views/LocationsView.vue
  - **Owner**: devops-engineer
- **The in-repo map image's licence is recorded next to the asset**
  - **Verification**: code-reviewer / tech-writer confirms a licence note file sits beside the image under apps/ye-olde-magic-shop/public/; no automated check can judge licence sufficiency
  - **Owner**: tech-writer

## Refactoring Opportunities

None

## Developer Handoff Notes

- LocationsView.vue line 74: L.map('map').setView([45.0, -120.0], 5) — do not add a crs option here or anywhere else in this file.
- LocationsView.vue line 80: the L.tileLayer('https://{s}.tile.openstreetmap.org/...') call is the exact line to remove and replace with L.imageOverlay(imageUrl, bounds, { className: 'map-tiles-parchment' }).addTo(map.value).
- Size the imageOverlay bounds generously beyond [41.0, 47.5] to [47.6, -121.1] / [-122.0, -118.0] so every marker, the Sword Coast polygon corners, and the zoom-6 geolocation flyTo centre all render inside the image with visible margin.
- Place the new image and its licence note together under apps/ye-olde-magic-shop/public/, following the existing flat-file convention (sits alongside aurelis.jpg, logo.png, favicon.svg) — no new subdirectory needed.
- Leave window.__SATURDAY_LEAFLET_MAP__ exposure (line 77) and the cleanup at line 157 untouched; they are outside this feature's scope.

## Open Architectural Questions

None
