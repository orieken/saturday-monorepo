---
intake: 1
id: leaflet-map-spec-gestures
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
depends-on: [leaflet-package-gestures, leaflet-local-tiles]
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 2
scope:
  paths: [apps/ye-olde-magic-shop/tests/map.spec.ts, docs/**]
---

# Fix the demo app's map tests to use real gestures

## Request

> Fix `tests/map.spec.ts` to use the gestures, and assert against the API where a shape is saved.

## Summary

The map tests stop drawing through Leaflet's API and say plainly which tests use test data.

## Acceptance Criteria

- [ ] Given `tests/map.spec.ts`, then no test calls the deprecated `drawLine`; a test that only needs
      a line to exist calls `addTestLine` and its name says so.
- [ ] Given a test that checks a marker or popup the user interacts with, then it uses `clickAt`.
- [ ] Given every test, then it still asserts at least what it asserted before.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/leaflet-map-spec-gestures/` excepted.

## Out of Scope

- Drawing tools in the app; there is nothing saved through the API on this view yet.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Playwright.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
