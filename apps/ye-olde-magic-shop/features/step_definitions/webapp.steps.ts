import { Given, When, Then, setDefaultTimeout } from '@cucumber/cucumber';
import { expect, Page } from '@playwright/test';
import * as assert from 'assert';
import { TestWorld } from '../support/world';
import { MagicShop } from '../../lib/site/MagicShop.ts';

setDefaultTimeout(60 * 10000);

// Prefer BASE_URL from environment (injected by the test-runner when running inside Docker)
let baseUrl = process.env.BASE_URL || 'http://localhost:8000';

Given('the demo web app is running', async function (this: TestWorld) {
  this.site = new MagicShop(this.page as Page);
  await this.site.visit()
  console.log({ baseUrl });
});

Given('I am on the demo shop home page', async function (this: TestWorld) {
  if (!this.site) throw new Error('Site not initialized');
  await this.site.homePage.visit();
});

Given('I open {string}', async function (this: TestWorld, url: string) {
    if (!this.page) throw new Error('Page not initialized');
    await this.page.goto(url);
});

Then('the page title should contain {string}', async function (this: TestWorld, titlePart: string) {
    if (!this.page) throw new Error('Page not initialized');
    const title = await this.page.title();
    expect(title).toContain(titlePart);
});
