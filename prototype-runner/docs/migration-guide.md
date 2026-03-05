Migration Guide – Consolidating Prototypes into `final/`

This guide describes a safe, incremental migration path to create the `final/` directory and move/consolidate code from the prototypes:

- `cuke_dashboard_full/`
- `cuke_playwright_dashboard_full/`
- `test_dashboard_full_with_steps/`

Read `docs/architecture.md` and `docs/project-structure.md` first.

Principles

- Keep commits small and reversible; prefer copy-then-adapt over in-place moves.
- Preserve working states at each step (backend can run, tests can run, UI can build).
- De-duplicate progressively, not up front.

Step 0 – Create Skeleton `final/`

1. Create the directory tree from `docs/project-structure.md`:
   - `final/backend/{cmd,internal/{httpserver,registry,runs,runner},data,reports}`
   - `final/frontend/src/{components,stores,canned}`
   - `final/tests/{features,steps,support,reports}`
   - `final/tools/cucumber-indexer/src`
   - `final/local-cluster`
   - `final/docs`
2. Add a `README.md` to `final/` pointing back to `docs/` for architecture and quality rules.

Step 1 – Backend (Go)

Source of truth: prefer `cuke_dashboard_full/backend` as base, cross-check with `cuke_playwright_dashboard_full/backend` for runner specifics.

1. Copy these packages to `final/backend`:
   - `cmd/server` (entrypoint)
   - `internal/httpserver` (REST API and routes)
   - `internal/registry` (in-memory index store)
   - `internal/runs` (in-memory run store)
   - `internal/runner` (Cucumber + Playwright runner)
2. Ensure the runner builds a command like:
   - `npx cucumber-js --config cucumber.mjs --name "<Scenario Name>" --format json:reports/<suite>/<runId>/cucumber.json --format html:reports/<suite>/<runId>/index.html tests/features/<file>.feature:<line>`
3. Copy `data/` and `reports/` directories; keep them git-ignored if appropriate.
4. Verify locally:
   - `cd final/backend && go run ./cmd/server`
   - Confirm API routes respond (see `docs/architecture.md`).

Step 2 – Tests (CucumberJS + Playwright)

Source of truth: prefer `cuke_playwright_dashboard_full/tests` for minimal runner compatibility; bring additional steps/examples from `test_dashboard_full_with_steps` as needed.

1. Copy `tests/` folder into `final/tests`:
   - `features/` – select representative features from the prototypes.
   - `steps/` – start minimal; add more from `test_dashboard_full_with_steps` later.
   - `support/world.ts`, `support/hooks.ts` – ensure Playwright lifecycle is correct.
   - `package.json`, `cucumber.mjs`, `tsconfig.json`, etc. (if present) – adapt paths if needed.
2. Run locally:
   - `cd final/tests && npm install && npm test`
   - Confirm HTML report is generated under `final/tests/reports/`.

Step 3 – Indexer Tool (TypeScript)

Source of truth: `cuke_dashboard_full/tools/cucumber-indexer`.

1. Copy to `final/tools/cucumber-indexer`.
2. Build and generate index:
   - `cd final/tools/cucumber-indexer && npm install && npm run build`
   - `node dist/index.js --features ../../tests/features --out ../../backend/data/cucumber_index.json`
3. POST the index to backend (optional):
   - `curl -X POST http://localhost:8080/api/cucumber/index -H 'Content-Type: application/json' --data-binary @../../backend/data/cucumber_index.json`

Step 4 – Wire Backend + Tests

1. Ensure backend’s runner executes in `final/tests` directory (either working directory change or absolute path resolution).
2. Verify run flow:
   - `POST /api/runs` with `{ "suite": "default", "scenarioId": "..." }`
   - `GET /api/runs/{id}` returns status and `reportUrl`.
   - Confirm report is written to `final/backend/reports/<suite>/<runId>/index.html`.

Step 5 – Frontend (Vue 3)

Source of truth: `cuke_dashboard_full/frontend`.

1. Copy to `final/frontend` and `npm install`.
2. Configure environment (proxy or CORS) so frontend can reach backend:
   - Option A: Vite proxy to `localhost:8080`.
   - Option B: Enable CORS in backend.
3. Implement/verify API integrations:
   - `GET /api/frameworks/cucumber/suites/{suite}/scenarios`
   - `POST /api/runs`
   - `GET /api/runs/{id}`
4. Run locally: `cd final/frontend && npm run dev`.

Step 6 – Local Dev via docker-compose

Source of truth: `cuke_dashboard_full/local-cluster`.

1. Copy `local-cluster/` into `final/local-cluster`.
2. Mount `../tests` into the backend container so the runner can execute `npx cucumber-js`.
3. Start: `cd final/local-cluster && docker compose up`.

Step 7 – De-duplication and Clean-up

1. Consolidate types/interfaces for `CucumberIndex` across backend, indexer, and tests (avoid divergence).
2. Remove duplicate features/steps once parity and coverage are achieved.
3. Update READMEs in `final/` to document commands and contribution rules.

Step 8 – Quality Gates and CI (Optional here, recommended)

1. Establish lint/format/test scripts per subproject.
2. Add minimal CI to run: Go build/test, Node lint/test, Playwright in headless mode (if feasible in CI).

Appendix – Mapping from Prototypes to Final

- Backend:
  - Base: `cuke_dashboard_full/backend`
  - Cross-check runner: `cuke_playwright_dashboard_full/backend/internal/runner`
- Frontend:
  - Base: `cuke_dashboard_full/frontend`
- Tests:
  - Base: `cuke_playwright_dashboard_full/tests`
  - Extras: `test_dashboard_full_with_steps/tests`
- Indexer:
  - Base: `cuke_dashboard_full/tools/cucumber-indexer`
- Local cluster:
  - Base: `cuke_dashboard_full/local-cluster`
