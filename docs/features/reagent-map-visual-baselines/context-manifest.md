# Context Manifest: reagent-map-visual-baselines

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` — Only existing 'chrome'-project spec for this view — shows the test fixture imports (`@orieken/saturday-playwright-cesium`), the SEEDED_SITES fixture, and the waitForSettled() pattern the new @visual spec must follow to avoid flakiness.
- `apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts` — Existing route-interception pattern (page.route(...).abort('failed')) for simulating /api/reagent-sites failure — the direct model for the new error-state @visual test the spec requires.
- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — The view under test. onMounted currently rethrows on fetch failure with no visible error UI rendered (just an unhandled pageerror) — the analyst must decide whether an error-state element needs to be added before it can be screenshotted, since none exists today.
- `apps/ye-olde-magic-shop/src/App.vue` — Defines the nav link (data-testid="reagent-map-link", lines 12-17) that the page-header screenshot AC covers, including its v-if on the feature flag.
- `apps/ye-olde-magic-shop/playwright.config.ts` — Must gain snapshotPathTemplate and the RUNNER_VARIANT-driven per-image baseline directory; already shows the existing 'chrome' project and resolveRunnerSlowMo()/resolveRunnerVideo() pattern for reading env vars at call time, which the new RUNNER_VARIANT skip logic should match.
- `scripts/ui-suite-in-chrome.sh` — Needs a new --update-baselines flag (and --ubi9 composability) plus RUNNER_VARIANT export into the container; shows the existing parse_args/exit-code contract this must extend.
- `scripts/ui-suite-in-chrome.test.sh` — Existing stubbed-docker test harness pattern the spec requires be extended to cover the new --update-baselines flag; shows the fake_docker dispatch style to follow.
- `docs/testing-cesium.md` — Contains the 'No visual baselines' Known Gaps entry (around line 309) the spec requires be closed, and the existing writing-style-governed doc structure to extend.
- `docs/runner-image.md` — Must document when to regenerate baselines (including after a browser lock bump) alongside the existing 'Updating the browsers (both locks)' section it should cross-reference.
- `docs/features/reagent-map-globe/visual-qa-engineer.md` — Direct precedent from the prior feature in this same bounded context: its Recommendations section already concluded the Cesium canvas is a poor pixel-diff target and that screenshotting the DOM-level chrome (site list, nav link, loading/error states) around it is more stable — exactly the approach this spec's ACs take.

## Token Budget

- **Status**: OK — ~18967 tokens, 15.8% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` | 7841 | 1960 |
| `apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts` | 2444 | 611 |
| `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` | 5626 | 1406 |
| `apps/ye-olde-magic-shop/src/App.vue` | 2908 | 727 |
| `apps/ye-olde-magic-shop/playwright.config.ts` | 2630 | 657 |
| `scripts/ui-suite-in-chrome.sh` | 5746 | 1436 |
| `scripts/ui-suite-in-chrome.test.sh` | 6948 | 1737 |
| `docs/testing-cesium.md` | 11645 | 2911 |
| `docs/runner-image.md` | 23746 | 5936 |
| `docs/features/reagent-map-globe/visual-qa-engineer.md` | 6344 | 1586 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `apps/ye-olde-magic-shop/tests/reagent-map-pan-zoom.spec.ts — drag/zoom gesture coverage, no visual-baseline relevance`
- `apps/ye-olde-magic-shop/tests/reagent-map-api.spec.ts — API contract tests; change-surface exclusion (no data model or API changes in this spec)`
- `apps/ye-olde-magic-shop/tests/reagent-map-a11y-check.spec.ts — accessibility checks, different concern; nav-link/list testids it uses are already visible via App.vue/ReagentMapView.vue`
- `apps/ye-olde-magic-shop/tests/reagent-map-feature-flag.spec.ts — pure feature-flag unit logic, unrelated to screenshot baselines`
- `apps/mock-api/** — backend; change-surface exclusion, no API changes in scope.paths`
- `containers/saturday-runner/Containerfile and lock-packages.sh — image build internals; out of scope.paths, and docs/runner-image.md already summarizes what's needed`
- `docs/features/local-fonts/visual-qa-engineer.md — same UNCONFIGURED-report pattern as reagent-map-globe's but a different, less visually-relevant bounded context (font swap, not the reagent map UI)`
- `apps/ye-olde-magic-shop/src/api/reagentSites.ts — thin fetch wrapper; fetchReagentSites() rejection behavior is already visible through ReagentMapView.vue's onMounted and reagent-map-resilience.spec.ts`
