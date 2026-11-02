# @orieken/saturday-playwright-cesium

A Playwright testing fixture for CesiumJS globes, following the same seam-contract pattern as
[`@orieken/saturday-playwright-leaflet`](../saturday-playwright-leaflet). It acts on a Cesium
`Viewer` the way a user does — reading back entities, projecting geographic points to screen
coordinates, clicking through `page.mouse`, and flying the camera — without depending on pixels.

## 1. Application Setup (the Seam Contract)

For this package to test your app, it must expose the Cesium `Viewer` instance (and the `Cesium`
module your app built it with) on `window`:

```typescript
import * as Cesium from 'cesium';

const viewer = new Cesium.Viewer('cesiumContainer', { /* ... */ });

// Expose the seam for testing
(window as any).__SATURDAY_CESIUM__ = { viewer, Cesium };
```

### Attach Metadata to Entities

To find and assert against specific entities, attach a `saturdayMetadata` object to the entity
before — or after — adding it to `viewer.entities`:

```typescript
const entity = viewer.entities.add({
  position: Cesium.Cartesian3.fromDegrees(longitude, latitude, 0),
  point: { pixelSize: 10 },
});

(entity as any).saturdayMetadata = { id: 'reagent-site-1', name: 'Reagent Site 1', kind: 'point' };
```

`kind` is passed through to test code unchanged — this package never branches on its value. It
dispatches on the entity's own Cesium shape (`position`, `polyline`, `ellipse`) to decide what
geometry to report, so `kind` can be any string your app finds useful.

Entities with no `saturdayMetadata` are invisible to `getEntities()` — exactly like an untagged
Leaflet layer is invisible to `saturday-playwright-leaflet`'s `getMarkers()`.

## 2. Test Setup

Import the extended `test` object from this package instead of `@playwright/test` directly. This
injects the `cesium` fixture into your tests.

```typescript
import { test } from '@orieken/saturday-playwright-cesium';
import { expect } from '@playwright/test';

test('reagent site is visible from the default view', async ({ page, cesium }) => {
  await page.goto('/map');
  await cesium.waitForSettled();

  const entities = await cesium.getEntities();
  expect(entities.find((e) => e.metadata.id === 'reagent-site-1')).toBeDefined();
});
```

## Available Helpers

- `cesium.waitForSettled(timeoutMs = 15000)`: Waits until the globe's tiles are loaded, the camera
  has stopped moving, and the first render has completed. Rejects with a message naming which
  condition (`tiles`, `camera`, `render`) is still pending if the globe never settles within the
  timeout.
- `cesium.getEntities()`: Returns every entity carrying `saturdayMetadata`, with its geometry in
  degrees and metres — point position, line endpoints, or circle centre and radius. Entities with
  no `saturdayMetadata` are excluded.
- `cesium.toScreen(lat, lng)`: Projects a geographic point to page coordinates, or returns `null`
  when the point is on the far side of the globe (occluded) or projects outside the canvas.
- `cesium.clickAt(lat, lng)`: Drives a real mouse click (via `page.mouse`, never a synthetic DOM
  event) at the screen position for the given point. Throws if the point is not visible.
- `cesium.flyTo(lat, lng, heightMeters)`: Moves the camera to the given point and height, then
  waits for the globe to settle at the destination.
- `cesium.metersPerPixelAt(lat, lng)`: Ground resolution, in metres per pixel, at the given point
  under the current camera position — used to size the pick-accuracy tolerance for `clickAt`.

## Domain Language

- **Seam contract**: what an app under test must expose (`window.__SATURDAY_CESIUM__`) for this
  package to test it.
- **Settled**: the globe's tiles are loaded, the camera has stopped moving, and nothing is waiting
  to render.

## Geometry Helpers

`distance`, `lineLength` and `pointInCircle` (exported from this package) compute great-circle
distance, polyline length, and point-in-circle membership on a sphere — all parameterized by an
explicit `radiusMeters`, never a hardcoded Earth radius. Pass your globe's own radius (read from
`viewer.scene.globe.ellipsoid` on the app side):

```typescript
import { distance } from '@orieken/saturday-playwright-cesium';

const radiusMeters = 6_000_000; // read from your app's own ellipsoid, not assumed
const metres = distance({ lat: 0, lng: 0 }, { lat: 0, lng: 1 }, radiusMeters);
```

These are pure functions with no dependency on `cesium` or `@playwright/test`, so they also run in
plain Node, with no browser.

## Testing

This package tests itself against a static fixture page under `fixtures/`. The page builds a
Cesium `Viewer` on an invented planet, **Aurelis** — a spherical ellipsoid of radius 6,000 km, not
WGS84 Earth — with `baseLayerPicker`, `geocoder` and ion imagery all off, one local equirectangular
image as its only imagery source, and known point, line and circle entities carrying
`saturdayMetadata`. The fixture server (`fixtures/serve.mjs`) binds to `127.0.0.1` only and serves
Cesium's own `Build/Cesium` assets and the planet texture directly from this package's
`node_modules` and `fixtures/assets/` — no CDN, no Cesium ion account, no network access at all.

Run the suite:

```bash
pnpm --filter @orieken/saturday-playwright-cesium test
```

This runs the Node-only `vitest` geometry suite first, then the Playwright fixture suite in real
Chrome (`channel: 'chrome'`, headed, with `--use-gl=angle --use-angle=swiftshader` — Cesium needs a
real `webgl2` context, which Leaflet never did; see `docs/runner-image.md`). A single non-zero exit
from either suite fails the command.

The suite also runs fully offline inside the Saturday runner image:

```bash
docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian \
  bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-cesium test'
```

Run `pnpm install` and `pnpm --filter @orieken/saturday-playwright-cesium build` before this
command, since it assumes dependencies are installed and the package is already built. The
`node_modules` directory must contain binaries resolved for the container's platform
(`linux-x64`) — if you installed on a different host OS/architecture, run `pnpm install` once more
inside a network-enabled container of the same image before the `--network none` run.
