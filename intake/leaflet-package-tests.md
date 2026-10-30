---
intake: 1
id: leaflet-package-tests
type: chore
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
  max-usd: 25
  max-hours: 3
scope:
  paths: [packages/saturday-playwright-leaflet/**, docs/**]
---

# Give the Leaflet package its own tests

## Request

> `saturday-playwright-leaflet` has no tests; its `test` script exits 1.

## Summary

The package is tested against a static fixture page it ships with, so it can be changed safely and
serves as the pattern for the Cesium package.

## Acceptance Criteria

- [ ] Given `packages/saturday-playwright-leaflet/fixtures/`, then a static page loads Leaflet from
      the package's own dependencies (no CDN), a local image, and a known set of markers with
      `saturdayMetadata`, a polygon and a layer control, and exposes `__SATURDAY_LEAFLET_MAP__`.
- [ ] Given the package's `test` script, then it serves that page and runs a Playwright suite with at
      least one test per helper in the README, each asserting a known value from the fixture.
- [ ] Given each helper's test, then it fails when the fixture's corresponding value is changed.
- [ ] Given the suite, then it runs in the real-Chrome image with no network.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/leaflet-package-tests/` excepted.

## Out of Scope

- New helpers; changing `drawLine`.

## Domain Language

- No new terms introduced.

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
