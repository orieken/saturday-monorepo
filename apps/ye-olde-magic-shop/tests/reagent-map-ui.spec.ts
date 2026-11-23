/**
 * @issue reagent-map-ui-tests
 * E1: the 15 @level-ui scenarios of features/reagent-map/, one test each, titled as in the
 * feature files. Spec package sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413.
 *
 * Gestures go through saturday-playwright-cesium (clickAt, flyTo, toScreen) and the reagent map's
 * own controls; state is read through getEntities and the seam (D-5). Drawing a line is the
 * DrawLineFlow (D-6). Positions are checked within 1 km on Aurelis, using the globe's radius from
 * the seam (D-10, D-14). AC11 reads the saved expedition back through the shop's API.
 *
 * Saving writes apps/mock-api/data/expeditions.json, and AC3 and AC4 count what it holds, so this
 * file runs in one worker under the shared expeditions lock and puts the file back after each test.
 * Runs only under the 'chrome' project: Cesium needs a real webgl2 context.
 */
import fs from 'node:fs';
import reagentSitesFixture from '../../mock-api/data/reagent-sites.json';
import expeditionsFixture from '../../mock-api/data/expeditions.json';
import {
  test,
  distance,
  type CesiumEntityResult,
  type CesiumFixture,
} from '@orieken/saturday-playwright-cesium';
import { expect, type APIRequestContext, type Page } from '@playwright/test';
import { MagicShop } from '../lib/site/MagicShop';
import { ShopApi, type GeoPoint } from '../lib/api/ShopApi';
import {
  EXPEDITIONS_PATH,
  LOCK_HOOK_TIMEOUT_MS,
  acquireExpeditionsLock,
  releaseExpeditionsLock,
} from '../lib/support/expeditions-lock';

type SeededSite = GeoPoint & { id: string; name: string };
type SeededExpedition = { id: string; name: string; owner: string; from: GeoPoint; to: GeoPoint };
type Line = { line: GeoPoint[] };

const SITES: SeededSite[] = reagentSitesFixture;
const EXPEDITIONS: SeededExpedition[] = expeditionsFixture;
const site = (id: string): SeededSite => SITES.find((candidate) => candidate.id === id)!;
const MOONWELL_SPRING = site('moonwell-spring');
const EMBER_CRAG = site('ember-crag');
const FROSTCAP_HOLLOW = site('frostcap-hollow');
const TIDEWRACK_SHOALS = site('tidewrack-shoals');
const THORNVEIL_GLADE = site('thornveil-glade');
const GLIMMERDEEP_CAVES = site('glimmerdeep-caves');
const SEEDED = EXPEDITIONS.filter((expedition) => expedition.owner === 'shop');
const ELARAS = EXPEDITIONS.find((expedition) => expedition.name === 'Moonwell Spring to Frostcap Hollow')!;
const BORINS = EXPEDITIONS.find((expedition) => expedition.name === 'Ember Crag to Saltmarsh Cairn')!;

const ELARA = 'elara@aurelis.example';
const TOLERANCE_METERS = 1_000;
const FLY_HEIGHT_METERS = 3_000_000;
const FIRST_POINT_STATUS = 'Line tool active. Click the globe to place the first point.';
const SIGN_IN_MESSAGE = 'Sign in to save this expedition.';
const ZERO_LENGTH_MESSAGE = 'Zero-length line: from and to are the same place';
const SAVED_STATUS = 'Expedition saved.';

const ofKind = (entities: CesiumEntityResult[], kind: string) =>
  entities.filter((entity) => entity.metadata.kind === kind);

function joins(entity: CesiumEntityResult, a: GeoPoint, b: GeoPoint, radius: number): boolean {
  const [first, second] = (entity.geometry as Line).line;
  const near = (p: GeoPoint, q: GeoPoint) => distance(p, q, radius) <= TOLERANCE_METERS;
  return (near(first, a) && near(second, b)) || (near(first, b) && near(second, a));
}

async function expectDrawnAt(cesium: CesiumFixture, expected: SeededExpedition[]): Promise<void> {
  const radius = await cesium.globeRadiusMeters();
  const drawn = ofKind(await cesium.getEntities(), 'expedition');
  expect(drawn.map((entity) => entity.metadata.name).sort()).toEqual(expected.map((e) => e.name).sort());
  for (const expedition of expected) {
    const entity = drawn.find((candidate) => candidate.metadata.id === expedition.id)!;
    expect(joins(entity, expedition.from, expedition.to, radius), `${expedition.name} ends`).toBe(true);
  }
}

async function signInBrowser(page: Page, request: APIRequestContext, email: string): Promise<void> {
  const session = await new ShopApi(request).signIn(email);
  await page.addInitScript(({ token, user }) => {
    localStorage.setItem('dnd_auth_token', token);
    localStorage.setItem('dnd_user', JSON.stringify(user));
  }, session);
}

const storedUnchanged = (seed: string) => expect(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')).toBe(seed);

// The screen point one pixel past the globe's edge, along the canvas's middle row, found with
// the camera's own pickEllipsoid: the edge, not open space beside the globe (that is AC9).
async function pastTheEdge(page: Page): Promise<{ x: number; y: number }> {
  return page.evaluate(() => {
    const { viewer, Cesium } = (window as any).__SATURDAY_CESIUM__;
    const canvas = viewer.scene.canvas;
    const rect = canvas.getBoundingClientRect();
    const [cx, cy] = [canvas.clientWidth / 2, canvas.clientHeight / 2];
    const hits = (dx: number) =>
      Boolean(viewer.camera.pickEllipsoid(new Cesium.Cartesian2(cx + dx, cy), viewer.scene.globe.ellipsoid));
    let [lo, hi] = [0, canvas.clientWidth];
    for (let i = 0; i < 40; i++) {
      const mid = (lo + hi) / 2;
      if (hits(mid)) lo = mid;
      else hi = mid;
    }
    return { x: rect.left + cx + lo + 1, y: rect.top + cy };
  });
}

// 'default', not 'serial': the tests still run in order in one worker, so the lock covers the
// whole file, but one failure does not skip the rest. Each scenario reports for itself.
test.describe.configure({ mode: 'default' });

test.describe('Reagent map UI scenarios (E1)', () => {
  let seed: string;
  let shop: MagicShop;

  test.beforeAll(async () => {
    test.setTimeout(LOCK_HOOK_TIMEOUT_MS);
    await acquireExpeditionsLock();
    seed = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');
  });

  test.afterEach(() => {
    fs.writeFileSync(EXPEDITIONS_PATH, seed, 'utf8');
  });

  test.afterAll(() => {
    releaseExpeditionsLock();
  });

  test.beforeEach(({ page }) => {
    shop = new MagicShop(page);
  });

  async function openMap(cesium: CesiumFixture): Promise<void> {
    await shop.reagentMapPage.open();
    await cesium.waitForSettled();
  }

  test('@AC1 Every reagent site is shown at its position', async ({ cesium }) => {
    await openMap(cesium);
    const radius = await cesium.globeRadiusMeters();
    const shown = ofKind(await cesium.getEntities(), 'site');

    for (const seeded of SITES) {
      await expect(shop.reagentMapPage.siteList).toContainText(seeded.name);
      const entity = shown.find((candidate) => candidate.metadata.id === seeded.id)!;
      const at = (entity.geometry as { point: GeoPoint }).point;
      expect(distance(at, seeded, radius), `${seeded.name} position`).toBeLessThanOrEqual(TOLERANCE_METERS);
    }
  });

  for (const seeded of [FROSTCAP_HOLLOW, TIDEWRACK_SHOALS]) {
    test(`@AC2 Sites near a pole and on the 180th meridian are shown once, where they are: ${seeded.name}`, async ({
      cesium,
    }) => {
      await openMap(cesium);
      await cesium.flyTo(seeded.lat, seeded.lng, FLY_HEIGHT_METERS);
      const radius = await cesium.globeRadiusMeters();

      const matches = ofKind(await cesium.getEntities(), 'site').filter((e) => e.metadata.id === seeded.id);
      expect(matches).toHaveLength(1);
      const at = (matches[0].geometry as { point: GeoPoint }).point;
      expect(distance(at, seeded, radius)).toBeLessThanOrEqual(TOLERANCE_METERS);
      expect(await cesium.toScreen(seeded.lat, seeded.lng), `${seeded.name} on screen`).not.toBeNull();
    });
  }

  test('@AC3 Seeded expeditions are drawn where they were seeded', async ({ cesium }) => {
    await openMap(cesium);
    await expectDrawnAt(cesium, SEEDED);
  });

  test("@AC4 A signed-in adventurer's own expedition is drawn with the seeded ones", async ({
    page,
    cesium,
    request,
  }) => {
    await signInBrowser(page, request, ELARA);
    await openMap(cesium);

    await expectDrawnAt(cesium, [...SEEDED, ELARAS]);
    const names = ofKind(await cesium.getEntities(), 'expedition').map((entity) => entity.metadata.name);
    expect(names).not.toContain(BORINS.name);
  });

  test('@AC5 Flying to a site on the far side brings it into view', async ({ cesium }) => {
    await openMap(cesium);
    expect(await cesium.toScreen(TIDEWRACK_SHOALS.lat, TIDEWRACK_SHOALS.lng), 'starts on the far side').toBeNull();

    await shop.reagentMapPage.flyTo(TIDEWRACK_SHOALS.id);
    await expect
      .poll(() => cesium.toScreen(TIDEWRACK_SHOALS.lat, TIDEWRACK_SHOALS.lng), { timeout: 15_000 })
      .not.toBeNull();
  });

  test('@AC6 Two points on the land make a line that can be saved', async ({ cesium }) => {
    await openMap(cesium);
    await shop.drawLine.execute({ cesium, from: MOONWELL_SPRING, to: EMBER_CRAG });
    const radius = await cesium.globeRadiusMeters();

    const lines = ofKind(await cesium.getEntities(), 'draft-line');
    expect(lines).toHaveLength(1);
    expect(joins(lines[0], MOONWELL_SPRING, EMBER_CRAG, radius)).toBe(true);
    await expect(shop.reagentMapPage.saveButton).toBeEnabled();
  });

  test('@AC7 One point is not a line', async ({ cesium }) => {
    await openMap(cesium);
    await shop.reagentMapPage.startLine();
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);

    expect(ofKind(await cesium.getEntities(), 'draft-point')).toHaveLength(1);
    expect(ofKind(await cesium.getEntities(), 'draft-line')).toHaveLength(0);
    await expect(shop.reagentMapPage.saveButton).toBeDisabled();
  });

  test('@AC8 A line follows the shortest path over the surface', async ({ page, cesium }) => {
    await openMap(cesium);
    await shop.reagentMapPage.startLine();
    // The two sites are never on screen together, so the camera goes to each before its click.
    for (const end of [THORNVEIL_GLADE, GLIMMERDEEP_CAVES]) {
      await cesium.flyTo(end.lat, end.lng, FLY_HEIGHT_METERS);
      await cesium.clickAt(end.lat, end.lng);
    }
    const radius = await cesium.globeRadiusMeters();

    const lines = ofKind(await cesium.getEntities(), 'draft-line');
    expect(lines).toHaveLength(1);
    expect(joins(lines[0], THORNVEIL_GLADE, GLIMMERDEEP_CAVES, radius)).toBe(true);
    const arc = await page.evaluate(() => {
      const { viewer, Cesium } = (window as any).__SATURDAY_CESIUM__;
      const line = [...viewer.entities.values].find((e: any) => e.saturdayMetadata?.kind === 'draft-line');
      return { actual: line.polyline.arcType.getValue(), geodesic: Cesium.ArcType.GEODESIC };
    });
    expect(arc.actual).toBe(arc.geodesic);
  });

  test('@AC9 A point on space beside the globe is not placed', async ({ page, cesium }) => {
    await openMap(cesium);
    await shop.reagentMapPage.startLine();
    const box = (await page.locator('#cesiumContainer').boundingBox())!;
    // The globe's disc is far from the canvas corners at the default view.
    await page.mouse.click(box.x + 10, box.y + 10);

    expect(ofKind(await cesium.getEntities(), 'draft-point')).toHaveLength(0);
    await expect(shop.reagentMapPage.lineToolStatus).toHaveText(FIRST_POINT_STATUS);
  });

  test('@AC10 A point on the edge of the globe, where the land faces away, is not placed', async ({
    page,
    cesium,
  }) => {
    await openMap(cesium);
    await shop.reagentMapPage.startLine();
    const edge = await pastTheEdge(page);
    await page.mouse.click(edge.x, edge.y);

    expect(ofKind(await cesium.getEntities(), 'draft-point')).toHaveLength(0);
  });

  test('@AC11 A saved expedition is shown again and held by the shop where it was drawn', async ({
    page,
    cesium,
    request,
  }) => {
    await signInBrowser(page, request, ELARA);
    await openMap(cesium);
    await shop.drawLine.execute({ cesium, from: MOONWELL_SPRING, to: EMBER_CRAG, save: true });
    await expect(shop.reagentMapPage.saveStatus).toHaveText(SAVED_STATUS);

    await openMap(cesium);
    const radius = await cesium.globeRadiusMeters();
    const known = new Set(EXPEDITIONS.map((expedition) => expedition.id));
    const fresh = ofKind(await cesium.getEntities(), 'expedition').filter((e) => !known.has(e.metadata.id));
    expect(fresh, 'her new expedition is shown').toHaveLength(1);
    expect(joins(fresh[0], MOONWELL_SPRING, EMBER_CRAG, radius)).toBe(true);

    const elara = new ShopApi(request);
    await elara.signIn(ELARA);
    const held = (await elara.visibleExpeditions()).find((e) => e.id === fresh[0].metadata.id);
    expect(held, 'the shop holds it for Elara').toBeDefined();
    expect(distance(held!.from, MOONWELL_SPRING, radius)).toBeLessThanOrEqual(TOLERANCE_METERS);
    expect(distance(held!.to, EMBER_CRAG, radius)).toBeLessThanOrEqual(TOLERANCE_METERS);
  });

  test('@AC12 A signed-out adventurer is asked to sign in when saving', async ({ cesium }) => {
    await openMap(cesium);
    await shop.drawLine.execute({ cesium, from: MOONWELL_SPRING, to: EMBER_CRAG, save: true });

    await expect(shop.reagentMapPage.saveMessage).toHaveText(SIGN_IN_MESSAGE);
    storedUnchanged(seed);
  });

  test('@AC13 Placing both points at the same place cannot be saved', async ({ page, cesium, request }) => {
    await signInBrowser(page, request, ELARA);
    await openMap(cesium);
    // Two clicks on the same screen point pick the same place, at any zoom.
    await shop.drawLine.execute({ cesium, from: MOONWELL_SPRING, to: MOONWELL_SPRING, save: true });

    await expect(shop.reagentMapPage.saveMessage).toHaveText(ZERO_LENGTH_MESSAGE);
    storedUnchanged(seed);
  });

  test('@AC14 Starting a new line abandons the unfinished one', async ({ page, cesium, request }) => {
    await signInBrowser(page, request, ELARA);
    await openMap(cesium);
    await shop.reagentMapPage.startLine();
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    expect(ofKind(await cesium.getEntities(), 'draft-point')).toHaveLength(1);

    await shop.reagentMapPage.startLine();

    expect(ofKind(await cesium.getEntities(), 'draft-point')).toHaveLength(0);
    expect(ofKind(await cesium.getEntities(), 'draft-line')).toHaveLength(0);
    storedUnchanged(seed);
  });

  test('@AC15 An unsaved line is not kept after leaving the map', async ({ page, cesium, request }) => {
    await signInBrowser(page, request, ELARA);
    await openMap(cesium);
    await shop.drawLine.execute({ cesium, from: MOONWELL_SPRING, to: EMBER_CRAG });
    expect(ofKind(await cesium.getEntities(), 'draft-line')).toHaveLength(1);

    // Leave through the app's own nav, so the view unmounts without a page reload.
    await page.getByTestId('shop-link').click();
    await page.getByTestId('reagent-map-link').click();
    await cesium.waitForSettled();

    expect(ofKind(await cesium.getEntities(), 'draft-point')).toHaveLength(0);
    expect(ofKind(await cesium.getEntities(), 'draft-line')).toHaveLength(0);
    storedUnchanged(seed);
  });
});
