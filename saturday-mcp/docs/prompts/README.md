# Agent Prompts

Ready-to-fire handoff prompts for outstanding work on saturday-mcp. Each file is self-contained — a fresh agent (Claude Code chat, subagent, or other MCP client) can pick up any prompt without prior conversation context.

## How to use

1. Open a fresh chat / spawn a fresh agent.
2. Copy the entire contents of the target prompt file into the message.
3. The agent will have everything it needs: repo path, prior state, scope, guardrails, deliverable format.

## Prompts

| File | Scope | Estimated size |
|---|---|---|
| [backfill-workflow-tool-tests.md](backfill-workflow-tool-tests.md) | Add unit tests for `internal/tools/workflow_tool.go` (currently 0% coverage) | Small — 1-2 tests, 1 commit |
| [refresh-docs-post-retrofit.md](refresh-docs-post-retrofit.md) | Update `README.md` and `docs/architecture.md` to reflect the post-retrofit trinity + adapter structure | Medium — 2-4 commits |
| [mcp-expand.md](mcp-expand.md) | Add framework-tooling capabilities (analytical tools, personas, KI search, etc.) to saturday-mcp | Large — spawn subagents in batches, similar cadence to the mcp-add retrofit |

## Convention

Every prompt in this directory:
- States the target repo path explicitly
- Summarizes prior state (what's done, where to look)
- Enumerates the scope (what to do, what NOT to touch)
- Lists guardrails (git hygiene — the parent monorepo has 100+ unrelated in-progress files, so `git add -A` is banned)
- Specifies escalation criteria (when to stop and ask instead of pushing through)
- Requests a specific report format

Add new prompts as new files. When a prompt is executed and committed, either delete it or move it to a `docs/prompts/done/` subdirectory (convention: TBD).
