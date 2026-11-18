# Visual QA Report: reagent-map-save

**Summary**: ⚠️ PASS WITH RECOMMENDATIONS

*(Template `shared/templates/visual-qa-report.template.md` does not exist in this repo — this report follows the section structure from the agent's own spec instead.)*

## Invocation Conditions

- `heatmap-data/` exists: `apps/ye-olde-magic-shop/heatmap-data/heatmap-spec-ts__heatmap-verification.json` — present.
- Playwright visual baselines exist: `apps/ye-olde-magic-shop/.playwright-baselines/{debian,ubi9}/*.png` (4 baselines per image) — present.

Both conditions are met, so this is a full analysis, not an UNCONFIGURED skip.

## Heatmap Coverage Analysis

`@orieken/saturday-ml-analyzer` is **not installed** for this app — it's not a dependency of `apps/ye-olde-magic-shop/package.json`, and `packages/saturday-ml-analyzer` has no built `dist/` (`npm run build` has never been run there). `node node_modules/@orieken/saturday-ml-analyzer/dist/index.js heatmap-data` cannot execute. No Coverage Score, Hotspots, or Cold Spots can be computed by tooling.

More importantly, the one file present in `heatmap-data/` is not evidence about this feature. It is `heatmap-spec-ts__heatmap-verification.json`, produced by `tests/heatmap.spec.ts` — a generic fixture smoke test ("heatmap verification") that clicks the site title on the home page (`/`). It contains a single click event and never navigates to `/reagent-map`. No scenario in `reagent-map-save.spec.ts` or `reagent-map-visual.spec.ts` captures heatmap data; the `@orieken/saturday-playwright-heatmap` fixture is wired into `tests/heatmap.spec.ts` only, not into the reagent-map suite.

**Finding (HIGH — manual read of the DOM, since tooling can't run):** the new primary-journey elements this feature ships —
- `save-line-tool` (the Save button itself, a primary CTA)
- `line-tool-status`, `save-line-message`, `save-line-status` (its feedback states)

— have zero heatmap interaction capture of any kind. This isn't a Cold Spot the analyzer flagged; it's an absence of instrumentation on the feature entirely. Recommendation: either wire the heatmap fixture into `reagent-map-save.spec.ts`'s existing scenarios, or explicitly accept this feature as outside heatmap coverage — but that should be a decision, not a silent gap.

## Visual Regression Results

Ran `npx playwright test --grep "@visual"` locally (outside `RUNNER_VARIANT`): all 4 tests in `tests/reagent-map-visual.spec.ts` correctly **skipped** — confirms the intended guard (`test.skip(!process.env.RUNNER_VARIANT, ...)`) behaves as documented; these only run under `scripts/ui-suite-in-chrome.sh`.

Did not re-run the full Docker-based runner-image suite — qa-engineer already ran `--update-baselines` for both `debian` and `ubi9` images this session and recorded in `qa-report.md`: all 4 `@visual` tests passed in both images, and `git diff` showed **zero baseline PNG bytes changed** in either image.

I verified why that's the expected outcome, not a coincidence: read `src/views/ReagentMapView.vue` and confirmed the four baselined elements —
- `header` (`reagent-map-header.png`)
- `reagent-site-list` (`reagent-site-list.png`)
- `reagent-map-error` (`reagent-map-error.png`)
- `#cesiumContainer` (`globe-default-view.png`)

— are siblings of, not ancestors of, the new save-tool toolbar (`start-line-tool`, `save-line-tool`, `line-tool-status`, `save-line-message`, `save-line-status`), which lives in its own `<div class="mb-8 ...">` block between the error banner and the globe/list row. None of the 4 baselined regions contain any of the new elements, so a 0-byte diff is exactly what should happen — this is confirmation, not a surprise finding.

**No screenshot diffs detected.** No regression.

**Finding (INFO — recommendation, not a failure):** no visual baseline exists yet for the new save-tool toolbar itself (default/empty state, mid-draw state with Save enabled, success state, error state). The feature shipped without visual regression coverage for its own new UI, even though the surrounding page is fully baselined. Nothing to diff against today, so this can't FAIL, but it's a gap worth closing before the next change to that toolbar ships blind.

## Cold Spots

| Element | Journey role | Severity |
|---|---|---|
| `save-line-tool` | Primary CTA (save an expedition) | HIGH — no heatmap capture exists for it at all |
| `line-tool-status` / `save-line-message` / `save-line-status` | Primary-journey feedback (draw/save state, error, success) | HIGH — same |

No decorative/read-only elements were assessed since no computable heatmap run exists to source them from.

## Recommendations

1. Add the `@orieken/saturday-playwright-heatmap` fixture to `reagent-map-save.spec.ts` (or a dedicated scenario) so the save flow's interaction coverage is actually measured, not just functionally tested.
2. Build `packages/saturday-ml-analyzer` (`pnpm --filter @orieken/saturday-ml-analyzer build`) and add it as a devDependency of `ye-olde-magic-shop` if heatmap Coverage Score is meant to gate this app going forward — right now the tool can't run here at all.
3. Add `@visual` baseline coverage for the save-tool toolbar's four states (idle, drawable, saving, saved/error) in a future pass, paired with `visual-qa-engineer` review before merge.

None of the above block this delivery — no regression was found in what IS baselined, and the new UI simply predates any visual/heatmap instrumentation being pointed at it.