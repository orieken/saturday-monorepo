---
intake: 1
id: map-heatmap-source
type: feature
plan: deliver-feature
status: draft
priority: 4
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [saturday-playwright-cesium, leaflet-package-gestures]
budget:
  max-attempts: 2
  max-usd: 30
  max-hours: 3
scope:
  paths: [packages/saturday-playwright-heatmap/**, packages/saturday-ml-analyzer/**, packages/saturday-playwright-cesium/**, packages/saturday-playwright-leaflet/**, docs/**]
---

# Report which map features no test touched

## Request

> Use the heatmap and the ML analyzer to track the items on the map.

## Summary

On a map, every click hits one canvas and the things to cover are map features, not DOM elements.
The heatmap gains a map source, so the analyzer can report which features no test came near.

## Acceptance Criteria

- [ ] Given a test using the Cesium or Leaflet fixture with the heatmap enabled, then each click on
      the map is recorded with its lat/lng as well as its page coordinates.
- [ ] Given the map's entities or markers with `saturdayMetadata`, then the heatmap records them as
      interactables, by id and name.
- [ ] Given that heatmap data, then `saturday-ml-analyzer` reports cold spots by feature name, with
      no change to its clustering, and a test on fixture data shows a feature with no nearby click
      reported cold.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/map-heatmap-source/` excepted.

## Out of Scope

- Finding the map on the page with ML; the seam already gives its container.

## Domain Language

- **Cold spot**: an interactable no test interaction came near.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Unit tests on recorded heatmap data; one Playwright test per map package.

## Open Questions

- "Near" on a map: Assumption if unresolved: a radius in metres from the analyzer's options.

## Infrastructure / Deploy Notes

- None
