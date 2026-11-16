/**
 * @issue reagent-map-line-tool
 * Covers draw-line.feature's acceptance scenarios: placing two points completes
 * a geodesic line ready to save (AC6), one point alone shows no completed line
 * (AC7), the completed line follows the shortest path over the surface rather
 * than a flat-map straight line (AC8), a point can be placed anywhere on the
 * visible surface including open sea, a click beside the globe places nothing
 * (AC9), a click on the globe's far edge where the surface faces away places
 * nothing (AC10), and starting the tool again discards any point already
 * placed (AC14). Every placed point and completed line carries saturdayMetadata
 * kind 'draft-point'/'draft-line', read through getEntities().
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 */
import reagentSitesFixture from '../../mock-api/data/reagent-sites.json';
import { test, type CesiumEntityResult } from '@orieken/saturday-playwright-cesium';
import { expect, type Page } from '@playwright/test';

type SeededSite = { id: string; name: string; lat: number; lng: number };

const SEEDED_SITES: SeededSite[] = reagentSitesFixture;
const MOONWELL_SPRING = SEEDED_SITES.find((site) => site.id === 'moonwell-spring')!;
const EMBER_CRAG = SEEDED_SITES.find((site) => site.id === 'ember-crag')!;
const THORNVEIL_GLADE = SEEDED_SITES.find((site) => site.id === 'thornveil-glade')!;
const GLIMMERDEEP_CAVES = SEEDED_SITES.find((site) => site.id === 'glimmerdeep-caves')!;
// Confirmed visible from the default camera view and over open sea in
// aurelis.jpg (blue channel clearly exceeds green at this pixel) — verified
// against the running app before writing this test, since the spec only
// requires "the land" to mean the planet's surface for now, sea included.
const OPEN_SEA_POINT = { lat: 0, lng: 0 };
// Matches ReagentMapView's own FLY_TO_HEIGHT_METERS — Thornveil Glade and
// Glimmerdeep Caves are not simultaneously visible from the default view (both
// sit beyond the visible hemisphere's ~70.5-degree horizon from Moonwell
// Spring), so each point is placed after flying the camera overhead of it.
const FLY_TO_HEIGHT_METERS = 3_000_000;

function draftPointsOf(entities: CesiumEntityResult[]): CesiumEntityResult[] {
  return entities.filter((entity) => entity.metadata.kind === 'draft-point');
}

function draftLinesOf(entities: CesiumEntityResult[]): CesiumEntityResult[] {
  return entities.filter((entity) => entity.metadata.kind === 'draft-line');
}

async function startLineTool(page: Page): Promise<void> {
  await page.getByTestId('start-line-tool').click();
}

// Binary-searches, with the current camera's own pickEllipsoid, for the
// screen radius (measured from the globe's on-screen center) beyond which a
// click no longer lands on the globe. Mirrors pickLandPoint's own check, so
// the found point sits right at the globe's rendered edge, not somewhere in
// open space beside it (that is the separate AC9 case below). Returns a
// page-absolute point, one pixel past the boundary, from the canvas element's
// own bounding rect. Not #cesiumContainer's: it has a 2px border on each side
// that the inner canvas doesn't, which was enough to miss the edge.
async function findLimbEdgeClickPoint(page: Page): Promise<{ x: number; y: number }> {
  return page.evaluate(() => {
    const seam = (window as any).__SATURDAY_CESIUM__;
    const { viewer, Cesium } = seam;
    const ellipsoid = viewer.scene.globe.ellipsoid;
    const canvas = viewer.scene.canvas;
    const rect = canvas.getBoundingClientRect();
    const centerX = canvas.clientWidth / 2;
    const centerY = canvas.clientHeight / 2;

    function isOnVisibleHemisphere(offset: number): boolean {
      const windowPosition = new Cesium.Cartesian2(centerX + offset, centerY);
      return Boolean(viewer.camera.pickEllipsoid(windowPosition, ellipsoid));
    }

    let lo = 0;
    let hi = canvas.clientWidth;
    for (let i = 0; i < 40; i++) {
      const mid = (lo + hi) / 2;
      if (isOnVisibleHemisphere(mid)) lo = mid;
      else hi = mid;
    }

    return { x: rect.left + centerX + lo + 1, y: rect.top + centerY };
  });
}

test.describe('Reagent map line tool (AC6, AC7, AC9, AC14)', () => {
  test.beforeEach(async ({ page, cesium }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();
  });

  test('placing two points on the land completes a line ready to save (AC6)', async ({ page, cesium }) => {
    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(EMBER_CRAG.lat, EMBER_CRAG.lng);

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(2);
    expect(draftLinesOf(entities)).toHaveLength(1);
  });

  test('placing only one point shows no completed line, and there is nothing to save (AC7)', async ({
    page,
    cesium,
  }) => {
    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(1);
    expect(draftLinesOf(entities)).toHaveLength(0);
  });

  test('every placed point and the completed line are marked as drafts, not a saved expedition', async ({
    page,
    cesium,
  }) => {
    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(EMBER_CRAG.lat, EMBER_CRAG.lng);

    const entities = await cesium.getEntities();
    const drafts = [...draftPointsOf(entities), ...draftLinesOf(entities)];
    expect(drafts).toHaveLength(3);
    for (const draft of drafts) {
      expect(['draft-point', 'draft-line']).toContain(draft.metadata.kind);
      // Draft entities are never persisted, unlike 'site'/'expedition' entities
      // which carry a real id and name.
      expect(draft.metadata.id).toBe('');
      expect(draft.metadata.name).toBe('');
    }
  });

  test('a point can be placed anywhere on the visible surface, sea included', async ({ page, cesium }) => {
    await startLineTool(page);
    await cesium.clickAt(OPEN_SEA_POINT.lat, OPEN_SEA_POINT.lng);

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(1);
  });

  test('clicking beside the globe places nothing and the tool still waits for its first point (AC9)', async ({
    page,
    cesium,
  }) => {
    await startLineTool(page);
    const canvasBox = await page.locator('#cesiumContainer').boundingBox();
    expect(canvasBox).not.toBeNull();

    // The default view's globe disc sits well clear of the canvas corners
    // (confirmed against the running app: the rendered disc's on-screen
    // radius is under 300px against a canvas several hundred pixels wider and
    // taller) — a corner pixel is unambiguously open space beside the globe.
    await page.mouse.click(canvasBox!.x + 10, canvasBox!.y + 10);

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(0);
    expect(draftLinesOf(entities)).toHaveLength(0);
    await expect(page.getByTestId('line-tool-status')).toHaveText(
      'Line tool active. Click the globe to place the first point.'
    );
  });

  test("clicking the globe's far edge, where the surface faces away, places nothing (AC10)", async ({
    page,
    cesium,
  }) => {
    await startLineTool(page);

    // One pixel past the globe's rendered edge, found with the camera's own
    // pickEllipsoid: the globe's edge, not open space (AC9's case, tested
    // above at a canvas corner far from the disc). pickEllipsoid returns only
    // the side of the surface facing the camera, so a point there never
    // reaches the far side.
    const { x: edgeX, y: edgeY } = await findLimbEdgeClickPoint(page);
    await page.mouse.click(edgeX, edgeY);

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(0);
  });

  test('starting the line tool again discards any point already placed (AC14)', async ({ page, cesium }) => {
    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    expect(draftPointsOf(await cesium.getEntities())).toHaveLength(1);

    await startLineTool(page);

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(0);
    expect(draftLinesOf(entities)).toHaveLength(0);
  });
});

test.describe("The completed line follows the shortest path over Aurelis's surface (AC8)", () => {
  test.beforeEach(async ({ page, cesium }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();
  });

  test('a line from Thornveil Glade to Glimmerdeep Caves is drawn as a geodesic arc, not a flat straight line', async ({
    page,
    cesium,
  }) => {
    await startLineTool(page);
    await cesium.flyTo(THORNVEIL_GLADE.lat, THORNVEIL_GLADE.lng, FLY_TO_HEIGHT_METERS);
    await cesium.clickAt(THORNVEIL_GLADE.lat, THORNVEIL_GLADE.lng);
    await cesium.flyTo(GLIMMERDEEP_CAVES.lat, GLIMMERDEEP_CAVES.lng, FLY_TO_HEIGHT_METERS);
    await cesium.clickAt(GLIMMERDEEP_CAVES.lat, GLIMMERDEEP_CAVES.lng);

    const entities = await cesium.getEntities();
    const lines = draftLinesOf(entities);
    expect(lines).toHaveLength(1);
    const geometry = lines[0].geometry as { line: Array<{ lat: number; lng: number }> };
    expect(geometry.line).toHaveLength(2);

    const arcType = await page.evaluate(() => {
      const { viewer, Cesium } = (window as any).__SATURDAY_CESIUM__;
      const entity = [...viewer.entities.values].find(
        (candidate: any) => candidate.saturdayMetadata?.kind === 'draft-line'
      );
      return { actual: entity.polyline.arcType.getValue(), geodesic: Cesium.ArcType.GEODESIC };
    });
    expect(arcType.actual).toBe(arcType.geodesic);
  });
});
