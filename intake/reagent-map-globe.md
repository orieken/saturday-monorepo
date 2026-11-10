---
intake: 1
id: reagent-map-globe
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
depends-on: [saturday-playwright-cesium, reagent-sites-api]
budget:
  max-attempts: 2
  max-usd: 30
  max-hours: 3
slice:
  of: reagent-map-spec
  spec: docs/specs/reagent-map-spec/acceptance-criteria.json
  digest: sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413
  manifest: docs/specs/reagent-map-spec/package.sha256
  covers: [AC1:ui, AC2]
scope:
  paths: [apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/public/**, apps/ye-olde-magic-shop/package.json, apps/ye-olde-magic-shop/vite.config.ts, apps/ye-olde-magic-shop/.env, pnpm-lock.yaml, apps/ye-olde-magic-shop/tests/reagent-map-*.spec.ts, apps/ye-olde-magic-shop/playwright.config.ts, scripts/ui-suite-in-chrome.sh, docs/features/reagent-map-globe/**]
---

# Show Aurelis and its reagent sites on a 3D globe

## Request

> Ye Olde Magic Shop gets a second map view: a 3D globe of an invented planet with the reagent sites marked on it.

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

A new view at `/reagent-map`, behind a feature flag that is off by default, shows Aurelis as an offline 3D Cesium globe with the reagent sites from `GET /api/reagent-sites`, a list of the sites beside it, and the seam contract `saturday-playwright-cesium` needs.

The reagent map is behind the build-time flag `VITE_FEATURE_REAGENT_MAP` (introduced by
`reagent-map-globe`), off by default until `reagent-map-save` turns it on. With the flag off, the
app behaves exactly as it does today.

## Acceptance Criteria

- [ ] Given `VITE_FEATURE_REAGENT_MAP` unset or not `true`, then there is no `/reagent-map` route and
      no link to it, and the app is otherwise unchanged. Given it is `true`, then the route and a nav
      link exist. `.env` sets it to `false`.
- [ ] Given the view, then Cesium runs with no network: no ion token or imagery; `baseLayerPicker`,
      `geocoder` and every ion default off; one equirectangular image of Aurelis from the app's own
      files; Cesium's static assets served by the app (`CESIUM_BASE_URL`); `SceneMode.SCENE3D`; and an
      ellipsoid that is a sphere of radius 6,000 km.
- [ ] Given the seam contract in `packages/saturday-playwright-cesium/README.md`, then the view sets
      `window.__SATURDAY_CESIUM__ = { viewer, Cesium }` and gives each site entity `saturdayMetadata`
      (`id`, `name`, `kind: 'site'`).
- [ ] Given the view, then every site is shown within 1 km of its position and listed by name beside
      the globe (spec AC1), including Frostcap Hollow near the pole and Tidewrack Shoals on the 180th
      meridian, each once (spec AC2).
- [ ] Given the view open, then no request leaves the app's and mock API's own origins.
- [ ] Given the flag, then whether the route exists is decided by a small pure function, tested for
      `true`, `false` and unset without a browser. `scripts/ui-suite-in-chrome.sh` starts the app with
      `VITE_FEATURE_REAGENT_MAP=true`, so the pre-ship check exercises the view, while `.env` keeps it
      off. Its new specs, named `tests/reagent-map-*.spec.ts`, run in the `chrome` project only (real
      Chrome with software WebGL), and `test:e2e` with bundled Chromium still passes without them.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/reagent-map-globe/` excepted.

## Out of Scope

- Expeditions, fly-to, the line tool, saving: the next three orders.
- The Leaflet locations view: unchanged.
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
