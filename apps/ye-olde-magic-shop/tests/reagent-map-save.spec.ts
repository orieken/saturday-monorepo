/**
 * @issue reagent-map-save
 * Covers reagent-map-spec AC11/AC12/AC13/AC15 at the UI level: saving a
 * completed line through the save-line-tool button draws it again after a
 * reload and the shop holds it within 1 km (AC11); saving while signed out
 * asks the adventurer to sign in and sends no POST (AC12); saving a line
 * whose two points are the same place shows the shop's own zero-length
 * message, never a client-recomputed one (AC13); an unsaved line leaves no
 * draft behind after the adventurer leaves the map and reopens it (AC15).
 * Also proves the double-click guard noted in implementation-notes.md sends
 * exactly one POST per completed line.
 * Shares apps/mock-api/data/expeditions.json with reagent-map-api.spec.ts and
 * reagent-map-expeditions.spec.ts — guarded by the same exclusive lock (same LOCK_PATH, own
 * local copy of acquire/release in each file; see the comment above them below and in
 * reagent-map-api.spec.ts for why these are duplicated rather than imported) so none of the
 * three files observes the file mid-write by another (architecture-notes.md).
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 */
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import reagentSitesFixture from '../../mock-api/data/reagent-sites.json';
import { test, distance, type CesiumEntityResult } from '@orieken/saturday-playwright-cesium';
import { expect, type APIRequestContext, type Page } from '@playwright/test';

// Keyed to the same path as reagent-map-api.spec.ts's and reagent-map-expeditions.spec.ts's
// own copies of this lock (architecture-notes.md). Not imported from either file: importing a
// sibling *.spec.ts file executes its top-level test.describe(...)/configure(...) calls as a
// side effect (Playwright's "currently loading file suite" is a single value a plain import
// does not change), re-registering that file's tests inside this one and colliding on a
// second describe.configure({mode: 'serial'}) call — confirmed by running the cross-import
// version architecture-notes.md originally specified, which failed with "serial" mode is
// already assigned for the enclosing scope. Each file keeps its own copy instead; mutual
// exclusion still comes from all three contending for the same filesystem path.
const LOCK_PATH = path.join(os.tmpdir(), 'reagent-map-expeditions.lock');
// See the matching constant's comment in reagent-map-api.spec.ts: measured locally at
// --workers=3 --repeat-each=3, the lock's hold time (one per file's whole serial describe
// block, not one test — this file's own block runs ~17-18s) summed to ~69s of exclusive
// contention across the 9 total blocks (3 files x 3 repeats); 30_000 measurably fired on a
// legitimate holder and caused a lock-theft double-release. 180_000 is comfortably above that
// measured worst case.
const STALE_LOCK_MS = 180_000;
const HOOK_TIMEOUT_MS = 200_000;

async function acquireExpeditionsLock(): Promise<void> {
  await expect(() => {
    try {
      fs.closeSync(fs.openSync(LOCK_PATH, 'wx'));
    } catch (err) {
      if (
        (err as NodeJS.ErrnoException).code === 'EEXIST' &&
        Date.now() - fs.statSync(LOCK_PATH).mtimeMs > STALE_LOCK_MS
      ) {
        fs.unlinkSync(LOCK_PATH);
      }
      throw err;
    }
  }).toPass({ timeout: STALE_LOCK_MS });
}

function releaseExpeditionsLock(): void {
  fs.unlinkSync(LOCK_PATH);
}

// Same reason as the lock above: a local copy, not an import from
// reagent-map-expeditions.spec.ts, which also defines this. Unlike that file's version, this
// one returns the user object too — saveLine()'s authStore.isAuthenticated check (unlike
// fetchExpeditions()'s getAuthHeaders(), which reads the token straight out of localStorage)
// is a Pinia getter that also needs authStore.user, hydrated from the 'dnd_user' key by
// App.vue's authStore.hydrate() — the token alone leaves isAuthenticated false.
async function signInAs(
  request: APIRequestContext,
  email: string
): Promise<{ token: string; user: unknown }> {
  const response = await request.post(`${MOCK_API_URL}/api/auth/login`, {
    data: { email, password: 'password123' },
  });
  return response.json();
}

type SeededSite = { id: string; name: string; lat: number; lng: number };
type StoredExpedition = { id: string; from: { lat: number; lng: number }; to: { lat: number; lng: number } };

const SEEDED_SITES: SeededSite[] = reagentSitesFixture;
const MOONWELL_SPRING = SEEDED_SITES.find((site) => site.id === 'moonwell-spring')!;
const EMBER_CRAG = SEEDED_SITES.find((site) => site.id === 'ember-crag')!;

const EXPEDITIONS_PATH = path.resolve(__dirname, '../../mock-api/data/expeditions.json');
const EXPEDITIONS_API_PATTERN = '**/api/expeditions';
const MOCK_API_URL = process.env.MOCK_API_URL || 'http://localhost:8001';
const ELARA_EMAIL = 'elara@aurelis.example';
const POSITION_TOLERANCE_METERS = 1_000;
const ZERO_LENGTH_MESSAGE = 'Zero-length line: from and to are the same place';
const SIGN_IN_MESSAGE = 'Sign in to save this expedition.';
const SAVED_MESSAGE = 'Expedition saved.';

function draftPointsOf(entities: CesiumEntityResult[]): CesiumEntityResult[] {
  return entities.filter((entity) => entity.metadata.kind === 'draft-point');
}

function draftLinesOf(entities: CesiumEntityResult[]): CesiumEntityResult[] {
  return entities.filter((entity) => entity.metadata.kind === 'draft-line');
}

function expeditionEntities(entities: CesiumEntityResult[]): CesiumEntityResult[] {
  return entities.filter((entity) => entity.metadata.kind === 'expedition');
}

async function startLineTool(page: Page): Promise<void> {
  await page.getByTestId('start-line-tool').click();
}

// Mirrors reagent-map-expeditions.spec.ts's own addInitScript pattern, extended to also set
// 'dnd_user': both keys must be in localStorage before the next navigation, since that is
// when App.vue's authStore.hydrate() reads them and ReagentMapView's saveLine() needs the
// resulting authStore.isAuthenticated to be true, not just the raw token.
async function signInBrowserAs(page: Page, request: APIRequestContext, email: string): Promise<void> {
  const { token, user } = await signInAs(request, email);
  await page.addInitScript(
    ({ storedToken, storedUser }) => {
      localStorage.setItem('dnd_auth_token', storedToken);
      localStorage.setItem('dnd_user', JSON.stringify(storedUser));
    },
    { storedToken: token, storedUser: user }
  );
}

async function readVisibleExpeditions(
  request: APIRequestContext,
  email?: string
): Promise<StoredExpedition[]> {
  const headers: Record<string, string> = {};
  if (email) {
    const { token } = await signInAs(request, email);
    headers.Authorization = `Bearer ${token}`;
  }
  const response = await request.get(`${MOCK_API_URL}/api/expeditions`, { headers });
  return response.json();
}

async function countPostsTo(page: Page, pattern: string): Promise<{ count: () => number }> {
  let count = 0;
  await page.route(pattern, async (route) => {
    if (route.request().method() === 'POST') count++;
    await route.continue();
  });
  return { count: () => count };
}

// This file writes to expeditions.json (every successful save), so it needs the same
// acquire-lock-before-reading-seed / restore-seed-after-each / release-lock-after-all
// scaffolding reagent-map-api.spec.ts already uses — see architecture-notes.md.
test.describe.configure({ mode: 'serial' });

test.describe('Reagent map save (AC11, AC12, AC13, AC15)', () => {
  let seed: string;

  test.beforeAll(async () => {
    test.setTimeout(HOOK_TIMEOUT_MS);
    await acquireExpeditionsLock();
    seed = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');
  });

  test.afterEach(() => {
    fs.writeFileSync(EXPEDITIONS_PATH, seed, 'utf8');
  });

  test.afterAll(() => {
    releaseExpeditionsLock();
  });

  test.beforeEach(async ({ page, cesium }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();
  });

  test('Elara, signed in, sees her saved expedition again after reopening the reagent map (AC11)', async ({
    page,
    cesium,
    request,
  }) => {
    await signInBrowserAs(page, request, ELARA_EMAIL);
    await page.reload();
    await cesium.waitForSettled();

    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(EMBER_CRAG.lat, EMBER_CRAG.lng);
    await page.getByTestId('save-line-tool').click();
    await expect(page.getByTestId('save-line-status')).toHaveText(SAVED_MESSAGE);

    await page.reload();
    await cesium.waitForSettled();

    const radiusMeters = await cesium.globeRadiusMeters();
    const entities = expeditionEntities(await cesium.getEntities());
    const match = entities.find((entity) => {
      const geometry = entity.geometry as { line: Array<{ lat: number; lng: number }> };
      if (geometry.line?.length !== 2) return false;
      const [first, second] = geometry.line;
      const matchesInOrder =
        distance(first, MOONWELL_SPRING, radiusMeters) <= POSITION_TOLERANCE_METERS &&
        distance(second, EMBER_CRAG, radiusMeters) <= POSITION_TOLERANCE_METERS;
      const matchesReversed =
        distance(first, EMBER_CRAG, radiusMeters) <= POSITION_TOLERANCE_METERS &&
        distance(second, MOONWELL_SPRING, radiusMeters) <= POSITION_TOLERANCE_METERS;
      return matchesInOrder || matchesReversed;
    });
    expect(match, 'expected the saved expedition to be drawn after reopening the map').toBeDefined();

    const held = await readVisibleExpeditions(request, ELARA_EMAIL);
    const heldMatch = held.find((expedition) => expedition.id === match!.metadata.id);
    expect(heldMatch, 'expected the shop to hold the saved expedition').toBeDefined();
    expect(distance(heldMatch!.from, MOONWELL_SPRING, radiusMeters)).toBeLessThanOrEqual(
      POSITION_TOLERANCE_METERS
    );
    expect(distance(heldMatch!.to, EMBER_CRAG, radiusMeters)).toBeLessThanOrEqual(
      POSITION_TOLERANCE_METERS
    );
  });

  test('saving without being signed in asks the adventurer to sign in and stores nothing (AC12)', async ({
    page,
    cesium,
  }) => {
    const posts = await countPostsTo(page, EXPEDITIONS_API_PATTERN);

    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(EMBER_CRAG.lat, EMBER_CRAG.lng);
    await page.getByTestId('save-line-tool').click();

    await expect(page.getByTestId('save-line-message')).toHaveText(SIGN_IN_MESSAGE);
    expect(posts.count()).toBe(0);
    expect(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')).toBe(seed);
  });

  test("saving a zero-length line shows the shop's own message and stores nothing (AC13)", async ({
    page,
    cesium,
    request,
  }) => {
    await signInBrowserAs(page, request, ELARA_EMAIL);
    await page.reload();
    await cesium.waitForSettled();

    await startLineTool(page);
    // Both points at the exact same place, never a client-recomputed near-but-under-1km
    // offset: clicking twice at the same lat/lng targets the same screen pixel both times, so
    // the two picked points are the same place regardless of camera zoom — deterministic,
    // unlike trying to land two distinct picks within 1km of each other by sub-pixel accuracy.
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await page.getByTestId('save-line-tool').click();

    await expect(page.getByTestId('save-line-message')).toHaveText(ZERO_LENGTH_MESSAGE);
    expect(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')).toBe(seed);
  });

  test('an unsaved line is gone after leaving the reagent map and opening it again (AC15)', async ({
    page,
    cesium,
  }) => {
    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(EMBER_CRAG.lat, EMBER_CRAG.lng);
    expect(draftPointsOf(await cesium.getEntities())).toHaveLength(2);
    expect(draftLinesOf(await cesium.getEntities())).toHaveLength(1);

    await page.goto('/');
    await page.goto('/reagent-map');
    await cesium.waitForSettled();

    const entities = await cesium.getEntities();
    expect(draftPointsOf(entities)).toHaveLength(0);
    expect(draftLinesOf(entities)).toHaveLength(0);
    expect(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')).toBe(seed);
  });

  test('rapid double-clicking save produces exactly one stored expedition', async ({
    page,
    cesium,
    request,
  }) => {
    await signInBrowserAs(page, request, ELARA_EMAIL);
    await page.reload();
    await cesium.waitForSettled();

    await startLineTool(page);
    await cesium.clickAt(MOONWELL_SPRING.lat, MOONWELL_SPRING.lng);
    await cesium.clickAt(EMBER_CRAG.lat, EMBER_CRAG.lng);

    const posts = await countPostsTo(page, EXPEDITIONS_API_PATTERN);
    const beforeCount = (JSON.parse(seed) as unknown[]).length;

    // One real double-click: both clicks reach the button before the first save returns.
    // Two parallel click() calls would not test that: the second waits for the button,
    // which the first save removes, and times out.
    await page.getByTestId('save-line-tool').dblclick();
    await expect(page.getByTestId('save-line-status')).toHaveText(SAVED_MESSAGE);

    expect(posts.count()).toBe(1);
    const afterCount = (JSON.parse(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')) as unknown[]).length;
    expect(afterCount).toBe(beforeCount + 1);
  });
});
