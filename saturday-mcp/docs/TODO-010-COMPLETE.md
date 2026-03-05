# TODO-010: Resource & Prompt Providers

## Status: ✅ Complete

We have enhanced the Prompt Provider for the Saturday MCP server to include a robust set of predefined prompts that assist developers in building and debugging with the Saturday Framework.

## Implemented Prompts

### 1. `plan_feature`
- **Purpose**: Helps developers plan the implementation of a new feature.
- **Inputs**: `feature` (description or name)
- **Output**: A step-by-step implementation plan covering Page Objects, Flows, Step Definitions, and Test scenarios using Saturday patterns.

### 2. `explain_framework`
- **Purpose**: Educational prompt for onboarding new developers.
- **Inputs**: None.
- **Output**: A concise explanation of the Saturday architecture (Site, Pages, Flows, Steps) and how they fit together.

### 3. `debug_error` (New)
- **Purpose**: Assists in root cause analysis of test failures.
- **Inputs**: 
    - `error`: The error message or stack trace.
    - `context` (optional): Relevant code snippets or step definition.
- **Output**: Root cause analysis, specific debugging steps, and fix suggestions tailored to common Saturday/Playwright pitfalls.

### 4. `generate_gherkin` (New)
- **Purpose**: Converts requirements into structured BDD scenarios.
- **Inputs**: `requirements` (user story).
- **Output**: A complete `.feature` file content with Feature, Background, Scenarios, and Tags, following best practices (declarative steps, independent scenarios).

## Testing

All prompts are covered by unit tests in `internal/prompts/provider_test.go`, verifying that they:
- Are correctly registered and listed.
- Accept the defined arguments.
- Return the expected text content in the prompt message.

We verified the implementation with `go test ./internal/prompts/...` and achieved success.
