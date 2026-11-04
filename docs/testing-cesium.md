# Testing a Cesium globe with Saturday

This page explains how we test the demo app's reagent map, what the app had to change, and what went
wrong on the way. It covers both sides, so a second app can do the same. For the runner image
itself, see [runner-image.md](runner-image.md); for the package's API, see
[its README](../packages/saturday-playwright-cesium/README.md).

## The short version

A Cesium globe is one `<canvas>`. It has no DOM elements for a marker, so a normal locator cannot
find one. Three things replace locators:

1. **A seam.** The app exposes its Cesium viewer on `window`, and tags each entity with metadata.
2. **Gestures through real mouse events.** A latitude and longitude become a page position, then a
   real `page.mouse` click.
3. **State reads.** Tests ask the viewer what is on the globe, in degrees and metres.

## Where is the marker? No ML needed

To check that something landed at the right coordinates, tests use geometry against the viewer's
own state, not image recognition. `getEntities()` returns each entity's position in degrees, and
`distance()` compares it to the expected position using the planet's radius:

```typescript
const radiusMeters = await cesium.globeRadiusMeters();
const errorMeters = distance(expected, site.geometry.point, radiusMeters);
expect(errorMeters).toBeLessThanOrEqual(1_000);
```

Each kind of check answers a different question:

| Question | How it's checked |
|---|---|
| Is the entity at the right coordinates? | `getEntities()` with `distance()`, within the spec's 1 km |
| Did a click land where intended? | `clickAt(lat, lng)`, then compare what the globe picked, within `metersPerPixelAt` |
| Did the shop store the right shape? | Read it back through the API. This is the one check independent of rendering (AC11) |
| Is the point on land? | Read the planet image's pixel at those coordinates |
| Does it look right? | Screenshots, which state reads can't replace (see "What these tests can't see") |

`@orieken/saturday-ml-analyzer` answers a different question: which parts of a page no test touched.
It doesn't check correctness.

## What the app had to do

All of this is in `apps/ye-olde-magic-shop`, mostly `src/views/ReagentMapView.vue`.

### For the tests

| Change | Why |
|---|---|
| `window.__SATURDAY_CESIUM__ = { viewer, Cesium }`, set on mount and cleared on unmount | The seam. Every helper reads through it, and tests need `Cesium` for its projection functions |
| `entity.saturdayMetadata = { id, name, kind }` on each entity | How tests find a specific entity. `getEntities()` ignores anything untagged |
| `data-testid` on the site list and nav link | They are ordinary DOM, so ordinary locators work |

### To run with no network

| Change | Why |
|---|---|
| `baseLayerPicker: false`, `geocoder: false`, and one `baseLayer` | Cesium's defaults reach for ion, which needs a token and the network |
| `SingleTileImageryProvider` with `public/aurelis.jpg` | One local image covers the planet |
| `cesiumStaticAssetsPlugin` in `vite.config.ts` serves `/cesium/`, and `CESIUM_BASE_URL = '/cesium/'` in `main.ts` | Cesium loads workers and assets at run time. They come from `node_modules`, and nothing generated is committed |
| Cesium loaded as a global `<script>`, not `import` | `cesium@1.146.0` could not be bundled by Rollup for a production build, because of a packaging mismatch with `@cesium/engine` |
| The Cesium credit hidden | The default logo says "Cesium ion", and this view doesn't use ion |

**Still open:** `src/styles.css` imports Google Fonts on every page, so the no-request test skips
those two hosts. See "Known gaps".

### For an invented planet

| Change | Why |
|---|---|
| `new Cesium.Ellipsoid(r, r, r)` with `r = 6_000_000`, passed to the `Globe` and to every position | Without it, positions are computed on Earth's ellipsoid |
| The radius read from the seam in tests, never hard-coded to Earth's | A helper assuming Earth's radius would be about 6% wrong here, and would fail |
| A camera set explicitly on load | Otherwise Cesium uses its Earth-sized default view (see "What went wrong") |
| `scripts/generate-aurelis-image.py` draws land under every site, at 1024×512 | The first image put four sites in the sea, and 2048×1024 made rendering slow enough to time out |

### Everything else

- **A feature flag.** `VITE_FEATURE_REAGENT_MAP`, off in `.env`. `isReagentMapEnabled` is a small pure
  function, so it's tested without a browser.
- **Lazy loading.** The route uses `() => import(...)`, so nothing Cesium-related loads with the flag off.
- **Accessibility.** The globe has `role="img"` and an `aria-label` pointing at the site list, which
  is the screen-reader equivalent of the canvas. `aria-hidden` won't work, because Cesium makes its
  container focusable.
- **The API.** `GET /api/reagent-sites`, plus `GET` and `POST /api/expeditions`, in `apps/mock-api`.

## How the tests work

### The package

`@orieken/saturday-playwright-cesium` has its own suite, against a fixture globe in its `fixtures/`
folder.

| Helper | What it does |
|---|---|
| `waitForSettled()` | Waits for tiles loaded, camera still, and a frame drawn since the camera last moved. Rejects after 15 seconds, naming what it was waiting on |
| `getEntities()` | Every tagged entity, with its geometry in degrees and metres |
| `toScreen(lat, lng)` | Page coordinates, or `null` on the far side or off screen |
| `clickAt(lat, lng)` | A real `page.mouse` click at that point |
| `flyTo(lat, lng, heightMeters)` | Moves the camera, then waits to settle |
| `metersPerPixelAt(lat, lng)` | Ground size of a pixel, to set a click's tolerance |
| `globeRadiusMeters()` | The radius, from the seam |
| `distance`, `lineLength`, `pointInCircle` | Pure geometry, tested in Node without a browser |

There is deliberately no helper that draws through Cesium's API. Leaflet's `drawLine` does that, so
its tests pass even when the app's drawing tool is broken. "Choose the line tool, place two points,
save" belongs in a Saturday Flow in the app's test site, built on these gestures.

### The browser

WebGL needs real, headed Chrome with a software renderer:

```typescript
{
  name: 'chrome',
  use: {
    channel: 'chrome',
    headless: false,
    launchOptions: { args: ['--use-gl=angle', '--use-angle=swiftshader'] },
  },
}
```

The renderer is ANGLE SwiftShader. Mesa llvmpipe gave a null `webgl2` context in the runner image.

`tests/reagent-map-*.spec.ts` run only in the `chrome` project. A `testIgnore` keeps them out of the
bundled `chromium` project, so `npm run test:e2e` still works on a machine with no Xvfb.

### Running them

```bash
scripts/ui-suite-in-chrome.sh          # debian
scripts/ui-suite-in-chrome.sh --ubi9   # ubi9
```

The script does these steps in order:

1. Builds the image.
2. Starts the container with `--shm-size=1g`, because Chrome crashes tabs with Docker's default 64 MB.
3. Installs dependencies and builds the workspace packages.
4. Starts the app with `VITE_FEATURE_REAGENT_MAP=true`.
5. Runs `runner-check`.
6. Runs the `chrome` project with `--workers=3`, because more pages at once starve tile loading.

| Exit code | Means |
|---|---|
| 2 | Docker is not installed |
| 3 | The image build failed |
| 4 | The app did not become healthy in 180 seconds |
| 5 | `runner-check` failed |
| 6 | The dependency install or workspace build failed |
| 64 | The script got an unknown argument |
| Other | Playwright's own exit code |

Use a clean checkout. An untracked workspace package makes `pnpm install --frozen-lockfile` fail.

To run the package's own suite with networking off:

```bash
docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian \
  bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-cesium test'
```

Install dependencies inside a networked container first. A macOS `node_modules` won't run in Linux.

## What went wrong

Each of these passed its tests until someone looked closer. They're the reason for some of the rules
above.

| Problem | Why tests missed it | What fixed it |
|---|---|---|
| `toScreen` returned canvas coordinates, so clicks would miss on any page with a header | The fixture canvas sat at (0, 0), where the two are equal | Add the canvas's page offset. The fixture canvas now sits at (64, 48), and the old code fails |
| `waitForSettled` only ever waited for the first frame | Nothing moved the camera and then checked straight away | A camera move now requires a fresh frame |
| The default view was Cesium's Earth-sized one. Moonwell Spring was off-view and the globe was small | Tests read state, not pixels. A screenshot found it | An explicit camera, and a test that Moonwell is visible and Tidewrack isn't |
| Four of the eight sites were in the sea | Same | A generated image, and a test reading the pixel under each site |
| A killed developer stage left the suite passing but unreviewed | The stage hit the 30-minute limit | The work was finished by hand, and bigger orders are now split |
| A stale pnpm store volume failed with a missing rollup binary | Unrelated to the code | Delete the `ui-suite-pnpm-store` volume |
| Tests crashed or timed out at random, with a different test each run | Software rendering under emulation ran short of resources: 64 MB of shared memory, too many pages at once, and a large image | `--shm-size=1g`, `--workers=3`, and a 1024×512 image. 10 of 10 runs passed after that |

**The lesson:** make the fixture as awkward as a real page. Offset the canvas, use a planet that isn't
Earth-sized, and look at a screenshot at least once.

## What these tests can't see

State reads prove what the viewer holds, not what a person sees. A marker can have the right
coordinates and still be hidden, off camera or the wrong colour. Check a screenshot too.

## Known gaps

- **Google Fonts.** The no-request test skips `fonts.googleapis.com` and `fonts.gstatic.com`, because
  `src/styles.css` imports them on every page. Serve the fonts locally, then remove the exception.
- **No visual baselines.** The visual QA stage reports itself as unconfigured for this view. Generate
  and compare baselines only in the runner image, and start with the DOM around the globe.
- **No network blocking in the script.** `scripts/ui-suite-in-chrome.sh` uses Docker's default network.
  The app's "no request" check comes from a test that watches requests. The package's own suite does
  run with `--network none`.
- **No coverage tool.** The app has none, so reviewers check coverage by hand.

## Doing this in another app

1. Add `@orieken/saturday-playwright-cesium` as a dev dependency, and pin `cesium` to one exact version.
2. Build the viewer with your ellipsoid, one local image, and pickers and ion off.
3. Serve Cesium's assets from your own origin, and set `CESIUM_BASE_URL`.
4. Set an explicit camera.
5. Expose the seam, and tag your entities with `saturdayMetadata`.
6. Add a `chrome` project with the launch args above, and keep your specs out of any bundled-Chromium project.
7. Write the app's own gestures as Flows, and check results with `getEntities()` and your API.
8. Run it in both images, and look at a screenshot.
