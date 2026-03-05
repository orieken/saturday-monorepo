import { setWorldConstructor, World } from '@cucumber/cucumber';
import * as messages from '@cucumber/messages';
import { BrowserContext, Browser, Page, APIRequestContext } from 'playwright';
import { SaturdayWorldOptions } from './saturday-world-options';
import { BaseSite, ConsoleLogger } from '@orieken/saturday-core';
import { SiteManager } from '../support/managers/site-manager';
import { TabManager } from '../support/managers/tab-manager';

export class SaturdayWorld extends World {
  context?: BrowserContext;
  feature?: messages.Pickle;
  page!: Page;
  browser?: Browser;
  request?: APIRequestContext;
  scenarioName?: string;
  consoleLogger!: ConsoleLogger;

  private _siteManager?: SiteManager;
  private _tabManager?: TabManager;

  constructor(options: SaturdayWorldOptions) {
    super(options);
  }

  get siteManager(): SiteManager {
    if (!this._siteManager) {
      if (!this.page) {
         // Fallback if accessed before page exists, though rare in tests
         throw new Error('Page must be initialized before utilizing SiteManager defaults');
      }
      this._siteManager = new SiteManager(this.page);
    }
    return this._siteManager;
  }

  get tabManager(): TabManager {
    if (!this._tabManager) {
      if (!this.context) {
        throw new Error('Browser context must be initialized before using TabManager');
      }
      // If we lazily init, we assume 'page' is the main one
      this._tabManager = new TabManager(this.context, this.page);
    }
    return this._tabManager;
  }

  initializeManagers(): void {
    if (this.context && this.page) {
      this._tabManager = new TabManager(this.context, this.page);
    }
  }

  async cleanupManagers(): Promise<void> {
    if (this._tabManager) {
      try {
        await this._tabManager.closeAll(true);
      } catch (error) {
        console.warn('Error closing tabs:', error);
      }
    }

    if (this._siteManager) {
      this._siteManager.clear();
    }
  }
}

// Optional: Export a method to install it, or let consumer do it
export function installSaturdayWorld() {
  setWorldConstructor(SaturdayWorld);
}
