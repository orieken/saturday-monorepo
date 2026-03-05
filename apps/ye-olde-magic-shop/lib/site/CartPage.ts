
import { BasePage, BaseSite, ButtonElement, BaseElement } from '@orieken/saturday-core';
import { Page } from '@playwright/test';

class TextElement extends BaseElement {}

export class CartPage extends BasePage {
  get containerSelector() { return 'body'; }
  get path() { return '/cart'; }

  constructor(page: Page, site: BaseSite) {
    super(page, site);
  }

  protected initializeElements(): void {
    this.registerElement('increaseQtyBtn', ButtonElement, '[data-testid="increase-qty-btn"]');
    this.registerElement('decreaseQtyBtn', ButtonElement, '[data-testid="decrease-qty-btn"]');
    this.registerElement('removeItemBtn', ButtonElement, '[data-testid="remove-item-btn"]');
    
    // Using a broad selector for now as we don't know the exact testid, but assuming it's a button
    // In a real app we'd ask the dev to add a testid
    this.registerElement('completeTransactionBtn', ButtonElement, 'button:has-text("COMPLETE TRANSACTION")');
    this.registerElement('cartItemQty', TextElement, '[data-testid="cart-item-qty"]');
  }
  
  protected initializeFilters(): void {}

  get increaseQtyBtn() { return this.getElement<ButtonElement>('increaseQtyBtn'); }
  get decreaseQtyBtn() { return this.getElement<ButtonElement>('decreaseQtyBtn'); }
  get removeItemBtn() { return this.getElement<ButtonElement>('removeItemBtn'); }
  get completeTransactionBtn() { return this.getElement<ButtonElement>('completeTransactionBtn'); }
  get cartItemQty() { return this.getElement<TextElement>('cartItemQty'); }

  async changeItemQuantity(targetQuantity: number) {
    const getCurrentQty = async () => {
      const qtyText = await this.cartItemQty.getText();
      return parseInt(qtyText || '0', 10);
    };

    const currentQty = await getCurrentQty();
    if (currentQty === targetQuantity) return;

    if (currentQty < targetQuantity) {
      // increasing
      await this.increaseQtyBtn.clickUntil(async () => {
         return (await getCurrentQty()) === targetQuantity;
      }, { maxRetries: targetQuantity - currentQty + 5 });
    } else {
      // decreasing
      await this.decreaseQtyBtn.clickUntil(async () => {
         return (await getCurrentQty()) === targetQuantity;
      }, { maxRetries: currentQty - targetQuantity + 5 });
    }
  }

  async emptyCart() {
    return this.page.locator('[data-testid="cart-empty"]');
  }

  async isItemVisible() {
    await this.page.locator('[data-testid="cart-item"]').first().waitFor();
  }

  async getItemQuantity() {
    return this.cartItemQty.getText();
  }

  async getItemUnitPrice() {
    return this.page.locator('[data-testid="cart-item-price"]').first().textContent();
  }

  async getItemSubtotal() {
    return this.page.locator('[data-testid="cart-item-total"]').first().textContent();
  }

  async getAllItemsSubtotals() {
    return this.page.locator('[data-testid="cart-item-total"]').allTextContents();
  }

  async getCartTotal() {
    return this.page.locator('[data-testid="cart-total"]').textContent();
  }

  async removeFirstItem() {
    return this.removeItemBtn.click();
  }

  async completeTransaction() {
    await this.completeTransactionBtn.click();
  }
}
