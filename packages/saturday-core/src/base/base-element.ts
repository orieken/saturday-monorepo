import { Locator, Page } from 'playwright';
import { TrainingData, ElementCaptureOptions, ElementValidationOptions, ElementAnomalyResult } from '../ml/interfaces/ml-types';

function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

export abstract class BaseElement {
  protected page: Page;
  protected selector: string;
  protected locator: Locator;

  constructor(page: Page, selector: string) {
    this.page = page;
    this.selector = selector;
    this.locator = page.locator(selector);
  }

  getLocator(): Locator {
    return this.locator;
  }

  getSelector(): string {
    return this.selector;
  }

  async waitForVisible(): Promise<void> {
    await this.locator.waitFor({ state: 'visible' });
  }

  async waitForHidden(): Promise<void> {
    await this.locator.waitFor({ state: 'hidden' });
  }

  async isVisible(): Promise<boolean> {
    return await this.locator.isVisible();
  }

  async exists(): Promise<boolean> {
    return await this.locator.count() > 0;
  }

  async getText(): Promise<string> {
    return await this.locator.textContent() || '';
  }

  async getInnerText(): Promise<string> {
    return await this.locator.innerText();
  }

  async getAttribute(name: string): Promise<string | null> {
    return await this.locator.getAttribute(name);
  }

  async click(): Promise<void> {
    await this.locator.click();
  }

  async clickUntil(
    condition: () => Promise<boolean>,
    options: { maxRetries?: number; interval?: number } = {}
  ): Promise<void> {
    const { maxRetries = 20, interval = 200 } = options;
    
    // Check initial condition first
    if (await condition()) return;

    for (let i = 0; i < maxRetries; i++) {
      await this.click();
      await this.page.waitForTimeout(interval);
      if (await condition()) return;
    }

    throw new Error(`Condition requested in clickUntil was not met after ${maxRetries} attempts`);
  }

  async doubleClick(): Promise<void> {
    await this.locator.dblclick();
  }

  async rightClick(): Promise<void> {
    await this.locator.click({ button: 'right' });
  }

  async hover(): Promise<void> {
    await this.locator.hover();
  }

  async scrollIntoView(): Promise<void> {
    await this.locator.scrollIntoViewIfNeeded();
  }

  async getBoundingBox(): Promise<{ x: number; y: number; width: number; height: number } | null> {
    return await this.locator.boundingBox();
  }

  // ML methods for element capture and validation
  async captureForTraining(label: string, options?: ElementCaptureOptions): Promise<TrainingData> {
    await this.waitForVisible();

    const boundingBox = await this.getBoundingBox();
    if (!boundingBox) {
      throw new Error(`Element ${this.selector} is not visible for capture`);
    }

    const screenshot = await this.page.screenshot({
      clip: boundingBox,
      ...options?.screenshotOptions
    });

    return {
      id: generateId(),
      timestamp: new Date(),
      label,
      imageBuffer: screenshot,
      metadata: {
        pageUrl: this.page.url(),
        elementSelector: this.selector,
        boundingBox,
        browserInfo: await this.getBrowserInfo(),
        viewportSize: await this.getViewportSize(),
        elementType: await this.locator.evaluate(el => el.tagName.toLowerCase()),
        ...options?.metadata
      }
    };
  }

  async detectAnomalies(modelLabel: string, options?: ElementValidationOptions): Promise<ElementAnomalyResult> {
    const boundingBox = await this.getBoundingBox();
    if (!boundingBox) {
      throw new Error(`Element ${this.selector} is not visible for detection`);
    }

    const screenshot = await this.page.screenshot({ clip: boundingBox });

    // This would use a site-specific element detector
    return await this.detectElementAnomalies(screenshot, modelLabel, options);
  }

  private async detectElementAnomalies(
    screenshot: Buffer,
    modelLabel: string,
    options?: ElementValidationOptions
  ): Promise<ElementAnomalyResult> {
    // Implementation would delegate to ML framework
    // This is a placeholder - actual implementation would come from the detectors facade
    return {
      isValid: true,
      confidence: 0.95,
      anomalies: [],
      elementSelector: this.selector,
      timestamp: new Date(),
      modelUsed: modelLabel
    };
  }

  private async getBrowserInfo(): Promise<{ userAgent: string; platform: string; cookieEnabled: boolean }> {
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