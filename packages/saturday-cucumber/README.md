# @orieken/saturday-cucumber

Cucumber.js integration for the Saturday automation framework. This package simplifies the setup of Playwright-based BDD tests by providing pre-configured Hooks, World, and browser management.

## Features

*   **Custom World**: Provides direct access to Playwright `page`, `browser`, and `context` in your step definitions.
*   **Lifecycle Hooks**: Automatic browser launch, context creation, and video/trace recording managed via `BeforeAll`/`AfterAll`/`Before`/`After` hooks.
*   **Report Attachments**: Automatically attaches screenshots and videos to Cucumber reports on failure.

## Installation

```bash
pnpm add @orieken/saturday-cucumber
```

## Usage

### 1. Initialize Hooks

In your `features/support/init.ts` (or similar setup file):

```typescript
import { installSaturdayHooks } from '@orieken/saturday-cucumber';

// Installs Before/After hooks for browser management
installSaturdayHooks();
```

### 2. Use SaturdayWorld in Steps

```typescript
import { Given } from '@cucumber/cucumber';
import { SaturdayWorld } from '@orieken/saturday-cucumber';

Given('I open the homepage', async function(this: SaturdayWorld) {
  // Access Playwright Page object directly
  await this.page.goto('https://example.com');
});
```
