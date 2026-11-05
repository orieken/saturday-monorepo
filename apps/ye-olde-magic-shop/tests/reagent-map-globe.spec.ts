/**
 * @issue reagent-map-globe
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 * Exercises AC2 (offline Cesium config), AC3 (Saturday seam + metadata), AC4
 * (all 8 sites shown once, within 1km, listed by name), and the NFR that no
 * request leaves the app's or the mock API's own origin.
 */
import reagentSitesFixture from '../../mock-api/data/reagent-sites.json';
import { test, distance } from '@orieken/saturday-playwright-cesium';
import { expect, type Page } from '@playwright/test';

type SeededSite = { id: string; name: string; lat: number; lng: number };

const SEEDED_SITES: SeededSite[] = reagentSitesFixture;
const FROSTCAP_HOLLOW = SEEDED_SITES.find((site) => site.id === 'frostcap-hollow')!;
const TIDEWRACK_SHOALS = SEEDED_SITES.find((site) => site.id === 'tidewrack-shoals')!;
const POSITION_TOLERANCE_METERS = 1_000;
const EXPECTED_RADIUS_METERS = 6_000_000;
const APP_ORIGIN_HOSTNAME = 'localhost';
const ALLOWED_PORTS = ['8000', '8001'];

async function collectOffOriginRequests(page: Page): Promise<string[]> {
  const offOrigin: string[] = [];
  page.on('request', (request) => {
    const url = new URL(request.url());
    // blob:/data: are in-memory, non-network resources — only http(s)
    // requests can actually leave the page's own origin.
    if (url.protocol !== 'http:' && url.protocol !== 'https:') return;
    if (url.hostname !== APP_ORIGIN_HOSTNAME || !ALLOWED_PORTS.includes(url.port)) {
      offOrigin.push(request.url());
    }
  });
  return offOrigin;
}

test.describe('Reagent map globe (AC2, AC3, AC4)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/reagent-map');
  });

  test('every reagent site sits on land in the planet image', async ({ cesium, page }) => {
    await cesium.waitForSettled();

    const landBySite = await page.evaluate(async (sites) => {
      const image = new Image();
      image.src = '/aurelis.jpg';
      await image.decode();
      const canvas = document.createElement('canvas');
      canvas.width = image.width;
      canvas.height = image.height;
      const context = canvas.getContext('2d')!;
      context.drawImage(image, 0, 0);
      return sites.map((site) => {
        const x = Math.min(image.width - 1, Math.floor(((site.lng + 180) / 360) * image.width));
        const y = Math.min(image.height - 1, Math.floor(((90 - site.lat) / 180) * image.height));
        const [, green, blue] = context.getImageData(x, y, 1, 1).data;
        return [site.name, green > blue];
      });
    }, SEEDED_SITES);

    expect(landBySite).toEqual(SEEDED_SITES.map((site) => [site.name, true]));
  });

  test('shows no Cesium ion credit, since the view does not use ion', async ({ cesium, page }) => {
    await cesium.waitForSettled();
    await expect(page.locator('.cesium-widget-credits')).toBeHidden();
  });

  test('the default view is centred over Moonwell Spring, with Tidewrack Shoals on the far side', async ({
    cesium,
  }) => {
    await cesium.waitForSettled();
    const moonwell = SEEDED_SITES.find((site) => site.id === 'moonwell-spring')!;

    expect(await cesium.toScreen(moonwell.lat, moonwell.lng)).not.toBeNull();
    expect(await cesium.toScreen(TIDEWRACK_SHOALS.lat, TIDEWRACK_SHOALS.lng)).toBeNull();
  });

  test('runs Cesium with no ion token/imagery, pickers off, SCENE3D, and a 6,000km-radius sphere (AC2)', async ({
    cesium,
    page,
  }) => {
    await cesium.waitForSettled();

    const config = await page.evaluate(() => {
      const { viewer, Cesium } = (window as any).__SATURDAY_CESIUM__;
      const ellipsoid = viewer.scene.globe.ellipsoid;
      return {
        sceneMode: viewer.scene.mode,
        scene3dValue: Cesium.SceneMode.SCENE3D,
        radii: [ellipsoid.radii.x, ellipsoid.radii.y, ellipsoid.radii.z],
        baseLayerPicker: viewer.baseLayerPicker,
        geocoder: viewer.geocoder,
        imageryLayerCount: viewer.imageryLayers.length,
      };
    });

    expect(config.sceneMode).toBe(config.scene3dValue);
    expect(config.radii).toEqual([EXPECTED_RADIUS_METERS, EXPECTED_RADIUS_METERS, EXPECTED_RADIUS_METERS]);
    expect(config.baseLayerPicker).toBeFalsy();
    expect(config.geocoder).toBeFalsy();
    expect(config.imageryLayerCount).toBe(1);
    expect(await cesium.globeRadiusMeters()).toBe(EXPECTED_RADIUS_METERS);
  });

  test('exposes window.__SATURDAY_CESIUM__ and tags every site entity with saturdayMetadata (AC3)', async ({
    cesium,
    page,
  }) => {
    await cesium.waitForSettled();

    const seamShape = await page.evaluate(() => {
      const seam = (window as any).__SATURDAY_CESIUM__;
      return {
        hasViewer: Boolean(seam?.viewer),
        hasCesium: Boolean(seam?.Cesium),
      };
    });
    expect(seamShape.hasViewer).toBe(true);
    expect(seamShape.hasCesium).toBe(true);

    const entities = await cesium.getEntities();
    expect(entities).toHaveLength(SEEDED_SITES.length);
    for (const entity of entities) {
      expect(entity.metadata.kind).toBe('site');
      expect(entity.metadata.id).toBeTruthy();
      expect(entity.metadata.name).toBeTruthy();
    }
  });

  test('shows all 8 reagent sites exactly once, each within 1km of its published position (AC4)', async ({
    cesium,
  }) => {
    await cesium.waitForSettled();
    const radiusMeters = await cesium.globeRadiusMeters();
    const entities = await cesium.getEntities();

    expect(entities).toHaveLength(SEEDED_SITES.length);

    for (const seeded of SEEDED_SITES) {
      const matches = entities.filter((entity) => entity.metadata.id === seeded.id);
      expect(matches, `expected exactly one entity for ${seeded.id}`).toHaveLength(1);

      const geometry = matches[0].geometry as { point: { lat: number; lng: number } };
      expect(matches[0].metadata.name).toBe(seeded.name);
      const errorMeters = distance(
        { lat: seeded.lat, lng: seeded.lng },
        geometry.point,
        radiusMeters
      );
      expect(errorMeters).toBeLessThanOrEqual(POSITION_TOLERANCE_METERS);
    }
  });

  test('lists every reagent site by name beside the globe (AC4)', async ({ cesium, page }) => {
    await cesium.waitForSettled();
    const list = page.getByTestId('reagent-site-list');

    for (const seeded of SEEDED_SITES) {
      await expect(list).toContainText(seeded.name);
    }
  });

  test('Frostcap Hollow (near the pole) and Tidewrack Shoals (on the antimeridian) stay single entities after flying to each (AC4 edge case)', async ({
    cesium,
  }) => {
    await cesium.waitForSettled();
    const radiusMeters = await cesium.globeRadiusMeters();

    for (const seeded of [FROSTCAP_HOLLOW, TIDEWRACK_SHOALS]) {
      await cesium.flyTo(seeded.lat, seeded.lng, 3_000_000);
      const entities = await cesium.getEntities();
      const matches = entities.filter((entity) => entity.metadata.id === seeded.id);

      expect(matches, `expected exactly one entity for ${seeded.id}`).toHaveLength(1);
      const geometry = matches[0].geometry as { point: { lat: number; lng: number } };
      const errorMeters = distance({ lat: seeded.lat, lng: seeded.lng }, geometry.point, radiusMeters);
      expect(errorMeters).toBeLessThanOrEqual(POSITION_TOLERANCE_METERS);
    }
  });
});

test('reagent map view makes no request off the app\'s own origin or the mock API\'s origin (NFR)', async ({
  page,
}) => {
  const offOrigin = await collectOffOriginRequests(page);

  await page.goto('/reagent-map');
  await page.waitForFunction(() => Boolean((window as any).__SATURDAY_CESIUM__), null, {
    timeout: 15000,
  });
  await page.waitForFunction(
    () => (window as any).__SATURDAY_CESIUM__.viewer.scene.globe.tilesLoaded === true,
    null,
    { timeout: 15000 }
  );

  expect(offOrigin).toEqual([]);
});
