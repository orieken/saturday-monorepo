---
intake: 1
id: leaflet-local-tiles
type: feature
plan: deliver-feature
status: draft
priority: 3
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [ui-suite-real-chrome]
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 2
scope:
  paths: [apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/public/**, apps/ye-olde-magic-shop/tests/**, apps/mock-api/public/**, docs/**]
---

# Serve the Leaflet map without the internet

## Request

> Make Leaflet local if possible. `LocationsView.vue` loads tiles from `tile.openstreetmap.org`.

## Summary

The locations map uses a local image of the Sword Coast instead of OpenStreetMap tiles, so the view
and its tests work with no network.

## Acceptance Criteria

- [ ] Given the locations view, when it loads, then the map's imagery comes only from the app's own
      origin; no request leaves it.
- [ ] Given the existing markers, territory polygon and popups, then each appears in the same place
      relative to the others as before.
- [ ] Given a Playwright test that blocks every request not to the app's origin, then the locations
      view loads, `leaflet.waitForMapLoad()` resolves and `leaflet.checkTileFailures()` is false.
- [ ] Given the imagery, then its licence is recorded next to it (made in-repo, or with a licence
      that permits redistribution).
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/leaflet-local-tiles/` excepted.

## Out of Scope

- Changing the Leaflet package.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- The view makes zero requests to other origins.

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Playwright, run through `scripts/ui-suite-in-chrome.sh`.

## Open Questions

- Image overlay on `CRS.Simple` or local tiles: Assumption if unresolved: one image with
  `L.imageOverlay` and `CRS.Simple`, with marker coordinates converted and the conversion tested.

## Infrastructure / Deploy Notes

- None
