import { test, expect } from '@playwright/test';
import { Site } from '../lib/site/MagicShop';

test.describe('Shopping Cart', () => {
  let site: Site;

  test.beforeEach(async ({ page }) => {
    site = new Site(page);
    await site.homePage.visit();
  });

  test('Add the same item twice increases quantity to 2', async () => {
    await site.homePage.addFirstProductToCart();
    await site.homePage.addFirstProductToCart();
    await site.homePage.goToCart();

    const qty = await site.cartPage.getItemQuantity();
    expect(qty).toBe('2');

    const unitPriceStr = await site.cartPage.getItemUnitPrice();
    const subtotalStr = await site.cartPage.getItemSubtotal();
    
    const unitPrice = parseFloat(unitPriceStr!.replace(' gp', ''));
    const subtotal = parseFloat(subtotalStr!.replace(' gp', ''));

    expect(subtotal).toBe(unitPrice * 2);
  });

  test('Increase quantity from the cart page', async () => {
    await site.homePage.addFirstProductToCart();
    await site.homePage.goToCart();

    await site.cartPage.changeItemQuantity(2);

    const qty = await site.cartPage.getItemQuantity();
    expect(qty).toBe('2');

    const unitPriceStr = await site.cartPage.getItemUnitPrice();
    const subtotalStr = await site.cartPage.getItemSubtotal();
    
    const unitPrice = parseFloat(unitPriceStr!.replace(' gp', ''));
    const subtotal = parseFloat(subtotalStr!.replace(' gp', ''));

    expect(subtotal).toBe(unitPrice * 2);
    
    const cartTotalStr = await site.cartPage.getCartTotal();
    const cartTotal = parseFloat(cartTotalStr!.replace(' gp', ''));
    expect(cartTotal).toBe(subtotal);
  });
});
