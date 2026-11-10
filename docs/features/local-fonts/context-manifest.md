# Context Manifest: local-fonts

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/ye-olde-magic-shop/src/styles.css` — Core change: remove the cross-origin @import of Google Fonts (Cinzel/Inter) — the AC requires styles.css to have no @import from another origin
- `apps/ye-olde-magic-shop/src/main.ts` — Single entry point loaded by every route; the self-hosted @fontsource/cinzel and @fontsource/inter CSS imports (400/700 and 400/600) belong here
- `apps/ye-olde-magic-shop/package.json` — Add @fontsource/cinzel and @fontsource/inter as exact-pinned dependencies, per the AC's 'pinned to exact versions' requirement
- `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` — Remove PRE_EXISTING_OFF_ORIGIN_HOSTS and the skip branch for fonts.googleapis.com/fonts.gstatic.com so the no-request test runs with no exception
- `apps/ye-olde-magic-shop/tests/local-fonts.spec.ts` — Target file for the new test (does not exist yet) asserting Cinzel/Inter load via document.fonts.load at the exact weight, per the AC's explicit document.fonts.check prohibition
- `docs/testing-cesium.md` — Remove the 'Still open' Google Fonts note near the Cesium-config table and the matching 'Known gaps' entry once fonts are self-hosted

## Token Budget

- **Status**: INCOMPLETE — ~5976 tokens, 5.0% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/src/styles.css` | 533 | 133 |
| `apps/ye-olde-magic-shop/src/main.ts` | 581 | 145 |
| `apps/ye-olde-magic-shop/package.json` | 2226 | 556 |
| `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` | 8589 | 2147 |
| `apps/ye-olde-magic-shop/tests/local-fonts.spec.ts` | — | not measured: cannot read apps/ye-olde-magic-shop/tests/local-fonts.spec.ts |
| `docs/testing-cesium.md` | 11980 | 2995 |

> 1 pinned file(s) could not be measured, so the total above is a floor, not a figure.

## Knowledge Items

None

## ADRs

None

## Pruned

- `pnpm-lock.yaml — in scope.paths but a 10,029-line generated lockfile; regenerate with `pnpm install` after adding the two @fontsource deps rather than reading or hand-editing it`
- `.claude/factory/claims/local-fonts.md — this is the feature spec itself, already handed to the analyst directly; not re-pinned here`
- `apps/ye-olde-magic-shop/scripts/generate-aurelis-image.py and scripts/ui-suite-in-chrome.sh — existing verification tooling the AC requires to pass, not to change; no edit needed`
- `apps/mock-api/** and reagent-sites-api files — different bounded context (Expedition Planning data layer); this chore stays inside the frontend's font-serving and its own tests`
- `blastRadius omitted — mcp__loom__impact/dependents_of/find_files were not available in this session; the spec's scope.paths and the known salvage-branch diff (83beb6e..3629f44) already enumerate the full touched-file set for this narrow chore, so a graph walk adds nothing here`
