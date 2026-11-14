/**
 * @issue reagent-map-pan-zoom
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 *
 * Pans and zooms the reagent map with real mouse gestures (page.mouse drag
 * and wheel), the way an adventurer does, never through Cesium's camera API.
 * State reads then check where the camera went: its height above Aurelis,
 * and where a site lands on screen.
 */
import reagentSitesFixture from '../../mock-api/data/reagent-sites.json';
import { test, distance, type CesiumFixture, type CesiumScreenPoint } from '@orieken/saturday-playwright-cesium';
import { expect, type Page } from '@playwright/test';

type SeededSite = { id: string; name: string; lat: number; lng: number };
type CameraState = { lat: number; lng: number; heightMeters: number };

const SEEDED_SITES: SeededSite[] = reagentSitesFixture;
const MOONWELL_SPRING = SEEDED_SITES.find((site) => site.id === 'moonwell-spring')!;
// The wheel tests zoom over Ember Crag, not Moonwell Spring. Moonwell sits at
// the exact centre of the default view, and a wheel there makes Cesium 1.146
// rotate the camera about a zero-length axis: the camera position becomes NaN
// and rendering stops. That is an open defect, not covered here.
const EMBER_CRAG = SEEDED_SITES.find((site) => site.id === 'ember-crag')!;

const DRAG = { dx: 160, dy: 90 };
const DRAG_STEPS = 20;
// Cesium adds inertia ("flick") to a drag released less than 0.4 s after the
// button went down, so the globe keeps spinning past the cursor. Holding the
// button longer makes this a plain grab-and-drag. This is gesture timing,
// not a wait for state: every state wait below is a wait-for-condition.
const HOLD_BEFORE_RELEASE_MS = 600;
// A grabbed point stays under the cursor. Measured error: 4.1 px in each of
// 6 runs. With the camera controller on Earth's ellipsoid instead of Aurelis's
// (the bug this test found), the error was 18.9 px in each of 3 runs.
const GRAB_TOLERANCE_PIXELS = 8;
// A drag pans; it must not zoom. 1% of 12,000 km is 120 km.
const PAN_HEIGHT_TOLERANCE_RATIO = 0.01;
const WHEEL_NOTCH_DELTA = 100;
const WHEEL_NOTCHES = 5;
// Five wheel notches move the camera far more than this. The floor stops a
// one-metre jitter from counting as a zoom.
const MINIMUM_ZOOM_RATIO = 0.1;
const CAMERA_MOVE_TIMEOUT_MS = 15_000;
const POSITION_TOLERANCE_METERS = 1_000;

function readCamera(): CameraState {
  type Seam = {
    viewer: { camera: { positionWC: unknown }; scene: { globe: { ellipsoid: unknown } } };
    Cesium: {
      Cartographic: { fromCartesian(cartesian: unknown, ellipsoid: unknown): { latitude: number; longitude: number; height: number } };
      Math: { toDegrees(radians: number): number };
    };
  };
  const seam = (window as unknown as { __SATURDAY_CESIUM__?: Seam }).__SATURDAY_CESIUM__;
  if (!seam) throw new Error('Cesium seam not found: expected window.__SATURDAY_CESIUM__ = { viewer, Cesium }');
  const { viewer, Cesium } = seam;
  // Aurelis's own ellipsoid: camera.positionCartographic would use Earth's.
  const carto = Cesium.Cartographic.fromCartesian(viewer.camera.positionWC, viewer.scene.globe.ellipsoid);
  return {
    lat: Cesium.Math.toDegrees(carto.latitude),
    lng: Cesium.Math.toDegrees(carto.longitude),
    heightMeters: carto.height,
  };
}

async function waitForCameraToLeave(page: Page, cesium: CesiumFixture, before: CameraState): Promise<void> {
  await expect
    .poll(() => page.evaluate(readCamera), { timeout: CAMERA_MOVE_TIMEOUT_MS })
    .not.toEqual(before);
  await cesium.waitForSettled();
}

async function dragGlobe(page: Page, from: CesiumScreenPoint, to: CesiumScreenPoint): Promise<void> {
  await page.mouse.move(from.x, from.y);
  await page.mouse.down();
  await page.mouse.move(to.x, to.y, { steps: DRAG_STEPS });
  await page.waitForTimeout(HOLD_BEFORE_RELEASE_MS);
  await page.mouse.up();
}

async function spinWheel(page: Page, at: CesiumScreenPoint, deltaPerNotch: number): Promise<void> {
  await page.mouse.move(at.x, at.y);
  for (let notch = 0; notch < WHEEL_NOTCHES; notch++) {
    await page.mouse.wheel(0, deltaPerNotch);
  }
}

async function siteOnScreen(cesium: CesiumFixture, site: SeededSite): Promise<CesiumScreenPoint> {
  const point = await cesium.toScreen(site.lat, site.lng);
  expect(point, `${site.name} should be on screen`).not.toBeNull();
  return point!;
}

test.describe('Reagent map pan and zoom', () => {
  test.beforeEach(async ({ page, cesium }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();
  });

  test('dragging the globe pans it: the grabbed site follows the cursor and the height stays the same', async ({
    page,
    cesium,
  }) => {
    const before = await page.evaluate(readCamera);
    const grabbed = await siteOnScreen(cesium, MOONWELL_SPRING);
    const dropAt = { x: grabbed.x + DRAG.dx, y: grabbed.y + DRAG.dy };

    await dragGlobe(page, grabbed, dropAt);
    await waitForCameraToLeave(page, cesium, before);

    const landed = await siteOnScreen(cesium, MOONWELL_SPRING);
    expect(Math.hypot(landed.x - dropAt.x, landed.y - dropAt.y)).toBeLessThanOrEqual(GRAB_TOLERANCE_PIXELS);

    const after = await page.evaluate(readCamera);
    // Dragging the globe right and down turns it east and south under the
    // camera, so the camera ends up west and north of where it started.
    expect(after.lng).toBeLessThan(before.lng);
    expect(after.lat).toBeGreaterThan(before.lat);
    expect(Math.abs(after.heightMeters - before.heightMeters) / before.heightMeters).toBeLessThanOrEqual(
      PAN_HEIGHT_TOLERANCE_RATIO
    );
  });

  test('scrolling the wheel up zooms in: the camera drops and each pixel covers less ground', async ({
    page,
    cesium,
  }) => {
    const before = await page.evaluate(readCamera);
    const metersPerPixelBefore = await cesium.metersPerPixelAt(EMBER_CRAG.lat, EMBER_CRAG.lng);
    const emberCrag = await siteOnScreen(cesium, EMBER_CRAG);

    await spinWheel(page, emberCrag, -WHEEL_NOTCH_DELTA);
    await waitForCameraToLeave(page, cesium, before);

    const after = await page.evaluate(readCamera);
    expect(after.heightMeters).toBeLessThan(before.heightMeters * (1 - MINIMUM_ZOOM_RATIO));
    expect(await cesium.metersPerPixelAt(EMBER_CRAG.lat, EMBER_CRAG.lng)).toBeLessThan(
      metersPerPixelBefore
    );
    // Cesium zooms towards the cursor, so the site under it stays in view.
    await siteOnScreen(cesium, EMBER_CRAG);
  });

  test('scrolling the wheel down zooms out: the camera rises and each pixel covers more ground', async ({
    page,
    cesium,
  }) => {
    const before = await page.evaluate(readCamera);
    const metersPerPixelBefore = await cesium.metersPerPixelAt(EMBER_CRAG.lat, EMBER_CRAG.lng);
    const emberCrag = await siteOnScreen(cesium, EMBER_CRAG);

    await spinWheel(page, emberCrag, WHEEL_NOTCH_DELTA);
    await waitForCameraToLeave(page, cesium, before);

    const after = await page.evaluate(readCamera);
    expect(after.heightMeters).toBeGreaterThan(before.heightMeters * (1 + MINIMUM_ZOOM_RATIO));
    expect(await cesium.metersPerPixelAt(EMBER_CRAG.lat, EMBER_CRAG.lng)).toBeGreaterThan(
      metersPerPixelBefore
    );
  });

  test('panning and zooming move the camera, never the sites', async ({ page, cesium }) => {
    const before = await page.evaluate(readCamera);
    const grabbed = await siteOnScreen(cesium, MOONWELL_SPRING);

    await dragGlobe(page, grabbed, { x: grabbed.x + DRAG.dx, y: grabbed.y + DRAG.dy });
    await waitForCameraToLeave(page, cesium, before);
    const afterPan = await page.evaluate(readCamera);
    await spinWheel(page, await siteOnScreen(cesium, EMBER_CRAG), -WHEEL_NOTCH_DELTA);
    await waitForCameraToLeave(page, cesium, afterPan);

    const radiusMeters = await cesium.globeRadiusMeters();
    // reagent-map-expeditions now also tags expedition polylines with
    // saturdayMetadata (kind: 'expedition'), so getEntities() returns both
    // kinds — filter to this test's own kind before counting.
    const entities = (await cesium.getEntities()).filter((entity) => entity.metadata.kind === 'site');
    expect(entities).toHaveLength(SEEDED_SITES.length);
    for (const seeded of SEEDED_SITES) {
      const entity = entities.find((candidate) => candidate.metadata.id === seeded.id);
      expect(entity, `expected an entity for ${seeded.id}`).toBeDefined();
      const geometry = entity!.geometry as { point: { lat: number; lng: number } };
      expect(distance({ lat: seeded.lat, lng: seeded.lng }, geometry.point, radiusMeters)).toBeLessThanOrEqual(
        POSITION_TOLERANCE_METERS
      );
    }
  });
});
