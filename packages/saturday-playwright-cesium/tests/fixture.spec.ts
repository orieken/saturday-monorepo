import { expect, type Page } from '@playwright/test';
import { test } from '../src/fixture';
import type { CesiumFixture } from '../src/fixture';
import { distance } from '../src/geometry';

// Every assertion below targets a literal, known value baked into
// fixtures/assets/fixture-globe.js — change that value and the matching test
// here must fail (shared/rules/test-repair-contract.md's "prove it can fail").
// Camera, canvas size (1024x512) and the geographic test points were chosen
// empirically against the real fixture globe (see implementation-notes.md).

const AURELIS_RADIUS_METERS = 6_000_000;

async function pickedLocationAfterClick(
  page: Page,
  cesium: CesiumFixture,
  lat: number,
  lng: number,
): Promise<{ lat: number; lng: number } | null> {
  await page.evaluate(() => {
    const { viewer, Cesium } = (window as any).__SATURDAY_CESIUM__;
    (window as any).__LAST_PICK__ = null;
    viewer.screenSpaceEventHandler.setInputAction((event: any) => {
      const ellipsoid = viewer.scene.globe.ellipsoid;
      const cartesian = viewer.camera.pickEllipsoid(event.position, ellipsoid);
      if (!cartesian) return;
      const carto = Cesium.Cartographic.fromCartesian(cartesian, ellipsoid);
      (window as any).__LAST_PICK__ = {
        lat: Cesium.Math.toDegrees(carto.latitude),
        lng: Cesium.Math.toDegrees(carto.longitude),
      };
    }, Cesium.ScreenSpaceEventType.LEFT_CLICK);
  });
  await cesium.clickAt(lat, lng);
  return page.evaluate(() => (window as any).__LAST_PICK__);
}

test.describe('CesiumFixture against the local fixture globe', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('waitForSettled resolves once the fixture globe’s tiles are loaded', async ({
    cesium,
    page,
  }) => {
    await expect(cesium.waitForSettled()).resolves.toBeUndefined();
    const tilesLoaded = await page.evaluate(
      () => (window as any).__SATURDAY_CESIUM__.viewer.scene.globe.tilesLoaded,
    );
    expect(tilesLoaded).toBe(true);
  });

  test('waitForSettled rejects with a message naming the stuck condition after a timeout', async ({
    cesium,
    page,
  }) => {
    await cesium.waitForSettled();
    await page.evaluate(() => {
      (window as any).__SATURDAY_CESIUM_SETTLE__.cameraMoving = true;
    });
    await expect(cesium.waitForSettled(200)).rejects.toThrow(/camera/);
    await page.evaluate(() => {
      (window as any).__SATURDAY_CESIUM_SETTLE__.cameraMoving = false;
    });
  });

  test('getEntities returns metadata-tagged entities with geometry in degrees and metres, excluding the untagged one', async ({
    cesium,
  }) => {
    await cesium.waitForSettled();
    const entities = await cesium.getEntities();
    // Five entities exist on the fixture globe; only four carry saturdayMetadata.
    expect(entities).toHaveLength(4);

    const byId = Object.fromEntries(entities.map((entity) => [entity.metadata.id, entity]));

    expect(byId['reagent-site-1'].metadata).toEqual({
      id: 'reagent-site-1',
      name: 'Reagent Site 1',
      kind: 'point',
    });
    const point = byId['reagent-site-1'].geometry as { point: { lat: number; lng: number } };
    expect(point.point.lat).toBeCloseTo(5, 4);
    expect(point.point.lng).toBeCloseTo(10, 4);

    const lineGeometry = byId['supply-route-1'].geometry as {
      line: Array<{ lat: number; lng: number }>;
    };
    expect(lineGeometry.line).toHaveLength(2);
    expect(lineGeometry.line[0].lat).toBeCloseTo(5, 4);
    expect(lineGeometry.line[0].lng).toBeCloseTo(10, 4);
    expect(lineGeometry.line[1].lat).toBeCloseTo(-10, 4);
    expect(lineGeometry.line[1].lng).toBeCloseTo(30, 4);

    const circleGeometry = byId['exclusion-zone-1'].geometry as {
      circle: { center: { lat: number; lng: number }; radiusMeters: number };
    };
    expect(circleGeometry.circle.center.lat).toBeCloseTo(0, 4);
    expect(circleGeometry.circle.center.lng).toBeCloseTo(0, 4);
    expect(circleGeometry.circle.radiusMeters).toBeCloseTo(400_000, 0);

    // mystery-site-1 carries a deliberately unexpected `kind` string —
    // getEntities() must still extract its geometry from Cesium's own
    // entity shape, never by branching on this value (architecture-notes.md).
    expect(byId['mystery-site-1'].metadata.kind).toBe('unexpected-kind-value');
    const mysteryPoint = byId['mystery-site-1'].geometry as { point: { lat: number; lng: number } };
    expect(mysteryPoint.point.lat).toBeCloseTo(-5, 4);
    expect(mysteryPoint.point.lng).toBeCloseTo(-5, 4);
  });

  test('toScreen returns page coordinates for a point visible on the near side of the globe', async ({
    cesium,
  }) => {
    await cesium.waitForSettled();
    const point = await cesium.toScreen(5, 10);
    expect(point).not.toBeNull();
    expect(point!.x).toBeGreaterThanOrEqual(0);
    expect(point!.x).toBeLessThanOrEqual(1024);
    expect(point!.y).toBeGreaterThanOrEqual(0);
    expect(point!.y).toBeLessThanOrEqual(512);
  });

  test('toScreen returns null for a point on the far side of the globe', async ({ cesium }) => {
    await cesium.waitForSettled();
    expect(await cesium.toScreen(0, 170)).toBeNull();
  });

  test('toScreen returns null for a point that projects outside the visible canvas', async ({
    cesium,
  }) => {
    await cesium.waitForSettled();
    // Not occluded by the globe (still on the near side) but far enough from
    // the camera's boresight that it falls outside the 1024x512 canvas.
    expect(await cesium.toScreen(20, 0)).toBeNull();
  });

  test('clickAt drives a real mouse click that picks the requested geographic location', async ({
    cesium,
    page,
  }) => {
    await cesium.waitForSettled();
    const lat = 5;
    const lng = 10;

    const picked = await pickedLocationAfterClick(page, cesium, lat, lng);
    expect(picked).not.toBeNull();

    const toleranceMeters = (await cesium.metersPerPixelAt(lat, lng)) * 2;
    const errorMeters = distance({ lat, lng }, picked!, AURELIS_RADIUS_METERS);
    expect(errorMeters).toBeLessThanOrEqual(toleranceMeters);
  });

  test('clickAt near the horizon still picks within the (wider) tolerance metersPerPixelAt allows', async ({
    cesium,
    page,
  }) => {
    await cesium.waitForSettled();
    const lat = 7;
    const lng = 0;

    const picked = await pickedLocationAfterClick(page, cesium, lat, lng);
    expect(picked).not.toBeNull();

    const toleranceMeters = (await cesium.metersPerPixelAt(lat, lng)) * 2;
    const errorMeters = distance({ lat, lng }, picked!, AURELIS_RADIUS_METERS);
    expect(errorMeters).toBeLessThanOrEqual(toleranceMeters);
  });

  test('flyTo moves the camera to the requested location and settles there', async ({
    cesium,
    page,
  }) => {
    await cesium.waitForSettled();
    await cesium.flyTo(40, 50, 4_000_000);

    const cameraLocation = await page.evaluate(() => {
      const { viewer, Cesium } = (window as any).__SATURDAY_CESIUM__;
      const ellipsoid = viewer.scene.globe.ellipsoid;
      const carto = Cesium.Cartographic.fromCartesian(viewer.camera.positionWC, ellipsoid);
      return {
        lat: Cesium.Math.toDegrees(carto.latitude),
        lng: Cesium.Math.toDegrees(carto.longitude),
        height: carto.height,
      };
    });

    expect(cameraLocation.lat).toBeCloseTo(40, 1);
    expect(cameraLocation.lng).toBeCloseTo(50, 1);
    expect(Math.abs(cameraLocation.height - 4_000_000)).toBeLessThan(50_000);
  });
});

test('fixture page makes no request off its own origin', async ({ page }) => {
  const offOrigin: string[] = [];
  page.on('request', (request) => {
    const url = new URL(request.url());
    // blob:/data: are in-memory, non-network resources (Cesium uses blob
    // URLs to bootstrap its worker pool) — only http(s) requests can
    // actually leave the page's own origin.
    if (url.protocol !== 'http:' && url.protocol !== 'https:') return;
    if (url.hostname !== '127.0.0.1') offOrigin.push(request.url());
  });

  await page.goto('/');
  await page.waitForFunction(() => Boolean((window as any).__SATURDAY_CESIUM__), null, {
    timeout: 15000,
  });
  await page.waitForFunction(
    () => (window as any).__SATURDAY_CESIUM__.viewer.scene.globe.tilesLoaded === true,
    null,
    { timeout: 15000 },
  );

  expect(offOrigin).toEqual([]);
});
