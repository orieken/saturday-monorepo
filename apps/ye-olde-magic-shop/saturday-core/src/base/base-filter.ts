import { Page } from 'playwright';


export abstract class BaseFilter {
  protected page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  abstract matches(condition: any): Promise<boolean>;
  abstract getValue(): Promise<any>;
  abstract apply(value: any): Promise<void>;
  abstract clear(): Promise<void>;
}
