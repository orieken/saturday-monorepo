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

---

## 2. Milestone 1: Core Analytical & Search Tools

Milestone 1 focuses on 5 high-leverage tools that operate on pure input → structured output contracts with zero side effects:

| Tool | Capability Category | Input | Output / Contract | Key Value |
|---|---|---|---|---|
| `analyze_complexity` | Analytical | `projectPath`, `maxComplexity` (default 7), `maxLines` (default 30) | `ComplexityAnalysisResult` (file/func violations breakdown) | Enforces Beck & Sandi Metz clean code thresholds |
| `check_accessibility` | Analytical | `filePath` or `projectPath` | `AccessibilityReportResult` (violations list with line & element) | Checks semantic HTML, ARIA, and missing labels |
| `check_ubiquitous_language` | Analytical | `projectPath`, `dictionaryPath` | `UbiquitousLanguageResult` (unapproved synonyms & violations) | Ensures code matches `DOMAIN_DICTIONARY.md` terms |
| `search_ki` | Search / Query | `query`, `tags`, `domain` | `KISearchResult` (ranked KI & ADR matches with paths) | Instant cross-project Knowledge Item discovery |
| `verify_dependencies` | Analytical | `projectPath` | `DependencyVerificationResult` (Clean Arch boundary breaches) | Prevents inner layers from importing outer layers |

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

### Phase C — Search & Knowledge Query Tools

- **Op 6 (`search_ki`)**: Create `internal/tools/search_ki_tool.go` implementing `domain.Tool`, add response struct `KISearchResult` in `responses.go`, wire into `tool_provider.go`, add unit tests in `internal/tools/search_ki_tool_test.go`.

### Phase D — Integration & Documentation

- **Op 7 (Integration & E2E)**: Update `internal/integration/e2e_test.go` and `internal/server/testing.go` to cover all 5 new tools.
- **Op 8 (Documentation Refresh)**: Update `README.md` and `docs/architecture.md` to reflect the 5 new framework tools (bringing total tool inventory from 16 to 21).

---

## 4. Future Roadmap (Milestones 2 & 3)

### Milestone 2: Advanced Analysis & Structured Generators
- `validate_migrations`: Expand/contract SQL migration validator.
- `saturday_test_advisor` & `sunday_test_advisor`: Coverage & gap auditing for Saturday UI tests and Sunday API tests.
- `create_ki` & `create_adr`: Generators for producing structured KI and ADR markdown files.
- `scaffold_docs`: Implementation guide documentation generator.

### Milestone 3: Framework Personas & Workflows
- **Personas**: `architect`, `code_reviewer`, `security_reviewer`, `accessibility_engineer`, `sre_engineer`, `performance_engineer`.
- **Workflows**: `review_pr_workflow` (combines code + security + accessibility review), `analyze_repo_workflow` (combines complexity + dependencies + accessibility + language checks).

---

## 5. Open Questions for User Approval

Before initiating execution of Phase B, the following decisions require user feedback/approval:

1. **Milestone 1 Tool Selection**: Confirm approval of the 5 proposed Milestone 1 tools (`analyze_complexity`, `check_accessibility`, `check_ubiquitous_language`, `search_ki`, `verify_dependencies`).
2. **Server Rename Sub-Question**: Recommend **deferring** any module/directory rename (e.g. `saturday-mcp` → `framework-mcp`) until Milestone 1 ships. Deferral prevents disrupting existing MCP client configs while immediately delivering the broad framework capabilities.

---

## 6. Verification & Discipline

- **Build & Test**: `go build ./... && go test ./...` must be green before every commit.
- **Coverage**: All new tools must achieve ≥ 85% unit test coverage.
- **Git Hygiene**: NEVER use `git add -A` or `git add .`. Stage explicit files under `saturday-mcp/` only.
- **Observability**: All tools automatically receive OTel span wrapping via `server/registration.go` and `tracing_middleware.go`.
