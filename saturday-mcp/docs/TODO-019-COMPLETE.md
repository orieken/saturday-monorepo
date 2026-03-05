# TODO-019: Test Execution Tools

## Status: ✅ Complete (Execution Phase)

We have implemented the ability for the MCP server to execute tests and capture their output. This is the first half of the "Test Execution & Self-Healing" goal.

## Implemented Components

### 1. Test Executor (`internal/executor/runner.go`)
- **Wrapper**: Wraps `os/exec` to run shell commands in a structured way.
- **Output Capture**: Captures both `stdout` and `stderr` combined.
- **Filtering**: Supports passing grep filters (e.g., `-g 'login'`) to Playwright.
- **Environment**: Supports injecting custom environment variables.

### 2. Run Tests Tool (`run_tests`)
Exposed via MCP to allow the agent to invoke tests.

**Input:**
- `projectPath`: Root directory of the project.
- `command`: (Optional) Custom command, defaults to `npx playwright test`.
- `filter`: (Optional) Test name filter.

**Output:**
- `success`: Boolean (exit code 0).
- `output`: Full console log.
- `summary`: Simple status string.

## Next Steps (Self-Healing)

Now that we can *run* tests, the next step (Phase 2 of TODO-019) is to build the loop:
1.  Run Test -> Fail.
2.  Pass `output` + `targetFile` to `analyze_error`.
3.  Use `analyze_impact` to see what else broke.
4.  Use `generate_page` or `replace_file_content` to fix.

This loop is likely an *agentic workflow* rather than a single MCP tool, so the infrastructure is now ready!
