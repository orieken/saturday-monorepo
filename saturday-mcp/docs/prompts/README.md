# Agent Prompts

Ready-to-fire handoff prompts for work on `saturday-mcp`. Each file is self-contained — a fresh agent (Claude Code chat, subagent, or other MCP client) can pick up any prompt without prior conversation context.

## How to use

1. Open a fresh chat / spawn a fresh agent.
2. Copy the entire contents of the target prompt file into the message.
3. The agent will have everything it needs: repo path, prior state, scope, guardrails, deliverable format.

## Completed Prompts (`docs/prompts/done/`)

| File | Scope | Status |
|---|---|---|
| [backfill-workflow-tool-tests.md](done/backfill-workflow-tool-tests.md) | Add unit tests for `internal/tools/workflow_tool.go` (100% coverage) | Completed |
| [refresh-docs-post-retrofit.md](done/refresh-docs-post-retrofit.md) | Update `README.md` and `docs/architecture.md` to reflect post-retrofit trinity structure | Completed |
| [mcp-expand.md](done/mcp-expand.md) | Framework-tooling expansion — Milestone 1 (6 analytical & search tools) | Completed |

## Future Expansion Roadmap (`mcp-expand-plan.md`)

The roadmap for upcoming work is documented in [`mcp-expand-plan.md`](../mcp-expand-plan.md):

- **Milestone 2**: Advanced Verification & Generators
  - `search_features` (semantic vector search via sqlite-vec over feature archives)
  - `validate_migrations` (expand/contract SQL migration validator)
  - `saturday_test_advisor` & `sunday_test_advisor` (test coverage gap auditors)
  - `create_ki`, `create_adr`, `scaffold_docs` (structured markdown generators)
- **Milestone 3**: Framework Personas & Composite Workflows
  - Personas: `architect`, `code_reviewer`, `security_reviewer`, `accessibility_engineer`, `sre_engineer`, `performance_engineer`
  - Workflows: `review_pr_workflow`, `analyze_repo_workflow`

## Convention

Every prompt in this directory:
- States the target repo path explicitly
- Summarizes prior state (what's done, where to look)
- Enumerates the scope (what to do, what NOT to touch)
- Lists guardrails (git hygiene — explicit paths only, `git add -A` is banned)
- Specifies escalation criteria (when to stop and ask)
- Requests a specific report format

When a prompt is executed and committed, move it to `docs/prompts/done/`.
