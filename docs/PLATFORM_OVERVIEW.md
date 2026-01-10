# Saturday Framework Platform Overview

**The AI-Native, Site-Centric Automation Framework for Modern Web Applications.**

---

## 🚀 The Elevator Pitch

Saturday is a next-generation test automation ecosystem designed to bridge the gap between traditional E2E testing and modern AI capabilities. Unlike fragile script-based approaches, Saturday treats your application as a structured "Site" composed of intelligent Pages, Flows, and Elements. It comes out-of-the-box with **Machine Learning integrations** for visual self-healing and anomaly detection, **multi-tab/multi-site management** for complex workflows, and **OpenTelemetry observability**. It's not just about checking if a button exists; it's about verifying the holistic user experience with the power of an agentic architecture.

---

## 🏛️ Core Pillars

### 1. **Site-Centric Architecture**
Move beyond loose collections of page objects. Saturday introduces the `BaseSite` facade pattern, providing a single, unified entry point for your application's Pages, Business Flows, and AI capabilities. 
- **Lazy-Loaded Components**: Pages and Flows are instantiated only when needed.
- **State-Aware Guards**: Use decorators like `@RequiresFilter('isLoggedIn')` to prevent tests from interacting with elements in the wrong state.

### 2. **AI & Machine Learning Native**
Saturday isn't just "AI-compatible"; it has ML built into its core.
- **Visual Baselines**: Establish "Golden Masters" of your UI with a single command.
- **Anomaly Detection**: Automatically detect visual regressions or unexpected layout shifts using `site.detectors`.
- **Self-Healing (Roadmap)**: Intelligent element location that adapts to UI changes.

### 3. **Structured BDD Integration**
Built on top of Cucumber.js and Playwright, Saturday provides a robust BDD environment from day one.
- **SaturdayWorld**: A pre-configured Cucumber World giving you instant access to `page`, `browser`, `siteManager`, and `tabManager`.
- **Manager Ecosystem**: Effortlessly handle multi-tab (popups) and multi-site (SSO, cross-app) scenarios without flaky window management code.

### 4. **Deep Observability**
Don't just guess why a test failed. Saturday integrates with OpenTelemetry and provides rich artifacts.
- **Heatmaps**: Automatically generate click/interaction heatmaps of your test runs.
- **Tracing**: Full integration with OTel reporters for distributed tracing of your test execution flow.
- **Rich Reports**: Automatic attachment of screenshots, video, logs, and active tab states on failure.

---

## 📦 The Ecosystem

### Core Framework
| Package | Purpose |
|---------|---------|
| **`@orieken/saturday-core`** | The heart of the framework. Contains `BaseSite`, `BasePage`, ML Facades, and Decorators. |
| **`@orieken/saturday-cucumber`** | The BDD glue. Provides `SaturdayWorld`, Hooks, `SiteManager`, and `TabManager`. |

### AI & Tools
| Package | Purpose |
|---------|---------|
| **`saturday-mcp`** | **AI Agent Server**. Enables LLMs (like Claude/Antigravity) to generate and analyze Saturday code instantly. |
| **`@orieken/saturday-ml-analyzer`** | Machine Learning libraries for anomaly detection and visual regression analysis. |
| **`@orieken/saturday-cucumber-indexer`** | Static analysis tool that indexes Gherkin feature files for metadata and semantic search. |

### Observability
| Package | Purpose |
|---------|---------|
| **`@orieken/saturday-playwright-heatmap`** | Visual analytics tool that overlays click/interaction heatmaps on your application. |
| **`@orieken/saturday-playwright-otel-reporter`** | Standardized OpenTelemetry reporter ensuring every test run is traceable. |
| **`@orieken/saturday-playwright-k6-exporter`** | Converts Playwright execution traces into k6 load testing scripts. |

### Applications
| App | Purpose |
|-----|---------|
| **`@orieken/cartridge`** | The Saturday UI Component System and Reference Implementation (Vue 3 + Tailwind). |
| **`saturday-console`** | Centralized Go-based dashboard for managing test runs and viewing aggregated results. |
| **`ye-olde-magic-shop`** | The official demo e-commerce application used for validating framework features. |

---

## 💡 Why Saturday?

*   **For Developers**: "It feels like writing application code, not standard test scripts."
*   **For QA Engineers**: "I can handle complex multi-tab flows and visual validation without importing five different libraries."
*   **For Managers**: "We have better visibility into failures and a framework that promotes reuse and maintainability."

---

## 🔗 Key Capabilities

### Multi-Context Management
Switching between an Admin Panel and a User Portal in the same test?
```typescript
// Register environments
this.siteManager.register('admin', new AdminSite(page));
this.siteManager.register('store', new StoreSite(page));

// Switch context seamlessly
await this.siteManager.get('admin').loginV2.execute();
```

### Visual Intelligence
Validate a page looks correct without writing 50 assertions.
```typescript
await site.validatePageAgainstBaseline('checkout_page');
```

---

*Saturday: Automation for the rest of us.*
