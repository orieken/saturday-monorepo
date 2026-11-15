/**
 * @issue reagent-map-globe
 * @issue reagent-map-expeditions
 * Covers the failure-path branches in ReagentMapView.vue's onMounted that
 * the happy-path specs never exercise (flagged in implementation-notes.md's
 * notesForQa): fetchReagentSites() rejecting after buildViewer() has already
 * resolved, fetchExpeditions() rejecting after buildViewer() and
 * fetchReagentSites() have both resolved (AC3/AC4's companion failure path,
 * added by reagent-map-expeditions), and buildViewer() itself rejecting.
 * In every case the viewer must be destroyed, leaving no canvas or WebGL
 * context behind, across repeated mounts, or (for buildViewer()) onMounted
 * must rethrow with no entities/sites ever shown.
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 */
import { test, expect } from '@playwright/test';

const REAGENT_SITES_API_PATTERN = '**/api/reagent-sites';
const EXPEDITIONS_API_PATTERN = '**/api/expeditions';
const CESIUM_SCRIPT_PATTERN = '**/cesium/Cesium.js';
const PAGE_ERROR_TIMEOUT_MS = 20_000;
const MOUNT_REPETITIONS = 3;

async function assertNoViewerArtifactsLeft(page: import('@playwright/test').Page) {
  const seamIsAbsent = await page.evaluate(() => (window as any).__SATURDAY_CESIUM__ === undefined);
  expect(seamIsAbsent).toBe(true);

  const canvasCount = await page.locator('#cesiumContainer canvas').count();
  expect(canvasCount).toBe(0);

  const siteListItemCount = await page.getByTestId('reagent-site-list').locator('li').count();
  expect(siteListItemCount).toBe(0);
}

test.describe('ReagentMapView onMounted failure paths', () => {
  test('destroys the viewer and leaves no canvas when fetchReagentSites() rejects after the viewer resolves, across repeated mounts', async ({
    page,
  }) => {
    await page.route(REAGENT_SITES_API_PATTERN, (route) => route.abort('failed'));

    for (let mount = 0; mount < MOUNT_REPETITIONS; mount++) {
      const pageErrorPromise = page.waitForEvent('pageerror', { timeout: PAGE_ERROR_TIMEOUT_MS });
      await page.goto('/reagent-map');
      const error = await pageErrorPromise;
      expect(error.message).toBeTruthy();

      await assertNoViewerArtifactsLeft(page);
    }
  });

  test('destroys the viewer and leaves no canvas when fetchExpeditions() rejects after the viewer and reagent sites both resolve, across repeated mounts', async ({
    page,
  }) => {
    await page.route(EXPEDITIONS_API_PATTERN, (route) => route.abort('failed'));

    for (let mount = 0; mount < MOUNT_REPETITIONS; mount++) {
      const pageErrorPromise = page.waitForEvent('pageerror', { timeout: PAGE_ERROR_TIMEOUT_MS });
      await page.goto('/reagent-map');
      const error = await pageErrorPromise;
      expect(error.message).toBeTruthy();

      await assertNoViewerArtifactsLeft(page);
    }
  });

  test('rethrows with no entities/sites shown when buildViewer() itself rejects', async ({ page }) => {
    await page.route(CESIUM_SCRIPT_PATTERN, (route) => route.abort('failed'));

    const pageErrorPromise = page.waitForEvent('pageerror', { timeout: PAGE_ERROR_TIMEOUT_MS });
    await page.goto('/reagent-map');
    const error = await pageErrorPromise;
    expect(error.message).toBeTruthy();

    await assertNoViewerArtifactsLeft(page);
  });
});
