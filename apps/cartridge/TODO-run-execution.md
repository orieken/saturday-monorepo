# TODO: Run execution architecture — Docker containers & Kubernetes Jobs

Purpose
- Design a single REST-driven run execution system that can trigger test runs either by launching a Docker container (local / single-host) or by creating a Kubernetes Job (k8s cluster).
- Provide clear API, reliable log delivery to clients (modal + dashboard), stored reports, and robust lifecycle handling.

Quick plan (what I'll deliver in this doc)
- Requirements & assumptions
- High-level architecture for both Docker and Kubernetes executors
- API contract and data model changes
- Log streaming / polling design (recommended approach + fallback)
- Backend responsibilities and runner behavior
- Client changes needed (UI + stores)
- Security, ops, and cleanup considerations
- Concrete implementation task list with priorities and acceptance criteria

Assumptions
- The existing test-runner-service exposes `POST /api/runs` and `GET /api/runs/:id` (and supports listing). We will extend this API.
- Reports are produced as static HTML files that the service will host under a public (or proxied) URL, e.g., `/reports/<runId>/index.html` or S3-compatible object store.
- The environment may be either a single host with Docker socket available to the service or a Kubernetes cluster with API access (service account + permissions) from the service.
- Runs produce logs (stdout/stderr) that must be persisted and streamable to clients.

Design overview

1) API contract (expanded)

- POST /api/runs
  - Request JSON:
    {
      "framework": "cucumber",
      "suiteId": "final-cucumber-project",
      "scenarioId": "item-details.feature:9",
      "executor": "docker" | "k8s",         // optional — server can choose default
      "executorOptions": { ... },             // optional runner-specific options
      "image": "node:18-alpine",            // optional, for docker/k8s
      "command": ["npm","run","test","--","--filter=item-details"], // command to run
      "env": { "FOO": "bar" },
      "timeoutSeconds": 1800,
      "resources": { "cpu": "500m", "memory": "512Mi" }
    }
  - Response: 202 Accepted with Run object { id, status: 'running', startedAt, executor, ... }

- GET /api/runs/:id
  - Returns current Run object, including: id, status (running|passed|failed|error), startedAt, finishedAt?, reportUrl?, logs?: string[] (or metadata only), executor
  - Support query params: `?includeLogs=true&fromLine=N&limit=M`

- GET /api/runs/:id/logs
  - Returns array of log lines (or chunks), with optional `fromLine` and `tail` semantics.
  - Useful for polling or backfill after reconnect.

- GET /api/runs/:id/stream  (optional, SSE)
  - If implemented, returns `text/event-stream` emitting log lines and status updates.

- GET /reports/:runId/*
  - Serve stored HTML reports (or proxy to object store).

Data model additions (Run)
- executor?: 'docker' | 'k8s'
- executorOptions?: object
- image?: string
- command?: string[]
- env?: Record<string,string>
- logs?: string[]  // persisted accumulated lines (useful for backfill)
- reportUrl?: string
- startedAt/finishedAt/status
- exitCode?/error

2) Backend runner architecture

A. Docker executor (single host)
- The service uses the Docker socket (or Docker client library) to create and start containers:
  - Create container with given image, command, env, mounts (if needed), network, and read-only workspace.
  - Attach to stdout/stderr streams and continuously capture lines.
  - Persist log lines to a store (local FS under `/app/reports/<runId>/logs.json` or DB) and push updates to in-memory store for SSE/poll responses.
  - When the container exits, collect exit code, gather report artifacts (e.g., copy from container), place them under `/app/reports/<runId>/index.html`, set `reportUrl` and final `status` in Run.
  - Clean up container and resources.

B. Kubernetes executor (cluster)
- Service acts as a controller that posts a Kubernetes Job manifest (or custom CRD) to the cluster via the Kubernetes API:
  - Job spec uses the image/command/env and sets restartPolicy: Never, backoffLimit, resource requests/limits, TTL.
  - Job writes logs to stdout/stderr and stores the report artifact to a shared location (e.g., a PVC, NFS, or object storage like S3/MinIO). Alternatively, the Job can `curl` the service to upload the report.
  - The service watches the Job status (via API polling or informers) and polls the Pod logs to capture output. Kubernetes API provides `pods/log` endpoint which can be streamed.
  - When the Job completes, collect artifacts from the storage and update the Run record with `reportUrl`.
  - Cleanup: delete job (and pods) or configure `ttlSecondsAfterFinished` on Job to auto-clean.

Executor selection and fallback
- API allows client to request `executor` preference.
- The server chooses default executor based on configuration (e.g., `PREFERRED_EXECUTOR=auto`):
  - If k8s config present -> use k8s, else if Docker available -> docker, else error.
- Validate that the requested executor is available and respond with 400 if not.

3) Logs handling and streaming (recommendation)
- Persist logs on the server as they arrive (append-only); store either in a DB or as JSON line files under `reports/<runId>/logs.jsonl`.
- Provide two mechanisms for clients:
  - SSE endpoint (`GET /api/runs/:id/stream`) that streams log lines + status events.
  - Polling `GET /api/runs/:id/logs?fromLine=N` for clients that can't use SSE.
- Client `RunModal` behavior:
  - On open, fetch `GET /api/runs/:id?includeLogs=true` to get existing history.
  - If SSE is available, open EventSource and append incoming lines.
  - Fallback to polling every ~1s while `run.status === 'running'`.
- Persist logs in the store `runsStore.runs[runId].logs` to support new windows and reloads.

4) Client changes (test-runner-ui)
- `POST /api/runs` UI uses the expanded request schema (executor, image, command, env, timeout).
- `RunModal`:
  - Add log area (left/top) with monospace lines appended as they arrive.
  - Show embedded report iframe when `run.reportUrl` exists (right/bottom or toggle between views).
  - Use SSE when supported; fallback to polling.
- `runs` store:
  - Already persists runs in localStorage. Extend to persist `logs` field (append as new lines arrive) and restore them on load.
  - Expose `getRunLogs(runId, fromLine?)` that calls backend logs endpoint and merges into local run model.

5) Security & operational concerns
- Authentication/authorization: Only authenticated clients should be allowed to start runs and access reports/logs.
- Network: If reports are stored on a private object store, the service must proxy reports or generate signed URLs.
- Resource limits & quotas: For Docker the host must limit concurrency; for k8s use Job resource requests/limits and cluster quotas.
- Isolation: k8s Jobs provide better multi-tenant isolation than running arbitrary containers via host Docker socket.
- Cleanup: enforce TTL and storage retention policies.

6) Observability & metrics
- Track metrics: runs started, running, passed, failed, avg duration, logs emitted, SSE connections.
- Expose Prometheus metrics for runs & runner health.

7) Testing & QA
- Unit tests for runner orchestration logic (mock Docker, mock k8s API).
- Integration tests that actually run a small test image and assert the report is stored and `reportUrl` set.
- E2E test for UI: start run, open modal, see logs and embedded report.

Implementation task list (prioritized)

Priority P0 — Core run lifecycle
- [ ] Add extended Run model in server: executor, image, command, env, logs[], reportUrl, startedAt, finishedAt, exitCode, error.
  - Acceptance: `GET /api/runs/:id` returns full run object; `POST /api/runs` accepts new fields.
- [ ] Implement Docker executor (server side) to:
  - Create container, capture stdout/stderr lines and append to persisted logs.
  - Collect artifacts (copy from container) and place under `reports/<runId>/`.
  - Update run status and reportUrl when finished.
  - Acceptance: a sample run started via `POST /api/runs` completes and `reportUrl` becomes available.
- [ ] Add logs endpoints: `GET /api/runs/:id/logs?fromLine=N` and `GET /api/runs/:id/stream` (SSE) stub.
  - Acceptance: logs endpoint returns appended lines; SSE returns new lines.

+Explicit logging implementation tasks (the 3 options)
+- [ ] Implement polling-based logs endpoint and client integration (fast, P0/P2).
+  - Acceptance: `RunModal` can poll `/api/runs/:id/logs?fromLine=N` and append missing lines.
+- [ ] Implement SSE streaming on server and client (P2).
+  - Acceptance: `RunModal` opens an EventSource and receives log lines and status events in near-real-time.
+- [ ] Prototype WebSocket-based streaming/broadcast (P2, optional).
+  - Acceptance: Server broadcasts run events via WS and client subscribes to runId channels.
+
Priority P1 — Kubernetes executor & reliability
- [ ] Implement Kubernetes executor that creates Jobs and captures Pod logs (streaming or tailing), supports resource requests and timeout.
  - Acceptance: same behavior as Docker executor, but via Job. Provide config to enable/disable k8s mode.
- [ ] Implement artifact publishing for k8s (either built-in upload to service or write to object store and return URL).
- [ ] Add executor selection logic and configuration flags.
+
+Priority P1.5 — K8s demo and docs
+- [ ] Add a simple `k8s-demo` with a sample Job manifest and README showing how to run a demo job locally on a cluster (kind/minikube) and retrieve logs/reports.
+  - Acceptance: `k8s-demo/demo-job.yaml` can be applied to a cluster and the README steps show how to view logs and fetch the generated report.

Priority P2 — Client & UX
- [ ] Update `RunModal` to display logs (monospace area) and embed report when available; use SSE if available, fallback to polling.
  - Acceptance: Open modal, see historical logs, see new logs in near-real-time, and embedded report when ready.
- [ ] Expose run creation UI options (executor selector, image, command overrides) behind an "advanced" toggle.
- [ ] Persist logs in `runsStore` and ensure new windows show existing logs and resume streaming/polling.

Priority P3 — Hardening & Ops
- [ ] Add metrics (Prometheus) and health checks.
- [ ] Add RBAC and auth checks around run creation and artifact access.
- [ ] Implement retention policy and automatic cleanup for old runs and artifacts.

Suggested implementation timeline (example)
- Day 1: Extend API & Run model; implement Docker executor basic flow; implement logs endpoint and persistence; update runs store to persist logs.
- Day 2: Improve Docker executor robustness, implement artifact retrieval and reportUrl assignment; wire `RunModal` polling to show logs and report.
- Day 3: Implement SSE streaming on server; update client to use SSE with polling fallback; add tests.
- Day 4+: Implement k8s executor and production hardening (metrics, RBAC, retention).

Edge cases & failure modes
- Executor fails to start: mark run as `error` with message and optionally job/container logs.
- Container/Job never produces report: mark finished but `reportUrl` blank and notify.
- Log storage failure (disk full): mark run with `error` and include available logs in the response.
- SSE connection broken: client should reconnect and fetch missed lines by polling using `fromLine`.

Open questions
- Where to store reports long-term (local FS, shared PV, or S3)?
- Do we allow arbitrary images and commands or restrict to pre-approved images for security?
- How many concurrent runs should be allowed per server/cluster? Implement queueing?

Next immediate steps I can take for you
- Implement the server `Run` model changes and `GET /api/runs/:id/logs` + `POST /api/runs` changes (P0 items), and a basic Docker executor that runs a sample test container and writes logs + report.
- Or, if you prefer, I can implement the client-side `RunModal` polling to display logs (assuming the backend `GET /api/runs/:id?includeLogs=true` is available).

+New immediate option:
+- Create the K8s demo Job + README so you can try running a job on a local cluster and observe logs/reports. I can produce the YAML and instructions now.
+
Tell me which immediate next step you'd like me to implement and I'll start coding and testing it.
