import { BaseElement } from '../base/base-element';

export class ButtonElement extends BaseElement {
  async isDisabled(): Promise<boolean> {
    return await this.locator.isDisabled();
  }

  async isEnabled(): Promise<boolean> {
    return await this.locator.isEnabled();
  }
}
