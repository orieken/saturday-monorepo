# QA report: reagent-map-line-tool

## How this order finished

The run got through security review. Its QA stage then failed twice. Each time, it started several
full suite runs at once, so three Cesium pages per run competed for software rendering and 10 tests
timed out. Then it ended its turn while waiting on a run, without returning its report, and the order
parked. The work was saved to `salvage/reagent-map-line-tool` and finished by hand.

## Change made during QA

`pickLandPoint` also checked `EllipsoidalOccluder.isPointVisible`. The QA stage showed no test can fail
on it: deleting it left every test green. `pickEllipsoid` returns the ray's first hit on the surface,
which always faces the camera, so the check can never return false. It was removed, and the comments in
the view and in `reagent-map-line-tool.spec.ts` now say why.

## Criteria

| Criterion | Test in `reagent-map-line-tool.spec.ts` |
|---|---|
| Two points complete a geodesic line, ready to save (AC6, AC8) | "placing two points on the land…", "a line from Thornveil Glade to Glimmerdeep Caves…" |
| One point: no completed line (AC7) | "placing only one point shows no completed line…" |
| Any near-side point can be placed, sea included | "a point can be placed anywhere on the visible surface, sea included" |
| A click beside the globe places nothing, and the tool still waits (AC9) | "clicking beside the globe places nothing…" |
| A click on the far edge places nothing (AC10) | "clicking the globe's far edge…" |
| Starting again clears the line (AC14) | "starting the line tool again discards any point already placed" |
| Draft entities carry `kind: 'draft-point'` and `'draft-line'` | "every placed point and the completed line are marked as drafts…" |
| No baseline changes unless on purpose | The four `@visual` tests pass unchanged. The new button sits above the globe, outside every captured element. |

## Test of the tests

Two runs in debian, each with a deliberate break in `ReagentMapView.vue`, then reverted:

| Break | Tests that failed |
|---|---|
| A: starting again keeps old points, and the line uses `ArcType.NONE` | AC14, AC8 |
| B: the second point never draws the line | AC6, the draft-kinds test, AC8 |

Run B also had one pan-zoom failure, the known flake below.

## Suite runs

Full suite, one run at a time, from the run's own code before the change above:

| Image | Run | Result |
|---|---|---|
| debian | 1 | 77 passed |
| debian | 2 | 76 passed, 1 failed: `reagent-map-pan-zoom.spec.ts` "panning and zooming move the camera, never the sites" |
| ubi9 | 1 | 77 passed |
| ubi9 | 2 | 77 passed |

The final runs with the change are in `pre-ship-checks.md`.

## Known gaps

- **No keyboard way to place a point.** The accessibility review found it: placing a point needs a
  mouse. It needs an interaction design, such as choosing sites from the list as ends, not a markup
  fix. Left for a follow-up order.
- **`#cesiumContainer` still has `role="img"`.** It now takes clicks, so "image" undersells it. Same
  follow-up.
- **Pan-zoom tests flake under 3 workers.** They are not part of this order. They get their own fix.
