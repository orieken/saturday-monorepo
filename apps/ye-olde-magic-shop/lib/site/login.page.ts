import { BasePage, BaseSite } from '@orieken/saturday-core';
import { Page } from '@playwright/test';

export class LoginPage extends BasePage {
  get containerSelector() { return 'body'; }
  get path() { return '/login'; }

  constructor(page: Page, site: BaseSite) {
    super(page, site);
  }

  protected initializeElements(): void {}
  protected initializeFilters(): void {}

  async navigate() {
    const url = `${(this.site as any).baseUrl}${this.path}`;
    await this.page.goto(url, { waitUntil: 'networkidle' });
  }

  async login(email: string, password: string) {
    await this.page.getByTestId('email-input').fill(email);
    await this.page.getByTestId('password-input').fill(password);
    await this.page.getByTestId('login-button').click();
    // Wait for navigation or error message
    await this.page.waitForTimeout(500);
  }

  async getErrorMessage() {
    const errorElement = this.page.getByTestId('login-error');
    return await errorElement.textContent();
  }

  async isOnLoginPage() {
    return this.page.url().includes('/login');
  }

  async isErrorVisible() {
    return await this.page.getByTestId('login-error').isVisible();
  }
}
