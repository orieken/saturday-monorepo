/**
 * @issue reagent-map-expeditions
 * Covers AC3/AC4 of expeditions-visibility.feature's two @level-ui scenarios:
 * every expedition GET /api/expeditions returns for the current adventurer is
 * drawn as a geodesic polyline, both ends within 1 km of its stored ends,
 * tagged saturdayMetadata { id, name, kind: 'expedition' } — signed out, only
 * the shop's seeded expeditions; signed in as Elara, the seeded ones plus her
 * own, never Borin's.
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 */
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import expeditionsFixture from '../../mock-api/data/expeditions.json';
import { test, distance, type CesiumEntityResult } from '@orieken/saturday-playwright-cesium';
import { expect, type APIRequestContext } from '@playwright/test';

// Keyed to the same path as reagent-map-api.spec.ts's and reagent-map-save.spec.ts's own
// copies of this lock (architecture-notes.md) — mutual exclusion comes from all three files
// contending for the same filesystem path, not from a shared module. Each file keeps its own
// copy because importing a sibling *.spec.ts file to reuse its functions would also execute
// that file's top-level test.describe(...)/configure(...) calls as a side effect of the
// import (Playwright tracks "the currently loading file suite" as a single value that a
// plain import does not change), re-registering its tests inside the importing file and
// colliding on a second describe.configure({mode: 'serial'}) call — confirmed by running the
// cross-import version architecture-notes.md originally specified, which failed with
// "serial" mode is already assigned for the enclosing scope.
const LOCK_PATH = path.join(os.tmpdir(), 'reagent-map-expeditions.lock');
// See the matching constant's comment in reagent-map-api.spec.ts: measured locally at
// --workers=3 --repeat-each=3, the lock's hold time (one per file's whole serial describe
// block, not one test) summed to ~69s of exclusive contention across the 9 total blocks (3
// files x 3 repeats); 30_000 measurably fired on a legitimate holder and caused a lock-theft
// double-release. 180_000 is comfortably above that measured worst case.
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

type SeededExpedition = {
  id: string;
  name: string;
  owner: string;
  from: { lat: number; lng: number };
  to: { lat: number; lng: number };
};

const EXPEDITIONS: SeededExpedition[] = expeditionsFixture;
const SPRING_TO_CRAG = EXPEDITIONS.find((expedition) => expedition.id === 'expedition-seed-1')!;
const GLADE_TO_CAVES = EXPEDITIONS.find((expedition) => expedition.id === 'expedition-seed-2')!;
const CAIRN_TO_SHOALS = EXPEDITIONS.find((expedition) => expedition.id === 'expedition-seed-3')!;
const ELARAS_EXPEDITION = EXPEDITIONS.find((expedition) => expedition.id === 'expedition-elara-1')!;
const BORINS_EXPEDITION = EXPEDITIONS.find((expedition) => expedition.id === 'expedition-borin-1')!;

const SEEDED_EXPEDITIONS = [SPRING_TO_CRAG, GLADE_TO_CAVES, CAIRN_TO_SHOALS];
const POSITION_TOLERANCE_METERS = 1_000;
const MOCK_API_URL = process.env.MOCK_API_URL || 'http://localhost:8001';
const ELARA_EMAIL = 'elara@aurelis.example';

async function signInAs(request: APIRequestContext, email: string): Promise<string> {
  const response = await request.post(`${MOCK_API_URL}/api/auth/login`, {
    data: { email, password: 'password123' },
  });
  return (await response.json()).token;
}

function expeditionEntities(entities: CesiumEntityResult[]): CesiumEntityResult[] {
  return entities.filter((entity) => entity.metadata.kind === 'expedition');
}

function assertDrawnBetween(
  entity: CesiumEntityResult,
  seeded: SeededExpedition,
  radiusMeters: number
): void {
  expect(entity.metadata.name).toBe(seeded.name);
  const geometry = entity.geometry as { line: Array<{ lat: number; lng: number }> };
  expect(geometry.line).toHaveLength(2);
  const [first, second] = geometry.line;
  const matchesInOrder =
    distance(first, seeded.from, radiusMeters) <= POSITION_TOLERANCE_METERS &&
    distance(second, seeded.to, radiusMeters) <= POSITION_TOLERANCE_METERS;
  const matchesReversed =
    distance(first, seeded.to, radiusMeters) <= POSITION_TOLERANCE_METERS &&
    distance(second, seeded.from, radiusMeters) <= POSITION_TOLERANCE_METERS;
  expect(
    matchesInOrder || matchesReversed,
    `expected "${seeded.name}" to have both ends within ${POSITION_TOLERANCE_METERS}m of ` +
      `its stored ends, got [${JSON.stringify(first)}, ${JSON.stringify(second)}]`
  ).toBe(true);
}

// This file only reads expeditions.json (via the live GET inside each test body), never
// writes it — but it still shares the cross-file lock with reagent-map-api.spec.ts and
// reagent-map-save.spec.ts so its exact-count assertions never observe the file mid-write by
// one of those two (architecture-notes.md).
test.describe.configure({ mode: 'serial' });

test.describe('Reagent map expeditions (AC3, AC4)', () => {
  test.beforeAll(async () => {
    test.setTimeout(HOOK_TIMEOUT_MS);
    await acquireExpeditionsLock();
  });

  test.afterAll(() => {
    releaseExpeditionsLock();
  });

  test('signed out, only the shop\'s three seeded expeditions are drawn', async ({ page, cesium }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();

    const radiusMeters = await cesium.globeRadiusMeters();
    const entities = expeditionEntities(await cesium.getEntities());

    expect(entities).toHaveLength(SEEDED_EXPEDITIONS.length);
    expect(entities.map((entity) => entity.metadata.name).sort()).toEqual(
      SEEDED_EXPEDITIONS.map((expedition) => expedition.name).sort()
    );

    for (const seeded of SEEDED_EXPEDITIONS) {
      const entity = entities.find((candidate) => candidate.metadata.id === seeded.id);
      expect(entity, `expected an expedition entity for ${seeded.id}`).toBeDefined();
      expect(entity!.metadata.kind).toBe('expedition');
      assertDrawnBetween(entity!, seeded, radiusMeters);
    }
  });

  test("Elara signed in sees the seeded expeditions plus her own, never Borin's", async ({
    page,
    cesium,
    request,
  }) => {
    const token = await signInAs(request, ELARA_EMAIL);
    await page.addInitScript((storedToken) => {
      localStorage.setItem('dnd_auth_token', storedToken);
    }, token);

    await page.goto('/reagent-map');
    await cesium.waitForSettled();

    const radiusMeters = await cesium.globeRadiusMeters();
    const entities = expeditionEntities(await cesium.getEntities());
    const expectedVisible = [...SEEDED_EXPEDITIONS, ELARAS_EXPEDITION];

    expect(entities).toHaveLength(expectedVisible.length);
    expect(entities.map((entity) => entity.metadata.name).sort()).toEqual(
      expectedVisible.map((expedition) => expedition.name).sort()
    );
    expect(entities.map((entity) => entity.metadata.name)).not.toContain(BORINS_EXPEDITION.name);

    for (const seeded of expectedVisible) {
      const entity = entities.find((candidate) => candidate.metadata.id === seeded.id);
      expect(entity, `expected an expedition entity for ${seeded.id}`).toBeDefined();
      assertDrawnBetween(entity!, seeded, radiusMeters);
    }
  });
});
