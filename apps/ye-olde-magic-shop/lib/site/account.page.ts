import { BasePage, BaseSite } from '@orieken/saturday-core';
import { Page } from '@playwright/test';

export class AccountPage extends BasePage {
  get containerSelector() { return 'body'; }
  get path() { return '/account'; }

  constructor(page: Page, site: BaseSite) {
    super(page, site);
  }

  protected initializeElements(): void {}
  protected initializeFilters(): void {}

  async navigate() {
    const url = `${(this.site as any).baseUrl}${this.path}`;
    await this.page.goto(url);
    await this.page.waitForLoadState('domcontentloaded');
  }

  async isOnAccountPage() {
    return this.page.url().includes('/account');
  }

  async switchToOrderHistoryTab() {
    await this.page.getByTestId('order-history-tab').click();
    await this.page.waitForTimeout(300);
  }

  async getOrderCount() {
    return await this.page.getByTestId('order-item').count();
  }

  async getOrders() {
    const orderElements = await this.page.getByTestId('order-item').all();
    return orderElements;
  }

  async clickFirstOrder() {
    await this.page.getByTestId('order-item').first().click();
  }

  async getOrderDetails() {
    // Get the first order item container
    const firstOrder = this.page.getByTestId('order-item').first();
    
    return {
      date: await firstOrder.getByTestId('order-date').textContent(),
      total: await firstOrder.getByTestId('order-total').textContent(),
      items: await firstOrder.getByTestId('order-line-item').count()
    };
  }

  async logout() {
    await this.page.getByTestId('logout-button').click();
    await this.page.waitForTimeout(500);
  }

  async isAccountInfoVisible() {
    // Check if user name is visible in header
    const userInfo = this.page.getByTestId('user-name');
    return await userInfo.isVisible().catch(() => false);
  }
}
