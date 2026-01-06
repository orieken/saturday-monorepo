# @orieken/saturday-core

Core abstractions and utilities for the Saturday automation framework. This package provides the foundational classes for building scalable, site-centric test automation architectures.

## Features

*   **Abstract Base Classes**: `BaseSite`, `BasePage`, `BaseFlow`, `BaseElement` for structured page object models.
*   **Site Facade Pattern**: Unified access point for Pages, Flows, Services, and ML tools.
*   **ML Integration**: Built-in facades for connecting with Machine Learning models (Visual validation, Anomaly detection).
*   **Utilities**: Common helpers like `ConsoleLogger`.

## Installation

```bash
pnpm add @orieken/saturday-core
```

## Usage

### Creating a Site

Extend `BaseSite` to create your application's entry point:

```typescript
import { BaseSite } from '@orieken/saturday-core';
import { HomePage } from './pages/home.page';

export class MySite extends BaseSite {
  public homePage: HomePage;

  constructor(page: Page) {
    super(page);
    this.homePage = new HomePage(page);
  }
}
```

### Creating a Page

Extend `BasePage` for page objects:

```typescript
import { BasePage } from '@orieken/saturday-core';

export class HomePage extends BasePage {
  constructor(page: Page) {
    super(page);
  }
  
  async visit() {
    await this.page.goto('/');
  }
}
```
