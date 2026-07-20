# ADR-001: Use invopop/jsonschema for Tool Output Schemas

Date: 2026-07-20
Status: Accepted
Deciders: Oscar Rieken
Technical Story: mcp-add retrofit Phase F (op 14) — every tool needed a declarative output schema for MCP tool discovery.

## Context
Every extracted tool in `internal/tools/` returns a typed response struct (`GenerationResult`, `DocumentationResult`, etc.) and needed a matching JSON schema surfaceable to MCP clients. Options ranged from hand-authoring schemas per tool to generating them from the Go types themselves.

## Decision
Use `github.com/invopop/jsonschema` to reflect Go response structs into `*jsonschema.Schema` values. The `domain.Tool` interface exposes `OutputSchema() *jsonschema.Schema`; `RegisterTools` marshals the result and assigns it to `mcp.Tool.RawOutputSchema` before calling `AddTool`.

## Alternatives Considered
| Option | Why Rejected |
|---|---|
| Hand-author per-tool JSON schemas | 14 tools × ~30 lines each = drift risk on every response-shape change |
| `xeipuuv/gojsonschema` | Validation only, no reflection/generation |
| OpenAPI-derived schemas | Over-engineered — no need for full OpenAPI, just per-tool output shapes |
| Continue with `map[string]interface{}` responses (pre-Phase-F state) | No machine-readable contract, no compile-time enforcement |

## Consequences
- **Easier**: The Go struct with its json tags IS the schema. Adding or removing a response field updates the advertised schema automatically. Response shape and schema cannot drift.
- **Harder**: One additional dependency (`invopop/jsonschema` + its transitive `santhosh-tekuri/jsonschema` and `dlclark/regexp2`). `go.mod`'s go directive bumped 1.23.2 → 1.24 to satisfy invopop's minimum.
- **Changed**: All typed response structs live in `internal/tools/responses.go`. Adding a new tool means adding both the tool file and a response type; the `OutputSchema()` implementation is one line of `jsonschema.Reflect(&MyResponse{})`.

## Fitness Function
Judgment-only for now. A test iterating `handler.Tools()` and asserting each tool's `OutputSchema()` returns a non-nil, non-empty schema would mechanize enforcement — not added because there are only 14 tools and the pattern is uniform. Add if the tool count grows or if drift is observed.
