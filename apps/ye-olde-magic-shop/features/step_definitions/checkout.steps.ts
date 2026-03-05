
import { When, Then } from '@cucumber/cucumber';
import { expect } from '@playwright/test';
import { TestWorld } from '../support/world';

When('I complete the transaction', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.cartPage.completeTransaction();
});


Then('I should be redirected to the order history page', async function (this: TestWorld) {
  if (!this.site || !this.page) throw new Error('Site or Page not initialized');
  await expect(this.page).toHaveURL(/.*\/account/);
  // Optionally verify "ORDER HISTORY" text or tab active
  await expect(this.page.getByTestId('order-history-tab')).toBeVisible();
});

Then('I should see the new order in the history', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  
  // Ensure we are on the order history tab
  await this.site.accountPage.switchToOrderHistoryTab();
  
  // Wait for orders to load
  const count = await this.site.accountPage.getOrderCount();
  expect(count).toBeGreaterThan(0);
  
  // Verify the top order is PENDING (just created)
  const firstOrder = this.page!.getByTestId('order-item').first();
  await expect(firstOrder).toContainText('PENDING');
});
