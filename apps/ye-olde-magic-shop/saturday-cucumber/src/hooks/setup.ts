import { chromium, firefox, webkit } from 'playwright';
import { browserOptions } from '../utils/browser-options';
import { SaturdayWorld } from '../world/saturday-world';
import { ITestCaseHookParameter } from '@cucumber/cucumber';
import { Browser } from 'playwright';
import fs from 'node:fs';
import { ConsoleLogger } from '@orieken/saturday-core';
import { getConfig } from '../config';

declare global {
  // eslint-disable-next-line no-var
  var browser: Browser;
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace NodeJS {
    interface Global {
      browser: Browser;
    }
  }
}

const defaultBrowsers: { [k: string]: () => Promise<Browser> } = {
  firefox: async (): Promise<Browser> => firefox.launch(browserOptions),
  webkit: async (): Promise<Browser> => webkit.launch(browserOptions),
  chrome: async (): Promise<Browser> => chromium.launch({ ...browserOptions, channel: 'chrome' }),
  chromium: async (): Promise<Browser> => chromium.launch(browserOptions),
};

export function createBrowser(): () => Promise<void> {
  return async function () {
    const config = getConfig();
    const browserName = process.env.BROWSER ?? 'chrome';
    
    // Check for custom browser definition in config
    if (config.customBrowsers && config.customBrowsers[browserName]) {
      global.browser = await config.customBrowsers[browserName]({ ...browserOptions, ...config.browser });
      return;
    }
    
    // Check default browsers
    if (defaultBrowsers[browserName]) {
       global.browser = await defaultBrowsers[browserName]();
       return;
    }
    
    throw new Error(`Browser '${browserName}' not found in default or custom browsers.`);
  };
}

export function closeBrowser(): () => Promise<void> {
  return async function () {
    await global.browser.close();
  };
}

export function createContext(): (this: SaturdayWorld, { pickle }: ITestCaseHookParameter) => Promise<void> {
  return async function (this: SaturdayWorld, { pickle }: ITestCaseHookParameter) {
    const config = getConfig();
    this.context = await global.browser.newContext({
      acceptDownloads: true,
      recordVideo: config.context.recordVideo,
      viewport: config.context.viewport,
      ...config.context
    });

    this.page = await this.context?.newPage();
    this.request = this.context?.request;
    this.feature = pickle;
    this.scenarioName = pickle.name.replace(/\s+/g, '_');
    this.consoleLogger = new ConsoleLogger(this.page);
    
    this.initializeManagers();
  };
}

async function attachTabsInfo(this: SaturdayWorld) {
  // Access property directly to avoid auto-init throw if not set (though hooks run after init)
  // We use the tabManager getter if context exists
  if (this.context && this.tabManager && this.tabManager.count() > 1) {
    const tabsInfo: string[] = ['=== Active Tabs ==='];
    
    await this.tabManager.forEach(async (page, name, metadata) => {
      const isActive = name === this.tabManager.getActiveName();
      const marker = isActive ? '→' : ' ';
      let url = 'unknown';
      try { url = page.url(); } catch {}

      tabsInfo.push(
        `${marker} ${name}: ${url}` +
        (metadata.purpose ? ` [${metadata.purpose}]` : '') +
        (metadata.openedFrom ? ` (from: ${metadata.openedFrom})` : '')
      );
    });

    this.attach(tabsInfo.join('\n'), 'text/plain');
  }
}

async function attachSitesInfo(this: SaturdayWorld) {
  // Check private prop or safe check to avoid instantiation side effects if not used
  const manager = (this as any)._siteManager;
  if (manager && manager.count() > 0) {
    const sitesInfo: string[] = ['=== Registered Sites ==='];
    
    const sites = manager.listSites();
    const activeName = manager.getActiveName();
    
    for (const name of sites) {
      const isActive = name === activeName;
      const marker = isActive ? '→' : ' ';
      const site = manager.get(name);
      let baseUrl = 'unknown';
      try { baseUrl = site.getBaseUrl(); } catch {}
      
      sitesInfo.push(`${marker} ${name}: ${baseUrl} (${site.constructor.name})`);
    }

    this.attach(sitesInfo.join('\n'), 'text/plain');
  }
}

async function attachScreenshot(this: SaturdayWorld) {
  // Take screenshot from active page defined by TabManager if available
  let activePage = this.page;
  if (this.context && (this as any)._tabManager) {
     activePage = this.tabManager.tryGetActive() || this.page;
  }
  
  try {
      const image = await activePage?.screenshot();
      if (image) {
        this.attach(image, 'image/png');
      }
  } catch (e) {
      console.warn('Failed to take screenshot', e);
  }
}

async function attachVideo(this: SaturdayWorld) {
  const videoPath = (await this.page.video()?.path()) ?? '';

  try {
    await fs.promises.access(videoPath!);
    const videoData = await fs.promises.readFile(videoPath!);
    this.attach(videoData, 'video/mp4');
  } catch (error) {
    console.error('Error accessing or reading video file:', error);
  }
}

async function createReport(this: SaturdayWorld, { result }: ITestCaseHookParameter) {
  if (result) {
    this.attach(`Status: ${result?.status}. Duration: ${result.duration?.seconds}s`);
    await attachScreenshot.call(this);
    await attachTabsInfo.call(this);
    await attachSitesInfo.call(this);
  }
}

export function closeContext(): (this: SaturdayWorld, hookParameter: ITestCaseHookParameter) => Promise<void> {
  return async function (this: SaturdayWorld, { result }: ITestCaseHookParameter) {
    await createReport.call(this, { result } as ITestCaseHookParameter);

    const logs = this.consoleLogger.getFormattedLogs();
    if (logs) {
      this.attach(logs, 'text/plain');
    }

    await this.cleanupManagers();

    await this.page?.close();
    await this.context?.close();
    await attachVideo.call(this);
  };
}
