# Agent Prompts

Ready-to-fire handoff prompts for work on `saturday-mcp`. Each file is self-contained — a fresh agent (Claude Code chat, subagent, or other MCP client) can pick up any prompt without prior conversation context.

## How to use

1. Open a fresh chat / spawn a fresh agent.
2. Copy the entire contents of the target prompt file into the message.
3. The agent will have everything it needs: repo path, prior state, scope, guardrails, deliverable format.

## Active Prompts (ready to fire in a fresh chat / subagent)

| File | Scope | Estimated size |
|---|---|---|
| [mcp-expand-milestone-2.md](mcp-expand-milestone-2.md) | Plan + execute M2: retrieval tier 3 (`search_features` via sqlite-vec), `validate_migrations`, saturday/sunday test advisors, KI/ADR/doc generators | Large — Phase A plan draft + ~9 execution ops |
| [mcp-expand-milestone-3.md](mcp-expand-milestone-3.md) | Plan + execute M3: 6 framework personas mirroring `shared/agents/*.md` + 2 composite workflows (review_pr, analyze_repo). **Depends on M2 landing first.** | Large — Phase A plan draft + ~10 execution ops |
| [post-m1-followups.md](post-m1-followups.md) | Two small M1 leftovers: `/reindex-docs` skill + `Handler.Shutdown()` graceful-close hook | Small — 2 independent ops, either can go first |

## Completed Prompts (`docs/prompts/done/`)

| File | Scope | Status |
|---|---|---|
| [backfill-workflow-tool-tests.md](done/backfill-workflow-tool-tests.md) | Add unit tests for `internal/tools/workflow_tool.go` (100% coverage) | Completed |
| [refresh-docs-post-retrofit.md](done/refresh-docs-post-retrofit.md) | Update `README.md` and `docs/architecture.md` to reflect post-retrofit trinity structure | Completed |
| [mcp-expand.md](done/mcp-expand.md) | Framework-tooling expansion — Milestone 1 (6 analytical & search tools) | Completed |

## Recommended execution order

1. **post-m1-followups.md** first — small, closes M1 loose ends before M2 starts
2. **mcp-expand-milestone-2.md** next — do Phase A (plan draft) → get user approval → Phase B execution
3. **mcp-expand-milestone-3.md** last — same Phase A → approval → Phase B pattern

## Convention

Every prompt in this directory:
- States the target repo path explicitly
- Summarizes prior state (what's done, where to look)
- Enumerates the scope (what to do, what NOT to touch)
- Lists guardrails (git hygiene — explicit paths only, `git add -A` is banned)
- Specifies escalation criteria (when to stop and ask)
- Requests a specific report format

When a prompt is executed and committed, move it to `docs/prompts/done/`.
