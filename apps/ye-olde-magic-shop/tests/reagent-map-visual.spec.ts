/**
 * @issue reagent-map-visual-baselines
 * @visual screenshot baselines for the reagent map's stable DOM chrome (page
 * header/nav link, site list, error state) and the Cesium globe at its fixed
 * default view. Baselines are per-runner-image (see playwright.config.ts's
 * snapshotPathTemplate) since the same page renders differently on a Mac and
 * in each runner image — these tests only run with RUNNER_VARIANT set (see
 * scripts/ui-suite-in-chrome.sh) and skip everywhere else.
 * Runs only under the 'chrome' Playwright project (see playwright.config.ts's
 * testIgnore on the chromium project) — Cesium needs a real webgl2 context.
 */
import { test } from '@orieken/saturday-playwright-cesium';
import { expect } from '@playwright/test';

const REAGENT_SITES_API_PATTERN = '**/api/reagent-sites';
// AC3's own ceiling: never raised past 1% to paper over globe flakiness —
// see the five-run repeat check recorded in the delivery record instead.
const GLOBE_MAX_DIFF_PIXEL_RATIO = 0.01;

test.describe('Reagent map visual baselines', () => {
  test.skip(
    !process.env.RUNNER_VARIANT,
    'RUNNER_VARIANT is unset: visual baselines render differently on a developer machine ' +
      'than in a runner image, so these tests only run inside scripts/ui-suite-in-chrome.sh.'
  );

  test('page header and nav link match the stored baseline', { tag: '@visual' }, async ({ page }) => {
    await page.goto('/reagent-map');

    const header = page.locator('header');
    await expect(page.getByTestId('reagent-map-link')).toBeVisible();
    await expect(header).toHaveScreenshot('reagent-map-header.png');
  });

  test('site list matches the stored baseline', { tag: '@visual' }, async ({ cesium, page }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();

    const list = page.getByTestId('reagent-site-list');
    await expect(list).toHaveScreenshot('reagent-site-list.png');
  });

  test('error state matches the stored baseline when fetching reagent sites fails', { tag: '@visual' }, async ({
    page,
  }) => {
    await page.route(REAGENT_SITES_API_PATTERN, (route) => route.abort('failed'));
    const pageErrorPromise = page.waitForEvent('pageerror');
    await page.goto('/reagent-map');
    await pageErrorPromise;

    const errorElement = page.getByTestId('reagent-map-error');
    await expect(errorElement).toBeVisible();
    await expect(errorElement).toHaveAttribute('role', 'alert');
    await expect(errorElement).toHaveScreenshot('reagent-map-error.png');
  });

  test('globe default view matches the stored baseline within tolerance', { tag: '@visual' }, async ({
    cesium,
    page,
  }) => {
    await page.goto('/reagent-map');
    await cesium.waitForSettled();

    const globe = page.locator('#cesiumContainer');
    await expect(globe).toHaveScreenshot('globe-default-view.png', {
      maxDiffPixelRatio: GLOBE_MAX_DIFF_PIXEL_RATIO,
    });
  });
});
