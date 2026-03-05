Unified Docs for the Final Project

These docs consolidate the key concepts from the prototypes in this repository and outline a clear path to a single, coherent codebase under a new `final/` directory at the repository root.

Prototypes included here:

- `cuke_dashboard_full/` – Full-stack slice: Go backend + Vue 3 frontend + CucumberJS/Playwright tests + indexer + local cluster
- `cuke_playwright_dashboard_full/` – Backend + runner + tests with a placeholder web app
- `test_dashboard_full_with_steps/` – Another full-stack slice with additional example steps and docs

Docs index:

- `architecture.md` – Target architecture and guiding principles
- `project-structure.md` – Target directory structure (with Mermaid diagram)
- `migration-guide.md` – Step-by-step plan to create `final/` and migrate code incrementally
- `adrs/` – Architecture Decision Records (start with ADR-0001)

If you’re starting now, begin with `migration-guide.md`.

Notes / TODOs:

- We plan to adopt the `@saturday` npm scope for published packages. Initial package names:
  - `@saturday/test-runner-ui` (frontend/dashboard)
  - `@saturday/test-runner-service` (backend/api)
  - `@saturday/cucumber-indexer` (indexer CLI)
  - `@saturday/demo-webapp` (demo storefront used for tests)
- Port mapping for local development (defaults):
  - `test-runner-service` (Go backend) — http://localhost:9001
  - `test-runner-ui` (frontend/UI) — http://localhost:9000 (proxies /api to 9001)
  - `web-app` (demo storefront) — http://localhost:8000 (proxies /api to mock-api:8001)
  - `mock-api` (demo items/orders) — http://localhost:8001

Docker

You can build and run all services together using Docker Compose from the `final/` directory:

```bash
cd final
docker-compose up --build
```

This will build images for `test-runner-service`, `test-runner-ui`, `mock-api`, and `web-app`, and expose them on the ports listed above.

- TODO: Add `@saturday/jest-indexer` and `@saturday/playwright-indexer` to support other test frameworks. These will be small CLIs that create or translate existing test metadata into the repository's index format.
