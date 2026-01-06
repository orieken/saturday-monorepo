import { BaseSite } from './base-site';
import { BasePage } from './base-page';

export abstract class BaseFlow {
  protected site: BaseSite;

  constructor(site: BaseSite) {
    this.site = site;
  }

  abstract execute(params?: Record<string, any>): Promise<void>;

  protected getPage<T extends BasePage>(pageName: string): T {
    return this.site.getPage<T>(pageName);
  }

  protected getFlow<T extends BaseFlow>(flowName: string): T {
    return this.site.getFlow<T>(flowName);
  }
}