---
intake: 1
id: reagent-map-view
type: feature
plan: deliver-feature
status: draft
priority: 2
created: 2026-11-02
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [saturday-playwright-cesium]
budget:
  max-attempts: 2
  max-usd: 40
  max-hours: 4
scope:
  paths: [apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/public/**, apps/ye-olde-magic-shop/package.json, apps/mock-api/**, pnpm-lock.yaml, docs/features/reagent-map-view/**]
---

# Plan reagent-gathering expeditions on a 3D globe

## Request

> Ye Olde Magic Shop gets a second map view: a 3D globe of an invented planet. Reagent sites are
> marked on it, and saved expeditions are drawn on it when the view loads. An adventurer draws a
> straight line to plan an expedition and saves it.

## Spec package (pinned)

Pilot wave 1, order D1. Built from spec package `reagent-map-spec`, accepted at every gate:

- `docs/specs/reagent-map-spec/` (criteria, seed data, domain model) and `features/reagent-map/`
- digest `sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413` over the files
  listed in `docs/specs/reagent-map-spec/package.sha256`. If that file's digest differs from this
  one, stop: the spec changed after this order was written.

`intake: 1` has no `inputs` field yet, so the pin is written here (D-30).

## Summary

A new view in the demo app shows Aurelis as a 3D CesiumJS globe (`SceneMode.SCENE3D`) with the
reagent sites, the expeditions the adventurer may see, a list of sites to fly to, and a line tool.
The mock API serves reagent sites and expeditions and stores new ones. The view implements the seam
contract `saturday-playwright-cesium` publishes, so the tests can act on it.

## Acceptance Criteria

- [ ] Given the spec package, then the view and API meet AC1–AC15 in
      `docs/specs/reagent-map-spec/acceptance-criteria.md`, with the settled answers there: 1 km
      tolerance, geodesic lines, nothing happens on a click off the land, each adventurer sees their
      own expeditions plus the seeded ones, and points closer than 1 km are the same place.
- [ ] Given `apps/mock-api`, then `GET /api/reagent-sites` returns the eight seeded sites;
      `GET /api/expeditions` returns the seeded expeditions plus, with a valid session, the caller's
      own and never another adventurer's; `POST /api/expeditions` requires a session, rejects a line
      whose ends are closer than 1 km on a sphere of radius 6,000 km, and stores the line for the
      caller. Seed data is exactly the spec's seed, including Elara's and Borin's expeditions and
      their users.
- [ ] Given the view, then Cesium runs with no network: no ion token or imagery, `baseLayerPicker`
      and `geocoder` off, one equirectangular image of Aurelis from the app's own files, and an
      ellipsoid that is a sphere of radius 6,000 km.
- [ ] Given the seam contract in `packages/saturday-playwright-cesium/README.md`, then the view
      exposes `window.__SATURDAY_CESIUM__ = { viewer, Cesium }` and gives every site and expedition
      entity its `saturdayMetadata` (`id`, `name`, `kind`).
- [ ] Given the existing Leaflet locations view and its suite, then they are unchanged and still pass.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under
      `docs/features/reagent-map-view/` excepted.

## Out of Scope

- The Saturday UI suite and the Sunday API suite for these scenarios: orders E1 and F1, run by hand
  from the same spec package (`docs/specs/reagent-map-spec/orders/`).
- Circles, enclosure and every wave 2 or 3 shape.
- Naming, editing or deleting an expedition.

## Domain Language

- Adventurer, Reagent Site, Expedition, Seeded Expedition, Aurelis: `DOMAIN_DICTIONARY.md`, section 7.

## Non-Functional Requirements

- No request leaves the app's and mock API's own origins while the view is open.

## Trust Boundaries

- Expeditions are per adventurer: the API returns another adventurer's expedition to no one, and
  stores an expedition only for the session's own adventurer. This is the security gate's check.

## Test Approach

Unit tests for the API's visibility and zero-length rules and for the view's pieces. The scenario
suites are E1 (Saturday, UI) and F1 (Sunday-style API, Playwright Test), written separately.

## Open Questions

- None: every question was answered at the spec's gates (`docs/specs/reagent-map-spec/corrections.md`).

## Infrastructure / Deploy Notes

- None
