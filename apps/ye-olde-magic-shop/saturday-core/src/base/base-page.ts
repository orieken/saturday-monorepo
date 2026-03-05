
import { Locator, Page } from 'playwright';
import { BaseSite } from './base-site';
import { BaseElement } from './base-element';
import { BaseFilter } from './base-filter';
import { TrainingData, ValidationResult, CaptureOptions, ValidationOptions, TrainingResult, BrowserInfo } from '../ml/interfaces/ml-types';

function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

export abstract class BasePage {
  abstract containerSelector: string;
  abstract path?: string;

  protected page: Page;
  protected site: BaseSite;
  protected container!: Locator;
  protected elements: Map<string, BaseElement> = new Map();
  protected filters: Map<string, BaseFilter> = new Map();

  constructor(page: Page, site: BaseSite) {
    this.page = page;
    this.site = site;
    this.createContainer();
    this.initializeElements();
    this.initializeFilters();
  }

  protected abstract initializeElements(): void;
  protected abstract initializeFilters(): void;

  protected createContainer(): void {
    this.container = this.page.locator(this.containerSelector);
  }

  protected registerElement<T extends BaseElement>(
    name: string,
    elementClass: new (page: Page, selector: string) => T,
    selector: string,
  ): void {
    const element = new elementClass(this.page, selector);
    this.elements.set(name, element);
  }

  protected registerFilter<T extends BaseFilter>(name: string, filterClass: new (page: Page) => T): void {
    const filter = new filterClass(this.page);
    this.filters.set(name, filter);
  }

  getElement<T extends BaseElement>(name: string): T {
    const element = this.elements.get(name);
    if (!element) {
      throw new Error(`Element '${name}' is not registered in ${this.constructor.name}`);
    }
    return element as T;
  }

  getFilter<T extends BaseFilter>(name: string): T {
    const filter = this.filters.get(name);
    if (!filter) {
      throw new Error(`Filter '${name}' is not registered in ${this.constructor.name}`);
    }
    return filter as T;
  }

  async visit(): Promise<void> {
    if (!this.path) {
      throw new Error(`Path is not defined for ${this.constructor.name}`);
    }
    const url = `${this.site.getBaseUrl()}${this.path}`;
    await this.page.goto(url);
    await this.waitForLoad();
  }

  async waitForLoad(): Promise<void> {
    await this.container.waitFor({ state: 'visible' });
  }

  async isLoaded(): Promise<boolean> {
    return await this.container.isVisible();
  }

  getContainer(): Locator {
    return this.container;
  }

  async scrollToTop(): Promise<void> {
    await this.page.evaluate(() => window.scrollTo(0, 0));
  }

  async scrollToBottom(): Promise<void> {
    await this.page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
  }

  // ML capture methods
  async captureForTraining(label: string, options?: CaptureOptions): Promise<TrainingData> {
    const screenshot = await this.page.screenshot({
      fullPage: options?.fullPage ?? true,
      clip: options?.clip
    });

    return {
      id: generateId(),
      timestamp: new Date(),
      label,
      imageBuffer: screenshot,
      metadata: {
        pageUrl: this.page.url(),
        containerSelector: this.containerSelector,
        browserInfo: await this.getBrowserInfo(),
        viewportSize: await this.getViewportSize(),
        ...options?.metadata
      }
    };
  }

  async validateWithML(modelLabel: string, options?: ValidationOptions): Promise<ValidationResult> {
    const currentScreenshot = await this.page.screenshot({
      fullPage: options?.fullPage ?? true,
      clip: options?.clip
    });

    // This would use the site's detector facade
    return await this.site.detectors.anomalyDetector.validateImage(
      currentScreenshot,
      modelLabel,
      options
    );
  }

  // Helper methods for ML
  private async getBrowserInfo(): Promise<BrowserInfo> {
    return await this.page.evaluate(() => ({
      userAgent: navigator.userAgent,
      platform: navigator.platform,
      cookieEnabled: navigator.cookieEnabled
    }));
  }

  private async getViewportSize(): Promise<{ width: number; height: number }> {
    return await this.page.evaluate(() => ({
      width: window.innerWidth,
      height: window.innerHeight
    }));
  }
}
