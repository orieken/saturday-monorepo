# mcp-expand Milestone 2 — Advanced Analysis, Retrieval Tier 3, Structured Generators

Plan + execute M2 of mcp-expand. M1 shipped 22 tools (16 pre-M1 + 6 M1). M2 adds 5-6 more: retrieval tier 3, migration validator, test advisors, KI/ADR/doc generators.

## Target repo

`/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp` (subdir in parent `saturday-monorepo` git repo).

## Source of truth

`saturday-mcp/mcp-expand-plan.md` — see §4 "Milestone 2" for the tool list. This handoff executes what §4 promises.

## Prior context (read first)

- `saturday-mcp/mcp-expand-plan.md` — the plan, M1 close-out section documents the pattern
- `saturday-mcp/mcp-add-plan.md` "Retrofit Complete" — the shape of a completed milestone doc
- `saturday-mcp/internal/tools/` — 22 existing tools; the shape is well-established (see `analyze_complexity_tool.go` for the reference)
- `saturday-mcp/internal/tools/retriever.go` + `bm25_retriever.go` — the retrieval adapter pattern search_features (M2's first tool) will extend
- Framework [ADR-002](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md) — corpus-aware retrieval; search_features is retrieval **tier 3 (vector via sqlite-vec)**

## Scope

### Phase A — Draft the M2 execution plan (do this FIRST, get approval, then execute)

Extend `mcp-expand-plan.md` §4 into a full §5 "Milestone 2 Execution" mirroring the shape §2 + §3 use for M1. Include:

- Tool inventory table (5-6 tools with capability, input, output/contract, key value — like §2)
- Per-op breakdown (like §3) with SHAs left blank for the executor to fill
- Milestone 2 open questions requiring user approval before code lands (e.g. sqlite-vec version choice, does search_features share DB path with search_docs, do saturday/sunday advisors need per-language variants)
- Ordered dependencies: sqlite-vec choice must resolve before search_features starts; personas out of scope (that's M3)

Commit as `docs(mcp): draft mcp-expand-plan §5 Milestone 2 execution plan`.

**Pause here. Do not write code until the user approves M2 open questions.**

### Phase B — Execute each M2 op as its own commit

M1 established the pattern: one commit per tool, per-op subagent-friendly scope. The following ops fit that shape:

- **Op M2.1 `search_features`** — vector retrieval tier 3. Uses sqlite-vec (verify pure-Go availability — modernc.org/sqlite doesn't ship the vec extension). If CGO is required, halt and re-raise; M1 committed to `CGO_ENABLED=0`. Shares `.claude/rag/` dir with docs-fts5.sqlite as `docs-vec.sqlite`. Implements `Retriever` interface. Corpus: installed project's `docs/features/` (per ADR-002).
- **Op M2.2 `validate_migrations`** — walks a project's migration directory (e.g. `db/migrations/`), flags expand/contract violations per `shared/rules/architecture-guardrails.md` #2 (DROP COLUMN, RENAME COLUMN, DROP TABLE, NOT NULL without DEFAULT). Analyzer + tool pattern from M1.
- **Op M2.3 `saturday_test_advisor`** — audits Saturday coverage (dead Page/Flow/Element/Partial types, scenarios that reference non-existent primitives). Cross-refs `shared/rules/testing-conventions.md` Saturday section. Analyzer + tool pattern.
- **Op M2.4 `sunday_test_advisor`** — audits Sunday API test coverage per `shared/rules/testing-conventions.md` Sunday section. Analyzer + tool pattern.
- **Op M2.5 `create_ki`** — generator tool taking structured input (name, tags, domain, body sections) → produces KI markdown to a temp path or writes it via `domain.FileWriter`. Must produce output matching `shared/schemas/ki-frontmatter.schema.json` from the framework repo.
- **Op M2.6 `create_adr`** — same shape for ADRs. Match the ecosystem's ADR-001 format (see `saturday-mcp/docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md` for local precedent).
- **Op M2.7 `scaffold_docs`** — comprehensive markdown implementation-guide generator. Takes a feature name + sections, produces a structured doc under `docs/features/<name>/README.md`.
- **Op M2.8 Integration + e2e** — extend `TestE2E_ToolInventory` (currently pins 22) to the new total. Add per-tool e2e like M1 Op 8.
- **Op M2.9 Docs refresh** — README + docs/architecture.md + close M2 in mcp-expand-plan.md.

## Discipline (non-negotiable — locked in by M1's execution)

- **One commit per op.** ~9 commits for M2.
- Conventional Commits: `feat(mcp): add <tool> (mcp-expand M2 Op X)`.
- **NEVER `git add -A`.** Parent monorepo has 100+ files of unrelated in-progress work. Always stage explicit paths under `saturday-mcp/`.
- `git status --short saturday-mcp/` before commit — verify only intended files.
- `go build ./... && go test ./...` must be green after every commit.
- Coverage ≥ 85% per new tool file.
- Follow the exact pattern established by M1 tools (`analyze_complexity_tool.go` is the reference shape).
- Reuse `internal/analyzers/walkutil.go`'s exported `SkipUninterestingDir` / `SkippedDirNames` for any new analyzer that walks a project tree.
- Reuse the `Retriever` interface for `search_features` — don't invent a new interface just because vector ≠ BM25.
- Reuse the `emptyResult`/`marshalToolResult` helpers in `search_ki_tool.go` for any new tool with the same nil-backend-degrades-gracefully pattern (search_features will need it).
- The `Handle*` wrappers in `internal/server/testing.go` use a `dispatch(name, request)` helper — add a new wrapper via that helper per new tool.

## Escalation criteria

Stop and report if:
- sqlite-vec requires CGO — halt, describe the tradeoff (pure-Go fallbacks + LightRAG-Python alternative per design pack `06-LightRAG-Strategy.md`)
- Any new tool would exceed the current `tools/` package's Sandi Metz limits (classes ≤ 100 lines etc.) — halt, propose the extract
- The M2 op count exceeds ~10 — halt, propose splitting M2 into M2a + M2b
- A user-approval-required decision surfaces mid-execution — halt, don't push past

## Report format (under 400 words for the whole M2 arc)

Per phase report separately. Each op reports:
```
STATUS: complete | stopped-at-<reason>
Commit: <sha> <message>
Coverage: <pct>
Test suite: all green | <details>
Notes for next op: <patterns established>
```

Milestone close-out report on Op M2.9:
```
M2 COMPLETE
Ops shipped: 9 / 9
New tool count: <n> (was 22)
Coverage snapshot: <per new tool>
Known gaps for M3: <list>
Recommended next step: draft M3 handoff
```

Go.
