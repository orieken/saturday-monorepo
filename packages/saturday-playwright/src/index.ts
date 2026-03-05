import { test as base } from '@playwright/test';
import { SiteManager, TabManager, ConsoleLogger } from '@orieken/saturday-core';

export type SaturdayFixtures = {
  consoleLogger: ConsoleLogger;
  siteManager: SiteManager;
  tabManager: TabManager;
};

export const test = base.extend<SaturdayFixtures>({
  consoleLogger: async ({ page }, use, testInfo) => {
    const logger = new ConsoleLogger(page);
    await use(logger);

    if (testInfo.status !== testInfo.expectedStatus) {
      const logs = logger.getFormattedLogs();
      if (logs) {
        await testInfo.attach('Console Logs', {
          contentType: 'text/plain',
          body: Buffer.from(logs),
        });
      }
    }
  },
  tabManager: async ({ context, page }, use, testInfo) => {
    const tabManager = new TabManager(context, page);
    await use(tabManager);

    if (testInfo.status !== testInfo.expectedStatus) {
      await testInfo.attach('Active Tab State', {
        contentType: 'application/json',
        body: Buffer.from(JSON.stringify({ 
          activeTabName: tabManager.getActiveName(), 
          tabCount: tabManager.count() 
        }, null, 2)),
      });
    }

    await tabManager.closeAll(true);
  },
  siteManager: async ({ page }, use, testInfo) => {
    const siteManager = new SiteManager(page);
    await use(siteManager);

    if (testInfo.status !== testInfo.expectedStatus) {
      await testInfo.attach('Site Manager State', {
        contentType: 'application/json',
        body: Buffer.from(JSON.stringify({ 
          activeSiteName: siteManager.getActiveName() ?? null, 
          sites: siteManager.listSites() 
        }, null, 2)),
      });
    }

    siteManager.clear();
  },
});

export { expect } from '@playwright/test';
