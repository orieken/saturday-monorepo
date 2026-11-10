---
intake: 1
id: reagent-sites-api
type: feature
plan: deliver-feature
status: ready
priority: 2
created: 2026-11-03
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: []
budget:
  max-attempts: 2
  max-usd: 15
  max-hours: 2
slice:
  of: reagent-map-spec
  spec: docs/specs/reagent-map-spec/acceptance-criteria.json
  digest: sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413
  manifest: docs/specs/reagent-map-spec/package.sha256
  covers: [AC1:api]
scope:
  paths: [apps/mock-api/**, docs/features/reagent-sites-api/**]
---

# Serve the reagent sites of Aurelis from the shop API

## Request

> The shop publishes the reagent sites of the invented planet Aurelis, so the reagent map can show them.

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

The mock API gains `GET /api/reagent-sites`, returning the spec's eight seeded sites. First of six slices of the reagent map view; API only, no UI.

The reagent map is behind the build-time flag `VITE_FEATURE_REAGENT_MAP` (introduced by
`reagent-map-globe`), off by default until `reagent-map-save` turns it on. With the flag off, the
app behaves exactly as it does today.

## Acceptance Criteria

- [ ] Given `apps/mock-api/data/reagent-sites.json`, then it holds exactly the eight sites in
      `docs/specs/reagent-map-spec/acceptance-criteria.json` (`seed.reagentSites`), each with an
      `id`, `name`, `lat` and `lng`.
- [ ] Given `GET /api/reagent-sites`, with or without a session, then it returns all eight sites
      with their positions unchanged (spec AC1, API scenario in `features/reagent-map/reagent-sites.feature`).
- [ ] Given `apps/mock-api` has no test runner, then its tests use Node's built-in `node:test`
      (no new dependency), run by a `test` script in `apps/mock-api/package.json`.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/reagent-sites-api/` excepted.

## Out of Scope

- Expeditions (`expeditions-api`). Any UI.
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
