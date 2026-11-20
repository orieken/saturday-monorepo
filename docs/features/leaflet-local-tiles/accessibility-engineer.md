---
feature: "leaflet-local-tiles"
bounded_context: "ye-olde-magic-shop frontend — shop-locations map (Leaflet)"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/src/views/LocationsView.vue", "apps/ye-olde-magic-shop/public/sword-coast-map.svg", "apps/ye-olde-magic-shop/public/sword-coast-map.LICENSE.txt"]
issue_refs: ["leaflet-local-tiles"]
linked_adrs: []
linked_kis: []
---

# Accessibility Report: leaflet-local-tiles

## Overall Status

**APPROVED**

## Scope Reviewed

- `LocationsView.vue`'s `onMounted` diff: the `L.tileLayer(OSM)` call is replaced with
  `L.imageOverlay('/sword-coast-map.svg', swordCoastImageBounds, { className: 'map-tiles-parchment' })`.
- The template (`<h2>`, `<p>`, the "Find Nearest Shop" `<button>`, `<div id="map">`) is unchanged by this diff.
- The new asset `apps/ye-olde-magic-shop/public/sword-coast-map.svg` and its adjacent licence note.

## Semantic HTML Review

No template markup was added, removed, or restructured by this change. The existing `<button>` for
"Find Nearest Shop", the heading/paragraph pair, and the `<div id="map">` map container are untouched
and remain correctly semantic (a real `<button>` triggers the geolocation action; no `onClick` was
attached to a non-interactive element by this diff).

## ARIA & Accessibility Review

- **New image overlay's accessible name**: Verified against the pinned Leaflet 1.9.4 source
  (`node_modules/.pnpm/leaflet@1.9.4/.../leaflet-src.js:9373-9375` and `:9554`) that `ImageOverlay`'s
  default `options.alt` is `''`, and `_initImage()` sets `img.alt = this.options.alt` from the merged
  options. The developer's call passes no `alt` option, so the default applies and the resulting
  `<img class="leaflet-image-layer">` renders with `alt=""`. This is the correct choice: the map's real
  information (shop names, inventory, rumors, the Sword Coast territory label) is carried by the marker
  popups and the GeoJSON popup, not by the base image, so marking it decorative is accurate and prevents
  assistive tech from announcing a meaningless image description.
- **Interactivity of the overlay**: `interactive` was not set, so it defaults to `false` — the overlay
  correctly stays out of the tab order and does not intercept pointer/keyboard focus intended for the
  markers and controls above it.
- No other ARIA attributes were added, removed, or needed a change as a result of this diff.

## Keyboard Navigation Review

No new interactive element was introduced. The existing "Find Nearest Shop" button and Leaflet's own
keyboard-pannable map/controls are unaffected by swapping the tile source. Nothing here regresses
Tab/Enter/Space/arrow-key reachability or `:focus-visible` behavior.

## Fixes Applied

## Escalations / Recommendations (non-blocking, pre-existing, outside this diff's scope)

- `<div id="map">` has no `aria-label`/`aria-labelledby` naming it as the shop-locations map for
  landmark/region navigation. This predates this change and is not something the current diff touches;
  fixing it would add an untested line against this feature's explicit scope constraints (no file
  outside `scope.paths` beyond what's needed, every changed line covered by a test). Recommend a
  follow-up ticket rather than folding it into this feature.
- `findNearestShop`'s geolocation failure path only does `console.error(...)` — a screen-reader or
  sighted user gets no visible/announced feedback when geolocation is denied or fails. Pre-existing,
  untouched by this diff; same follow-up recommendation applies.

## Feedback for the Developer

None — the implementation's only accessibility-relevant surface (the new image overlay) relies on
Leaflet's own correct default (`alt=""`, non-interactive), and no template or ARIA regression was
introduced.
