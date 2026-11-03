---
intake: 1
id: reagent-map-expeditions
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
depends-on: [reagent-map-globe, expeditions-api]
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/src/**, docs/features/reagent-map-expeditions/**]
---

# Draw visible expeditions and fly to a site on the reagent map

## Request

> Saved expeditions are drawn on the globe when the view loads, and the camera can be moved to any reagent site.

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

Still behind the flag: the reagent map draws the expeditions `GET /api/expeditions` returns for the current adventurer, and choosing a site from the list flies the camera to it.

The reagent map is behind the build-time flag `VITE_FEATURE_REAGENT_MAP` (introduced by
`reagent-map-globe`), off by default until `reagent-map-save` turns it on. With the flag off, the
app behaves exactly as it does today.

## Acceptance Criteria

- [ ] Given the view, then each expedition from `GET /api/expeditions` is drawn as a geodesic line with
      both ends within 1 km of its stored ends, with `saturdayMetadata` (`id`, `name`, `kind:
      'expedition'`): the three seeded ones signed out (spec AC3), plus the adventurer's own when
      signed in (spec AC4).
- [ ] Given a site chosen from the list, including Tidewrack Shoals from the default view where it
      is on the far side, then the camera flies so the site is on screen on the near side of the globe
      (spec AC5).
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/reagent-map-expeditions/` excepted.

## Out of Scope

- The line tool and saving: the next two orders.
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
