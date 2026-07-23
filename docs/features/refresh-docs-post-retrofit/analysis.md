# Technical Analysis: refresh-docs-post-retrofit

## 1. Feature Overview
Refresh `saturday-mcp/README.md` and `saturday-mcp/docs/architecture.md` post `mcp-add` retrofit to accurately reflect the Trinity domain model (`Tool`, `Persona`, `Workflow`), the adapter layer (`internal/adapters/`), OTel tracing middleware, declarative output schemas (`invopop/jsonschema`), and ADR cross-references.

## 2. Requirements & Acceptance Criteria
- **AC-1**: `README.md` updated with Trinity architecture, 14 tools + 2 workflows + resources + 10 prompts inventory, OTel configuration env vars, development guide, testing instructions (≥85% coverage target), and ADR links.
- **AC-2**: `docs/architecture.md` updated with layer diagram, Trinity concepts, adapter breakdown, MCP protocol layer, step-by-step extension guides, observability, and links to `mcp-add-plan.md`, `ADR-001`, and `ADR-002`.
- **AC-3**: Remove fragile specific line numbers and LOC count references (e.g. "880-LOC", "~1200-LOC") to prevent documentation rot.
- **AC-4**: Cross-references verified:
  - `docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md` linked from `README.md` and `docs/architecture.md`.
  - `docs/adrs/ADR-002-default-otlp-grpc-otel-trace-export.md` linked from `README.md` and `docs/architecture.md`.
  - `mcp-add-plan.md` linked from `README.md` and `docs/architecture.md`.

## 3. Scope & Architectural Boundaries
- **Owning Context**: Documentation (`README.md`, `docs/architecture.md`)
- **Architectural Flags**: None (documentation refresh only)
