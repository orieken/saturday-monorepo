# @orieken/saturday-playwright

Playwright integration for the Saturday automation framework.

This package provides a seamless bridge between Playwright's native test runner and the Saturday framework's core architecture by exposing Saturday's management utilities directly as Playwright fixtures.

## Features

- **Saturday Fixtures**: Automatically initializes and injects `siteManager`, `tabManager`, and `consoleLogger` into your Playwright tests.
- **State Cleanup**: Automatically isolates test state by closing open tabs and wiping manager registries after every test run.
- **Enhanced Telemetry on Failure**: If a test fails, this package automatically intercepts the failure and attaches rich context directly to the Playwright report:
  - Text dumps of all browser console logs.
  - Active tab tracking and window state.
  - Site Manager registry state.

## Installation

```bash
pnpm add -D @orieken/saturday-playwright @orieken/saturday-core @playwright/test
```

## Usage

Instead of importing `test` and `expect` from `@playwright/test`, import them from `@orieken/saturday-playwright`.

```typescript
import { test, expect } from '@orieken/saturday-playwright';

test('my automated test', async ({ page, siteManager, tabManager, consoleLogger }) => {
  // Use tab manager to open a new tab
  const newTab = await tabManager.openInNewTab('dashboard', 'https://example.com/dashboard');
  
  // Register a Saturday Core Site Model
  const mySite = new MySiteModel(newTab);
  siteManager.register('mainSite', mySite);

  await mySite.loginFlow.execute();
  await expect(newTab).toHaveTitle(/Dashboard/);

  // You can easily pull console logs explicitly if needed
  const logs = consoleLogger.getLogs();
  console.log(`Intercepted ${logs.length} browser logs`);
});
```

## Composition with Other Plugins

Because `@orieken/saturday-playwright` exports a standard extended `@playwright/test` object, it composes beautifully with other Playwright fixtures like the Saturday Heatmap tracker.

Use Playwright's `mergeTests` utility to combine custom extensions:

```typescript
import { test as base, expect } from '@orieken/saturday-playwright';
import { test as heatmapTest } from '@orieken/saturday-playwright-heatmap';
import { mergeTests } from '@playwright/test';

export const test = mergeTests(base, heatmapTest);

// 'test' now has access to { siteManager, tabManager, consoleLogger, heatmap }
```
