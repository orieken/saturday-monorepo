# Saturday Monorepo - Comprehensive Project Overview

> **The AI-Native, Site-Centric Test Automation Ecosystem for Modern Web Applications**

---

## 📋 Table of Contents

- [Executive Summary](#executive-summary)
- [Vision & Mission](#vision--mission)
- [Architecture Overview](#architecture-overview)
- [Core Components](#core-components)
- [Applications](#applications)
- [Packages & Libraries](#packages--libraries)
- [Infrastructure & Deployment](#infrastructure--deployment)
- [Technology Stack](#technology-stack)
- [Development Workflow](#development-workflow)
- [Key Features & Capabilities](#key-features--capabilities)
- [Roadmap & Future Direction](#roadmap--future-direction)
- [Getting Started](#getting-started)
- [Contributing](#contributing)
- [License](#license)

---

## Executive Summary

The **Saturday Monorepo** is a unified ecosystem that revolutionizes test automation by bridging the gap between **functional testing** (Playwright/Cucumber), **performance testing** (k6), and **observability** (OpenTelemetry). It's designed to make testing smarter, faster, and AI-ready through the integration of Machine Learning capabilities and the Model Context Protocol (MCP).

### What Makes Saturday Different?

Unlike traditional test automation frameworks that rely on fragile, script-based approaches, Saturday treats your application as a structured **Site** composed of intelligent Pages, Flows, and Elements. It comes with:

- **Built-in Machine Learning** for visual self-healing and anomaly detection
- **Multi-tab/multi-site management** for complex workflows
- **OpenTelemetry observability** for full traceability
- **AI-powered code generation** via MCP server integration
- **Automatic k6 script generation** from Playwright tests
- **Site-centric architecture** promoting reusability and maintainability

---

## Vision & Mission

### Mission Statement

To bridge the gap between **functional testing** (Playwright/Cucumber), **performance testing** (k6), and **observability** (OpenTelemetry), while enabling **AI-driven development** through the Model Context Protocol (MCP).

### Core Principles

1. **AI-Native Design**: Machine Learning isn't an add-on—it's built into the core architecture
2. **Observability First**: Every test execution emits rich OpenTelemetry data for full traceability
3. **Developer Experience**: Unified CLI and robust Console for managing the entire testing lifecycle
4. **Site-Centric Architecture**: Move beyond loose collections of page objects to a structured, maintainable approach
5. **BDD Integration**: Structured Behavior-Driven Development with Cucumber.js and Playwright

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Saturday Ecosystem                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   Console    │  │  Cartridge   │  │  Saturday    │          │
│  │  (Go/K8s)    │  │   (Vue 3)    │  │  MCP Server  │          │
│  │              │  │              │  │   (Go)       │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         │                  │                  │                  │
│         └──────────────────┴──────────────────┘                  │
│                            │                                     │
│  ┌─────────────────────────┴─────────────────────────┐          │
│  │           Core Framework Layer                     │          │
│  ├────────────────────────────────────────────────────┤          │
│  │  saturday-core  │  saturday-cucumber               │          │
│  └────────────────────────────────────────────────────┘          │
│                            │                                     │
│  ┌─────────────────────────┴─────────────────────────┐          │
│  │           Observability & Tools Layer             │          │
│  ├────────────────────────────────────────────────────┤          │
│  │  • Playwright OTel Reporter                        │          │
│  │  • Cucumber OTel Formatter                         │          │
│  │  • Playwright k6 Exporter                          │          │
│  │  • Heatmap Generator                               │          │
│  │  • ML Analyzer                                     │          │
│  └────────────────────────────────────────────────────┘          │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Site-Centric Pattern

The framework follows a **Site-Centric** architecture pattern:

```
Site (Facade)
  ├── Pages (Specific Views)
  │     └── Elements (Interactive Components)
  ├── Flows (Multi-Page Business Logic)
  ├── Filters (State-Dependent Conditions)
  └── ML Tools (Trainers, Detectors, Analyzers)
```

**Benefits:**
- Single entry point for your application
- Lazy-loaded components for efficiency
- State-aware guards to prevent invalid interactions
- Built-in ML capabilities for visual validation

---

## Core Components

### 1. Saturday Framework Core

The heart of the ecosystem, providing fundamental abstractions and utilities.

#### Key Classes

- **`BaseSite`**: The facade pattern entry point for your application
  - Holds references to Pages, Flows, and ML tools
  - Provides lazy-loading and state management
  - Integrates visual baselines and anomaly detection

- **`BasePage`**: Represents specific views or pages in your application
  - Manages page-specific elements
  - Handles navigation and waiting strategies
  - Supports conditional element access via filters

- **`BaseElement`**: Interactive components (buttons, inputs, links)
  - Built-in elements: `ButtonElement`, `InputElement`, `LinkElement`
  - Extensible for custom components
  - Smart waiting and interaction handling

- **`BaseFlow`**: Multi-page business logic encapsulation
  - Orchestrates complex user journeys
  - Promotes code reuse across tests
  - Maintains clean separation of concerns

- **`Filters`**: State-dependent conditions
  - Use `@RequiresFilter` decorator to protect elements
  - Prevents invalid interactions based on application state
  - Example: Accessing admin panel only when logged in

### 2. Saturday MCP Server

A **Model Context Protocol (MCP)** server that empowers AI assistants (like Claude, GitHub Copilot, Cursor) to generate and analyze Saturday framework code.

#### Capabilities

**Code Generation:**
- Sites, Pages, Flows, Elements
- Services (API clients)
- Cucumber step definitions
- Full project scaffolding

**Analysis Tools:**
- Framework structure analysis
- Pattern compliance validation
- Performance bottleneck detection
- Dependency graph and impact analysis

**AI Workflows:**
- Visual intelligence (generate code from screenshots)
- Test prioritization based on production metrics
- Autonomous QA workflow (implement features from requirements)
- Self-healing test execution

#### Supported IDEs

- ✅ Claude Desktop (recommended)
- ✅ VS Code (with GitHub Copilot, v1.102+)
- ✅ Cursor IDE
- ✅ JetBrains IDEs (IntelliJ, PyCharm, WebStorm, v2025.2+)
- ✅ Windsurf
- ⚠️ Antigravity (likely supported)

### 3. Console (Go-based Orchestration)

The central control plane for the Saturday ecosystem, built with Go and designed for Kubernetes.

**Features:**
- REST API for run management and reporting
- Efficient in-memory registry for Cucumber scenarios
- Dynamic test runner execution
- WebSocket support for real-time log streaming
- PostgreSQL persistence for historical data

**Status:** In active development

### 4. Cartridge Dashboard

A Vue 3-based dashboard UI for visualizing test results and real-time execution status.

**Features:**
- Real-time test execution monitoring
- Build and test trend visualization
- Environment comparison views
- Failure analysis and flaky test detection
- Tag-based test result filtering

---

## Applications

### 1. **Console** (`apps/console`)
- **Purpose**: Central control plane for orchestrating test execution
- **Stack**: Go, Kubernetes
- **Status**: Active development
- **Key Features**:
  - REST API for test management
  - Scenario registry
  - Runner integration
  - Live log streaming (WebSocket)

### 2. **Saturday MCP** (`saturday-mcp`)
- **Purpose**: AI assistant integration for code generation and analysis
- **Stack**: Go, Model Context Protocol
- **Status**: Production-ready
- **Key Features**:
  - 15+ code generation tools
  - Framework analysis and validation
  - Visual intelligence (screenshot to code)
  - Test execution and self-healing
  - CLI wrapper for standalone use

### 3. **Ye Olde Magic Shop** (`apps/ye-olde-magic-shop`)
- **Purpose**: Demo e-commerce application for framework validation
- **Stack**: Vue 3, Vite
- **Features**:
  - Complete e-commerce flow (browse, cart, checkout)
  - Login/authentication
  - Used for E2E and component test validation

### 4. **Cartridge** (`apps/cartridge`)
- **Purpose**: Dashboard UI for test results and execution status
- **Stack**: Vue 3, Vite, Tailwind CSS
- **Features**:
  - Test statistics and trends
  - Build tracking
  - Environment comparisons
  - Failure analysis

### 5. **Mock API** (`apps/mock-api`)
- **Purpose**: Lightweight mock backend for demo application
- **Stack**: Node.js, Express
- **Features**:
  - Product catalog endpoints
  - User authentication simulation
  - Order processing

### 6. **Saturday Node CLI** (`apps/saturday-node-cli`)
- **Purpose**: Internal CLI for Node.js-based tasks
- **Stack**: Node.js, TypeScript
- **Features**:
  - k6 exporter management
  - Build and deployment utilities

### 7. **Saturday Go CLI** (`apps/saturday-go-cli`)
- **Purpose**: Comprehensive CLI for scaffolding and framework management
- **Stack**: Go
- **Status**: Planned

---

## Packages & Libraries

### Core Framework

#### `@orieken/saturday-core`
**The foundation of the framework**

- Provides `BaseSite`, `BasePage`, `BaseElement`, `BaseFlow`
- ML facades for visual validation and anomaly detection
- Decorators for state-dependent element access
- Console logger for browser log capture

**Installation:**
```bash
pnpm add @orieken/saturday-core
```

#### `@orieken/saturday-cucumber`
**Cucumber.js integration**

- Custom `SaturdayWorld` with Playwright integration
- Lifecycle hooks for browser management
- `SiteManager` for multi-site scenarios
- `TabManager` for multi-tab/popup handling
- Automatic screenshot/video attachment on failure

**Installation:**
```bash
pnpm add @orieken/saturday-cucumber
```

### Observability & Reporting

#### `@orieken/saturday-playwright-otel-reporter`
**OpenTelemetry reporter for Playwright**

- Emits OTel traces for every test
- Integrates with distributed tracing systems
- Provides test execution visibility
- Supports custom span attributes

#### `@orieken/saturday-cucumber-otel-formatter`
**OpenTelemetry formatter for Cucumber**

- Emits OTel traces and metrics for scenarios and steps
- Correlates BDD scenarios with test execution
- Provides step-level timing and status

#### `@orieken/saturday-playwright-heatmap`
**Visual analytics tool**

- Generates click/interaction heatmaps
- Overlays heatmaps on application screenshots
- Identifies high-traffic UI areas
- Helps optimize test coverage

#### `@orieken/saturday-ml-analyzer`
**Machine Learning libraries**

- Statistical analysis for heatmap generation
- Test pattern recognition
- Anomaly detection algorithms
- Visual regression analysis

### Performance Testing

#### `@orieken/saturday-playwright-k6-exporter`
**Automatic k6 script generation**

- Records Playwright API calls during tests
- Generates k6 performance test scripts
- Supports redaction policies for secrets
- Environment variable management

**Key Features:**
- Pluggable redaction policies
- Automatic `.env.apis` generation
- Custom k6 naming via `k6Name` option
- Idempotent script updates

**Installation:**
```bash
npm i -D @saturday/playwright-k6-exporter
```

#### `@orieken/saturday-k6-redaction-basic`
**Security policies for k6 scripts**

- Redacts sensitive data (tokens, passwords, API keys)
- Replaces secrets with environment variable placeholders
- Generates `.env.example` for team collaboration
- Customizable redaction rules

### Utilities

#### `@orieken/saturday-cucumber-indexer`
**CLI tool for feature file indexing**

- Indexes Cucumber feature files for metadata
- Enables semantic search
- Supports reporting and analytics
- Generates feature catalogs

---

## Infrastructure & Deployment

### Local Kubernetes Cluster

The project includes a complete local Kubernetes setup for development and testing.

**Location:** `local-cluster/`

**Components:**
- Test runner infrastructure
- Friday platform (AI-powered test analysis)
- Observability stack (Grafana, Prometheus, OpenTelemetry Collector)
- Database services (PostgreSQL, Qdrant vector DB)

### Friday Platform

An AI-powered test analysis platform that integrates with the Saturday ecosystem.

**Components:**
- **Friday Service**: Backend API for test result processing
- **Friday Dashboard**: UI for viewing analysis results
- **PostgreSQL**: Test result storage
- **Qdrant**: Vector database for semantic search
- **Ollama**: Local LLM service

**Integration:**
- Accepts Cucumber JSON reports
- Processes build information
- Provides AI-powered failure analysis
- Generates test improvement suggestions

**Deployment:**
```bash
# Build images
docker build -t friday-service:latest -f friday-platform/friday-service/Dockerfile friday-platform/friday-service
docker build -t friday-dashboard:latest -f friday-platform/friday-dashboard/Dockerfile friday-platform/friday-dashboard

# Deploy to Kubernetes
kubectl apply -f local-cluster/friday/
```

### Prototype Runner

A containerized test execution environment.

**Location:** `prototype-runner/`

**Features:**
- Isolated test execution
- Kubernetes job-based runners
- Artifact collection and storage
- Integration with Console API

---

## Technology Stack

### Languages
- **TypeScript/JavaScript**: Core framework, packages, demo apps
- **Go**: MCP server, Console, CLI tools
- **Python**: Test examples and utilities

### Frontend
- **Vue 3**: Dashboard and demo applications
- **Vite**: Build tooling
- **Tailwind CSS**: Styling (Cartridge)

### Testing
- **Playwright**: Browser automation
- **Cucumber.js**: BDD framework
- **k6**: Performance testing
- **Jest/Vitest**: Unit testing (planned OTel integration)

### Observability
- **OpenTelemetry**: Distributed tracing and metrics
- **Grafana**: Dashboards and visualization
- **Prometheus**: Metrics collection

### Infrastructure
- **Kubernetes**: Container orchestration
- **Docker**: Containerization
- **Kind**: Local Kubernetes clusters
- **PostgreSQL**: Data persistence
- **Qdrant**: Vector database for AI features

### AI/ML
- **Ollama**: Local LLM service
- **Model Context Protocol (MCP)**: AI assistant integration
- **Custom ML analyzers**: Visual regression, anomaly detection

---

## Development Workflow

### Prerequisites

- **Node.js**: 20+
- **pnpm**: Package manager
- **Go**: 1.21+ (for MCP and Console)
- **Docker**: For containerization
- **kubectl**: For Kubernetes management

### Installation

```bash
# Clone the repository
git clone https://github.com/orieken/saturday-monorepo.git
cd saturday-monorepo

# Install dependencies
pnpm install

# Build all packages
pnpm run build
```

### Common Commands

```bash
# Build all packages
pnpm run build

# Run tests
pnpm run test

# Lint code
pnpm run lint

# Format code
pnpm run format

# List all workspaces
pnpm run ls:ws
```

### Working with Specific Packages

```bash
# Start the demo app
pnpm --filter @orieken/ye-olde-magic-shop run dev

# Run E2E tests
pnpm --filter @orieken/ye-olde-magic-shop run test:e2e

# Build a specific package
pnpm --filter @orieken/saturday-core run build
```

### MCP Server Development

```bash
cd saturday-mcp

# Install Go dependencies
go mod download

# Build the server
go build -o bin/saturday-mcp ./cmd/saturday-mcp

# Run the server
./bin/saturday-mcp

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### CLI Development

```bash
cd saturday-mcp

# Build the CLI
go build -o bin/saturday ./cmd/cli

# Use the CLI
./bin/saturday generate page LoginPage --path /login
./bin/saturday analyze framework ./my-project
./bin/saturday validate ./my-project
```

### Agent Workflows

The project includes predefined workflows for common tasks:

- **`/build-playwright-k6-exporter`**: Build and verify the k6 exporter package
- **`/implement-feature`**: Autonomous "QA Engineer" workflow for feature implementation
- **`/otel-metrics-pipeline`**: Understanding metrics flow from code to Grafana
- **`/self-heal-test`**: Self-healing test failure workflow using Saturday MCP

**Location:** `.agent/workflows/`

---

## Key Features & Capabilities

### 1. Multi-Context Management

Seamlessly switch between different application contexts (e.g., Admin Panel and User Portal) in the same test.

```typescript
// Register environments
this.siteManager.register('admin', new AdminSite(page));
this.siteManager.register('store', new StoreSite(page));

// Switch context seamlessly
await this.siteManager.get('admin').loginV2.execute();
await this.siteManager.get('store').homePage.visit();
```

### 2. Visual Intelligence

Validate pages against visual baselines without writing dozens of assertions.

```typescript
// Establish baseline
await site.establishPageBaseline('checkout_page_v1');

// Validate against baseline
const result = await site.validatePageAgainstBaseline('checkout_page_v1');

if (!result.isValid) {
  console.log('Anomalies detected:', result.anomalies);
}
```

### 3. Automatic k6 Script Generation

Record Playwright API calls and generate k6 performance test scripts automatically.

```typescript
import { test } from '@saturday/playwright-k6-exporter/fixture';
import { createK6Recorder } from '@saturday/playwright-k6-exporter';

test('users flow @k6', async ({}, testInfo) => {
  const setup = await createK6Recorder(testInfo.title);
  const { ctx, recorder } = setup!;

  const res = await ctx.get('https://api.example.com/users', { 
    k6Name: 'list users' 
  });
  expect(res.status()).toBe(200);

  await recorder.flushToK6();
  await ctx.dispose();
});
```

### 4. State-Dependent Element Access

Protect elements that require specific application states using filters.

```typescript
export class HomePage extends BasePage {
  // Define the condition
  async isLoggedIn(): Promise<boolean> {
    return await this.page.isVisible('#user-avatar');
  }

  // Protect the element
  @RequiresFilter('isLoggedIn')
  public accountSettings: ButtonElement;
}
```

### 5. OpenTelemetry Integration

Every test execution emits rich OpenTelemetry data for full traceability.

- Playwright tests emit traces via custom reporter
- Cucumber scenarios emit traces via custom formatter
- Distributed tracing across test execution
- Integration with Grafana and Prometheus

### 6. AI-Powered Code Generation

Use the Saturday MCP server with AI assistants to generate framework code.

**Available Tools:**
- `generate_site`: Create Site classes
- `generate_page`: Create Page classes
- `generate_flow`: Create Flow classes
- `generate_steps`: Create Cucumber step definitions
- `generate_element`: Create custom Element classes
- `generate_service`: Create API Service classes
- `analyze_framework`: Analyze project structure
- `validate_patterns`: Validate framework patterns
- `suggest_improvements`: Get improvement suggestions
- `run_tests`: Execute tests and capture output
- `prioritize_tests`: Rank tests by production usage

**Available Prompts:**
- `plan_feature`: Plan feature implementation
- `explain_framework`: Explain Saturday concepts
- `debug_error`: Analyze test failures
- `generate_gherkin`: Generate BDD scenarios
- `visual_page_object`: Generate code from screenshots
- `implement_feature`: Autonomous QA workflow

### 7. Multi-Tab and Multi-Site Testing

Handle complex scenarios involving multiple browser tabs or different applications.

**Multi-Tab Example:**
```typescript
// Open a popup
await this.tabManager.openTab('popup', 'https://example.com/popup');

// Switch to popup
await this.tabManager.switchTo('popup');

// Interact with popup
await this.page.click('#confirm');

// Close popup and return to main tab
await this.tabManager.closeTab('popup');
```

**Multi-Site Example:**
```typescript
// Register multiple sites
this.siteManager.register('app1', new App1Site(page));
this.siteManager.register('app2', new App2Site(page));

// Use different sites in the same test
await this.siteManager.get('app1').login.execute();
await this.siteManager.get('app2').dashboard.visit();
```

### 8. Heatmap Generation

Visualize test coverage and user interactions with automatic heatmap generation.

```typescript
// Heatmaps are generated automatically during test execution
// View them in the test artifacts directory
```

### 9. Self-Healing Tests (Roadmap)

Intelligent element location that adapts to UI changes, powered by ML.

---

## Roadmap & Future Direction

### Strategic Focus

1. **AI Integration**: Expand MCP server capabilities for deeper insights and automated refactoring
2. **Observability First**: Ensure every test execution emits rich OpenTelemetry data
3. **Developer Experience**: Unify CLI experience and provide robust Console

### Saturday MCP Server

- [x] Framework Analyzer
- [x] Analysis Validation
- [x] Resource Providers
- [x] Prompt Library
- [x] Knowledge Graph
- [x] Test Execution & Self-Healing
- [x] Visual Intelligence
- [x] Agentic Workflow
- [x] Observability Integration

### Saturday Console

- [ ] REST API for run management
- [ ] Scenario registry
- [ ] Runner integration
- [ ] PostgreSQL persistence
- [ ] Live log streaming (WebSocket)

### Observability & Metrics

- [ ] Jest/Vitest OTel instrumentation
- [ ] Standard Grafana dashboards
- [ ] End-to-end trace propagation

### Playwright & Cucumber Integration

- [ ] Enhanced heatmap ML analyzer
- [ ] Smart selector suggestions based on stability

### Performance Engineering (k6)

- [ ] Advanced redaction (JSON traversal, JWT)
- [ ] Custom enterprise header patterns
- [ ] Improved scenario converter

### Wishlist & Ideas

- **IDE Plugins**: Native VS Code/IntelliJ extensions
- **Self-Healing Tests**: Automated failure analysis and fixing
- **Distributed Runner Grid**: k8s-native scalable test execution

**Full Roadmap:** See [ROADMAP.md](./ROADMAP.md)

---

## Getting Started

### Quick Start

1. **Clone and Install**
   ```bash
   git clone https://github.com/orieken/saturday-monorepo.git
   cd saturday-monorepo
   pnpm install
   ```

2. **Build Packages**
   ```bash
   pnpm run build
   ```

3. **Run the Demo App**
   ```bash
   pnpm --filter @orieken/ye-olde-magic-shop run dev
   ```

4. **Run E2E Tests**
   ```bash
   pnpm --filter @orieken/ye-olde-magic-shop run test:e2e
   ```

### Using the MCP Server

1. **Build the MCP Server**
   ```bash
   cd saturday-mcp
   go build -o bin/saturday-mcp ./cmd/saturday-mcp
   ```

2. **Configure Your IDE**

   **Claude Desktop:**
   Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:
   ```json
   {
     "mcpServers": {
       "saturday": {
         "command": "/path/to/saturday-mcp/bin/saturday-mcp"
       }
     }
   }
   ```

   **VS Code:**
   Add to `.vscode/settings.json`:
   ```json
   {
     "mcp.servers": {
       "saturday": {
         "command": "/path/to/saturday-mcp/bin/saturday-mcp"
       }
     }
   }
   ```

3. **Use AI to Generate Code**
   - Ask your AI assistant to generate a Page class
   - Request framework analysis
   - Get improvement suggestions

### Creating Your First Test

1. **Create a Site**
   ```typescript
   import { BaseSite } from '@orieken/saturday-core';
   
   export class MySite extends BaseSite {
     constructor(page: Page) {
       super(page, 'https://myapp.com');
     }
   }
   ```

2. **Create a Page**
   ```typescript
   import { BasePage, InputElement, ButtonElement } from '@orieken/saturday-core';
   
   export class LoginPage extends BasePage {
     public usernameInput = new InputElement(this.page, '#username');
     public passwordInput = new InputElement(this.page, '#password');
     public submitButton = new ButtonElement(this.page, 'button[type="submit"]');
   }
   ```

3. **Create a Flow**
   ```typescript
   import { BaseFlow } from '@orieken/saturday-core';
   
   export class LoginFlow extends BaseFlow {
     async execute(user: string, pass: string) {
       await this.site.loginPage.usernameInput.fill(user);
       await this.site.loginPage.passwordInput.fill(pass);
       await this.site.loginPage.submitButton.click();
     }
   }
   ```

4. **Write a Cucumber Scenario**
   ```gherkin
   Feature: Login
     Scenario: Successful login
       Given I am on the login page
       When I enter "user@example.com" and "password123"
       Then I should see the dashboard
   ```

5. **Implement Step Definitions**
   ```typescript
   import { Given, When, Then } from '@cucumber/cucumber';
   import { SaturdayWorld } from '@orieken/saturday-cucumber';
   
   Given('I am on the login page', async function(this: SaturdayWorld) {
     await this.site.loginPage.visit();
   });
   
   When('I enter {string} and {string}', async function(
     this: SaturdayWorld, 
     username: string, 
     password: string
   ) {
     await this.site.loginFlow.execute(username, password);
   });
   ```

---

## Contributing

We welcome contributions! Here's how you can help:

1. **Check the Roadmap**: See [ROADMAP.md](./ROADMAP.md) for current priorities
2. **Review Open Issues**: Look for issues tagged with `good first issue`
3. **Follow Coding Standards**:
   - Use TypeScript for framework code
   - Follow existing patterns and conventions
   - Write tests for new features
   - Update documentation

4. **Submit Pull Requests**:
   - Fork the repository
   - Create a feature branch
   - Make your changes
   - Write/update tests
   - Submit a PR with a clear description

### Code Quality Standards

- **Cyclomatic Complexity**: < 7
- **Function Length**: Keep functions short and focused
- **Test Coverage**: Aim for high coverage
- **TDD/BDD**: Test-Driven Development is standard practice

---

## Documentation

### Core Documentation

- **[README.md](./README.md)**: Quick start and overview
- **[ROADMAP.md](./ROADMAP.md)**: Future plans and priorities
- **[PLATFORM_OVERVIEW.md](./docs/PLATFORM_OVERVIEW.md)**: Detailed platform architecture
- **[SITE-TAB-MANAGEMENT.md](./docs/SITE-TAB-MANAGEMENT.md)**: Multi-site and multi-tab testing guide

### Package Documentation

- **[saturday-core README](./packages/saturday-core/README.md)**: Core framework documentation
- **[saturday-cucumber README](./packages/saturday-cucumber/README.md)**: Cucumber integration guide
- **[saturday-playwright-k6-exporter README](./packages/saturday-playwright-k6-exporter/README.md)**: k6 exporter documentation
- **[saturday-mcp README](./saturday-mcp/README.md)**: MCP server documentation

### Integration Guides

- **[Friday Integration](./friday/INTEGRATION.md)**: Integrating with Friday platform
- **[Friday Workflow](./friday/WORKFLOW.md)**: Friday platform workflow

---

## License

**Apache License 2.0**

Copyright 2025 Oscar Rieken and the Saturday Framework Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

See [LICENSE](./LICENSE) for full details.

---

## Support & Community

- **Issues**: [GitHub Issues](https://github.com/orieken/saturday-monorepo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/orieken/saturday-monorepo/discussions)
- **Author**: Oscar Rieken (@orieken)

---

## Acknowledgments

The Saturday Framework is built on the shoulders of giants:

- **Playwright**: Modern browser automation
- **Cucumber.js**: BDD framework
- **k6**: Performance testing
- **OpenTelemetry**: Observability standard
- **Model Context Protocol**: AI assistant integration
- **Vue.js**: Progressive JavaScript framework
- **Go**: Efficient, concurrent programming

---

*Saturday: Automation for the rest of us.*

**Last Updated**: January 2026
