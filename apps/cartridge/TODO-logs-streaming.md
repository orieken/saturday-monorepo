# TODO: Logs streaming and modal improvements

Goal: provide near-real-time run logs inside the `RunModal` while a run is executing, and ensure runs/reports are available when opening the dashboard in a new window.

Current state
- `runs` are polled by the `runs` store and persisted to `localStorage` so multiple windows can see previously started runs and report links.
- `RunModal` embeds the HTML report (iframe) once `run.reportUrl` is available.
- We currently do not stream incremental logs into the modal; users must rely on the final report or the "View report" link.

Problems to solve
1. Provide live/near-live logs in the modal while a run is "running".
2. Avoid excessive polling or client/server load.
3. Ensure logs are available in new windows and across reloads when possible.

Possible solutions

Option A — Polling (existing pattern)
- Client periodically (e.g., every 1s) calls `GET /api/runs/:id/logs` or `GET /api/runs/:id` to fetch the latest logs chunk.
- Pros:
  - Simple to implement with existing REST endpoints and CORS.
  - Works behind proxies and load balancers without connection upgrades.
  - Easy to resume after reconnect or page reload.
- Cons:
  - Increased request volume (many requests for active runs).
  - Delay is bound by polling interval; shorter interval increases load.
- Implementation notes:
  - Extend runs API to return logs (or new `/logs` endpoint that supports `fromLine` query param).
  - In `RunModal`, start a short-interval timer while run.status === 'running' to fetch new lines and append to a local array.
  - Persist logs in `runsStore.runs[runId].logs` so new windows can read them from localStorage.

Option B — Server-Sent Events (SSE)
- Server exposes `GET /api/runs/:id/stream` as an SSE stream that emits log lines and final status.
- Pros:
  - Single long-lived connection per client, low overhead for many messages.
  - Works well for one-way streaming (server -> client).
  - No manual polling required on the client.
- Cons:
  - Not bi-directional (ok for logs); can't send control messages back on same connection.
  - Some corporate proxies can kill SSE connections; reconnect logic needed.
  - Needs server implementation (and reverse-proxy config) to support `text/event-stream` responses.
- Implementation notes:
  - When opening `RunModal`, create an EventSource tied to the run ID and append incoming events to visible logs.
  - On reconnect/resume, request recent history (or rely on persisted logs in the store).

Option C — WebSocket
- Server exposes a WebSocket endpoint (e.g., `/ws`) and sends messages identifying runId + payload.
- Pros:
  - Full-duplex, low-latency, robust for many kinds of events.
  - Efficient for high-frequency updates.
- Cons:
  - More complex server-side implementation and routing through proxies.
  - Need to implement authentication/authorization on the socket.
- Implementation notes:
  - On client open, subscribe to the run ID; server pushes log lines and status updates.
  - Persist logs in `runsStore` when received so other windows can access them.

Other considerations
- Persisting logs: store them on the server (recommended) and stream deltas to clients. Client-side localStorage persistence is only best-effort and limited.
- Security: if reports are behind auth, the iframe approach might need authenticated cookies or a proxy endpoint.
- Resource cleanup: ensure server stops streaming and frees resources when run completes.

Recommended immediate approach (fastest):
1. Add `logs?: string[]` to the `Run` model and expose `GET /api/runs/:id` that includes accumulated logs.
2. Implement short-interval polling (e.g., 1s) in `RunModal` while `run.status === 'running'` to fetch incremental logs using `fromLine` param.
3. Persist logs into `runsStore.runs[runId].logs` so other windows see them (already persisted via localStorage).

Recommended long-term approach:
- Implement SSE on the server for run log streaming, with proper reconnect/backfill logic and a small fallback to polling for environments where SSE is blocked.
- Optionally implement a WebSocket gateway if you need bi-directional control messages or many simultaneous subscribers.

Next tasks
- [ ] Add `logs` support to server `Run` model and API endpoints.
- [ ] Wire polling-based log fetching into `RunModal` as a first-pass UX.
- [ ] Implement SSE server and client logic as a follow-up.
- [ ] Add tests and integration checks to ensure logs appear in new windows and after reload.

