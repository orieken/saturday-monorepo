
# Architecture — @saturday/playwright-k6-exporter

This package captures **Playwright API calls** during tests and generates **k6** scripts that mirror those calls.

```mermaid
flowchart TD
    A["Test code<br/>Playwright test(@k6)"] -->|uses| B(createK6Recorder)
    B --> C["Wrapped APIRequestContext<br/>Proxy for get/post/..."]
    C -->|executes| D["Real HTTP via Playwright"]
    C -->|records| E["Call Log<br/>(method, url, headers, body, status)"]
    E --> F[K6 Generator]
    F -->|emit| G["perf/k6/<test-slug>.k6.js"]
    G --> H[k6 run / Git CI perf stage]
```

## Key Components

- **K6Recorder**: collects calls and converts them to a k6 script.
- **createK6Recorder**: helper to create a wrapped `APIRequestContext` per test.
- **Fixture**: `@saturday/playwright-k6-exporter/fixture` adds a Playwright `test` with tag-aware exporting.
- **CLI**: `k6-exporter-init` scaffolds an example spec and adds helpful npm scripts.

## Data Flow

1. Test calls `createK6Recorder(testInfo.title)` → gets wrapped `APIRequestContext`.
2. All HTTP verbs (`get`, `post`, `put`, `patch`, `delete`, `head`, `fetch`) are proxied.
3. Each call is appended to the recorder's call log (normalized headers/body).
4. On `flushToK6()` (or fixture teardown), a k6 script is written:
   - `default` function reproduces calls with `http.request`.
   - Each call has a `check` assertion and small `sleep` to space events.
5. Generated scripts are idempotent (filename derived from test title slug).

## Boundaries & Responsibilities (DDD-lite)

- **Application**: test orchestration (Playwright test files).
- **Domain**: `RecordedCall`, slugging rules, idempotent write policy.
- **Infrastructure**: Playwright request context proxy, file IO for k6 output.

## Extensibility

- **Redaction**: plug a sanitizer into the recorder before `flushToK6()`.
- **Scenarios**: inject k6 `options` via per-test metadata or env.
- **Batching**: detect concurrency and replace sequences with `http.batch`.


## Redaction & Env Emission

```mermaid
flowchart LR
    P["Policy (plugin)"] -->|redact header/body| R[K6Recorder]
    R -->|collect findings| E["(.env.apis)"]
    R -->|update| X[".env.example"]
    R -->|emit| K["perf/k6/*.k6.js"]
```
- Policies replace secrets with `${ENV_VAR}` placeholders.
- The recorder writes `.env.apis` (local values) and updates `.env.example` (keys only).
- Teams map secrets via CI/Secrets Manager to env vars for k6 runtime.


## Policy Auto-Discovery

```mermaid
flowchart LR
    Env[K6_REDACTION_POLICY] --> R{resolvePolicy}
    NPM["@saturday/k6-redaction-basic"] --> R
    Opt[options.policy] --> R
    R --> Recorder
```
Resolution order:
1) `options.policy` if provided by the caller
2) `import(process.env.K6_REDACTION_POLICY)` if set
3) `import('@saturday/k6-redaction-basic')` as a sensible default
