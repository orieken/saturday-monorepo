import { expect } from '@playwright/test';
import { test } from '@orieken/saturday-playwright-k6-exporter/fixture';
import { createK6Recorder } from '@orieken/saturday-playwright-k6-exporter';

test('homepage smoke test @k6', async ({ page }, testInfo) => {
  const setup = await createK6Recorder(testInfo.title);
  // Optional: check setup if you want to skip tests when K6_EXPORT is not set, 
  // but usually we just let it run as a normal playwright test if setup is null.
  
  const { recorder, ctx } = setup || { recorder: null, ctx: page.request };

  // Note: if using recorder, use 'ctx' for API calls you want recorded.
  // For standard page interactions, use 'page'.
  
  await page.goto('/');
  await expect(page).toHaveTitle(/Ye Olde Magic Shop/);

  // Example API call if we had one
  // await ctx.get('/api/health');

  if (recorder) {
    await recorder.flushToK6();
  }
});
