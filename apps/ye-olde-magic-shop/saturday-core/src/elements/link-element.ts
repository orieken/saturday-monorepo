import { BaseElement } from '../base/base-element';

export class LinkElement extends BaseElement {
  async getHref(): Promise<string | null> {
    return await this.getAttribute('href');
  }

  async getTarget(): Promise<string | null> {
    return await this.getAttribute('target');
  }

  async opensInNewTab(): Promise<boolean> {
    const target = await this.getTarget();
    return target === '_blank';
  }
}