---
intake: 1
id: saturday-playwright-cesium
type: feature
plan: deliver-feature
status: draft
priority: 2
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [ui-suite-real-chrome, leaflet-package-tests]
budget:
  max-attempts: 2
  max-usd: 40
  max-hours: 4
scope:
  paths: [packages/saturday-playwright-cesium/**, packages/saturday-geo/**, pnpm-lock.yaml, docs/**]
---

# Test CesiumJS maps with a Saturday package

## Request

> Build a Saturday package to test CesiumJS maps, like the one we have for Leaflet. Keep Cesium local
> with no account.

## Summary

`@orieken/saturday-playwright-cesium` lets Playwright tests act on a Cesium globe the way a user does
and read back what is on it, without depending on pixels. It is tested against its own fixture
globe: an invented planet, with no network and no Cesium ion account.

## Acceptance Criteria

- [ ] Given the package's seam contract in its README, then an app under test exposes
      `window.__SATURDAY_CESIUM__ = { viewer, Cesium }` and attaches `saturdayMetadata` (`id`,
      `name`, `kind`) to each entity it wants tests to find.
- [ ] Given `packages/saturday-playwright-cesium/fixtures/`, then a static page shows a Cesium
      `Viewer` with: a spherical ellipsoid of a declared radius; one equirectangular image covering
      the whole globe from the package's own files; `baseLayerPicker`, `geocoder` and ion imagery
      off; and known point, line and circle entities with metadata. It makes no request off its own
      origin.
- [ ] Given `cesium.waitForSettled()`, then it resolves once the globe's tiles are loaded, the
      camera has stopped and no render is pending, and rejects with a clear message after a timeout.
- [ ] Given `cesium.getEntities()`, then it returns each entity's metadata and geometry in degrees
      and metres: point position, line endpoints, circle centre and radius.
- [ ] Given `cesium.toScreen(lat, lng)`, then it returns page coordinates, or null when the point is
      on the far side of the globe or off-screen.
- [ ] Given `cesium.clickAt(lat, lng)`, then it clicks through `page.mouse`, and a test on the fixture
      shows the globe's pick returned that location within the tolerance `metersPerPixelAt` gives.
- [ ] Given `cesium.flyTo(lat, lng, heightMeters)`, then it moves the camera and waits until settled.
- [ ] Given distance, line-length and point-in-circle checks, then they take the globe's radius from
      the seam rather than assuming Earth's, and are unit-tested in Node without a browser. Whether
      they live in this package or a shared `saturday-geo` is decided at the design gate.
- [ ] Given the package's `test` script, then it runs the fixture suite in the real-Chrome image
      with no network, and every helper has a test that fails when the fixture's value changes.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/saturday-playwright-cesium/` excepted.

## Out of Scope

- Any change to the demo app. The reagent map view is a separate order, written from the reagent-map
  spec.
- Drawing-tool UX; that is a Flow in each app's tests.
- Terrain, 3D tiles and time-dynamic entities.

## Domain Language

- **Seam contract**: what an app must expose for this package to test it.
- **Settled**: tiles loaded, camera still, nothing waiting to render.

## Non-Functional Requirements

- The fixture page and suite make zero requests off their own origin.
- `waitForSettled` has an explicit timeout, default 15 seconds.

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Playwright Test against the package's fixture page, in the real-Chrome image; Node unit tests for
the geometry.

## Open Questions

- Planet radius and texture: Assumption if unresolved: a sphere of radius 6,371,008.8 m, Turf's
  default, with an in-repo equirectangular image of an invented land at 4096×2048.
- Pinned Cesium version supports a non-WGS84 ellipsoid for the whole viewer: Assumption if
  unresolved: pass the ellipsoid to `Globe` and to every geometry, and record any place it is not
  honoured.

## Infrastructure / Deploy Notes

Cesium's static assets (`Workers`, `Assets`, `Widgets`) are served from the package's fixture
server, with `CESIUM_BASE_URL` set to it.
