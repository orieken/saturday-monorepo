# mcp-expand Milestone 3 — Framework Personas & Composite Workflows

Plan + execute M3 of mcp-expand. Adds MCP personas (agent-style prompts callable from any MCP client) and composite workflows (chained tool executions).

## Prerequisite

**M2 must be complete before starting M3.** M3 personas may compose M2 tools (e.g., `code_reviewer` persona might call `analyze_complexity` + `verify_dependencies` internally). Confirm M2 exit criterion in `saturday-mcp/mcp-expand-plan.md`'s M2 close-out section before starting.

## Target repo

`/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp` (subdir in parent `saturday-monorepo` git repo).

## Source of truth

- `saturday-mcp/mcp-expand-plan.md` §4 "Milestone 3" — the persona + workflow list
- `saturday-mcp/internal/prompts/provider.go` — existing MCP prompt provider (personas hook here)
- `saturday-mcp/internal/domain/persona.go` — the Persona interface established in mcp-add
- `saturday-mcp/internal/workflows/` — existing workflow pattern (`RunTestsWorkflow`, `PrioritizeTestsWorkflow`) — M3 composite workflows follow the same shape
- Framework `shared/agents/*.md` — the 20+ agents whose prompts these personas mirror. **Personas SHOULD match framework agent semantics** so a user familiar with `/architect` in Claude Code gets the same behavior from the `architect` MCP persona.

## Scope

### Phase A — Draft M3 execution plan (do this FIRST, get approval)

Extend `mcp-expand-plan.md` with §6 "Milestone 3 Execution" — same shape as M1's §3 or M2's execution section (draft in Phase A of the M2 handoff). Include:

- Persona inventory: 6 personas (architect, code_reviewer, security_reviewer, accessibility_engineer, sre_engineer, performance_engineer). For each: which framework agent it mirrors, what M2 tools it composes internally (if any).
- Workflow inventory: 2 composites (review_pr_workflow, analyze_repo_workflow).
- Per-op breakdown.
- Open questions: how does a persona invocation from MCP receive project context (does the calling client pass it? do we introspect env vars?), do personas share the domain.Persona interface with MCP prompt-provider integration or does M3 introduce a new abstraction, do composite workflows use the existing `WorkflowTool` adapter or need a new `CompositeWorkflowTool` shape.

Commit as `docs(mcp): draft mcp-expand-plan §6 Milestone 3 execution plan`.

**Pause. Do not write code until user approves M3 open questions.**

### Phase B — Execute M3 ops per approved plan

Rough ordering:

- **Ops M3.1-M3.6** — 6 personas. Each: `internal/personas/<name>.go` implementing `domain.Persona`, prompt content mirroring the framework agent's system-prompt shape, unit tests. Register through `internal/prompts/provider.go`. Any exported persona metadata **MUST validate against `shared/schemas/agent-frontmatter.schema.json`** from ai-assistant-dot-files (per M2's frontmatter contract cross-ref rule).
- **Op M3.7 `review_pr_workflow`** — composes code_reviewer + security_reviewer + accessibility_engineer personas + tools (analyze_complexity, verify_dependencies, check_accessibility). Uses `WorkflowTool` adapter over `domain.Workflow`.
- **Op M3.8 `analyze_repo_workflow`** — composes analyze_complexity + verify_dependencies + check_accessibility + check_ubiquitous_language into one call, returns a merged report.
- **Op M3.9 Integration + e2e** — extend `TestE2E_ToolInventory` to new total.
- **Op M3.10 Docs refresh** — README + docs/architecture.md + M3 close-out in plan.

## Discipline (non-negotiable)

Same as M1 + M2:
- One commit per op.
- `feat(mcp): add <name> (mcp-expand M3 Op X)`.
- **NEVER `git add -A`.**
- Green build + tests per commit.
- Coverage ≥ 85% per new file.
- Follow M1/M2 patterns exactly — don't invent new abstractions unless plan Phase A explicitly approved them.

## Escalation criteria

Stop and report if:
- A persona needs runtime access to information the MCP client can't reasonably pass (e.g., the whole codebase) — halt, describe the constraint. Retrieval tools from M1/M2 may be the intended path.
- The domain.Persona interface needs extending — halt, propose the change as an ADR-worthy design decision.
- A workflow's steps span both personas AND tools with unclear ordering semantics — halt, ask for the sequencing decision.
- Any persona's prompt would drift from its mirrored framework agent's semantics — halt, describe the drift.

## Report format

Per-op:
```
STATUS: complete | stopped-at-<reason>
Commit: <sha> <message>
Coverage: <pct>
Persona/workflow: <name>
Frontmatter schema validation: <pass/fail>
```

M3 close-out on Op M3.10:
```
M3 COMPLETE
mcp-expand full arc: M1 (9 ops) + M2 (~9 ops) + M3 (~10 ops) shipped
Total MCP surface: 22 M1 tools + <n> M2 tools + <n> M3 personas + <n> M3 workflows
Recommended next step: rename saturday-mcp (deferred per mcp-expand scope decision — post-M1 was too early; M3 is a natural rename inflection point)
```

Go.
