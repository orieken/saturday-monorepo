import { test, expect } from '@playwright/test';

test('magic shop loads and displays wares', async ({ page }) => {
  // Setup the wait for response BEFORE navigating to avoid race conditions
  const itemsResponsePromise = page.waitForResponse(response => 
    response.url().includes('/api/items') && response.status() === 200
  );

  await page.goto('/');
  
  // Expect the page to have a title
  await expect(page).toHaveTitle(/Ye Olde Magic Shop/);

  // Wait for the API call to succeed
  await itemsResponsePromise;

  // Verify UI content confirms successful loading
  await expect(page.getByText('Wares for Sale')).toBeVisible();
  
  // Optionally verify that we have items
  // Assuming ItemCard renders something recognizable, but header is enough for "smoke"
});
