# DnD Magical Items Web App (demo)

A tiny Vue 3 + Vite app used as a demo storefront for testing the frontend/backend test runner.

Quick start:

1. `pnpm install` (at monorepo root)
2. `pnpm --filter @orieken/saturday-web-app run dev`

This runs the Vite dev server (port 8000) and the mock API (port 8001) concurrently.

## Running Tests

### End-to-End (Playwright)
Validates the app and generates k6 performance scripts via `@orieken/saturday-playwright-k6-exporter`.

```bash
pnpm --filter @orieken/saturday-web-app run test:e2e
```

### BDD (Cucumber)
Runs feature files using `@orieken/saturday-cucumber`.

```bash
pnpm --filter @orieken/saturday-web-app run test:bdd
```

## 🤖 AI-Powered Code Generation

Use the Saturday MCP Server with Claude Desktop to generate test code:

```bash
# See docs/USING_SATURDAY_MCP.md for full setup guide
```

**Quick Example:**
> "Use Saturday MCP to generate a Site class for the Magic Shop with pages: home, products, cart, checkout"

[📖 Full MCP Guide](./docs/USING_SATURDAY_MCP.md)

