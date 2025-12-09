
# Copilot / Agent Instructions — @saturday/playwright-k6-exporter

## North Star
Produce clean, testable code that adheres to SOLID and Clean Architecture. Favor small, composable functions with **cyclomatic complexity < 7**.

## Quality Gates (must pass)
- **Cyclomatic complexity:** < 7 per function/method
- **Max function length:** ≤ 30 LOC (excl. blanks/comments)
- **Params per function:** ≤ 4 (prefer value objects)
- **Duplication:** < 5% per module
- **Coverage:** ≥ 85% (unit + integration)
- **Lint:** 0 errors, 0 new warnings
- **Types:** No `any` unless documented with rationale

## Development Rules
1. **TDD/BDD**: Red–Green–Refactor, keep specs clear and focused.
2. **DDD-lite**: Separate domain (models, policies) from infrastructure (Playwright proxy, fs).
3. **Side-effects last**: Pure transforms first, IO at the edges.
4. **Logs over prints**: Use structured messages if needed; avoid console noise.
5. **No secrets**: Never write tokens/PII to files. Provide a `sanitize()` step before emit.

## Architectural Contracts
- Recorder must be **stateless aside from call log**; reset per test.
- File writes are **idempotent** (slug by test title).
- Public API surface:
  - `createK6Recorder(testTitle: string, outDir?: string)`
  - `class K6Recorder { wrap(ctx), record(), hasCalls(), flushToK6() }`
- Fixture path: `@saturday/playwright-k6-exporter/fixture`

## Commit Hygiene
- Conventional commits (feat, fix, docs, refactor, test, chore).
- Keep PRs small; update docs and diagrams when behavior changes.

## References
- Kent Beck (TDD, Simple Design)
- Martin Fowler (Refactoring, Architecture)
- Neal Ford (Evolutionary Architecture)
- Robert C. Martin / Uncle Bob (SOLID, Clean Code)
