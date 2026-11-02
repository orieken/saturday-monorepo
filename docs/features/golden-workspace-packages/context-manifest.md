# Context Manifest: golden-workspace-packages

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/ye-olde-magic-shop/package.json` — AC1 target — three devDependencies currently point at file:./saturday-core, file:./saturday-cucumber, file:./saturday-playwright-heatmap and must become workspace:*.
- `apps/ye-olde-magic-shop/Dockerfile` — AC4 — COPY's the vendored saturday-core/saturday-cucumber/saturday-playwright-heatmap dirs and strips their node_modules/lockfiles before npm install; must build from the workspace or be removed with the reason recorded.
- `apps/ye-olde-magic-shop/Dockerfile.runner` — AC4 — same vendored-copy coupling as Dockerfile, used by scripts/ui-suite-in-chrome.sh's runner image.
- `apps/ye-olde-magic-shop/saturday-cucumber/src/config/index.ts` — Confirmed behavior drift named in the spec's Open Questions: this vendored copy reads process.env.HEADLESS for the browser default; the workspace package version does not.
- `packages/saturday-cucumber/src/config/index.ts` — The workspace package's current config module — target for the HEADLESS default move the spec's starting point describes, with tests/config.spec.ts as its coverage.
- `apps/ye-olde-magic-shop/.env` — Sets HEADLESS=false — the environment value the drifted vendored config reads, confirming the divergence is live and exercised, not theoretical.
- `pnpm-workspace.yaml` — Confirms apps/* and packages/* are already both workspace members, so workspace:* resolution for the three dependencies needs no workspace-config change.
- `scripts/ui-suite-in-chrome.sh` — AC5's verification command — the task breakdown needs to know what it builds and runs to schedule a before/after run and record its output.
- `apps/ye-olde-magic-shop/tests/cart.spec.ts` — The pre-existing golden-test bug already fixed on main (#7) that originally blocked this work — confirms the starting point's premise before the task list assumes it's resolved.

## Token Budget

- **Status**: OK — ~3178 tokens, 2.6% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/package.json` | 2188 | 547 |
| `apps/ye-olde-magic-shop/Dockerfile` | 871 | 217 |
| `apps/ye-olde-magic-shop/Dockerfile.runner` | 704 | 176 |
| `apps/ye-olde-magic-shop/saturday-cucumber/src/config/index.ts` | 1259 | 314 |
| `packages/saturday-cucumber/src/config/index.ts` | 1233 | 308 |
| `apps/ye-olde-magic-shop/.env` | 133 | 33 |
| `pnpm-workspace.yaml` | 105 | 26 |
| `scripts/ui-suite-in-chrome.sh` | 4579 | 1144 |
| `apps/ye-olde-magic-shop/tests/cart.spec.ts` | 1653 | 413 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `apps/ye-olde-magic-shop/saturday-core/** — vendored framework copy; comparing it file-by-file against packages/saturday-core/** is a developer-level diffing task, not an analyst-level one`
- `apps/ye-olde-magic-shop/saturday-playwright-heatmap/** — same reason as saturday-core`
- `pnpm-lock.yaml — large generated lockfile, regenerated mechanically by pnpm install, not meaningful to read for analysis`
- `apps/ye-olde-magic-shop/tests/*.spec.ts (example.spec.ts, heatmap.spec.ts, magic-shop.spec.ts, map.spec.ts, playwright-config.spec.ts, webgl.spec.ts) and features/*.feature — the existing golden suite AC3 must keep green; these are QA's before/after baseline, not analyst scope, and none of them changes under Out of Scope`
- `docs/features/ui-suite-real-chrome/** — only prior delivery found in the same bounded context (Test Automation); checked for a retrospective.md, none exists, so there is no lesson to carry forward`
