
# Contributing

## Setup
```bash
npm i
npm run build
```
## Test
Add/extend tests in your host Playwright project that consumes this package. Keep unit helpers small; aim for ≥85% coverage.

## Style & Quality
- ESLint + TypeScript strict.
- Cyclomatic complexity < 7 per function—extract until it is.
- Keep functions ≤ 30 LOC (excl. blanks/comments).
- No more than 4 parameters per function (prefer objects).

## Docs
- Update `docs/ARCHITECTURE.md` when behavior or flow changes.
- If adding features (redaction, batching, thresholds), include a Mermaid diagram snippet.

## Releasing
- `npm version (patch|minor|major)`
- `npm publish` (public package) or your internal registry equivalent.
