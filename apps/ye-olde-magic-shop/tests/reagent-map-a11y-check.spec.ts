import { test, expect } from '@playwright/test';

test('reagent map a11y check', async ({ page }) => {
  await page.goto('/reagent-map');
  await page.waitForFunction(() => Boolean((window as any).__SATURDAY_CESIUM__), null, { timeout: 15000 });

  const globe = page.locator('#cesiumContainer');
  await expect(globe).toHaveAttribute('role', 'img');
  await expect(globe).toHaveAttribute('aria-label', /.+/);

  const list = page.getByTestId('reagent-site-list');
  await expect(list).toHaveAttribute('aria-label', 'Reagent sites');

  const navLink = page.getByTestId('reagent-map-link');
  await expect(navLink).toBeVisible();
  await navLink.focus();
  await expect(navLink).toBeFocused();
});
