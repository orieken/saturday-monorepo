import { BaseSite } from '@orieken/saturday-core';
import { Page } from 'playwright';

export class SiteManager {
  private sites: Map<string, BaseSite>;
  private activeSiteName: string | null;
  private page: Page;

  constructor(page: Page) {
    this.sites = new Map();
    this.activeSiteName = null;
    this.page = page;
  }

  register(name: string, site: BaseSite): void {
    if (this.sites.has(name)) {
      throw new Error(`Site '${name}' is already registered.`);
    }
    this.sites.set(name, site);
    if (this.activeSiteName === null) {
      this.activeSiteName = name;
    }
  }

  get<T extends BaseSite>(name: string): T {
    const site = this.sites.get(name);
    if (!site) {
      throw new Error(`Site '${name}' is not registered.`);
    }
    return site as T;
  }

  getActive<T extends BaseSite>(): T {
    if (!this.activeSiteName) {
      throw new Error('No active site is set.');
    }
    return this.get<T>(this.activeSiteName);
  }

  setActive(name: string): void {
    if (!this.sites.has(name)) {
      throw new Error(`Cannot set active site. Site '${name}' is not registered.`);
    }
    this.activeSiteName = name;
    // Note: Switching active site in the manager doesn't necessarily mean switching the active tab.
    // That coordination should happen in the World/Steps if needed, or by ensuring the site uses the correct page.
  }

  getActiveName(): string | null {
    return this.activeSiteName;
  }

  listSites(): string[] {
    return Array.from(this.sites.keys());
  }

  count(): number {
    return this.sites.size;
  }

  clear(): void {
    this.sites.clear();
    this.activeSiteName = null;
  }
}
