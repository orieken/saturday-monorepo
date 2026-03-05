import { Then, When } from '@cucumber/cucumber';
import { TestWorld } from '../support/world';
import assert from 'assert';
import { expect } from '@playwright/test';

When('I add the first product to the cart', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.homePage.addFirstProductToCart();
});

When('I go to the cart page', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.homePage.goToCart();
});

Then('I should see the product in the cart with quantity 1', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.cartPage.isItemVisible();
  const qty = await this.site.cartPage.getItemQuantity();
  assert.strictEqual(qty, '1');
});

When('I add the second product to the cart', function () {
  if (!this.site) throw new Error('Site not initialized');
  return this.site.homePage.addSecondProductToCart();
});

Then('I should see {int} items in the cart', async function (cartCount: number) {
  if (!this.site) throw new Error('Site not initialized');
  const count = await this.site.cartPage.page.locator('[data-testid="cart-item"]').count();
  assert.strictEqual(count, cartCount);
});

When('I remove the second cart item', function () {
  if (!this.site) throw new Error('Site not initialized');
  return this.site.cartPage.page.locator('[data-testid="remove-item-btn"]').nth(1).click();
});

Then('I should see {int} item in the cart', async function (itemCount: number) {
  if (!this.site) throw new Error('Site not initialized');
  const count = await this.site.cartPage.page.locator('[data-testid="cart-item"]').count();
  assert.strictEqual(count, itemCount);
});

Then('the cart total should equal the remaining item subtotal', async function () {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.page.getByTestId('cart-total').click();
  await expect(this.site.page.getByTestId('cart-total')).toContainText('120.00 gp');
  await expect(this.site.page.getByTestId('cart-item-total')).toContainText('120.00 gp');
});

When('I remove the first cart item', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.cartPage.removeFirstItem();
});

Then('the cart should be empty', async function () {
  if (!this.site) throw new Error('Site not initialized');

  const cart = await this.site.cartPage.emptyCart();
  await expect(cart).toBeVisible();
  await expect(cart).toContainText('Your stash is empty');

});

When('I increase the quantity of the first cart item to {int}', async function (this: TestWorld, quantity: number) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.cartPage.changeItemQuantity(quantity);
});

Then('I should see quantity {int} for the first cart item', function (this: TestWorld, quantity: number) {
  if (!this.site) throw new Error('Site not initialized');
  return this.site.cartPage.getItemQuantity().then((qtyText: string | null) => {
    const qty = parseInt(qtyText as string, 10);
    assert.strictEqual(qty, quantity);
  });

});

Then('the subtotal for the first cart item should be the unit price times {int}', function (this: TestWorld, quantity: number) {
  if (!this.site) throw new Error('Site not initialized');
  return Promise.all([
    this.site.cartPage.getItemUnitPrice(),
    this.site.cartPage.getItemSubtotal()
  ]).then(([unitPriceText, subtotalText]) => {
    const unitPrice = parseFloat(unitPriceText!.replace(' gp', ''));
    const subtotal = parseFloat(subtotalText!.replace(' gp', ''));
    assert.strictEqual(subtotal, unitPrice * quantity);
  });
});

Then('the cart total should reflect the updated subtotal', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  const subtotals = await this.site.cartPage.getAllItemsSubtotals()
  const cartTotal = await this.site.cartPage.getCartTotal()
  const sum = subtotals.reduce((acc: any, val: any) => {
    const num = parseFloat(val.replace(' gp', ''));
    return acc + num;
  }, 0);
  const totalNum = parseFloat((cartTotal as string).replace(' gp', ''));
  assert.strictEqual(totalNum, sum);
});

Then('the cart total should equal the subtotal of all items', async function () {
  if (!this.site) throw new Error('Site not initialized');
  const subtotals = await this.site.cartPage.getAllItemsSubtotals()
  const cartTotal = await this.site.cartPage.getCartTotal()
  const sum = subtotals.reduce((acc: any, val: any) => {
    const num = parseFloat(val.replace(' gp', ''));
    return acc + num;
  }, 0);
  const totalNum = parseFloat((cartTotal as string).replace(' gp', ''));
  assert.strictEqual(totalNum, sum);
});