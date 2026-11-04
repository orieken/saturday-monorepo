import fs from 'node:fs';
import path from 'node:path';
import { test, expect } from '@playwright/test';
import { ShopApi, Expedition, GeoPoint } from '../lib/api/ShopApi';
// The pure geometry module: no dependency on Cesium or a browser.
import { distance } from '../../../packages/saturday-playwright-cesium/src/geometry';

// Order F1, from features/reagent-map/ (@level-api). Spec tolerances are on Aurelis, never Earth.
const AURELIS_RADIUS_METERS = 6_000_000;
const TOLERANCE_METERS = 1000;

const MOONWELL_SPRING = { lat: 10, lng: 20 };
const EMBER_CRAG = { lat: 25, lng: 60 };
const ELARA = 'elara@aurelis.example';

const SEEDED_NAMES = ['Spring to Crag', 'Glade to Caves', 'Cairn to Shoals'];
const ELARAS_NAME = 'Moonwell Spring to Frostcap Hollow';
const BORINS_NAME = 'Ember Crag to Saltmarsh Cairn';

const EXPEDITIONS_PATH = path.resolve(__dirname, '../../mock-api/data/expeditions.json');

function metersBetween(a: GeoPoint, b: GeoPoint): number {
  return distance(a, b, AURELIS_RADIUS_METERS);
}

function names(expeditions: Expedition[]): (string | undefined)[] {
  return expeditions.map((expedition) => expedition.name);
}

// Saving writes to the mock API's data file, so these run one at a time and put it back.
test.describe.configure({ mode: 'serial' });

test.describe('Reagent map API', () => {
  let seed: string;

  test.beforeAll(() => {
    seed = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');
  });

  test.afterEach(() => {
    fs.writeFileSync(EXPEDITIONS_PATH, seed, 'utf8');
  });

  test('@AC1 The shop publishes every reagent site', async ({ request }) => {
    const sites = await new ShopApi(request).reagentSites();

    expect(sites).toHaveLength(8);
    const tidewrack = sites.find((site) => site.name === 'Tidewrack Shoals');
    expect(tidewrack).toMatchObject({ lat: 0, lng: 180 });
    const frostcap = sites.find((site) => site.name === 'Frostcap Hollow');
    expect(frostcap).toMatchObject({ lat: 85, lng: 0 });
  });

  test('@AC3 Signed out, only the seeded expeditions are visible', async ({ request }) => {
    const expeditions = await new ShopApi(request).visibleExpeditions();

    expect(names(expeditions).sort()).toEqual([...SEEDED_NAMES].sort());
  });

  test("@AC4 An adventurer sees their own expeditions and not another adventurer's", async ({
    request,
  }) => {
    const api = new ShopApi(request);
    await api.signIn(ELARA);

    const visible = names(await api.visibleExpeditions());

    expect(visible.sort()).toEqual([...SEEDED_NAMES, ELARAS_NAME].sort());
    expect(visible).not.toContain(BORINS_NAME);
  });

  test('@AC11 The shop keeps an expedition for the adventurer who saved it', async ({ request }) => {
    const elara = new ShopApi(request);
    await elara.signIn(ELARA);

    const saved = await (await elara.saveExpedition(MOONWELL_SPRING, EMBER_CRAG)).json();

    const held = (await elara.visibleExpeditions()).find((e) => e.id === saved.id);
    expect(held).toBeDefined();
    expect(metersBetween(held!.from, MOONWELL_SPRING)).toBeLessThanOrEqual(TOLERANCE_METERS);
    expect(metersBetween(held!.to, EMBER_CRAG)).toBeLessThanOrEqual(TOLERANCE_METERS);

    const signedOut = await new ShopApi(request).visibleExpeditions();
    expect(signedOut.map((e) => e.id)).not.toContain(saved.id);
  });

  test('@AC12 Saving while signed out is refused', async ({ request }) => {
    const response = await new ShopApi(request).saveExpedition(MOONWELL_SPRING, EMBER_CRAG);

    expect(response.status()).toBe(401);
    expect(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')).toBe(seed);
  });

  test('@AC13 A zero-length expedition is refused', async ({ request }) => {
    const elara = new ShopApi(request);
    await elara.signIn(ELARA);
    // 0.005 degrees is about 524 m on Aurelis, which is less than 1 km.
    const nearby = { lat: MOONWELL_SPRING.lat + 0.005, lng: MOONWELL_SPRING.lng };
    expect(metersBetween(MOONWELL_SPRING, nearby)).toBeLessThan(TOLERANCE_METERS);

    const response = await elara.saveExpedition(MOONWELL_SPRING, nearby);

    expect(response.status()).toBe(422);
    expect((await response.json()).error).toMatch(/zero-length/i);
    expect(fs.readFileSync(EXPEDITIONS_PATH, 'utf8')).toBe(seed);
  });
});
