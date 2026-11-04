---
intake: 1
id: reagent-map-line-tool
type: feature
plan: deliver-feature
status: draft
priority: 2
created: 2026-11-03
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [reagent-map-expeditions]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/tests/reagent-map-*.spec.ts, docs/features/reagent-map-line-tool/**]
---

# Draw a line on the reagent map

## Request

> An adventurer draws a straight-line expedition by placing two points on the land.

## Spec package (pinned)

Part of pilot wave 1's D1, split into six orders so each developer stage fits loom's 30-minute
stage limit (L3.78). Built from spec package `reagent-map-spec`:

- `docs/specs/reagent-map-spec/` (criteria, seed data, domain model) and `features/reagent-map/`
- digest `sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413` over the files
  listed in `docs/specs/reagent-map-spec/package.sha256`. If that file's digest differs, stop: the
  spec changed after this order was written.

Settled at the spec's gates: 1 km tolerance on the planet's surface; geodesic lines; a click off the
land places nothing and shows nothing; an adventurer sees the shop's seeded expeditions plus their
own, signed out only the seeded ones; two points closer than 1 km are the same place; Aurelis is a
sphere of radius 6,000 km.

## Summary

Still behind the flag: a line tool on the reagent map. Placing two points on the land draws a geodesic line between them; anything else places nothing. No saving yet.

The reagent map is behind the build-time flag `VITE_FEATURE_REAGENT_MAP` (introduced by
`reagent-map-globe`), off by default until `reagent-map-save` turns it on. With the flag off, the
app behaves exactly as it does today.

## Acceptance Criteria

- [ ] Given the line tool started, when two points are placed on the land, then a completed geodesic
      line joins them and is ready to save (spec AC6, AC8); after one point, no line is complete and it
      cannot be saved (spec AC7).
- [ ] Given the line tool started, then any point on the near side of the globe can be placed, sea
      included; "the land" in the spec means the planet's surface for now (see
      `docs/specs/reagent-map-spec/corrections.md`).
- [ ] Given the line tool started, then a click on space beside the globe, or on the edge where the
      land faces away, places nothing, shows nothing, and the tool still waits (spec AC9, AC10).
- [ ] Given one point placed, when the line tool is started again, then the new line has no points
      (spec AC14).
- [ ] Given the line being drawn, then its points and line carry `saturdayMetadata` (`kind:
      'draft-point'`, `'draft-line'`) so the tests can read them through `getEntities()`.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/reagent-map-line-tool/` excepted.

## Out of Scope

- Saving the line: `reagent-map-save`.
- The scenario suites: E1 and F1, run by hand from the spec package.

## Domain Language

- Adventurer, Reagent Site, Expedition, Seeded Expedition, Aurelis: `DOMAIN_DICTIONARY.md`, section 7.

## Non-Functional Requirements

- No request leaves the app's and mock API's own origins.

## Trust Boundaries

- None beyond what `expeditions-api` enforces.

## Test Approach

Unit tests for each change. The UI and API scenario suites are E1 and F1.

## Open Questions

- None: every question was answered at the spec's gates (`docs/specs/reagent-map-spec/corrections.md`).

## Infrastructure / Deploy Notes

- None
