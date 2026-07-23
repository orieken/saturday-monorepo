# MCP Expand Plan: Framework Tooling Expansion for saturday-mcp

**Server path**: `/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp`
**Language**: Go 1.24
**MCP SDK**: `github.com/mark3labs/mcp-go v0.56.0`
**Transport**: stdio (`server.ServeStdio`)
**Sub-flow**: `expand` — Framework-wide tooling expansion
**Date**: 2026-07-22

---

## 1. Executive Summary & Strategy

Following the completion of the `mcp-add` retrofit (59 commits), `saturday-mcp` has been established as a layered Clean Architecture Go server utilizing the Trinity domain model (`Tool`, `Persona`, `Workflow`), `domain.Tracer` middleware wrapping every tool invocation, and typed `invopop/jsonschema` output contracts ([ADR-001](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-001-use-invopop-jsonschema-tool-output-schemas.md), [ADR-002](file:///Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp/docs/adrs/ADR-002-default-otlp-grpc-otel-trace-export.md)).

**Chosen Scope**: **Broad Framework-MCP**. `saturday-mcp` will evolve into the reference MCP server for the entire context-engineering framework (`ai-assistant-dot-files`), surfacing analytical verification, Knowledge Item search, framework rule validation, personas, and workflows across any language/codebase.

### Recently-landed framework context (2026-07-22)

Two updates in `ai-assistant-dot-files` materially affect this expansion — the plan honors both:

1. **[ADR-002: Corpus-Aware Retrieval Strategy (Graduated RAG)](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md)** — a framework-level ADR (distinct from saturday-mcp's own ADR-002 on OTLP export). It prescribes a graduated retrieval strategy with three implementations behind one `Retrieve(query, corpus) → []Reference` adapter:

    | Corpus | Retrieval | Milestone in this plan |
    |---|---|---|
    | Framework KIs + ADRs (`shared/knowledge/`, `docs/adrs/`) | **LLM-as-retriever** — corpus fits in context; no index maintenance | `search_ki` (M1) |
    | Installed project's `docs/` (features, ADRs, KIs) | **BM25 via sqlite-fts5** — prose retrieval, mature, no ML | `search_docs` (M2, or promotion candidate for M1 — see §5) |
    | Installed project's feature-archive semantic similarity | **Vector via sqlite-vec** — "have we built this before?" is inherently semantic | `search_features` (M2) |
    | Installed project's source code | Lean on client (Grep/Glob/Read); optional vector backend later | Deferred indefinitely |

    ADR-002 explicitly names saturday-mcp as the transport for all four retrieval capabilities and calls the BM25 `search_docs` tool "the highest-leverage single addition" that "can ship in Milestone 1 of mcp-expand." That recommendation is surfaced as an open question in §5 rather than silently absorbed here.

2. **Frontmatter contracts + JSON schemas** landed in `shared/schemas/` (`agent-frontmatter.schema.json`, `skill-frontmatter.schema.json`, `ki-frontmatter.schema.json`) with paired contracts in `shared/contracts/`. Milestone 1 of this plan is Go-only and does not introduce agent or skill markdown files, so no immediate action is required. Milestone 3's personas — once implemented — SHOULD validate their exported metadata against these schemas so the framework-MCP's personas can be discovered by the same tooling that reads local agent files. Recorded here so it isn't lost when M3 planning begins.

---

## 2. Milestone 1: Core Analytical & Search Tools

Milestone 1 focuses on 6 high-leverage tools: 4 pure analytical tools (input → structured output, zero side effects), 1 LLM-as-retriever query tool over the framework corpus, and 1 BM25-index retrieval tool over the installed project's docs (per framework ADR-002's recommendation that this is "the highest-leverage single addition"):

| Tool | Capability Category | Input | Output / Contract | Key Value |
|---|---|---|---|---|
| `analyze_complexity` | Analytical | `projectPath`, `maxComplexity` (default 7), `maxLines` (default 30) | `ComplexityAnalysisResult` (file/func violations breakdown) | Enforces Beck & Sandi Metz clean code thresholds |
| `check_accessibility` | Analytical | `filePath` or `projectPath` | `AccessibilityReportResult` (violations list with line & element) | Checks semantic HTML, ARIA, and missing labels |
| `check_ubiquitous_language` | Analytical | `projectPath`, `dictionaryPath` | `UbiquitousLanguageResult` (unapproved synonyms & violations) | Ensures code matches `DOMAIN_DICTIONARY.md` terms |
| `verify_dependencies` | Analytical | `projectPath` | `DependencyVerificationResult` (Clean Arch boundary breaches) | Prevents inner layers from importing outer layers |
| `search_ki` | Retrieval tier 1 — LLM-as-retriever (framework corpus) | `query`, `tags`, `domain` | `KISearchResult` (ranked KI & ADR matches with paths, references only — never content copies) | Instant cross-project Knowledge Item discovery over the framework corpus |
| `search_docs` | Retrieval tier 2 — BM25 (installed-project corpus) | `query`, `docsPath` (default `docs/`), optional `section` filter | `DocSearchResult` (ranked doc matches with paths + snippet, references only) | Agent + human queries against the actual installed project's `docs/features/`, `docs/adrs/`, `docs/patterns/`, `docs/runbooks/` |

---

## 3. Operations Breakdown for Milestone 1

Each operation is one discrete commit adhering to conventional commits (`feat(mcp): ...`) and strict git hygiene.

### Phase A — Planning & Open Questions

- **Op 1**: Write and commit `mcp-expand-plan.md` planning artifact (`docs(mcp): add mcp-expand-plan.md for framework expansion`).

### Phase B — Analytical Tools Implementation

- **Op 2 (`analyze_complexity`)**: Create `internal/tools/analyze_complexity_tool.go` implementing `domain.Tool`, add response struct `ComplexityAnalysisResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/analyze_complexity_tool_test.go`.
- **Op 3 (`check_accessibility`)**: Create `internal/tools/check_accessibility_tool.go` implementing `domain.Tool`, add response struct `AccessibilityReportResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/check_accessibility_tool_test.go`.
- **Op 4 (`check_ubiquitous_language`)**: Create `internal/tools/check_ubiquitous_language_tool.go` implementing `domain.Tool`, add response struct `UbiquitousLanguageResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/check_ubiquitous_language_tool_test.go`.
- **Op 5 (`verify_dependencies`)**: Create `internal/tools/verify_dependencies_tool.go` implementing `domain.Tool`, add response struct `DependencyVerificationResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/verify_dependencies_tool_test.go`.

### Phase C — Search & Retrieval Tools

- **Op 6 (`search_ki`)**: Create `internal/tools/search_ki_tool.go` implementing `domain.Tool`, add response struct `KISearchResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/search_ki_tool_test.go`. Retrieval strategy follows framework [ADR-002](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md): **LLM-as-retriever** over the framework corpus (KIs + ADRs) — no index, no embedding, corpus fits in context. Adapter interface: `Retrieve(query, corpus) → []Reference`. Returns pointers to canonical markdown, never content copies.
- **Op 7 (`search_docs`)**: Create `internal/tools/search_docs_tool.go` implementing `domain.Tool`, add response struct `DocSearchResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/search_docs_tool_test.go`. Retrieval strategy per framework [ADR-002](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md): **BM25 via sqlite-fts5** over the installed project's `docs/features/`, `docs/adrs/`, `docs/patterns/`, `docs/runbooks/`. Sub-ops:
    - Add `github.com/mattn/go-sqlite3` (or comparable fts5-enabled driver) to `go.mod`
    - Introduce a `shared/rag/` adapter interface locally in the tool for now (`Retrieve(query, corpus) → []Reference`); when the framework's `shared/rag/` layer lands in AOS Phase 3, refactor to import it
    - Index lives at `.claude/rag/docs-fts5.sqlite` per installed project — created lazily on first `search_docs` call if absent
    - Add `/reindex-docs` skill spec in `internal/tools/search_docs_tool.go`'s package docs for manual rebuild after doc changes (skill implementation deferred to a follow-up; tool works without it via lazy first-call build)
    - Tests must use a temp directory + fixture docs; no reliance on the real installed project's `docs/`

### Phase D — Integration & Documentation

- **Op 8 (Integration & E2E)**: Update `internal/integration/e2e_test.go` and `internal/server/testing.go` to cover all 6 new tools.
- **Op 9 (Documentation Refresh)**: Update `README.md` and `docs/architecture.md` to reflect the 6 new framework tools (bringing total tool inventory from 16 to 22).

---

## 4. Future Roadmap (Milestones 2 & 3)

### Milestone 2: Advanced Analysis, Retrieval Tier 3 & Structured Generators
- **Retrieval tier 3 — `search_features`**: vector similarity via sqlite-vec over the installed project's feature archive. Semantic "have we built something like this before?" query. Reuses the `.claude/rag/` index directory established in M1's `search_docs` — adds a second sqlite file for the vector store.
- `validate_migrations`: Expand/contract SQL migration validator.
- `saturday_test_advisor` & `sunday_test_advisor`: Coverage & gap auditing for Saturday UI tests and Sunday API tests.
- `create_ki` & `create_adr`: Generators for producing structured KI and ADR markdown files.
- `scaffold_docs`: Implementation guide documentation generator.

### Milestone 3: Framework Personas & Workflows
- **Personas**: `architect`, `code_reviewer`, `security_reviewer`, `accessibility_engineer`, `sre_engineer`, `performance_engineer`. Any persona metadata exported by these MUST validate against `shared/schemas/agent-frontmatter.schema.json` (contract in `shared/contracts/agent-frontmatter-contract.md`) so personas surfaced by the framework-MCP are discoverable by the same tooling that reads local `.claude/agents/` files.
- **Workflows**: `review_pr_workflow` (combines code + security + accessibility review), `analyze_repo_workflow` (combines complexity + dependencies + accessibility + language checks).

---

## 5. Open Questions for User Approval

**Resolved 2026-07-23.** Original questions kept for historical context; resolutions inline.

1. ~~**Milestone 1 Tool Selection**~~: **RESOLVED** — 6 tools approved (see §2). Added `search_docs` per Q2 below.

2. ~~**Promote `search_docs` (BM25) into Milestone 1?**~~ **RESOLVED — YES (option b)**. Ship `search_docs` as the 6th M1 tool. Rationale: framework [ADR-002](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md) explicitly names it "the highest-leverage single addition." Original three options kept below for historical context:
   - **(a)** Keep M1 as-is, ship `search_docs` in Milestone 2.
   - **(b)** ← chosen. Add `search_docs` as a 6th M1 tool, accepting the extra sqlite-fts5 dependency and reindex-flow scope.
   - **(c)** Substitute `search_docs` for one of the M1 five (e.g., swap for `verify_dependencies`, which has the highest per-language implementation cost of the analytical tools).
   - **Recommendation**: **(a)** — keep M1 focused on pure input→output tools with no persistent state; ship `search_docs` first in M2 so the index-rebuild story gets its own dedicated milestone.

3. ~~**Server Rename Sub-Question**~~: **RESOLVED — DEFERRED**. M1 ships as `saturday-mcp`. Rename discussion revisited post-M1 when the broad-scope reality is visible and can inform the decision with data. Candidate names (`framework-mcp`, `context-mcp`, `craftsmanship-mcp`) preserved above.

---

## 6. Pre-existing Uncommitted Work

**As of 2026-07-23**, the repo has uncommitted Go files that appear to be an early start on Op 2 (`analyze_complexity`):

- `internal/tools/analyze_complexity_tool.go` (new)
- `internal/tools/analyze_complexity_tool_test.go` (new)
- `internal/analyzers/complexity_analyzer.go` (new)
- `internal/server/tool_provider.go` (modified — adds `complexityAnalyzer` wiring)
- `internal/tools/responses.go` (modified — adds 5 response structs for all originally-planned M1 tools; `DocSearchResult` for the newly-added `search_docs` still needs to be added)

**Decision (2026-07-23)**: keep the WIP. The execution subagent for M1 must:
1. Review the WIP against this plan (Op 2 spec above)
2. Adapt as needed to match the approved shape
3. Commit as Op 2 in a proper conventional-commit format
4. Add `DocSearchResult` to `responses.go` for the new `search_docs` tool (Op 7)

---

## 6. Verification & Discipline

- **Build & Test**: `go build ./... && go test ./...` must be green before every commit.
- **Coverage**: All new tools must achieve ≥ 85% unit test coverage.
- **Git Hygiene**: NEVER use `git add -A` or `git add .`. Stage explicit files under `saturday-mcp/` only.
- **Observability**: All tools automatically receive OTel span wrapping via `server/registration.go` and `tracing_middleware.go`.
