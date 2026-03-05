# Phase 6: The AI Agent Evolution (Brainstorming)

We have completed the foundational tools (Generation, Analysis, Validation). The next phase focuses on **Intelligence, Execution, and Context**.

## 🧠 TODO-018: Knowledge Graph (The Brain)
Current analysis is file-local (regex). We need deep understanding of relationships.
- **Symbol Indexing**: Parse TypeScript AST to understand classes, methods, and types.
- **Dependency Graph**: Map `Flows` -> `Pages` -> `Elements`.
- **Impact Analysis Tool**: `analyze_impact(file: string)` -> "Changing this selector affects 3 Flows and 12 Scenarios."

## 🛠️ TODO-019: Test Execution & Self-Healing (The Hands)
Giving the Agent the ability to run and fix code.
- **Test Runner Tool**: `run_test(filter: string)` -> Executes `npx playwright test ...` and returns structured JSON results/logs.
- **Self-Healing Loop**:
    1. Agent runs test -> implementation fails.
    2. Agent reads error logs + DOM snapshot.
    3. Agent updates Selector in Page Object.
    4. Agent re-runs test to verify.

## 👁️ TODO-020: Visual Intelligence (The Eyes)
Using Vision models to bridge Design -> Code.
- **Screenshot to Page Object**: `generate_from_image(path)` -> Analyzes UI screenshot, identifies interactables, generates Saturday Page Object.
- **Visual Diff Analysis**: Compare visual regression screenshots and explain *why* it failed (e.g., "The button moved 5px right").

## 📊 TODO-021: Data-Driven Insights (The Nerves)
Connecting code to production reality (Observability).
- **Heatmap Integration**: `get_page_metrics(path)` -> "70% of users click the 'Sign Up' button, but we don't have a test for it."
- **Flakiness History**: `analyze_stability(test_name)` -> "This test has failed 15% of the time over the last week."
