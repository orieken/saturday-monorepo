# ADR-002: Default to OTLP over gRPC for OTel Trace Export

Date: 2026-07-20
Status: Accepted
Deciders: Oscar Rieken
Technical Story: mcp-add retrofit Phase G — tool invocations are wrapped in an OTel span at registration; the tracer needs a real exporter for production and a no-op fallback for dev.

## Context
Phase G introduced `domain.Tracer` as an interface used by the tracing middleware in `internal/server/tracing_middleware.go`. Every OTel setup must pick a wire protocol. The Saturday-stack collectors already speak OTLP/gRPC natively, and the MCP server runs in the same operational environment as the rest of the Saturday tooling.

## Decision
Default to OTLP over gRPC via `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`. `cmd/saturday-mcp/main.go` reads `OTEL_EXPORTER_OTLP_ENDPOINT`; when unset, Handler uses a private no-op tracer. When set but the exporter setup errors during startup, the entrypoint also falls back to no-op rather than refusing to start.

## Alternatives Considered
| Option | Why Rejected |
|---|---|
| OTLP over HTTP | Works over 443 (easier through corporate firewalls) but slower per-span; not the Saturday-stack default |
| stdout exporter as default | Useful for debug but noisy in production; better as an opt-in via env |
| Direct Jaeger export | Legacy — modern integration point is the OTel collector, not the backend directly |
| No default (require explicit config) | Surprising for developers; forces every dev environment to configure endpoints even when it doesn't care about traces |

## Consequences
- **Easier**: Zero-config dev runs (noop when unset). Production wiring is one env var. Matches the rest of the Saturday collector configuration.
- **Harder**: Requires network reachability of a collector on port 4317 (default OTLP gRPC) in environments where tracing is enabled. Firewalls that only allow 443 need OTLP-HTTP instead — currently a code change (constructor swap in adapters/otel/).
- **Changed**: Adds `go.opentelemetry.io/otel`, `otel/sdk`, `otlptrace`, and `otlptracegrpc` to `go.mod`. Handler construction now takes a functional option (`WithTracer(...)`) so tests and dev inject the no-op tracer explicitly.

## Fitness Function
Judgment-only. `internal/adapters/otel/otel_tracer_test.go` covers the tracer contract with an in-memory span recorder; the exporter choice itself is a constructor argument in `main.go` and not enforced by CI.
