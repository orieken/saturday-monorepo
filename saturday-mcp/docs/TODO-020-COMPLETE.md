# TODO-020: Visual Intelligence

## Status: ✅ Complete (Prompt Phase)

We have implemented "Visual Intelligence" by leveraging the capabilities of MCP Clients (like Claude) through a specialized **Prompt**.

## Approach

Since the Go server itself does not embed a Visual LLM, we use the **MCP Prompt** mechanism to guide the connected agent.

*   **Prompt Name**: `visual_page_object`
*   **Workflow**:
    1.  User selects the `visual_page_object` prompt in Claude Desktop.
    2.  User attaches a screenshot of the UI.
    3.  User (optionally) provides the `componentName` (e.g., `LoginPage`) and `path` (e.g., `/login`).
    4.  The Prompt injects detailed system instructions on how to analyze the image and map it to Saturday Framework patterns.
    5.  Claude analyzes the image and generates the robust Page Object code.

## Why this approach?

*   **Leverage**: Uses the state-of-the-art vision capabilities of the host model (e.g., GPT-4o, Claude 3.5 Sonnet).
*   **Zero Dependency**: No need to add heavy ML libraries to the Go MCP server.
*   **Agentic**: Fits perfectly into the "Human-in-the-loop" workflow where the user provides the visual context.

## Usage

In Claude Desktop:
1.  Type `/` and select `visual_page_object`.
2.  Fill in the arguments (e.g., `LoginPage`, `/login`).
3.  Click the "Attach Image" button and select your UI mockup or screenshot.
4.  Hit Enter.
