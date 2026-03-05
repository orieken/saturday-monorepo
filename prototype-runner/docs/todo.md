Context and high-level plan

This document captures the current state of the demo test-runner project, the decisions we've made so far, what has been implemented, and a prioritized, actionable checklist for the next agent to pick up tomorrow.

Start here: the goal is to have a server-side job runner that can execute CucumberJS tests either by launching Docker containers or by creating Kubernetes Jobs. The UI should allow selecting executor (docker/k8s), opening a Run modal that streams logs (SSE), and viewing the generated HTML report.

Quick plan for the next agent
- Verify the cluster demo environment (kind) and the docker-compose demo both work locally.
- Finish operational pieces and remaining code/UI work so runs can be triggered reliably in both docker and k8s modes.

Current state (what we implemented and discovered)

Files / implementation notes
- Server-side k8s executor: implemented in
  - `final/test-runner-service/internal/runner/cucumber_runner.go`
    - This file already contains a full k8s executor implementation using client-go: it creates a Job, polls for the pod, streams pod logs into a host-side run.log, exec+tar to copy `/app/reports` from the job pod back into the service's `/app/reports/<suiteId>/<runId>/`, and cleans up the Job.
- SSE streaming endpoint for run logs:
  - `GET /api/runs/{runId}/stream` implemented in `final/test-runner-service/internal/httpserver/handlers.go` (StreamRunLogs).
  - Streams JSON messages for log lines and final status.
- Run creation API:
  - `POST /api/runs` is implemented; it accepts RunRequest (framework, suiteId, scenarioId, optional executor). If executor is not provided, it falls back to `DEFAULT_EXECUTOR` env or "docker".
- Cucumber index ingestion
  - `POST /api/cucumber/index` persists an index file under `data/<suite>-index.json` and registers it with the in-memory registry.
- Helper / operational assets added during the debugging session:
  - `final/local-cluster/test-runner-rbac.yaml` — Role & RoleBinding granting the `default` ServiceAccount in namespace `test-runner` permissions to create/manage Jobs and read pods/logs. (Temporary; security note below.)
  - `final/scripts/kind-prepare.sh` — Script to build `cucumber-project:local`, load it into kind, and apply RBAC. Use when running the demo in kind.

What I tested (manual run checks)
- Posted runs via POST /api/runs with executor: "docker" and executor: "k8s" and observed the service behavior.
- Diagnosed two runtime issues when running inside kind:
  1. RBAC: service account lacked permission to create Jobs — fixed by applying `test-runner-rbac.yaml`.
  2. Image availability: job pods were stuck in ContainerCreating until `cucumber-project:local` was built and loaded into the kind cluster — resolved by `kind load docker-image` (helper script automates this).

Decisions made
- Use a server-side k8s executor implemented with client-go to create Jobs and stream logs back. This is cluster-native and preferred for the kind demo and future k8s deployments.
- Keep the docker executor for local docker-compose demo; UI should allow switching between executors (per-run or global default).
- Use SSE (Server-Sent Events) for live streaming run logs to the UI (implemented). The Run modal in the UI should connect to `/api/runs/{runId}/stream` to display live logs.
- For now we applied RBAC binding to the `default` service account in namespace `test-runner` to unblock development; this should be hardened later to use a dedicated ServiceAccount.

Security note
- Current RoleBinding binds to `ServiceAccount default` in `test-runner`. This is acceptable for local demos but not production. Next steps should include creating a dedicated `test-runner-sa` ServiceAccount and update the deployment to use it and restrict permissions as narrowly as possible.

Where things are incomplete (high-level)
- UI
  - Run modal: not implemented in the UI yet. The backend SSE exists; the frontend must open a modal, subscribe to SSE, and render streaming logs and status.
  - Steps presentation: currently steps are numbered and showing line numbers — we want to remove numbers and show a plain list with no decoration.
  - Run list bug: clicking "view report" opens a new window but the new window's store doesn't contain runs or logs; runs must be persisted or made retrievable by runId so the dashboard shows the run.
  - Executor dropdown: ensure the UI dropdown switches between docker and k8s and the run request sends the selected executor.
- Backend / infra
  - RBAC: create a dedicated ServiceAccount and RoleBinding; do not use default SA for security.
  - Persistent storage for reports and posted index: consider mounting a PVC to `/app/reports` and `/app/data` so posted indexes and reports survive pod restarts.
  - Streaming improvements: SSE is implemented; ensure UI reconnect/backfill works (client should fetch history via GET /api/runs/{runId}/logs before opening SSE follow).
  - Error handling / retries: job creation should return helpful errors; consider exposing job creation errors to users earlier in the UI.

Detailed prioritized TODO (actionable items for the next agent)

1) Stabilize k8s demo (priority: high)
- Create a dedicated ServiceAccount and Role (narrow permissions), update the `test-runner-service` Deployment to use it.
  - Files: `final/local-cluster/test-runner-rbac.yaml` -> split into `Role` + `RoleBinding` referencing `ServiceAccount: test-runner-sa`
  - Update deployment manifest (if any) or k8s YAML to set `serviceAccountName: test-runner-sa`.
  - Acceptance: job creation works and events show no Forbidden errors.
- Ensure `cucumber-project:local` image is built and loaded into kind in CI or local setup.
  - Use `final/scripts/kind-prepare.sh` to build and load the image. Improve script to detect cluster name automatically and accept flags.
  - Acceptance: Job pods reach Running or CrashLoopBackOff (i.e., they are scheduled and container created).

2) Persist posted indexes and reports (priority: high)
- Add a PVC in the kind demo that is mounted into the `test-runner-service` pod at `/app/reports` and `/app/data`.
  - This ensures reports and ingest index files survive restarts.
  - Update `kind` setup (local-cluster manifests) to create a PVC and mount it.
- Acceptance: after posting index and running reports then restarting the service pod, the posted index and generated reports still exist.

3) UI: Run modal + SSE streaming + UX fixes (priority: high)
- Implement `RunModal` in `final/test-runner-ui` (or equivalent frontend) that:
  - Opens when clicking Run, shows the runName, suite and scenario, executor choice.
  - Calls POST /api/runs with the selected executor.
  - After POST returns, do a GET /api/runs/{runId}/logs to fetch existing lines, render them, then open SSE to `/api/runs/{runId}/stream` to show live lines and final status event.
  - On status event (passed/failed), show a link to `reportUrl` from run object or from final SSE status. Provide a button `Open report` that opens report URL in new tab.
- Modify step rendering: remove numbering and line decorators for steps. Show plain list. (File(s): `final/test-runner-ui/src/components/...` — search for step renderer.)
- Fix run list bug: ensure runs are persisted in a store that survives page reload or the report view loads data by runId from backend. Approaches:
  - Persist runs in backend (e.g., write run store to `./data/runs.json`) and GET /api/runs to return them, or
  - Keep runs in `reports` folder and fetch by listing `/reports/<suite>/` or provide GET /api/runs endpoint to list runs (currently runs are in memory only).
- Acceptance: modal shows logs live, report link works in new window, run list persists across reloads.

4) Backend: improve run persistence and listability (priority: medium)
- Persist RunStore to disk (e.g., `data/runs.json`) so run metadata survives restarts.
- Add GET /api/runs (list all) endpoint if absent.
- Acceptance: new tab view shows runs even after UI reload or service restart.

5) Logging and UX improvements (priority: medium)
- Improve log rotation/size handling for run.log (avoid unbounded growth): rotate or limit file size and keep indexes of lines.
- Make SSE messages include structured metadata (timestamp, level) to allow nicer formatting in the UI.
- Acceptance: UI displays timestamps and supports toggling structured log view.

6) CI/automation for demo (priority: low -> medium)
- Add a `make demo-kind` or `scripts/run-kind-demo.sh` that:
  - Creates the kind cluster (using `final/local-cluster` manifests),
  - Builds and loads images (using `scripts/kind-prepare.sh`),
  - Deploys `test-runner-service`, `web-app`, `mock-api`, and `test-runner-ui` into the `test-runner` namespace.
- Acceptance: a single script stands up the demo in kind.

7) Tests and quality gates (priority: low)
- Add unit tests around RunStore (StartRun, CompleteRun, GetRun) and runner k8s creation (mock client-go) if possible.
- Add an integration smoke test that posts an index, triggers a run with executor "k8s" and verifies the report is produced in `/app/reports` (requires kind preset in CI).

Operational runbook / commands (copy/paste)

Prepare kind and images (local):

```bash
# from repo root
# build & load image + apply RBAC
final/scripts/kind-prepare.sh [kind-cluster-name]

# verify service is running
kubectl --kubeconfig=/tmp/kind-kubeconfig.yaml -n test-runner get pods -o wide
```

Ingest example cucumber index (if you have index file):

```bash
# provided helper in repo
final/scripts/post-cucumber-index.sh /tmp/cuke-index.json
# verify
curl -sS http://localhost:9001/api/cucumber/index | jq .
```

Trigger a k8s run:

```bash
curl -sS -X POST http://localhost:9001/api/runs \
  -H 'Content-Type: application/json' \
  -d '{"framework":"cucumber","suiteId":"final-cucumber-project","scenarioId":"inventory-list.feature:9","executor":"k8s"}' | jq .
```

Follow pod and logs:

```bash
# list jobs/pods for run
kubectl --kubeconfig=/tmp/kind-kubeconfig.yaml -n test-runner get jobs -l run-id=<runId>
kubectl --kubeconfig=/tmp/kind-kubeconfig.yaml -n test-runner get pods -l run-id=<runId> -w
# tail logs
kubectl --kubeconfig=/tmp/kind-kubeconfig.yaml -n test-runner logs -f <podName> -c runner
# or use SSE in the UI: GET /api/runs/{runId}/stream
```

Where to look in the codebase
- Backend runner: `final/test-runner-service/internal/runner/cucumber_runner.go`
- HTTP handlers: `final/test-runner-service/internal/httpserver/handlers.go`
- Runs data model / store: `final/test-runner-service/internal/runs/` (store.go, model.go)
- UI project: `final/test-runner-ui/` (search for step rendering, run modal components)
- Local-cluster resources and helper scripts: `final/local-cluster/` and `final/scripts/` (we added `kind-prepare.sh`)
- RBAC we applied: `final/local-cluster/test-runner-rbac.yaml`

Testing checklist (for the agent picking this up)
- [ ] Run `final/scripts/kind-prepare.sh` and confirm it builds and loads cucumber-project image into kind and applies RBAC without errors.
- [ ] POST a cucumber index and confirm GET /api/cucumber/index shows features.
- [ ] Trigger a k8s run; confirm:
  - A Job is created labelled `run-id=<id>`.
  - A Pod is created, reaches Running or completes, and writes reports to `/app/reports/<suite>/<runId>/` in the service pod.
  - SSE streaming to `/api/runs/{id}/stream` emits log messages and a final status.
- [ ] Implement and test the Run modal in the UI: logs live, and report link opens the generated report.
- [ ] Harden RBAC and switch the Deployment to a dedicated ServiceAccount.

Acceptance criteria (minimum for closing the job)
- The UI can trigger a run with executor set to "k8s" and see live logs in the Run modal via SSE.
- The backend creates a Kubernetes Job, the Job Pod runs, the service copies the report back to `/app/reports/<suite>/<runId>/index.html`, and the UI can open that report.
- Runs list and view-report behavior work when opening the dashboard in a new window.

Notes / possible refinements
- Consider streaming logs from the Job pod to the UI directly (via Kubernetes API proxy + WS) for lower latency; current approach (service streams run.log as generated by its runner) is good for now and keeps a consistent log copy on disk.
- For scale, implement queueing and rate limiting on the service if many concurrent runs are expected.
- Think about authenticated APIs if demos will be shared outside local dev; add auth later.

If you want, next I can:
- (A) Create the dedicated ServiceAccount and update deployment YAML + RoleBinding (recommended next step).
- (B) Implement the RunModal + SSE wiring in the UI and the step-list style change.
- (C) Add persistence for RunStore and create GET /api/runs (list) endpoint.

Pick one or more (A/B/C) or tell me to proceed with the entire prioritized list and I will continue working through items and test each change end-to-end.

----
Document created: `docs/todo.md` — use this as the single source of context for tomorrow's agent to continue work.

