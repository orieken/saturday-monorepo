# Backfill unit tests for `internal/tools/workflow_tool.go`

## Target repo

`/Users/oscarrieken/Projects/Rieken/saturday-monorepo/saturday-mcp` (subdirectory in parent `saturday-monorepo` git repo — commits go to the parent).

## Prior context

`workflow_tool.go` is a thin generic adapter that wraps a `domain.Workflow` so it can be registered with the MCP server as a `domain.Tool`. It was introduced in Phase D of the mcp-add retrofit (commit `17f515a` — `refactor(mcp): wire workflows via tools.WorkflowTool adapter`).

Currently at **0% coverage**. The two workflows it wraps (`RunTestsWorkflow`, `PrioritizeTestsWorkflow`) are covered by `internal/workflows/*_test.go`, but the adapter itself is not.

## Scope

Add one test file: `internal/tools/workflow_tool_test.go` covering `WorkflowTool.Execute` and any other exported methods.

**Minimum coverage:**
1. **Happy path** — `Execute` delegates to `workflow.Run(ctx, args)` and returns its result unchanged. Use a hand-rolled `fakeWorkflow` struct that captures the call and returns a canned result.
2. **Error propagation** — when `workflow.Run` returns an error, `WorkflowTool.Execute` propagates it (or wraps it — read the actual code to see the pattern before writing the assertion).
3. **Metadata** — `Name()`, `Description()`, `InputSchema()`, `OutputSchema()` return the values passed to the constructor.

Target: ≥ 85% coverage on `workflow_tool.go`.

## Reference

Read before starting:
- `internal/tools/workflow_tool.go` — the source under test
- `internal/tools/testfixtures_test.go` — shared helpers (`buildRequest`, `extractText`, etc.)
- `internal/workflows/run_tests_workflow_test.go` — shows the fake-workflow-collaborator pattern; adapt it to fake the whole `domain.Workflow`

## Discipline

- **One commit** for this work.
- Conventional Commits: `test(mcp): add unit tests for WorkflowTool`
- **NEVER `git add -A` or `git add .`** — the parent monorepo has 100+ unrelated in-progress files. Stage explicit paths only: `git add saturday-mcp/internal/tools/workflow_tool_test.go`
- `git status --short` after staging, before commit
- `go build ./... && go test ./...` must be green before committing

## Escalation

Stop and report if:
- The `WorkflowTool` API doesn't cleanly accept a fake — its constructor may bind concrete types that need Extract Interface first (would be out of scope for this prompt)
- Coverage stays below 70% after honest effort — the file may have defensive code branches that are unreachable

## Report format (under 100 words)

```
STATUS: complete | stopped-at-<reason>
Commit: <sha> <message>
Coverage: workflow_tool.go: <pct> (was 0%)
Test suite: all green | <details>
```

Go.
