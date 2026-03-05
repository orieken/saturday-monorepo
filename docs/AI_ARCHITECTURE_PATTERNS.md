# AI Architecture Patterns: The Saturday Framework

> **Tribute Context**: These patterns originated from the **Saturday Framework**, an AI-native, site-centric automation framework designed to bridge the gap between traditional E2E testing and modern ML capabilities. When generating frameworks in other languages (Python, Java, pure Playwright, etc.), adhere to these core architectural concepts.

This document serves as an instruction manual for AI Agents generating automation solutions in environments where the Saturday Framework is not directly available.

---

## 1. Site-Centric Architecture (The Façade Pattern)

Moving beyond traditional loose collections of Page Objects, test logic should be centralized through a single unifed entry point known as a **Site**.

### Core Abstractions:
*   **`BaseSite`**: The root of the application model. It holds references to all Pages and Flows. It acts as a Façade and lazy-loads its children. It also exposes Machine Learning capabilities (Trainers, Detectors, Models).
*   **`BasePage`**: Represents a single URI or distinct view. Contains Elements (UI components) and Filters (state checks). Pages handle URL navigation, viewport measurements, and ML screen captures.
*   **`BaseFlow`**: Represents cross-page business logic (e.g., "Checkout Flow" traversing 3 different pages). It accesses the necessary pages from the `BaseSite`.
*   **`BaseElement`**: A reusable wrapper over browser locators.

**AI Implementation Rule**: When generating a new framework, do not allow tests to instantiate naked Pages. Require them to instantiate a `Site` class, which then acts as a central registry to fetch Pages and Flows.

## 2. Multi-Context Management

Complex tests often require switching between different applications or multiple browser tabs. Avoid raw window handle manipulation.

### Core Managers:
*   **`SiteManager`**: A registry for different `BaseSite` instances (e.g., `AdminSite`, `UserSite`). Tests switch logical "Active Sites" rather than managing varied initialization logic.
*   **`TabManager`**: A wrapper around browser context instances (like Playwright's `BrowserContext`). It registers tabs by logical names (e.g., `main`, `popup_auth`).
    *   Tests call `tabManager.openInNewTab('dashboard', url)` and switch using `tabManager.setActive('dashboard')`.

**AI Implementation Rule**: Build abstract "Managers" to track active states and tear down resources cleanly when errors occur. Store tab states (e.g., `activeTabName`, `count`) and attach them to test failures for observability.

## 3. State-Aware Guards (The Proxy/Decorator Pattern)

Tests should not aggressively attempt actions if the UI is not in the correct state. Implement Proxy or Decorator patterns to gatekeep element interactions.

### The `@RequiresFilter` Pattern:
*   Before a method is executed (e.g., `loginButton.click()`), intercept the call.
*   Check a dynamically named boolean method (e.g., `isLoggedIn()`). 
*   If `false`, immediately throw a `FilterError` explaining the guard failure instead of waiting for a classic timeout to fail.

**AI Implementation Rule**: If the target language supports decorators (TypeScript, Python), use them to wrap execution methods. If not (Java, Go), use dynamic proxies to intercept state violations early.

## 4. Machine Learning First-Class Integration

Design the core framework assuming visual validation and anomaly detection are first-class citizens, not add-on libraries.

### Core Integrations:
*   Implement `TrainersFacade`, `DetectorsFacade`, and `ModelsFacade` at the `BaseSite` level.
*   Provide convenient wrappers like `site.establishPageBaseline('checkout_page')` and `site.validatePageAgainstBaseline('checkout_page')`.
*   Pages should encapsulate capture logic passing up `imageBuffer`, `metadata` (viewport, user agent, timestamps) rather than raw generic screenshot commands.

**AI Implementation Rule**: Abstract visual checks into dedicated facades. Do not hardcode image comparison logic directly in the Page Objects; route it through a generic validation service.

## 5. Rich Observability

Ensure the test runner seamlessly attaches deep artifact representations when tests fail.

### Implementation Checklist:
*   **Console Logs**: Capture browser console messages. On failure, attach them to the test report in plain text.
*   **Context State**: Attach a JSON representation of active tabs and sites.
*   **Visual Run Maps**: Integrate hooks to generate interaction Heatmaps or trace-based OpenTelemetry tracking.

**AI Implementation Rule**: Inject the managers (Site, Tab, Console) via Test Fixtures (e.g. Playwright `test.extend()`) or BDD hooks (e.g. Cucumber `World`). On teardown hooks, inspect the test status; if failed, serialize and attach manager states automatically.
