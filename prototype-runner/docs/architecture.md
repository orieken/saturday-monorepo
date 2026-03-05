Target Architecture – Cucumber Scenario Dashboard (Final Project)

This document defines the target architecture we will converge on under `final/`.
It consolidates the best parts of the existing prototypes:

- `cuke_dashboard_full/` – complete slice (backend + frontend + tests + indexer + local cluster)
- `cuke_playwright_dashboard_full/` – focused on backend + runner + tests
- `test_dashboard_full_with_steps/` – variant with additional step examples and docs

Goals

- Single, coherent monorepo (`final/`) with clear module boundaries.
- Backed by a small Go HTTP API exposing Cucumber index and run management.
- Vue 3 frontend providing a scenario dashboard and one-click execution.
- CucumberJS + Playwright test workspace executed per-scenario via the backend.
- A TypeScript indexer CLI to parse `.feature` files and upload an index to the backend.
- Local docker-compose for easy end-to-end development.
- Documented quality rules, ADRs, and contribution guidelines.

High-level Components

- Backend (Go)
  - `httpserver` – REST API, minimal CORS.
  - `registry` – in-memory store of `CucumberIndex` per suite.
  - `runs` – in-memory run registry.
  - `runner` – executes `npx cucumber-js` commands with Playwright.
  - `data/` – index JSON storage; `reports/` – HTML/JSON report outputs.

- Frontend (Vue 3 + Vite + Pinia)
  - Tree view: Feature → Scenario → Steps.
  - Actions: trigger run, view run status, open HTML report.

- Tests (CucumberJS + Playwright)
  - `features/` – source `.feature` files.
  - `steps/` – step definitions.
  - `support/` – world/hooks wiring Playwright lifecycle.
  - `reports/` – test outputs.

- Tools (TypeScript)
  - `cucumber-indexer` – Parse features and emit a `CucumberIndex` JSON; optionally POST to backend.

- Local Development
  - `local-cluster` – docker-compose with backend (and optionally frontend) for a quick demo loop.

Data and Control Flow

1. Indexer parses `.feature` files → writes `CucumberIndex` JSON.
2. Backend accepts index via `POST /api/cucumber/index` or reads it from `data/`.
3. Frontend queries backend for scenarios, displays tree.
4. User triggers a run → `POST /api/runs` with `{ suite, scenarioId }`.
5. Backend resolves file + line from index, executes `npx cucumber-js` in `tests/`.
6. Runner writes logs + reports; backend tracks run status and exposes `GET /api/runs/{id}`.
7. Frontend polls or subscribes (future WS) to run status and links to `/reports/<suite>/<runId>/index.html`.

Public API (initial)

- `GET /api/frameworks/cucumber/suites/{suite}/scenarios`
- `POST /api/cucumber/index` – body: `CucumberIndex` JSON
- `POST /api/runs` – body: `{ suite, scenarioId }`
- `GET /api/runs/{id}` – returns `{ id, suite, scenarioId, status, reportUrl? }`
- `GET /reports/<suite>/<runId>/index.html` – static files

Quality Gates

- Cyclomatic complexity < 7 per function.
- Function length ≤ 30 LOC.
- Tests ≥ 85% coverage per subproject.
- TDD/BDD, SOLID, clean architecture.

Notes on Differences Between Prototypes

- Naming and folder layout differ slightly; the `final/` structure (see `project-structure.md`) is the single source of truth.
- Keep runner behavior from `cuke_playwright_dashboard_full` (simple and focused) and UI/Pinia patterns from `cuke_dashboard_full`.
- Prefer the richer step examples from `test_dashboard_full_with_steps` where applicable.
