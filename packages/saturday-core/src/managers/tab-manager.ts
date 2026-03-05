import { BrowserContext, Page } from 'playwright';

export interface TabMetadata {
  purpose?: string;
  openedFrom?: string;
}

export class TabManager {
  private tabs: Map<string, Page>;
  private metadata: Map<string, TabMetadata>;
  private context: BrowserContext;
  private activeTabName: string;

  constructor(context: BrowserContext, initialPage: Page, initialName: string = 'main') {
    this.context = context;
    this.tabs = new Map();
    this.metadata = new Map();
    
    // Register the initial page
    this.tabs.set(initialName, initialPage);
    this.metadata.set(initialName, { purpose: 'Initial Tab' });
    this.activeTabName = initialName;

    // Listen for new pages to auto-register popup windows
    context.on('page', (page) => {
      // Find if this page is already registered
      // This listener is useful for un-managed popups, but for explicit 'openInNewTab', we handle it manually.
    });
  }

  async openInNewTab(name: string, url: string, metadata?: TabMetadata): Promise<Page> {
    if (this.tabs.has(name)) {
      throw new Error(`Tab '${name}' is already registered.`);
    }

    const newPage = await this.context.newPage();
    await newPage.goto(url);
    
    this.tabs.set(name, newPage);
    this.metadata.set(name, metadata || {});
    this.activeTabName = name;
    
    // Bring to front
    await newPage.bringToFront();
    
    return newPage;
  }

  get(name: string): Page {
    const page = this.tabs.get(name);
    if (!page) {
      throw new Error(`Tab '${name}' is not registered.`);
    }
    return page;
  }

  tryGetActive(): Page | undefined {
    return this.tabs.get(this.activeTabName);
  }
  
  getActive(): Page {
    return this.get(this.activeTabName);
  }

  getActiveName(): string {
    return this.activeTabName;
  }

  setActive(name: string): void {
    const page = this.get(name); // Throws if not found
    this.activeTabName = name;
    page.bringToFront().catch(() => {}); // Fire and forget
  }

  async close(name: string): Promise<void> {
    const page = this.tabs.get(name);
    if (page) {
      await page.close();
      this.tabs.delete(name);
      this.metadata.delete(name);
      
      if (this.activeTabName === name) {
        // Fallback to main or first available
        const first = this.tabs.keys().next().value;
        this.activeTabName = first || '';
      }
    }
  }

  async closeAll(keepMain: boolean = true): Promise<void> {
    for (const [name, page] of this.tabs) {
      if (keepMain && name === 'main') continue;
      
      try {
        await page.close();
      } catch (e) {
        console.warn(`Failed to close tab ${name}`, e);
      }
      this.tabs.delete(name);
      this.metadata.delete(name);
    }
    
    if (keepMain && this.tabs.has('main')) {
      this.activeTabName = 'main';
    } else {
        this.activeTabName = '';
    }
  }

  count(): number {
    return this.tabs.size;
  }

  async forEach(callback: (page: Page, name: string, metadata: TabMetadata) => Promise<void>): Promise<void> {
    for (const [name, page] of this.tabs) {
        await callback(page, name, this.metadata.get(name) || {});
    }
  }
}
