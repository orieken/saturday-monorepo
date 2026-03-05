import { test } from '@orieken/saturday-playwright-heatmap';
import { expect } from '@playwright/test';

test('heatmap verification', async ({ page, heatmap }) => {
  await page.goto('/');
  await expect(page).toHaveTitle(/Ye Olde Magic Shop/);
  
  // Perform some interactions
  const buyButtons = page.locator('button').filter({ hasText: 'Buy' });
  if (await buyButtons.count() > 0) {
      await buyButtons.first().click();
  }
  
  // Click a link if available
  const links = page.locator('a');
  if (await links.count() > 0) {
       await links.first().click();
  }
});
