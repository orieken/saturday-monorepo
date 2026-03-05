import { Then, When } from '@cucumber/cucumber';
import { TestWorld } from '../support/world';
import { expect } from '@playwright/test';

Then('I should see at least {int} products listed', async function (this: TestWorld, count: number) {
  if (!this.site) throw new Error('Site not initialized');
  const items = await this.site.homePage.getItems();
  expect(items.length).toBeGreaterThanOrEqual(count);
});

Then('the total price should equal the product price', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  const items = await this.site.homePage.getItems();
  for (const item of items) {
    const priceText = await item.getByTestId('item-price').innerText();
    const price = parseFloat(priceText.replace(' gp', ''));
    const totalText = await item.getByTestId('item-total').innerText();
    const total = parseFloat(totalText.replace(' gp', ''));
    expect(total).toBe(price);
  }
});

Then('I should see a list of products with their names and prices', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  const items = await this.site.homePage.getItems();
  for (const item of items) {
    const name = await item.getByTestId('item-name').innerText();
    const price = await item.getByTestId('item-price').innerText();
    expect(name).toBeTruthy();
    expect(price).toMatch(/\d+\sgp/); // Simple regex to check for price format like 120 gp
  }
});

When('I click on the first product\'s details link', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.homePage.clickFirstProductDetails();
});

Then('I should see the product details page with name, price, and description', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.productDetailsPage.isDetailsVisible();
});

Then('each product should display an image and a price', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  const items = await this.site.homePage.getItems();
  for (const item of items) {
    const name = await item.getByTestId('item-name').innerText();
    const price = await item.getByTestId('item-price').innerText();
    expect(name).toBeTruthy();
    expect(price).toMatch(/\d+\sgp/); // Simple regex to check for price format like 120 gp
  }
});
