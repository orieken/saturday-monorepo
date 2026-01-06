import { setWorldConstructor, World } from '@cucumber/cucumber';
import * as messages from '@cucumber/messages';
import { BrowserContext, Browser, Page, APIRequestContext } from 'playwright';
import { CustomWorldOptions } from './custom-world-options';
import { BaseSite, ConsoleLogger } from '@orieken/saturday-core';

export class CustomWorld extends World {
  context?: BrowserContext;
  feature?: messages.Pickle;
  page!: Page;
  browser?: Browser;
  request?: APIRequestContext;
  scenarioName?: string;
  consoleLogger!: ConsoleLogger;

  constructor(options: CustomWorldOptions) {
    super(options);
  }
}

// Optional: Export a method to install it, or let consumer do it
export function installSaturdayWorld() {
  setWorldConstructor(CustomWorld);
}
