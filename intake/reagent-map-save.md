---
intake: 1
id: reagent-map-save
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
depends-on: [reagent-map-line-tool]
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 3
slice:
  of: reagent-map-spec
  spec: docs/specs/reagent-map-spec/acceptance-criteria.json
  digest: sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413
  manifest: docs/specs/reagent-map-spec/package.sha256
  covers: [AC11:ui, AC12:ui, AC13:ui, AC15]
scope:
  paths: [apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/.env, apps/ye-olde-magic-shop/tests/reagent-map-*.spec.ts, docs/features/reagent-map-save/**]
---

# Save expeditions from the reagent map and turn it on

## Request

> An adventurer saves an expedition and sees it again after reopening the map.

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

Saving from the line tool through `POST /api/expeditions`, with the sign-in and zero-length messages, then the feature flag on by default: the reagent map is complete for wave 1.

The reagent map is behind the build-time flag `VITE_FEATURE_REAGENT_MAP` (introduced by
`reagent-map-globe`), off by default until `reagent-map-save` turns it on. With the flag off, the
app behaves exactly as it does today.

## Acceptance Criteria

- [ ] Given Elara signed in with a completed line, when she saves and reopens the map, then her
      expedition is shown and the API holds it for her within 1 km of where she placed it (spec AC11).
- [ ] Given no adventurer signed in, when a completed line is saved, then they are asked to sign in and
      nothing is stored (spec AC12).
- [ ] Given both points less than 1 km apart, when saving, then the adventurer is told it is a
      zero-length line and nothing is stored (spec AC13).
- [ ] Given an unsaved line, when the adventurer leaves the map and opens it again, then the line is
      not shown and nothing is stored (spec AC15).
- [ ] Given the change is complete, then `.env` sets `VITE_FEATURE_REAGENT_MAP=true` and the nav link
      is shown, with the flag still able to turn the view off.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/reagent-map-save/` excepted.

## Out of Scope

- Naming, editing or deleting an expedition. Circles and later waves.
- The scenario suites: E1 and F1, run by hand from the spec package.

## Domain Language

- Adventurer, Reagent Site, Expedition, Seeded Expedition, Aurelis: `DOMAIN_DICTIONARY.md`, section 7.

## Non-Functional Requirements

- No request leaves the app's and mock API's own origins.

## Trust Boundaries

- Expeditions are per adventurer: never return another adventurer's expedition, and store one only for the session's own adventurer.

## Test Approach

Unit tests for each change. The UI and API scenario suites are E1 and F1.

## Open Questions

- None: every question was answered at the spec's gates (`docs/specs/reagent-map-spec/corrections.md`).

## Infrastructure / Deploy Notes

- None
