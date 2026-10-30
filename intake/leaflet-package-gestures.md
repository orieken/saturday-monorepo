---
intake: 1
id: leaflet-package-gestures
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
depends-on: [leaflet-package-tests]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [packages/saturday-playwright-leaflet/**, docs/**]
---

# Let Leaflet tests draw the way a user does

## Request

> `drawLine` adds a polyline through Leaflet's API, so a test using it passes when the app's drawing
> tool is broken.

## Summary

The package gains gestures that act through the mouse on the map, so tests prove what a user can do.

## Acceptance Criteria

- [ ] Given `leaflet.toScreen(lat, lng)`, then it returns the page coordinates of that point, or
      null when the point is outside the visible map.
- [ ] Given `leaflet.clickAt(lat, lng)` and `leaflet.dragFrom(from, to)`, then each acts through
      `page.mouse` at the converted coordinates, and a test on the fixture page shows the map's
      click handler received the expected lat/lng within one pixel's distance.
- [ ] Given `drawLine`, then it is renamed `addTestLine`, documented as a test-data helper that
      bypasses the UI, and the old name is kept as a deprecated alias for one release.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/leaflet-package-gestures/` excepted.

## Out of Scope

- Drawing tools in the demo app.

## Domain Language

- **Gesture**: a helper that acts through the mouse or keyboard, the way a user does.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Playwright Test against the package's fixture page.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
