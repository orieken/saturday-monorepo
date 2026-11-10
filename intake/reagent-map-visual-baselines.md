---
intake: 1
id: reagent-map-visual-baselines
type: chore
plan: deliver-feature
status: ready
priority: 2
created: 2026-11-04
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [reagent-map-globe, local-fonts]
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts, apps/ye-olde-magic-shop/.playwright-baselines/**, apps/ye-olde-magic-shop/playwright.config.ts, scripts/ui-suite-in-chrome.sh, scripts/ui-suite-in-chrome.test.sh, docs/testing-cesium.md, docs/runner-image.md, docs/features/reagent-map-visual-baselines/**]
---

# Give the reagent map screenshot baselines

## Request

> The visual QA stage reports itself as unconfigured, so nothing looks at the reagent map. A screenshot
> found a wrong default view and sites in the sea after every test had passed.

## Summary

Add Playwright screenshot baselines for the reagent map, so the `visual-qa-engineer` stage has
something to compare on every later slice. The same page renders differently on a Mac and in each
runner image, so baselines belong to an image and are only made and compared there.

## Acceptance Criteria

- [ ] Given `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts`, then it holds `toHaveScreenshot()`
      tests tagged `@visual` for what is stable: the site list, the page header with its nav link, and
      the view's error state when `/api/reagent-sites` fails (through route interception).
- [ ] Given the globe, then one `@visual` test covers `#cesiumContainer` at the fixed default view with a
      stated `maxDiffPixelRatio`, once settled. If it is flaky under the repeat-run criterion below,
      remove it, mask the canvas in the other tests, and record why in the delivery record. Do not
      raise the tolerance past 1% to make it pass.
- [ ] Given `playwright.config.ts`, then baselines are stored per image under
      `.playwright-baselines/{debian,ubi9}/` through `snapshotPathTemplate`, using a `RUNNER_VARIANT`
      environment variable that `scripts/ui-suite-in-chrome.sh` sets.
- [ ] Given `RUNNER_VARIANT` unset, as on a developer's Mac, then the `@visual` tests skip with a message
      saying why. Given it set, they never skip: a missing baseline fails.
- [ ] Given `scripts/ui-suite-in-chrome.sh --update-baselines`, alone or with `--ubi9`, then it regenerates
      that image's baselines. Without the flag, the script never writes one. Its test script covers the flag.
- [ ] Given both images, then the suite passes five runs in a row in each. Record the output.
- [ ] Given a one-line style change to the site list, then the matching `@visual` test fails in both
      images, and the same for a moved default camera on the globe test if it stays. Record the failing
      output, and revert the change.
- [ ] Given `docs/testing-cesium.md` and `docs/runner-image.md`, then they say how and when to regenerate
      baselines, including after a browser bump, and `testing-cesium.md`'s "No visual baselines" gap is
      gone. Both still follow `.claude/rules/writing-style.md`.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under
      `docs/features/reagent-map-visual-baselines/` excepted.

## Out of Scope

- Baselines for the rest of the app.
- Heatmap data for the map: order `map-heatmap-source`.
- Changing the view itself.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- The suite still runs with `--workers=3`, `--shm-size=1g`, and the same 180-second app health limit.

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Playwright `@visual` tests, through `scripts/ui-suite-in-chrome.sh`, in both images.

## Open Questions

- A browser bump changes rendering, so baselines must move with the locks: Assumption if unresolved:
  regenerate them in the same PR as the lock change, as `docs/runner-image.md` will say.

## Infrastructure / Deploy Notes

- After this merges, `visual-qa-engineer` finds `.playwright-baselines/` and stops reporting UNCONFIGURED.
