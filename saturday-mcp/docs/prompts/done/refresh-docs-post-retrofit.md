# Refresh saturday-mcp docs after the mcp-add retrofit

## Target repo

`/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp` (subdirectory in parent `saturday-monorepo` git repo — commits go to parent).

## Prior context

The mcp-add retrofit (59 commits, closed out in `mcp-add-plan.md`) restructured the entire server. Two doc files still describe the **pre-retrofit** shape and are misleading to any reader:

- `README.md` — the entry-point doc a new contributor reads first
- `docs/architecture.md` — the layer/module walkthrough

Phase I of the retrofit plan called for a tech-writer pass to align these with the trinity vocabulary and adapter structure. It never happened.

## What actually shipped that the docs should reflect

- **Trinity vocabulary**: domain-layer types are `Tool`, `Persona`, `Workflow` — not `Handler`, `Manager`, `Service`. Interfaces live in `internal/domain/`.
- **`handler.go` shrunk from 880 → 93 LOC** — it's now a thin composite that owns three provider slices (tools, personas, resources) and a `domain.Tracer`. Tool registration is a loop in `registration.go`.
- **Adapter layer**: `internal/adapters/{testrunner,filesystem,metricsfile,otel}/` — every external dep behind a domain interface.
- **OTel instrumentation**: every tool invocation wraps in a span at `RegisterTools` via `tracing_middleware.go`. Configurable via `OTEL_EXPORTER_OTLP_ENDPOINT`; falls back to a private noopTracer when unset. See `ADR-002`.
- **Output schemas**: every tool has a typed response struct in `internal/tools/responses.go` with an `OutputSchema() *jsonschema.Schema` method (via `invopop/jsonschema`). Surfaced to MCP clients via `mcp.Tool.RawOutputSchema` after the `mcp-go v0.56.0` upgrade. See `ADR-001`.
- **Adding a new tool** is now: drop a file in `internal/tools/`, register in the Handler's tool slice, done. No touching `RegisterTools`.

## Scope

Update these two files:

### 1. `README.md`

Sections to add or rewrite:
- **Architecture**: trinity + adapter overview (short — a diagram or bullet list, not a novel)
- **Tool inventory**: current list of the 14 tools + 2 workflows + resources + prompts
- **Configuration**: env vars — `OTEL_EXPORTER_OTLP_ENDPOINT` (optional), any others
- **Development**: how to add a new tool (link to `docs/architecture.md`)
- **Testing**: how to run tests, coverage target (≥85%), where fixtures live

Leave existing content that's still accurate (CLI usage, install instructions).

### 2. `docs/architecture.md`

Full rewrite (or heavy refresh). Should cover:
- Layer diagram: Domain (interfaces) → Use-cases (tools + workflows) → Adapters (transport, filesystem, metrics, otel, testrunner) → Frameworks (mcp-go, cobra, main.go)
- Trinity concepts: what a Tool is, what a Persona is, what a Workflow is, when to use each
- Extension points: how to add a tool (with a step-by-step), how to add a workflow, how to add an adapter
- Cross-references to `mcp-add-plan.md` (historical retrofit) and the two ADRs

Both files should reference the two ADRs in `docs/adrs/` where relevant.

## Discipline

- **One commit per file, minimum** — `docs(mcp): refresh README for post-retrofit trinity structure` + `docs(mcp): rewrite docs/architecture.md for adapter layer`. Or combine if the changes are small enough to review together.
- Conventional Commits: `docs(mcp): ...`
- **NEVER `git add -A`** — parent monorepo has 100+ unrelated in-progress files. Stage explicit paths only.
- `git status --short` before commit to verify only the intended files are staged

## Escalation

Stop and report if:
- You find the README or architecture.md already reflects the retrofit (someone else refreshed it) — verify and note
- You discover a section that references code that no longer exists (e.g., a mention of `internal/executor/` which was deleted) — flag any such staleness
- The docs cite specific line numbers or LOC counts — those should be removed, they'll rot fast

## Report format (under 150 words)

```
STATUS: complete | stopped-at-<reason>

Commits:
  <sha> <message>
  <sha> <message>

Files updated:
  README.md — <what changed>
  docs/architecture.md — <what changed>

Cross-references verified:
  - docs/adrs/ADR-001 linked from: <where>
  - docs/adrs/ADR-002 linked from: <where>
  - mcp-add-plan.md linked from: <where>

Test suite: N/A (docs-only)
```

Go.
