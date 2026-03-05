# Task: Implement "Autonomous QA Engineer" Workflow

## Status: 📋 Planning

**Objective**: Combine the capabilities of Visual Intelligence (TODO-020), Test Execution (TODO-019), and Knowledge Graph (TODO-018) into a unified, high-level agentic workflow. This workflow mimics the behavior of a human QA engineer implementing a new feature.

## The Workflow: "Implement & Verify"

We want to enable a single command/prompt that triggers the following autonomous loop:

### 1. Visualization & Modeling
*   **Input**: User requirement ("Login to the site") + UI Screenshot.
*   **Action**: Use `visual_page_object` prompt.
*   **Output**: New Page Object file (e.g., `LoginPage.ts`).

### 2. Test Generation
*   **Input**: The new Page Object.
*   **Action**: Generate a Playwright test spec (e.g., `login.spec.ts`) using standard MCP code generation.
*   **Output**: A new test file.

### 3. Verification Loop (The "Self-Healing" Core)
*   **Action**: Run `run_tests` targeting the new spec.
*   **Decision**:
    *   **Pass**: Proceed to Analysis.
    *   **Fail**:
        1.  Call `parse_test_failure`.
        2.  Read the specific error (e.g., "Selector not found").
        3.  Call `start_debugging` (or generic AI reasoning).
        4.  Edit the code (Page Object or Test) to fix.
        5.  **RETRY** (Max 3 attempts).

### 4. Impact Safety Check
*   **Input**: The modified/created files.
*   **Action**: Call `analyze_impact`.
*   **Output**: A list of other files potentially broken by these changes.
*   **Action (Optional)**: Run tests for those impacted files to ensure no regressions.

## Deliverables

1.  **Workflow Definition**: A regular `.agent/workflows/implement-feature.md` file that strictly defines this process for the agent to follow.
2.  **Orchestrator Prompt**: A new MCP prompt `implement_feature_loop` that sets the context for the agent to behave this way.
3.  **Documentation**: Example usage guide.

## Technical Requirements

*   No new Go code required (likely). This is about *orchestration* of existing tools.
*   We rely on the Agent's ability to "Chain of Thought" through these tools.
