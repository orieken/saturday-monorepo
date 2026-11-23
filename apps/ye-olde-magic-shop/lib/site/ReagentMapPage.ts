import { BasePage, BaseSite } from '@orieken/saturday-core';
import { Locator, Page } from '@playwright/test';

/**
 * The reagent map, as an adventurer sees it: the globe, the site list and the
 * line tool's controls. Gestures on the globe itself go through the
 * saturday-playwright-cesium fixture, never through Cesium's own API (D-5).
 */
export class ReagentMapPage extends BasePage {
  get containerSelector() { return 'body'; }
  get path() { return '/reagent-map'; }

  constructor(page: Page, site: BaseSite) {
    super(page, site);
  }

  protected initializeElements(): void {}
  protected initializeFilters(): void {}

  async open(): Promise<void> {
    await this.page.goto(this.path);
  }

  async startLine(): Promise<void> {
    await this.page.getByTestId('start-line-tool').click();
  }

  async save(): Promise<void> {
    await this.page.getByTestId('save-line-tool').click();
  }

  async flyTo(siteId: string): Promise<void> {
    await this.page.getByTestId(`fly-to-site-${siteId}`).click();
  }

  get saveButton(): Locator { return this.page.getByTestId('save-line-tool'); }
  get siteList(): Locator { return this.page.getByTestId('reagent-site-list'); }
  get lineToolStatus(): Locator { return this.page.getByTestId('line-tool-status'); }
  get saveMessage(): Locator { return this.page.getByTestId('save-line-message'); }
  get saveStatus(): Locator { return this.page.getByTestId('save-line-status'); }
}
