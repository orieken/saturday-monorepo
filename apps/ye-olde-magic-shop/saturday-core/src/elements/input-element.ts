import { BaseElement } from '../base/base-element';

export class InputElement extends BaseElement {
  async type(text: string, options?: { delay?: number }): Promise<void> {
    await this.locator.type(text, options);
  }

  async fill(text: string): Promise<void> {
    await this.locator.fill(text);
  }

  async clear(): Promise<void> {
    await this.locator.clear();
  }

  async getValue(): Promise<string> {
    return await this.locator.inputValue();
  }
}