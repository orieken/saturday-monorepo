---
intake: 1
id: expeditions-api
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
depends-on: []
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 3
scope:
  paths: [apps/mock-api/**, docs/features/expeditions-api/**]
---

# Store and serve expeditions per adventurer

## Request

> An adventurer saves an expedition and finds it again; the shop's seeded expeditions are visible to everyone; nobody sees another adventurer's expeditions.

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

The mock API gains `GET /api/expeditions` and `POST /api/expeditions`, with the spec's visibility and zero-length rules and its seed data. API only, no UI.

The reagent map is behind the build-time flag `VITE_FEATURE_REAGENT_MAP` (introduced by
`reagent-map-globe`), off by default until `reagent-map-save` turns it on. With the flag off, the
app behaves exactly as it does today.

## Acceptance Criteria

- [ ] Given the seed, then the three shop-owned expeditions (`seed.seededExpeditions`) exist, and
      adventurers Elara and Borin exist as users with the expeditions `seed.adventurers` gives them.
- [ ] Given `GET /api/expeditions` with no session, then exactly the three seeded expeditions are
      returned (spec AC3).
- [ ] Given `GET /api/expeditions` with Elara's session, then the seeded expeditions and Elara's
      are returned, and Borin's never is (spec AC4). This is the trust boundary: the security gate checks it.
- [ ] Given `POST /api/expeditions` with a session and a line between two points at least 1 km
      apart on a sphere of radius 6,000 km, then it is stored for the session's adventurer (never an
      owner from the request body) and returned to them by `GET` (spec AC11).
- [ ] Given `POST /api/expeditions` with no session, then it is refused with 401 and nothing is
      stored (spec AC12); with ends less than 1 km apart, it is refused with 422 naming a zero-length
      line and nothing is stored (spec AC13).
- [ ] Given the tests, then they use Node's built-in `node:test`, as `reagent-sites-api` does
      (whichever ships first adds the `test` script).
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/expeditions-api/` excepted.

## Out of Scope

- Reagent sites (`reagent-sites-api`). Any UI. Editing or deleting an expedition.
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
