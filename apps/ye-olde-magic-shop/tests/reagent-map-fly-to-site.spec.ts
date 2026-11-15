/**
 * @issue reagent-map-expeditions
 * Covers AC5 of fly-to-site.feature: choosing a site from the list flies the
 * camera so that site lands on screen on the near side of the globe, even
 * one that starts on the far side — Tidewrack Shoals from the default view.
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 */
import reagentSitesFixture from '../../mock-api/data/reagent-sites.json';
import { test, type CesiumFixture } from '@orieken/saturday-playwright-cesium';
import { expect, type Page } from '@playwright/test';

type SeededSite = { id: string; name: string; lat: number; lng: number };
type CameraState = { lat: number; lng: number; heightMeters: number };

const SEEDED_SITES: SeededSite[] = reagentSitesFixture;
const TIDEWRACK_SHOALS = SEEDED_SITES.find((site) => site.id === 'tidewrack-shoals')!;
// Cesium's default flyTo animation takes several seconds for a near-antipodal
// destination (measured by the developer when this was implemented), and the
// destination can flicker on screen mid-arc before the camera actually stops
// there. Calling waitForSettled() right after the click can also race its
// moveStart event (see implementation-notes.md's notesForQa) — so this waits
// for the camera to actually start, then leans on waitForSettled() to know
// when it has actually stopped, before reading the final on-screen position.
const CAMERA_MOVE_TIMEOUT_MS = 20_000;

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
  const carto = Cesium.Cartographic.fromCartesian(viewer.camera.positionWC, viewer.scene.globe.ellipsoid);
  return {
    lat: Cesium.Math.toDegrees(carto.latitude),
    lng: Cesium.Math.toDegrees(carto.longitude),
    heightMeters: carto.height,
  };
}

async function waitForFlightToSettle(page: Page, cesium: CesiumFixture, before: CameraState): Promise<void> {
  await expect
    .poll(() => page.evaluate(readCamera), { timeout: CAMERA_MOVE_TIMEOUT_MS })
    .not.toEqual(before);
  await cesium.waitForSettled();
}

test.describe('Flying to a reagent site (AC5)', () => {
  test.beforeEach(async ({ page, cesium }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();
  });

  test('Tidewrack Shoals starts on the far side, off screen', async ({ cesium }) => {
    expect(await cesium.toScreen(TIDEWRACK_SHOALS.lat, TIDEWRACK_SHOALS.lng)).toBeNull();
  });

  test('choosing Tidewrack Shoals from the list flies the camera so it lands on screen on the near side', async ({
    page,
    cesium,
  }) => {
    expect(await cesium.toScreen(TIDEWRACK_SHOALS.lat, TIDEWRACK_SHOALS.lng)).toBeNull();
    const before = await page.evaluate(readCamera);

    await page.getByTestId(`fly-to-site-${TIDEWRACK_SHOALS.id}`).click();
    await waitForFlightToSettle(page, cesium, before);

    const canvasBox = await page.locator('#cesiumContainer').boundingBox();
    const point = await cesium.toScreen(TIDEWRACK_SHOALS.lat, TIDEWRACK_SHOALS.lng);
    expect(canvasBox).not.toBeNull();
    expect(point).not.toBeNull();
    expect(point!.x).toBeGreaterThanOrEqual(canvasBox!.x);
    expect(point!.x).toBeLessThanOrEqual(canvasBox!.x + canvasBox!.width);
    expect(point!.y).toBeGreaterThanOrEqual(canvasBox!.y);
    expect(point!.y).toBeLessThanOrEqual(canvasBox!.y + canvasBox!.height);
  });
});
