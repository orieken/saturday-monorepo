# Technical Analysis: mcp-expand

## 1. Overview & Strategy
Expand `saturday-mcp` from a Saturday-only code generator into the broad framework-MCP surface for the entire context-engineering ecosystem (`ai-assistant-dot-files`).

## 2. Milestone 1 Recommendation (4-5 Tools)
Select high-leverage analytical & search tools with zero side effects:
- `analyze_complexity`: Calculates function cyclomatic complexity & line count against framework rules (< 7 complexity, < 30 LOC).
- `check_accessibility`: Scans HTML/Vue/JSX UI templates for semantic HTML violations (e.g. `div` click handlers, missing ARIA/labels).
- `check_ubiquitous_language`: Checks source files against domain terms in `DOMAIN_DICTIONARY.md`.
- `search_ki`: Searches Knowledge Items in `shared/knowledge/`, `.claude/knowledge/`, and `docs/adrs/` by tag, domain, or keyword.
- `verify_dependencies`: Enforces Clean Architecture dependency boundaries by inspecting Go/TS import graphs.

## 3. Milestones 2 & 3 Roadmap
- **Milestone 2 (Expanded Analysis & Structured Artifact Generators)**: `validate_migrations`, `saturday_test_advisor`, `sunday_test_advisor`, `create_ki`, `create_adr`, `scaffold_docs`.
- **Milestone 3 (Personas & Composite Workflows)**: `architect`, `code_reviewer`, `security_reviewer`, `accessibility_engineer`, `sre_engineer`, `performance_engineer`, `review_pr_workflow`, `analyze_repo_workflow`.

## 4. Open Questions to Raise for Approval
1. **Milestone 1 Scope**: Approve the recommended 5 tools (`analyze_complexity`, `check_accessibility`, `check_ubiquitous_language`, `search_ki`, `verify_dependencies`) for implementation.
2. **Server Naming**: Confirm deferring any repo/module rename (`saturday-mcp` → `framework-mcp` / `context-mcp`) until after Milestone 1 ships.
