
import { BasePage, BaseSite } from '@orieken/saturday-core';
import { Page } from '@playwright/test';

export class HomePage extends BasePage {
  get containerSelector() { return 'body'; }
  get path() { return '/'; }

  constructor(page: Page, site: BaseSite) {
    super(page, site);
  }

  protected initializeElements(): void {}
  protected initializeFilters(): void {}

  async addFirstProductToCart() {
    await this.page.locator('[data-testid="add-to-cart-btn"]').first().click();
  }
  async addSecondProductToCart() {
    await this.page.locator('[data-testid="add-to-cart-btn"]').nth(1).click();
  }

  async addRandomProductToCart() {
    const items = this.page.locator('[data-testid="item-card"]');
    const count = await items.count();
    const randomIndex = Math.floor(Math.random() * count);
    await items.nth(randomIndex).locator('[data-testid="add-to-cart-btn"]').click();
  }

  async goToCart() {
    await this.page.locator('[data-testid="cart-link"]').click();
  }
  
  async getItems() {
    return this.page.locator('[data-testid="item-card"]').all();
  }

  async clickFirstProductDetails() {
    await this.page.locator('[data-testid="item-card"] a').first().click();
  }
}
