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
import expeditionsFixture from '../../mock-api/data/expeditions.json';
import { test, distance, type CesiumEntityResult } from '@orieken/saturday-playwright-cesium';
import { expect, type APIRequestContext } from '@playwright/test';

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

test.describe('Reagent map expeditions (AC3, AC4)', () => {
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
