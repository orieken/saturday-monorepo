# @orieken/saturday-core

Core abstractions and utilities for the Saturday automation framework. This package provides the foundational classes for building scalable, site-centric test automation architectures and integrating advanced Machine Learning capabilities.

## Architecture

The framework follows a **Site-Centric** architecture pattern:

1.  **Site**: The entry point (Facade) for your application. It holds references to Pages, Flows, and ML tools.
2.  **Pages**: represent specific views or pages in your app.
3.  **Elements**: represent interactive parts of a page (buttons, inputs).
4.  **Flows**: represent complex business logic spanning multiple pages (e.g., "Login Flow", "Checkout Flow").
5.  **Filters**: represent state-dependent conditions (e.g., "User is logged in").

## Installation

```bash
pnpm add @orieken/saturday-core
```

## Usage

### 1. Creating a Site

Extend `BaseSite` to create your application's entry point.

```typescript
import { BaseSite } from '@orieken/saturday-core';
import { HomePage } from './pages/home.page';
import { LoginFlow } from './flows/login.flow';

export class MySite extends BaseSite {
  // Public properties for easy access
  public homePage: HomePage;
  public loginFlow: LoginFlow;

  constructor(page: Page) {
    super(page, 'https://mysite.com'); // Base URL is optional but recommended
    
    // Initialize components
    this.homePage = new HomePage(page, this);
    this.loginFlow = new LoginFlow(this);
  }

  // Optional: Override lifecycle methods
  protected initializePages(): void {
    // efficient lazy-loading can be implemented here if needed
  }
}
```

### 2. Creating a Page

Extend `BasePage` and register your elements.

```typescript
import { BasePage, InputElement, ButtonElement } from '@orieken/saturday-core';

export class LoginPage extends BasePage {
  public usernameInput: InputElement;
  public passwordInput: InputElement;
  public submitButton: ButtonElement;

  constructor(page: Page, site: BaseSite) {
    super(page, site);
    
    this.usernameInput = new InputElement(page, '#username');
    this.passwordInput = new InputElement(page, '#password');
    this.submitButton = new ButtonElement(page, 'button[type="submit"]');
  }

  async visit() {
    await this.page.goto('/login');
    await this.waitForLoad();
  }
}
```

### 3. Creating Elements

Use built-in elements (`ButtonElement`, `InputElement`, `LinkElement`) or creating custom ones by extending `BaseElement`.

```typescript
import { BaseElement } from '@orieken/saturday-core';

export class DropdownElement extends BaseElement {
  async select(value: string) {
    await this.click();
    await this.page.click(`text=${value}`);
  }
}
```

### 4. Creating Flows

Flows encapsulate multi-page business logic.

```typescript
import { BaseFlow } from '@orieken/saturday-core';

export class LoginFlow extends BaseFlow {
  async execute(user: string, pass: string) {
    await this.site.loginPage.visit();
    await this.site.loginPage.usernameInput.fill(user);
    await this.site.loginPage.passwordInput.fill(pass);
    await this.site.loginPage.submitButton.click();
  }
}
```

### 5. Machine Learning Integration

`BaseSite` provides built-in tools for AI-driven visual validation and anomaly detection.

#### Visual Baselines
Capture the current state of a page as a "golden master".

```typescript
// In your test setup or a dedicated training script
await site.establishPageBaseline('home_page_v1');
```

#### Visual Validation
Validate the current page against the trained baseline.

```typescript
const result = await site.validatePageAgainstBaseline('home_page_v1');

if (!result.isValid) {
  console.log('Anomalies detected:', result.anomalies);
}
```

#### Advanced ML Access
You can access granular ML tools via the site's facades:

```typescript
// Access trainer directly
await site.trainers.trainCurrentPage('custom_model_label');

// Access detectors
const anomalies = await site.detectors.detectAnomalies(screenshot);
```

### 6. Conditional Element Access (Filters)

Use the `@RequiresFilter` decorator to protect elements that require a specific state (like being logged in).

```typescript
import { BasePage, RequiresFilter } from '@orieken/saturday-core';

export class HomePage extends BasePage {
  
  // 1. Define the condition
  async isLoggedIn(): Promise<boolean> {
     return await this.page.isVisible('#user-avatar');
  }

  // 2. Protect the element
  @RequiresFilter('isLoggedIn')
  public accountSettings: ButtonElement;

  constructor(page: Page, site: BaseSite) {
    super(page, site);
    this.accountSettings = new ButtonElement(page, '#settings');
  }
}
```

If you try to access `accountSettings` when `isLoggedIn()` is false, a `FilterError` will be thrown.

## Utilities

### ConsoleLogger
A helper to neatly format and capture browser console logs during tests.

```typescript
import { ConsoleLogger } from '@orieken/saturday-core';

const logger = new ConsoleLogger(page);
// ... run test ...
const logs = logger.getFormattedLogs();
```
