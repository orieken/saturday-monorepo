---
intake: 1
id: local-fonts
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
depends-on: [reagent-map-globe]
budget:
  max-attempts: 2
  max-usd: 15
  max-hours: 2
scope:
  paths: [apps/ye-olde-magic-shop/src/styles.css, apps/ye-olde-magic-shop/src/main.ts, apps/ye-olde-magic-shop/package.json, apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts, pnpm-lock.yaml, docs/testing-cesium.md, docs/features/local-fonts/**]
---

# Serve the shop's fonts from its own origin

## Request

> The no-request test skips Google Fonts because every page imports them. Serve the fonts locally and
> remove the exception.

## Summary

`apps/ye-olde-magic-shop/src/styles.css` imports Cinzel and Inter from `fonts.googleapis.com` on every
page. So the app reaches the network on every page, and `tests/reagent-map-globe.spec.ts` has to
skip two hosts. Fonts loaded over the network can also make screenshots differ between runs, so this
comes before `reagent-map-visual-baselines`.

## Acceptance Criteria

- [ ] Given the app, then Cinzel (400, 700) and Inter (400, 600) load from the app's own origin,
      through `@fontsource/cinzel` and `@fontsource/inter` pinned to exact versions, and
      `src/styles.css` has no `@import` from another origin.
- [ ] Given `tests/reagent-map-globe.spec.ts`, then `PRE_EXISTING_OFF_ORIGIN_HOSTS` and the code that
      skips those hosts are gone, and the no-request test passes with no exception.
- [ ] Given the no-request test with the old `@import` put back, then it fails. Record the failing
      output in the delivery record.
- [ ] Given a page opened in the runner image, then headings still render in Cinzel and body text in
      Inter, checked through `document.fonts.check` for each family and weight in a test.
- [ ] Given `docs/testing-cesium.md`, then its Google Fonts entry under "Known gaps" is gone, and it
      still follows `.claude/rules/writing-style.md`.
- [ ] Given `scripts/ui-suite-in-chrome.sh` and `scripts/ui-suite-in-chrome.sh --ubi9`, then both pass.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/local-fonts/` excepted.

## Out of Scope

- Changing which fonts the shop uses, or any other styling.
- The map tiles the Leaflet locations view loads from OpenStreetMap: order `leaflet-local-tiles`.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- No page in the app requests a font from another origin.

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Playwright, through `scripts/ui-suite-in-chrome.sh`, in both images.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
