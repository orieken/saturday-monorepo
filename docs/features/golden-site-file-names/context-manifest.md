# Context Manifest: golden-site-file-names

## Scope

- **Tier**: developer
- **Built for**: developer

## Pinned Files

- `apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts` — Renamed from MagicShop.ts and already carries the six rewritten page imports (./home.page, ./cart.page, ./product-details.page, ./login.page, ./account.page, ./reagent-map.page) — verified clean by grep this session; developer must confirm this state and report it as a rename (filesCreated), not re-edit it.
- `apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts` — Not renamed (already name.type.ts-shaped) but its ReagentMapPage import is already rewritten to '../reagent-map.page' — verify and report in filesModified, not filesCreated.
- `apps/ye-olde-magic-shop/features/support/world.ts` — Cucumber World import already rewritten to '../../lib/site/magic-shop.site' — distance-1 consumer of the renamed site file; verify and report in filesModified.
- `apps/ye-olde-magic-shop/features/step_definitions/webapp.steps.ts` — Import already rewritten to '../../lib/site/magic-shop.site.ts', preserving this file's existing explicit .ts extension style — verify and report in filesModified.
- `apps/ye-olde-magic-shop/tests/cart.spec.ts` — Playwright spec import already rewritten to '../lib/site/magic-shop.site' — verify and report in filesModified.
- `apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts` — Playwright spec import already rewritten to '../lib/site/magic-shop.site' — verify and report in filesModified.
- `.claude/feature-workspace/golden-site-file-names/state/developer.rejected.json` — The prior developer attempt's rejected payload — it did the renames and import edits correctly on disk but reported them as 'old -> new' pairs in filesModified, which loom rejected as unverifiable. Read this before reporting again: list each renamed file once by its new path in filesCreated per the spec's 'Reporting the renames' section, never an old path or an old -> new pair.

## Token Budget

- **Status**: OK — ~7544 tokens, 4.7% of the developer tier budget (160000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts` | 1673 | 418 |
| `apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts` | 924 | 231 |
| `apps/ye-olde-magic-shop/features/support/world.ts` | 828 | 207 |
| `apps/ye-olde-magic-shop/features/step_definitions/webapp.steps.ts` | 1231 | 307 |
| `apps/ye-olde-magic-shop/tests/cart.spec.ts` | 1659 | 414 |
| `apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts` | 14845 | 3711 |
| `.claude/feature-workspace/golden-site-file-names/state/developer.rejected.json` | 9025 | 2256 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `apps/ye-olde-magic-shop/lib/site/home.page.ts, cart.page.ts, login.page.ts, account.page.ts, product-details.page.ts, reagent-map.page.ts — renamed only, no internal imports to rewrite (confirmed in analysis.md); existence and filename already verified this session via `ls` and `git status`, so reading full file content adds nothing a developer doesn't already have.`
- `No shared/knowledge, .claude/knowledge, or docs/adrs directory exists in this repo — confirmed this session — so there is no KI/ADR corpus to search.`
- `apps/ye-olde-magic-shop/features/step_definitions/login.steps.ts, cart.steps.ts, checkout.steps.ts, inventory.steps.ts, sample.steps.ts — reference `this.site.<page>` by instance property name only, not by lib/site file path; unaffected by the rename.`
- `apps/ye-olde-magic-shop/lib/api/ShopApi.ts, lib/support/expeditions-lock.ts — different directories/concerns, out of scope for this rename.`
- `docs/features/golden-site-file-names/ — this run's own delivery record, owned by tech-writer downstream, not a developer input.`
