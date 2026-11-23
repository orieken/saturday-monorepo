import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { expect } from '@playwright/test';

/**
 * One lock over apps/mock-api/data/expeditions.json for every spec that
 * writes it or counts what it holds. It uses the same path as the copies in
 * reagent-map-api.spec.ts, reagent-map-expeditions.spec.ts and
 * reagent-map-save.spec.ts, so all of them take turns. A spec holds it from
 * beforeAll to afterAll and runs its tests serially.
 */
export const EXPEDITIONS_PATH = path.resolve(__dirname, '../../../mock-api/data/expeditions.json');
const LOCK_PATH = path.join(os.tmpdir(), 'reagent-map-expeditions.lock');
const STALE_LOCK_MS = 180_000;
export const LOCK_HOOK_TIMEOUT_MS = 200_000;

function removeIfStale(): void {
  if (Date.now() - fs.statSync(LOCK_PATH).mtimeMs > STALE_LOCK_MS) fs.unlinkSync(LOCK_PATH);
}

export async function acquireExpeditionsLock(): Promise<void> {
  await expect(() => {
    try {
      fs.closeSync(fs.openSync(LOCK_PATH, 'wx'));
    } catch (err) {
      if ((err as NodeJS.ErrnoException).code === 'EEXIST') removeIfStale();
      throw err;
    }
  }).toPass({ timeout: STALE_LOCK_MS });
}

export function releaseExpeditionsLock(): void {
  fs.unlinkSync(LOCK_PATH);
}
