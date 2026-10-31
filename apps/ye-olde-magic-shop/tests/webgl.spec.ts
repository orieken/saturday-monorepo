import { test, expect } from '@playwright/test';

test.describe('WebGL support in the runner browser', () => {
  test('provides a webgl2 context on the running app page', async ({ page }) => {
    await page.goto('/');

    const hasWebgl2Context = await page.evaluate(() => {
      const canvas = document.createElement('canvas');
      return canvas.getContext('webgl2') !== null;
    });

    expect(hasWebgl2Context).toBe(true);
  });
});
