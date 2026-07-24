# mcp-expand M1 post-milestone follow-ups

Two small deferred items from M1 execution. Both cosmetic-quality — nothing broken today, both were noted with "deferred to follow-up" in their source commits.

## Target repo

`/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp`.

## Scope: two independent ops

### Op F1 — `/reindex-docs` skill for `search_docs`

**Deferred from M1 Op 7b (commit `5a47441`).** Today, `search_docs.Execute` calls `retriever.EnsureIndex(...)` on every invocation to catch new/changed docs. Fine for M1's corpus sizes, but wastes work when the caller knows the corpus hasn't changed.

Add an MCP skill (not a tool — this is a maintenance command) at `internal/skills/reindex_docs.go` that:
- Takes optional `corpusPaths []string` input (defaults to same discovery as the tool)
- Calls `retriever.EnsureIndex(corpusPaths)` once
- Returns a simple result: `{indexed, corpusPaths, durationMs}`
- Docs update in README noting the skill exists as an opt-in performance optimization

Then update `search_docs_tool.go` to add an optional `skipEnsureIndex bool` input arg — callers can flip it to true to skip the auto-rebuild and rely on the manual reindex flow.

**One commit**: `feat(mcp): add /reindex-docs skill + skipEnsureIndex option (M1 follow-up F1)`.

### Op F2 — `Handler.Shutdown()` hook for graceful sqlite close

**Deferred from M1 Op 7b (commit `5a47441`).** `BM25Retriever.Close()` exists but nothing calls it — sqlite is durable so ungraceful shutdown works fine, but a proper hook is preferred.

Add:
- `Handler.Shutdown(ctx context.Context) error` method that iterates the tool list and calls `Close()` on any tool that implements a `Closer` interface (or specifically on the BM25 retriever if simpler)
- Update `cmd/saturday-mcp/main.go` to install a signal handler (SIGTERM, SIGINT) that calls `Handler.Shutdown()` before `os.Exit`
- Update every test that instantiates `Handler` — add a `defer handler.Shutdown(ctx)` where the test would otherwise leak DB handles (usually only matters for tests that override `SATURDAY_MCP_DOCS_FTS_PATH` to a temp DB)

**One commit**: `feat(mcp): add Handler.Shutdown() hook for graceful shutdown (M1 follow-up F2)`.

## Discipline

- Two commits, one per op. F1 and F2 are independent — either can go first.
- Conventional Commits.
- **NEVER `git add -A`.**
- Green build + tests per commit.

## Escalation

- If `Handler.Shutdown()` requires restructuring the Handler struct significantly — halt, this was expected to be small. Propose the smaller shape.
- If the `skipEnsureIndex` flag on `search_docs` breaks any existing test's assumptions — halt, describe.

## Report

Per-op single-liner:
```
F1: <sha> <message>
F2: <sha> <message>
```

Go.
