import { test as base } from '@playwright/test';
import type { Page } from '@playwright/test';

export interface CesiumEntityMetadata {
  id: string;
  name: string;
  kind: string;
}

export type CesiumEntityGeometry =
  | { point: { lat: number; lng: number } }
  | { line: Array<{ lat: number; lng: number }> }
  | { circle: { center: { lat: number; lng: number }; radiusMeters: number } };

export interface CesiumEntityResult {
  metadata: CesiumEntityMetadata;
  geometry: CesiumEntityGeometry;
}

export interface CesiumScreenPoint {
  x: number;
  y: number;
}

const DEFAULT_SETTLE_TIMEOUT_MS = 15000;

// Every function below that is passed directly to page.evaluate() is
// serialized via Function.prototype.toString() and re-run inside the page —
// it cannot close over anything from this module's scope. The seam-missing
// message is therefore inlined as a literal in each one rather than shared
// from a module-level constant, and any helper one of them needs (see
// isWithinCanvasBounds below) must be nested inside it, never a sibling.

export interface CesiumFixture {
  /**
   * Waits until the globe's tiles are loaded, the camera has stopped moving,
   * and no render is pending. Rejects with a message naming which condition
   * timed out if the globe never reaches that state within `timeoutMs`.
   */
  waitForSettled(timeoutMs?: number): Promise<void>;

  /**
   * Returns every entity carrying `saturdayMetadata`, with its geometry in
   * degrees and metres. Entities with no `saturdayMetadata` are excluded.
   */
  getEntities(): Promise<CesiumEntityResult[]>;

  /**
   * Projects a geographic point to page coordinates, or returns null when
   * the point is on the far side of the globe or outside the canvas.
   */
  toScreen(lat: number, lng: number): Promise<CesiumScreenPoint | null>;

  /**
   * Drives a real mouse click (via page.mouse) at the screen position for
   * the given geographic point. Throws if the point is not visible.
   */
  clickAt(lat: number, lng: number): Promise<void>;

  /**
   * Flies the camera to the given geographic point and height, then waits
   * for the globe to settle at the destination.
   */
  flyTo(lat: number, lng: number, heightMeters: number): Promise<void>;

  /**
   * Ground resolution, in metres per pixel, at the given geographic point
   * under the current camera position. Used to size the tolerance for
   * clickAt's pick-accuracy assertions.
   */
  metersPerPixelAt(lat: number, lng: number): Promise<number>;
}

function installSettleTracking(): void {
  const seam = (window as any).__SATURDAY_CESIUM__;
  if (!seam) throw new Error('Cesium seam not found: expected window.__SATURDAY_CESIUM__ = { viewer, Cesium }');
  if ((window as any).__SATURDAY_CESIUM_SETTLE__) return;

  const state = { cameraMoving: false, renderPending: true };
  const { viewer } = seam;
  viewer.camera.moveStart.addEventListener(() => {
    state.cameraMoving = true;
  });
  viewer.camera.moveEnd.addEventListener(() => {
    state.cameraMoving = false;
  });
  viewer.scene.postRender.addEventListener(() => {
    state.renderPending = false;
  });
  (window as any).__SATURDAY_CESIUM_SETTLE__ = state;
}

function isSettled(): boolean {
  const seam = (window as any).__SATURDAY_CESIUM__;
  const state = (window as any).__SATURDAY_CESIUM_SETTLE__;
  if (!seam || !state) return false;
  return seam.viewer.scene.globe.tilesLoaded === true && !state.cameraMoving && !state.renderPending;
}

function describeUnsettledConditions(): string {
  const seam = (window as any).__SATURDAY_CESIUM__;
  const state = (window as any).__SATURDAY_CESIUM_SETTLE__;
  if (!seam || !state) return 'seam';

  const pending: string[] = [];
  if (seam.viewer.scene.globe.tilesLoaded !== true) pending.push('tiles');
  if (state.cameraMoving) pending.push('camera');
  if (state.renderPending) pending.push('render');
  return pending.join(', ');
}

async function pollUntilSettled(page: Page, timeoutMs: number): Promise<void> {
  try {
    // The app under test may attach window.__SATURDAY_CESIUM__ asynchronously
    // (as the fixture page itself does, after its imagery provider resolves),
    // so this can't assume the seam already exists just because page.goto()
    // has returned.
    await page.waitForFunction(() => Boolean((window as any).__SATURDAY_CESIUM__), null, {
      timeout: timeoutMs,
    });
  } catch {
    throw new Error(`cesium.waitForSettled() timed out after ${timeoutMs}ms waiting on: seam`);
  }

  await page.evaluate(installSettleTracking);
  try {
    await page.waitForFunction(isSettled, null, { timeout: timeoutMs });
  } catch {
    const pending = await page.evaluate(describeUnsettledConditions);
    throw new Error(`cesium.waitForSettled() timed out after ${timeoutMs}ms waiting on: ${pending}`);
  }
}

// geometryOf is nested inside extractEntities, not a sibling: this whole
// function is serialized via Function.prototype.toString() for
// page.evaluate(), and a sibling function reference would not exist in the
// page at runtime (see the module comment above).
function extractEntities(): CesiumEntityResult[] {
  const seam = (window as any).__SATURDAY_CESIUM__;
  if (!seam) throw new Error('Cesium seam not found: expected window.__SATURDAY_CESIUM__ = { viewer, Cesium }');
  const { viewer, Cesium } = seam;
  const ellipsoid = viewer.scene.globe.ellipsoid;
  const time = viewer.clock.currentTime;
  const toDegrees = (cartesian: any) => {
    const carto = Cesium.Cartographic.fromCartesian(cartesian, ellipsoid);
    return { lat: Cesium.Math.toDegrees(carto.latitude), lng: Cesium.Math.toDegrees(carto.longitude) };
  };
  const geometryOf = (entity: any): CesiumEntityGeometry | null => {
    if (entity.ellipse) {
      const center = entity.position?.getValue(time);
      if (!center) return null;
      const radiusMeters = entity.ellipse.semiMajorAxis.getValue(time);
      return { circle: { center: toDegrees(center), radiusMeters } };
    }
    if (entity.polyline) {
      const positions = entity.polyline.positions.getValue(time) ?? [];
      return { line: positions.map(toDegrees) };
    }
    if (entity.position) {
      const cartesian = entity.position.getValue(time);
      return cartesian ? { point: toDegrees(cartesian) } : null;
    }
    return null;
  };

  const results: CesiumEntityResult[] = [];
  for (const entity of viewer.entities.values) {
    if (!entity.saturdayMetadata) continue;
    const geometry = geometryOf(entity);
    if (!geometry) continue;
    results.push({ metadata: entity.saturdayMetadata, geometry });
  }
  return results;
}

function projectToScreen({ lat, lng }: { lat: number; lng: number }): CesiumScreenPoint | null {
  const seam = (window as any).__SATURDAY_CESIUM__;
  if (!seam) throw new Error('Cesium seam not found: expected window.__SATURDAY_CESIUM__ = { viewer, Cesium }');
  const { viewer, Cesium } = seam;
  const ellipsoid = viewer.scene.globe.ellipsoid;
  const isWithinCanvasBounds = (x: number, y: number): boolean => {
    if (x < 0 || y < 0) return false;
    const canvas = viewer.scene.canvas;
    return x <= canvas.clientWidth && y <= canvas.clientHeight;
  };

  const cartesian = Cesium.Cartesian3.fromRadians(
    Cesium.Math.toRadians(lng),
    Cesium.Math.toRadians(lat),
    0,
    ellipsoid,
  );
  const occluder = new Cesium.EllipsoidalOccluder(ellipsoid, viewer.camera.positionWC);
  if (!occluder.isPointVisible(cartesian)) return null;

  const canvasPosition = viewer.scene.cartesianToCanvasCoordinates(cartesian);
  if (!canvasPosition) return null;

  if (!isWithinCanvasBounds(canvasPosition.x, canvasPosition.y)) return null;
  return { x: canvasPosition.x, y: canvasPosition.y };
}

function pixelSizeAt({ lat, lng }: { lat: number; lng: number }): number {
  const seam = (window as any).__SATURDAY_CESIUM__;
  if (!seam) throw new Error('Cesium seam not found: expected window.__SATURDAY_CESIUM__ = { viewer, Cesium }');
  const { viewer, Cesium } = seam;
  const ellipsoid = viewer.scene.globe.ellipsoid;
  const cartesian = Cesium.Cartesian3.fromRadians(
    Cesium.Math.toRadians(lng),
    Cesium.Math.toRadians(lat),
    0,
    ellipsoid,
  );
  const boundingSphere = new Cesium.BoundingSphere(cartesian, 0);
  const canvas = viewer.scene.canvas;
  return viewer.camera.getPixelSize(boundingSphere, canvas.clientWidth, canvas.clientHeight);
}

function flyCameraTo({
  lat,
  lng,
  heightMeters,
}: {
  lat: number;
  lng: number;
  heightMeters: number;
}): void {
  const seam = (window as any).__SATURDAY_CESIUM__;
  if (!seam) throw new Error('Cesium seam not found: expected window.__SATURDAY_CESIUM__ = { viewer, Cesium }');
  const { viewer, Cesium } = seam;
  const ellipsoid = viewer.scene.globe.ellipsoid;
  const destination = Cesium.Cartesian3.fromRadians(
    Cesium.Math.toRadians(lng),
    Cesium.Math.toRadians(lat),
    heightMeters,
    ellipsoid,
  );
  viewer.camera.flyTo({ destination, duration: 0 });
}

export const test = base.extend<{ cesium: CesiumFixture }>({
  cesium: async ({ page }, use) => {
    const fixture: CesiumFixture = {
      async waitForSettled(timeoutMs = DEFAULT_SETTLE_TIMEOUT_MS) {
        await pollUntilSettled(page, timeoutMs);
      },

      async getEntities() {
        return page.evaluate(extractEntities);
      },

      async toScreen(lat, lng) {
        return page.evaluate(projectToScreen, { lat, lng });
      },

      async clickAt(lat, lng) {
        const point = await fixture.toScreen(lat, lng);
        if (!point) {
          throw new Error(`cesium.clickAt(${lat}, ${lng}): point is not visible on screen`);
        }
        await page.mouse.click(point.x, point.y);
      },

      async flyTo(lat, lng, heightMeters) {
        await page.evaluate(flyCameraTo, { lat, lng, heightMeters });
        await fixture.waitForSettled();
      },

      async metersPerPixelAt(lat, lng) {
        return page.evaluate(pixelSizeAt, { lat, lng });
      },
    };

    await use(fixture);
  },
});
