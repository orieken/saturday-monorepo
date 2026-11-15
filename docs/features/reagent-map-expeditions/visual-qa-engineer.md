Clean. Now producing the final artifact.

```markdown
---
feature: "reagent-map-expeditions"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/.playwright-baselines/debian/globe-default-view.png", "apps/ye-olde-magic-shop/.playwright-baselines/ubi9/globe-default-view.png"]
issue_refs: ["reagent-map-expeditions"]
linked_adrs: []
linked_kis: []
---

# Visual QA Report: reagent-map-expeditions

## Summary

**PASS** (visual regression) — the one changed baseline is an intentional, verified-correct update, not a regression. **UNCONFIGURED** (interaction-heatmap coverage) — `heatmap-data/` exists but holds no capture for any scenario in this feature; see below.

## Visual Regression Results

`apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` carries the feature's 4 `@visual`-tagged scenarios (page header/nav, site list, error state, globe default view), each stored per runner image under `.playwright-baselines/{debian,ubi9}/` (Cesium needs a real WebGL2 context, so these only run inside `scripts/ui-suite-in-chrome.sh`, gated on `RUNNER_VARIANT`).

- **Baseline change**: `globe-default-view.png` changed in both `debian/` and `ubi9/` baseline directories; `reagent-map-header.png`, `reagent-site-list.png`, and `reagent-map-error.png` are untouched in either — confirmed directly via `git diff --stat` on `.playwright-baselines/`, independent of qa-report.md's own account.
- **What changed, visually**: a single new white geodesic line now connects two site markers (Moonwell Spring and Ember Crag) on the globe's near-side hemisphere — the "Spring to Crag" seeded expedition, the only one of the three seeded expeditions whose both endpoints fall inside this fixed default view. Opened the new `debian/globe-default-view.png` directly and confirmed the line renders cleanly: no z-fighting, no broken segments, no stray artifacts elsewhere on the sphere.
- **Cross-runner consistency**: `sha256sum` on both new baselines returns the identical digest (`ddd67d85...`) — the debian and ubi9 variants produced byte-identical PNGs, consistent with both rendering through the same SwiftShader software GL driver, and consistent with the pre-change baselines also having been byte-identical.
- **Tolerance discipline**: `GLOBE_MAX_DIFF_PIXEL_RATIO` in the spec is `0.01` and was not raised to accommodate this change — the new baseline itself absorbs the new line, not a widened threshold.
- **Independent live run — debian**: ran the full suite (`scripts/ui-suite-in-chrome.sh`, no `--update-baselines`) against the already-committed baselines. 69/69 passed, including all 4 `@visual` tests — the current implementation's globe render matches the committed baseline within the 1% tolerance.
- **Independent live run — ubi9**: a first attempt run concurrently with the debian suite failed at the dependency-install step (`ENOENT` on a `vitest` binary) — both containers' `pnpm install --frozen-lockfile` were racing against the same bind-mounted checkout at once, not a defect in this feature's code. Stopped that container, re-ran ubi9 alone end-to-end (fresh install → build → app health → `runner-check` → suite), scoped this time to `--grep '@visual'` only (`--workers=1`) to isolate the visual signal from `tests/map.spec.ts`'s already-documented resource-contention flakiness (qa-report.md's Notes for Tech Writer section). Result: **4/4 `@visual` tests passed cleanly** in 19.1s, including the globe test.
- **No new or unexpected diffs** were found on either runner image. The only visual change present anywhere in this delivery is the intended one, and it reproduces correctly from the current code.

## Interaction Heatmap Coverage Analysis

`heatmap-data/` exists (`@orieken/saturday-playwright-heatmap` is wired into the app), and `@orieken/saturday-ml-analyzer` ran successfully against it after building it from source (its `dist/` wasn't checked in):

```
Test: heatmap verification (passed)
  Coverage Score: 1.5%
  Interactables: 66
  Hotspots (Clusters): 1
  Cold Spots (Untested): 65
```

This number is **not meaningful for this feature** and is reported here only to explain why it's being discounted rather than scored. The single capture file (`heatmap-spec-ts__heatmap-verification.json`) comes from `tests/heatmap.spec.ts` — a generic, pre-existing fixture-smoke-test that clicks the shop header logo and the first link on the storefront homepage (`/`). It predates this feature, never navigates to `/reagent-map`, and isn't one of this delivery's test files. None of the reagent-map-expeditions test files (`reagent-map-expeditions.spec.ts`, `reagent-map-fly-to-site.spec.ts`, `reagent-map-visual.spec.ts`, `reagent-map-globe.spec.ts`, `reagent-map-pan-zoom.spec.ts`) import the `@orieken/saturday-playwright-heatmap` fixture — confirmed by grep across `tests/reagent-map*.spec.ts`.

In plain terms: this feature has not adopted heatmap instrumentation in its own scenarios, so there is no cold-spot or coverage-gap signal to report for it. This is the expected "project hasn't adopted the fixture for this feature yet" case, not a coverage failure — it does not block the pipeline.

## Cold Spots

None scoped to this feature (see above — no heatmap capture exists for any reagent-map scenario). The fly-to-site buttons and site-list — the interactive primary-journey elements this feature adds — were separately verified keyboard-accessible by qa-engineer's ad-hoc check (qa-report.md's Accessibility Check section); that is a functional/a11y verification, not a heatmap-derived coverage signal, and is out of this report's scope.

## Recommendations

- If heatmap-derived coverage gap detection is wanted for this feature going forward, wire the `@orieken/saturday-playwright-heatmap` fixture into one or more of the `reagent-map-*.spec.ts` files (the existing `tests/heatmap.spec.ts` is a reasonable pattern to copy). Not required for this delivery.
- `packages/saturday-ml-analyzer` ships without a built `dist/` in this checkout — this report only had a usable analyzer because it was built locally (`node_modules/.bin/tsc`) for this run. Worth a build step somewhere in CI or a `prepare` script so the next person running this agent doesn't have to rediscover that.

## Known Issues / Out of Scope

- `tests/map.spec.ts` (Leaflet `/locations` page, unrelated to this feature) and `tests/reagent-map-pan-zoom.spec.ts`'s camera-move assertions showed intermittent timeouts under `--workers=3` on one of this session's ubi9 runs — this is the same resource-contention flakiness qa-engineer already diagnosed and documented in qa-report.md (software-rendered WebGL pages contending under parallel workers), reproduced again here, not a visual-regression finding and not caused by this feature's changes.
```