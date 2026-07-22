# Add framework-tooling capabilities to saturday-mcp

## Target repo

`/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp` (Go MCP server, mcp-go v0.56.0, stdio transport, in the parent `saturday-monorepo` git repo — commits go there).

## Prior context

The mcp-add retrofit is complete (59 commits, closed out in `mcp-add-plan.md`). Current structure:

- `internal/domain/{tool,persona,workflow,tracer,testrunner,filesystem,metrics}.go` — trinity interfaces
- `internal/tools/` — 14 Saturday-scaffolding tools (`generate_site`, `generate_page`, etc.) + `WorkflowTool` adapter + shared `schemas.go`/`responses.go`/`testfixtures_test.go`
- `internal/workflows/` — `RunTestsWorkflow`, `PrioritizeTestsWorkflow`
- `internal/adapters/{testrunner,filesystem,metricsfile,otel}/` — adapter implementations
- `internal/server/{handler,registration,tracing_middleware}.go` — thin composite + registration loop

Adding a new tool is now: drop a file in `internal/tools/`, register in the Handler's tool slice, done. Every tool auto-wraps in an OTel span. Every tool has typed input + output schemas via `invopop/jsonschema`.

## Framework context (another repo)

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` has the context-engineering framework: rules, agents, skills, blueprints, KIs. Many of its capabilities are the shape MCP tools want.

## Scope decision (resolved 2026-07-22)

**Chosen scope**: **Broad framework-MCP** — saturday-mcp becomes the framework's MCP surface, exposing analytical tools, personas, and search capabilities useful for any language/codebase using the framework. NOT just Saturday-specific tooling.

Rationale: better long-term positioning; makes saturday-mcp the reference implementation of an MCP that surfaces framework capabilities to any MCP client (Claude Code, Claude Desktop, Cursor, etc.).

**Open sub-question for the fresh chat to raise with the user** (not yet decided): rename `saturday-mcp` → something more accurate to the broader scope (e.g., `context-mcp`, `framework-mcp`, `craftsmanship-mcp`). Real cost — affects MCP client configs, the parent monorepo directory, and the go.mod module path. Reasonable to defer the rename until Milestone 1 ships and the broad-scope reality is visible; the current name doesn't functionally block anything.

## Candidate MCP capabilities to add — categorized by fit

### Strong fit (analytical tools — pure input → structured output)

- `analyze_complexity` — cyclomatic complexity check (framework rule: < 7)
- `check_accessibility` — semantic HTML violations in Vue/HTML/React
- `check_ubiquitous_language` — DOMAIN_DICTIONARY.md violations
- `verify_dependencies` — Clean Arch layer boundary check (import graph)
- `validate_migrations` — expand/contract SQL migration check
- `saturday_test_advisor` — audit Saturday coverage (dead pages/flows, missing scenarios)
- `sunday_test_advisor` — audit Sunday API test coverage
- `health_check` — validate framework install integrity

### Strong fit (search/query tools)

- `search_ki` — search Knowledge Items by tag/domain/keyword
- `query_memory` — registry-aware search across memory sources

### Strong fit (personas — MCP prompts)

- `architect` — Clean Architecture layer reviewer
- `code_reviewer` — SOLID + complexity + naming discipline
- `security_reviewer` — STRIDE threat modeling
- `accessibility_engineer` — semantic HTML, ARIA, keyboard nav
- `sre_engineer` — OTel span design, SLI cardinality
- `performance_engineer` — N+1 detection, timeout budgets, pagination

### Medium fit (workflows — composed tool chains)

- `review_pr_workflow` — code_reviewer + security_reviewer + accessibility_engineer in sequence
- `analyze_repo_workflow` — complexity + dependencies + accessibility + tests in one call
- `retrofit_scorecard_workflow` — like mcp-add Phase 2 gap analysis, generic

### Medium fit (generation tools)

- `create_ki` — structured input → Knowledge Item markdown
- `create_adr` — structured input → ADR markdown (Deciders, Alternatives, Fitness Function)
- `scaffold_docs` — comprehensive markdown implementation guide

### Poor fit (skip — interviews don't map to MCP)

- `bootstrap_project`, `spec_writer`, `event_storm`, `new_feature`, `mcp_add` — need turn-taking that MCP tools don't support

## Additional ideas beyond "port existing skills"

- Saturday-specific persona: `site_centric_reviewer` — reviews a `.feature` or Page/Flow against Site-Centric pattern
- Chained tool: `generate_page_and_verify` — runs `generate_page` then `check_accessibility` + `analyze_complexity` in one call
- Improvement to `generate_documentation`: emit framework-compliant docs (with layer tables, testing pyramid, ADR seeds)

## What to do

1. ~~Decide the scope question~~ — resolved above (broad framework-MCP)
2. Prune the candidate list to what's actually going to ship in Milestone 1 — recommend 4-6 tools, not all of them. Suggest starting with the analytical tools since they're purest input→output (analyze_complexity, check_accessibility, check_ubiquitous_language) plus search_ki (highest-leverage generic tool).
3. For each Milestone 1 tool: draft the same structure Phase B/C used in mcp-add (Tool interface implementation, typed input/output structs via invopop, shared helpers, unit tests, OTel span emission via existing middleware)
4. Produce an `mcp-expand-plan.md` (mirroring `mcp-add-plan.md`'s structure) at the saturday-mcp repo root before writing any code
5. Raise the rename sub-question (see scope-decision section above) — recommend deferring until after Milestone 1
6. Follow the same commit discipline: one op per commit, conventional commits, NEVER `git add -A` (parent monorepo has 100+ unrelated files), verify `go build && go test` after each commit

## Repo hygiene

- Commit style: `docs(mcp)` for planning artifacts, `feat(mcp)` for new tools/personas, `refactor(mcp)` for structural changes
- Parent git repo is `saturday-monorepo`, not `saturday-mcp` — commits land there
- ALWAYS stage explicit paths under `saturday-mcp/`, never use `git add -A`

## Start here

Scope is already resolved (broad framework-MCP — see section above). Proceed directly to drafting `mcp-expand-plan.md`. Do not write any code before the plan is approved. When the plan is ready, raise the rename sub-question and the Milestone 1 tool selection for user approval before starting execution.
