
# Saturday k6 Monorepo

## Layout
- `apps/cli` — TypeScript CLI (`saturday-k6`) for export workflows
- `packages/playwright-k6-exporter` — core exporter (Playwright → k6)
- `packages/k6-redaction-basic` — default redaction policy
- `docs/` — all design docs, ADRs, and agent instructions

## Workflows
- `npm run build` builds all workspaces
- `npm run test` runs tests across workspaces
- `npm run lint` lints packages and apps
