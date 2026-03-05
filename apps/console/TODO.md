Backend TODO – Final Project

This file tracks the work needed to stand up `saturday/console` by consolidating concepts from `cuke_dashboard_full` and `cuke_playwright_dashboard_full` (runner specifics) and wiring it to the future frontend and tests.

Foundations

- [ ] Create Go module and scaffolding
  - [ ] `go.mod` with module name `saturday/console`
  - [ ] `cmd/server/main.go` bootstrap with graceful shutdown
  - [ ] Basic configuration via env vars/flags (port, paths)
  - [ ] Logging setup (structured, leveled)

HTTP API

- [ ] Implement minimal REST API with CORS
  - [ ] `GET /api/frameworks/cucumber/suites/{suite}/scenarios` – serve scenarios from registry
  - [ ] `POST /api/cucumber/index` – accept and validate `CucumberIndex` JSON
  - [ ] `POST /api/runs` – body: `{ suite, scenarioId }` → start run, return run id
  - [ ] `GET /api/runs/{id}` – status + `reportUrl?`
  - [ ] Serve static reports at `/reports/<suite>/<runId>/...`
  - [ ] CORS for localhost dev; configurable origins
  - [ ] OpenAPI/Swagger description (optional but recommended)

Domain/Stores

- [ ] `internal/registry` – in-memory store of `CucumberIndex` per suite
  - [ ] Types for `CucumberIndex`, feature/scenario references
  - [ ] Validation on load (unique ids, file/line presence)
  - [ ] Optional persistence (disk) for cold start
- [ ] `internal/runs` – in-memory run store
  - [ ] Run lifecycle: queued → running → passed/failed
  - [ ] Concurrency-safe maps; background cleanup of old runs
  - [ ] Cancellation hook (optional)

Runner

- [ ] `internal/runner.CucumberRunner`
  - [ ] Resolve scenario → feature file + line via registry
  - [ ] Build command:
        `npx cucumber-js --config cucumber.mjs --name "<Scenario Name>" --format json:reports/<suite>/<runId>/cucumber.json --format html:reports/<suite>/<runId>/index.html tests/features/<file>.feature:<line>`
  - [ ] Execute with working dir set to `final/tests`
  - [ ] Stream stdout/stderr to `backend/reports/<suite>/<runId>/run.log`
  - [ ] Return exit code mapped to run status; populate `reportUrl`
  - [ ] Timeouts and max output size safeguards
  - [ ] Optional Docker-based runner variant (future)

Filesystem Layout

- [ ] Ensure directories exist at startup:
  - [ ] `backend/data/` for `cucumber_index.json`
  - [ ] `backend/reports/<suite>/<runId>/`
  - [ ] Configurable tests workspace path (default `../tests`)

Error Handling & Observability

- [ ] Consistent error responses with codes
- [ ] Structured logs with run ids/suite ids for traceability
- [ ] Basic metrics (counters for runs, durations) – optional

Security & Limits

- [ ] Validate inputs (suite, scenarioId)
- [ ] Constrain parallel runs (configurable)
- [ ] Prevent path traversal when resolving files

Indexer Compatibility

- [ ] Confirm compatibility with `final/tools/cucumber-indexer`
- [ ] Support loading index from disk and via `POST /api/cucumber/index`
- [ ] Add small sample `data/cucumber_index.json` for local demo (optional)

Testing

- [ ] Unit tests
  - [ ] Registry validation and lookups
  - [ ] Run store lifecycle
  - [ ] HTTP handlers (table-driven)
- [ ] Integration tests
  - [ ] Spin up server and hit API endpoints (without runner)
  - [ ] Runner happy-path executing a trivial scenario (requires `final/tests` deps)
- [ ] Test fixtures
  - [ ] Minimal `CucumberIndex` with one feature/scenario
  - [ ] Temporary workspace for reports in tests

Operations / DX

- [ ] Makefile with common targets: `run`, `test`, `lint`, `fmt`
- [ ] `golangci-lint` configuration (optional but recommended)
- [ ] `README.md` with setup, env vars, commands, API overview
- [ ] Example `.env` (port, paths)

Local Cluster

- [ ] `final/local-cluster/docker-compose.yml`
  - [ ] Backend service exposing 8080
  - [ ] Mount `../tests` and `./reports`
  - [ ] Optional frontend service proxying to backend

Future Enhancements (stretch)

- [ ] WebSockets or SSE for live run logs
- [ ] Persistent store (PostgreSQL) for historical runs
- [ ] Run queue with workers and retries
- [ ] Sharded report storage and retention policies

Acceptance Checklist (Definition of Done)

- [ ] `go run ./cmd/server` starts and serves documented endpoints
- [ ] `POST /api/cucumber/index` loads index; `GET scenarios` returns data
- [ ] `POST /api/runs` triggers a scenario; report written under `backend/reports/...`
- [ ] Frontend can consume endpoints (CORS/proxy ok)
- [ ] Basic unit + integration tests pass; coverage reported
