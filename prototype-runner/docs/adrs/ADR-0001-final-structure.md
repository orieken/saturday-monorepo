ADR-0001: Adopt Unified Final Project Structure under `final/`

Status: Accepted
Date: 2025-11-23

Context

Multiple prototypes exist in this repository:

- `cuke_dashboard_full/` (full-stack slice)
- `cuke_playwright_dashboard_full/` (backend + runner + tests)
- `test_dashboard_full_with_steps/` (full-stack slice with additional step examples)

They share the same domain but diverge in naming, layout, and scope. Contributors need a single, coherent structure to reduce confusion and enable steady progress.

Decision

Create a root-level `final/` directory that standardizes the layout for backend, frontend, tests, tools, local cluster, and docs. The target structure is documented in `docs/project-structure.md` and the architecture in `docs/architecture.md`.

Consequences

- Pros:
  - Clear single source of truth for folder layout and module responsibilities.
  - Easier onboarding; consistent documentation and commands.
  - Simplifies CI and future automation.
- Cons:
  - Temporary duplication during migration until prototypes are retired.
  - Requires careful synchronization to avoid divergence during transition.

Notes

- Migration will be incremental, as outlined in `docs/migration-guide.md`.
- Prototypes remain read-only during the transition; new work targets `final/`.
