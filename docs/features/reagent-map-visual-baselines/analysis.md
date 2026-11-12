---
feature: "reagent-map-visual-baselines"
bounded_context: "Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)"
domain_terms: ["Reagent Site", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts", "apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/playwright.config.ts", "apps/ye-olde-magic-shop/.playwright-baselines/", "scripts/ui-suite-in-chrome.sh", "scripts/ui-suite-in-chrome.test.sh", "docs/testing-cesium.md", "docs/runner-image.md", "docs/features/reagent-map-visual-baselines/"]
issue_refs: ["reagent-map-visual-baselines"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: reagent-map-visual-baselines

## Summary

The visual-qa-engineer stage reports itself UNCONFIGURED for the reagent map because no Playwright screenshot baselines exist for it. This feature adds `@visual`-tagged `toHaveScreenshot()` tests for the map's stable DOM chrome (site list, page header/nav link, and a new error-state element) and one tolerant test of the Cesium canvas at its fixed default view, storing baselines per runner image (debian/ubi9) via a `RUNNER_VARIANT`-keyed `snapshotPathTemplate`. `scripts/ui-suite-in-chrome.sh` gains an `--update-baselines` flag to regenerate them explicitly; the tests skip on a developer Mac (no `RUNNER_VARIANT`) and never skip in the runner image. `ReagentMapView.vue` currently has no visible error-state element, so one must be added for the error-state screenshot to mean anything, without weakening the existing resilience test's assertions. Docs close the 'No visual baselines' known gap and describe when to regenerate baselines, including after a browser-lock bump.

### Acceptance Criteria

- `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` MUST hold `@visual`-tagged `toHaveScreenshot()` tests for the site list, the page header with its nav link, and the view's error state when `/api/reagent-sites` fails through route interception.
  - Example: page.route('**/api/reagent-sites', route => route.abort('failed')) — same pattern as reagent-map-resilience.spec.ts's REAGENT_SITES_API_PATTERN
  - Example: test('...', { tag: '@visual' }, async ({ page }) => { ... })
- One `@visual` test MUST cover `#cesiumContainer` at the fixed default view, after `cesium.waitForSettled()`, with a stated `maxDiffPixelRatio` that never exceeds 1%.
  - Example: await expect(page.locator('#cesiumContainer')).toHaveScreenshot('globe-default-view.png', { maxDiffPixelRatio: 0.01 })
- If the globe `@visual` test is flaky under the five-run repeat criterion, the developer MUST remove it, mask the canvas in the other `@visual` tests, and record why in the delivery record, instead of raising the tolerance past 1%.
- `playwright.config.ts` MUST store baselines per runner image under `.playwright-baselines/{debian,ubi9}/` using `snapshotPathTemplate`, keyed by a `RUNNER_VARIANT` environment variable that `scripts/ui-suite-in-chrome.sh` sets.
  - Example: snapshotPathTemplate: '.playwright-baselines/{arg}{ext}' where {arg} is built from resolveRunnerVariant() + the test's own snapshot name, mirroring resolveRunnerSlowMo()/resolveRunnerVideo()'s read-env-at-call-time pattern
- When `RUNNER_VARIANT` is unset, as on a developer's Mac, every `@visual` test MUST skip with a message stating why. When it is set, the tests MUST never skip, and a missing baseline MUST fail the test.
- `scripts/ui-suite-in-chrome.sh --update-baselines`, alone or combined with `--ubi9`, MUST regenerate that image's baselines by passing `--update-snapshots` to Playwright. Without the flag, the script MUST never write a baseline. `ui-suite-in-chrome.test.sh` MUST cover the flag.
- The suite MUST pass five consecutive runs in each image (debian and ubi9), and the run output MUST be recorded in the delivery record.
- A one-line style change to the site list MUST fail the matching `@visual` test in both images. A moved default camera on the globe test, if that test stays, MUST do the same. The failing output MUST be recorded, then the change reverted.
- `docs/testing-cesium.md` and `docs/runner-image.md` MUST document how and when to regenerate baselines, including after a browser lock bump. `testing-cesium.md`'s 'No visual baselines' Known Gaps entry MUST be closed. Both documents MUST follow `.claude/rules/writing-style.md`.
- Every new or changed function MUST have cyclomatic complexity below 7. Every changed line MUST be covered by a test that can fail. No file outside `scope.paths` MUST change, except this run's own delivery record under `docs/features/reagent-map-visual-baselines/`.
- `ReagentMapView.vue` MUST render a visible, stable error-state element when `fetchReagentSites()` rejects, since none exists today and the error-state screenshot in the first acceptance criterion needs something stable to capture.
  - Example: A `data-testid="reagent-map-error"` element with `role="alert"`, set via a new ref just before the existing `throw sitesResult.reason`, so the existing pageerror event, `newViewer.destroy()` call, and the zero-canvas/zero-list-item assertions in `reagent-map-resilience.spec.ts` stay intact.

### Non-Functional Requirements

- **performance**: The suite MUST keep running with `--workers=3`, `--shm-size=1g`, and the existing app-health timeout; adding visual baselines MUST NOT change any of these. (app health timeout 180seconds)
- **accessibility**: The new error-state element MUST be exposed to assistive technology consistently with the existing `role`/`aria-label` pattern already used on `#cesiumContainer` and the site list.
- **other**: Baseline PNGs MUST be generated and compared only inside the Saturday runner image (debian or ubi9), never on a developer's Mac, since the same page renders differently on each.

## Proposed Fitness Functions

- **Visual regression on the reagent map's stable UI (site list, header/nav link, error state)**
  - **Verification**: Playwright toHaveScreenshot() under --project=chrome, run via scripts/ui-suite-in-chrome.sh
  - **Owner**: DevOps
- **Baselines never silently regenerate outside an explicit --update-baselines run**
  - **Verification**: ui-suite-in-chrome.test.sh asserts --update-snapshots is passed to playwright only when --update-baselines is given
  - **Owner**: DevOps
- **Cyclomatic complexity below 7 on every changed function**
  - **Verification**: analyze-complexity tooling against the diff
  - **Owner**: Architect
- **Globe canvas screenshot stability across 5 repeated runs per image**
  - **Verification**: Manually recorded five-run pass in the delivery record per image; a one-time stability proof, not a per-commit CI gate
  - **Owner**: QA

## Out of Scope

- Baselines for the rest of the app.
- Heatmap data for the map (ordered separately as map-heatmap-source).
- Changing the Cesium view's behavior — camera, globe rendering, or interaction logic — beyond the minimal visible error-state element the error-state acceptance criterion requires.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)
- **Context Crossings**: Test Automation (Saturday Framework / Playwright @visual baselines), Craftsmanship Governance (visual-qa-engineer's UNCONFIGURED→configured transition), Infrastructure (Saturday runner image's debian/ubi9 variants — docs/runner-image.md)

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` — New spec file holding all @visual-tagged toHaveScreenshot() tests this feature adds.
- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Needs a visible, stable error-state element so the error-state acceptance criterion has something meaningful to screenshot; none exists today.
- `apps/ye-olde-magic-shop/playwright.config.ts` — Needs snapshotPathTemplate and a RUNNER_VARIANT-driven per-image baseline directory, plus the @visual skip-when-unset logic.
- `apps/ye-olde-magic-shop/.playwright-baselines/` — New committed baseline PNGs per runner image (debian, ubi9), written only via --update-baselines.
- `scripts/ui-suite-in-chrome.sh` — Needs an --update-baselines flag (composable with --ubi9) and must export RUNNER_VARIANT into the container.
- `scripts/ui-suite-in-chrome.test.sh` — Needs new assertions covering the --update-baselines flag, alone and combined with --ubi9.
- `docs/testing-cesium.md` — Must close the 'No visual baselines' Known Gaps entry and describe the regeneration procedure.
- `docs/runner-image.md` — Must document regenerating baselines, cross-referenced with the existing 'Updating the browsers (both locks)' section.
- `docs/features/reagent-map-visual-baselines/` — This run's own delivery record — the one path outside the literal scope.paths list that the final acceptance criterion explicitly permits to change.

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Add a visible, stable error-state element to ReagentMapView.vue (data-testid="reagent-map-error", role="alert"), set via a new ref just before the existing `throw sitesResult.reason`, preserving the existing pageerror event, newViewer.destroy() call, and the zero-canvas/zero-list-item assertions in reagent-map-resilience.spec.ts.
2. Create apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts with @visual-tagged toHaveScreenshot() tests: page header + nav link, site list, the new error-state element (via the same route-interception pattern as reagent-map-resilience.spec.ts), and #cesiumContainer at the default view after cesium.waitForSettled() with an explicit maxDiffPixelRatio <= 0.01.
3. Gate every @visual test on process.env.RUNNER_VARIANT: skip with a stated reason when unset, never skip (missing baseline fails) when set.
4. Add snapshotPathTemplate to playwright.config.ts storing baselines under .playwright-baselines/{debian,ubi9}/, reading RUNNER_VARIANT at call time the same way resolveRunnerSlowMo()/resolveRunnerVideo() already do.
5. Add --update-baselines to scripts/ui-suite-in-chrome.sh (composable with --ubi9): export RUNNER_VARIANT into the container in start_container(), and pass --update-snapshots to the chrome suite run only when the flag is set; never write a baseline otherwise.
6. Extend scripts/ui-suite-in-chrome.test.sh to cover --update-baselines alone and combined with --ubi9, asserting the docker exec log contains --update-snapshots only when the flag is passed and RUNNER_VARIANT is exported in both variants.
7. Run scripts/ui-suite-in-chrome.sh --update-baselines and --ubi9 --update-baselines once each to generate and commit the initial baselines under .playwright-baselines/{debian,ubi9}/.
8. Run the suite five consecutive times in each image and record the output for the delivery record.
9. Make and revert a one-line site-list style change (and, if the globe test stays, a moved default camera) in each image, recording the failing output for the delivery record.

### QA Tasks

1. Verify @visual tests skip (not fail) when RUNNER_VARIANT is unset, locally on a Mac.
2. Verify a missing baseline fails (not skips) when RUNNER_VARIANT is set, using a disposable run with a baseline file temporarily renamed.
3. Verify ui-suite-in-chrome.test.sh's new --update-baselines assertions pass.
4. Confirm no changed file falls outside scope.paths except the run's own delivery record.

### Tech Writer Tasks

1. Update docs/testing-cesium.md: close the 'No visual baselines' Known Gaps entry, describing what is covered and the regeneration procedure.
2. Update docs/runner-image.md: cross-reference the 'Updating the browsers (both locks)' section with an instruction to regenerate baselines in the same PR as a browser-lock bump.

### DevOps Tasks

None

## Edge Cases and Risks

- RUNNER_VARIANT is set to a value other than 'debian' or 'ubi9'.: Treat as an unrecognized variant and fail loudly rather than silently falling back to one baseline directory, since a silent fallback would compare against the wrong image's baselines.
- The globe @visual test is flaky across the five-run check.: Remove it, mask the canvas region in the other @visual tests, and record why — per the acceptance criterion's own fallback, not a tolerance increase past 1%.
- A baseline file is missing in the runner image (RUNNER_VARIANT set).: The test MUST fail, never skip — skipping is reserved for RUNNER_VARIANT being unset entirely.
- Both buildViewer() and fetchReagentSites() reject at once.: Existing onMounted logic already rethrows viewerResult's rejection first and never reaches the sitesResult branch; the new error-state element only renders on a sitesResult rejection after buildViewer() succeeds, matching the resilience spec's covered branch and the error-state acceptance criterion's literal scope (/api/reagent-sites failing).

## Definition of Done

- reagent-map-visual.spec.ts exists with @visual-tagged toHaveScreenshot() tests for the site list, the header/nav link, the error state, and (if kept) the globe default view.
- playwright.config.ts's snapshotPathTemplate stores baselines under .playwright-baselines/{debian,ubi9}/, keyed by RUNNER_VARIANT.
- @visual tests skip with a stated reason when RUNNER_VARIANT is unset, and never skip (missing baseline fails) when it is set.
- scripts/ui-suite-in-chrome.sh supports --update-baselines (alone and with --ubi9), writing a baseline only when passed; ui-suite-in-chrome.test.sh covers it.
- Both debian and ubi9 images pass five consecutive runs; the output is recorded in the delivery record.
- A one-line site-list style change, and a moved default camera if the globe test stays, each fail their matching @visual test in both images; the failing output is recorded and the change reverted.
- docs/testing-cesium.md's 'No visual baselines' Known Gaps entry is closed, and docs/runner-image.md documents regenerating baselines after a browser lock bump; both pass writing-style review.
- Every new or changed function has cyclomatic complexity below 7; every changed line is covered by a test that can fail; no file outside scope.paths changes except this run's own delivery record.
- visual-qa-engineer, run after this merges, finds .playwright-baselines/ and stops reporting UNCONFIGURED for this feature.
