# Implementation Notes: refresh-docs-post-retrofit

## 1. Summary of Changes
Updated `saturday-mcp/README.md` and `saturday-mcp/docs/architecture.md` to accurately document the post-retrofit Trinity architecture (`Tool`, `Persona`, `Workflow`), Clean Architecture adapter layer, OTel tracing middleware, declarative JSON schemas (`invopop/jsonschema`), and cross-references to `mcp-add-plan.md`, `ADR-001`, and `ADR-002`. Removed fragile line numbers and LOC counts to prevent doc rot.

## 2. Files Modified
- `[MODIFY] saturday-mcp/README.md` -- Architecture summary, Trinity vocabulary, tool/workflow/persona inventory, OTel env vars, development guide, and ADR links.
- `[MODIFY] saturday-mcp/docs/architecture.md` -- Full architecture walkthrough, layer diagram, Trinity concepts, adapter overview, MCP protocol layer details, extension guides, observability, and ADR links.

## 3. Self-Review Checklist
- [x] Trinity domain concepts correctly explained.
- [x] All 14 tools + 2 workflows + resources + 10 personas accurately cataloged.
- [x] Cross-references to `docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md`, `docs/adrs/ADR-002-default-otlp-grpc-otel-trace-export.md`, and `mcp-add-plan.md` present and verified.
- [x] Specific fragile LOC counts and line numbers removed.
