---
feature: "leaflet-local-tiles"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: ["leaflet-local-tiles"]
linked_adrs: []
linked_kis: []
---

# Code Review: leaflet-local-tiles

## Overall Status

**APPROVED**

## Design Narrative

LocationsView.vue's onMounted swaps the network-dependent L.tileLayer(OSM) call for a single L.imageOverlay() over a new in-repo SVG, kept inside Leaflet's default EPSG:3857 CRS so the real-geography assertions in map.spec.ts stay valid; the asset and its CC BY 4.0 licence note sit flat in public/ alongside the app's other static assets, and a new Playwright test blocks every off-origin request and directly verifies the image overlay's own load success. Behaviour change: yes — the map's visual tile source moves from live internet tiles to a static local image.

## Design Score

- **Clarity** 4/5
- **Cohesion** 4/5
- **Coupling** 4/5
- **Craft** 4/5

## Security Surface

- Removes the view's only external network call (tile.openstreetmap.org) — reduces exfiltration/SSRF-adjacent surface; no new trust boundary is introduced by the image overlay, which is served same-origin from public/.

## Performance Surface

- No per-tile network fetches on pan/zoom — eliminates external request latency and flakiness entirely, which is the point of the change.
- Static single-image overlay cannot stream additional detail on zoom the way tiles could, but architecture-notes.md already priced this trade-off as a non-issue for a fixed-purpose test fixture view.

## Test Design Review

- 'locations view loads with every off-origin request blocked (NFR)' (apps/ye-olde-magic-shop/tests/map.spec.ts:139): YES. Verified Leaflet 1.9.4's own source (node_modules/.pnpm/leaflet@1.9.4/.../leaflet-src.js:9528) adds the 'leaflet-image-layer' class directly to the <img> element itself, so document.querySelector('.leaflet-image-layer') returns the real image node and the complete/naturalWidth>0 check genuinely fails if the overlay's own same-origin request is blocked or the image fails to decode. This resolves the prior iteration's vacuous-assertion finding (checkTileFailures() alone only inspected the now-absent .leaflet-tile elements); checkTileFailures() is correctly retained only as a secondary guard.
- No other test in tests/map.spec.ts was added or modified by this diff — the 9 pre-existing tests are unchanged text, consistent with the spec's 'passes unchanged' requirement.

## Verification of Developer Self-Review

- Developer's Self-Review Checklist and Simple Design Verification hold up against the diff: the new assertion is a single page.evaluate() with an intention-revealing local (imageOverlayLoaded), no new abstraction, no function over 30 LOC, no magic numbers (bounds and APP_ORIGIN are named constants). No discrepancy found.
- Named Refactoring Log entry ('Replace Vacuous Check with Direct Assertion') accurately describes what changed and matches my own prior finding's instruction — scoped to the test file only, left checkTileFailures()/fixture.ts untouched as instructed.

## Feedback for the Developer

None
