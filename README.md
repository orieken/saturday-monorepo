# Saturday Monorepo

Welcome to the **Saturday Monorepo**. This repository creates a unified ecosystem for next-generation test automation, observability, and performance engineering. It houses the **Saturday Framework**, a collection of tools designed to make testing smarter, faster, and AI-ready.

## 🌟 Mission

To bridge the gap between **functional testing** (Playwright/Cucumber), **performance testing** (k6), and **observability** (OpenTelemetry), while enabling **AI-driven development** through the Model Context Protocol (MCP).

## 📂 Project Structure

### � Applications (`apps/*`)

| Application | Description | Stack |
|-------------|-------------|-------|
| **[Console](./apps/console)** | The central control plane for the Saturday ecosystem. Orchestrates test execution and manages runners. | Go, Kubernetes |
| **[Saturday MCP](./saturday-mcp)** | A Model Context Protocol server that empowers AI assistants to generate and analyze Saturday framework code. | Go, MCP |
| **[Ye Olde Magic Shop](./apps/ye-olde-magic-shop)** | A Vue 3 e-commerce demo application used to validate framework features (E2E + Component tests). | Vue 3, Vite |
| **[Cartridge](./apps/cartridge)** | A dashboard UI for visualizing test results and real-time execution status. | Vue 3, Vite |
| **[Mock API](./apps/mock-api)** | Lightweight Express.js server providing mock data for the demo application. | Node.js, Express |
| **[Saturday Node CLI](./apps/saturday-node-cli)** | Internal CLI tool for managing Node.js-based tasks like the k6 exporter. | Node.js, TypeScript |
| **[Saturday Go CLI](./apps/saturday-go-cli)** | (Planned) A comprehensive CLI for scaffolding and managing the Saturday framework. | Go |

### 🛠️ Core Packages (`packages/*`)

| Package | Description |
|---------|-------------|
| **[@orieken/saturday-core](./packages/saturday-core)** | fundamental abstractions (`BasePage`, `BaseFlow`, etc.) and utilities for the framework. |
| **[@orieken/saturday-cucumber](./packages/saturday-cucumber)** | Cucumber.js integration with custom World and Hooks for Playwright. |
| **[@orieken/saturday-playwright-k6-exporter](./packages/saturday-playwright-k6-exporter)** | Records Playwright API calls and automatically generates k6 performance test scripts. |
| **[@orieken/saturday-k6-redaction-basic](./packages/saturday-k6-redaction-basic)** | Security policies to redact sensitive data (tokens, passwords) from generated k6 scripts. |
| **[@orieken/saturday-playwright-otel-reporter](./packages/saturday-playwright-otel-reporter)** | Custom Playwright reporter that emits OpenTelemetry traces for every test. |
| **[@orieken/saturday-cucumber-otel-formatter](./packages/saturday-cucumber-otel-formatter)** | Cucumber formatter that emits OTel traces and metrics for scenarios and steps. |
| **[@orieken/saturday-playwright-heatmap](./packages/saturday-playwright-heatmap)** | Generates visual heatmaps of test coverage based on interaction coordinates. |
| **[@orieken/saturday-ml-analyzer](./packages/saturday-ml-analyzer)** | ML-powered statistical analysis for heatmap generation and test pattern recognition. |
| **[@orieken/saturday-cucumber-indexer](./packages/saturday-cucumber-indexer)** | CLI tool to index Cucumber feature files for metadata and reporting. |

### 🚀 Starters

- **[Saturday Cucumber Starter](./saturday-cucumber-starter)**: A boilerplate template for starting new projects with the Saturday Framework.

## ⚡ Getting Started

This repository is a **pnpm workspace**.

### Prerequisites
- Node.js 20+
- pnpm
- Go 1.21+ (for MCP and Console)

### Installation

```bash
# Install dependencies
pnpm install

# Build all packages
pnpm run build
```

### Development Workflow

1.  **Framework Development**:
    - Work in `packages/`. Changes are instantly available to apps via workspace symlinking.
    - Run unit tests: `pnpm run test`

2.  **Web App Integration**:
    - Start the demo app: `pnpm --filter @orieken/ye-olde-magic-shop run dev`
    - Run E2E tests: `pnpm --filter @orieken/ye-olde-magic-shop run test:e2e`

3.  **MCP Server**:
    - Explore the API: `cd saturday-mcp && go run ./cmd/saturday-mcp`

## Advanced Features

- **Multi-Site & Multi-Tab Testing**: The framework supports managing multiple sites and tabs within a single test scenario. See [Site and Tab Management](docs/SITE-TAB-MANAGEMENT.md) for details.

## 🤝 Contributing

We welcome contributions! Please check the [ROADMAP.md](./ROADMAP.md) to see what we are working on next.
