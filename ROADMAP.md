# 🗺️ Saturday Monorepo Roadmap

This document outlines the high-level goals and next steps for the Saturday ecosystem. We are focused on building a fully integrated, AI-driven testing platform.

## 🎯 Strategic Focus

1.  **AI Integration**: Expand the capabilities of the Saturday MCP server to provide deeper insights, automated refactoring, and intelligent code generation.
2.  **Observability First**: Ensure every test execution (Playwright, Cucumber, K6) emits rich OpenTelemetry data for full traceability.
3.  **Developer Experience**: Unify the CLI experience and provide a robust Console for managing the testing lifecycle.

---

## 🏗️ Project Roadmaps

### 🤖 Saturday MCP Server
*Enabling AI agents to control the testing framework.*
- [ ] **Framework Analyzer (TODO-008)**: Implement deep analysis of existing project structures to understand page objects and flows.
- [ ] **Analysis Validation**: Add tools to validate adherence to framework patterns.
- [ ] **Resource Providers**: Expose documentation and templates as MCP resources.
- [ ] **Prompt Library**: Add `explain_framework` and `plan_feature` prompts.

### 🎛️ Saturday Console
*The central nervous system for test orchestration.*
- [ ] **REST API**: Implement the core API for run management and reporting.
- [ ] **Registry**: Build an efficient in-memory registry for Cucumber scenarios.
- [ ] **Runner Integration**: Wire up the logic to execute `cucumber-js` processes dynamically.
- [ ] **Persistence**: Add database support (PostgreSQL) for historical run data.
- [ ] **Live Logs**: Implement WebSockets for real-time log streaming from runners.

### 📊 Observability & Metrics
*Making tests visible.*
- [ ] **Jest/Vitest OTel**: Complete OpenTelemetry instrumentation for unit test runners.
- [ ] **Dashboarding**: Create standard Grafana dashboards for test metrics (Flakiness, Duration, Pass Rate).
- [ ] **Trace Linking**: Ensure end-to-end trace propagation between frontend user actions and backend API calls during tests.

### 🎭 Playwright & Cucumber Integration
- [ ] **Heatmap Enhancements**: Improve the ML analyzer for cluster detection in click maps.
- [ ] **Smart Selectors**: Use ML to suggest more robust selectors based on historical stability.

### ⚡ Performance Engineering (k6)
- [ ] **Advanced Redaction**: Add support for redaction of complex body structures (JSON traversal) and JWTs.
- [ ] **Custom Headers**: Allow configuration of custom enterprise header patterns for redaction.
- [ ] **Scenario Converter**: Improve the transpilation of complex Playwright logic into k6 scenarios.

---

## 💡 Wishlist & Ideas

- **IDE Plugins**: Native VS Code / IntelliJ extensions that bundle the MCP server.
- **Self-Healing Tests**: A feedback loop where failed tests are analyzed by the MCP server, fixed, and re-run automatically.
- **Distributed Runner Grid**: A k8s-native scalable grid for massive parallel test execution.

---

*This roadmap is a living document. Priorities may shift as we learn more from using the framework.*
