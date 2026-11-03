# Draft order E1: reagent-map-ui-tests (by hand)

| Field | Value |
|---|---|
| plan | `deliver-tests-saturday` (by hand; the plan does not exist yet, so not in `intake/`, D-27) |
| status | draft |
| depends on | `reagent-map-save` (the last D1 slice), `saturday-playwright-cesium` (C1) |
| spec package | `sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413` (`package.sha256`) |
| gates | testability → developer; model, ship → tester |

## Scope

The 15 `@level-ui` scenarios in `features/reagent-map/`, as Playwright Test specs in
`apps/ye-olde-magic-shop/tests/`, running in the runner image through `scripts/ui-suite-in-chrome.sh`.

## Rules

- Act through `saturday-playwright-cesium` gestures (`clickAt`, `flyTo`, `toScreen`) and read through
  its state (`getEntities`, `waitForSettled`); never a helper that draws through Cesium's API (D-5).
- "Choose the line tool, place A, place B, save" is a Saturday Flow, `draw-line.flow.ts`, in the
  app's test site; the reagent map is a Page (D-6).
- Positions are checked against the spec's 1 km tolerance using the planet's radius from the seam,
  never Earth's (D-10, D-14).
- AC11 reads the saved expedition back through the API: the cross-check.

## Judge

Every scenario fails on at least one injected fault (prompt 05, Judges), injected through mock-API
route interception: a line saved with an endpoint moved 5 km, a seeded expedition at the wrong place,
a site's metadata at the wrong position, and the line tool ignoring the second point.
