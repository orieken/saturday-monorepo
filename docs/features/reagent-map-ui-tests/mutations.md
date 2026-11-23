# E1 judge: reagent-map-ui-tests

Order `docs/specs/reagent-map-spec/orders/reagent-map-ui-tests.md`, run by hand. Spec package
`sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413`.

Suite: `apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts`. It has one test for each of the 15
`@level-ui` scenarios in `features/reagent-map/`, and AC2's outline gives two, so 16 tests. Each test
title starts with its criterion tag and uses the scenario's own name.

## How the tests act and read

- **Gestures:** `cesium.clickAt`, `cesium.flyTo` and `cesium.toScreen` from
  `saturday-playwright-cesium`, and the map's own buttons. Nothing draws through Cesium's API (D-5).
- **Drawing a line:** `DrawLineFlow` in `lib/site/flows/draw-line.flow.ts`. It chooses the line tool,
  places A, places B, and saves if asked. The map is `ReagentMapPage` (D-6). Both are registered in
  `MagicShop`.
- **Positions:** checked within 1 km, with the radius read from the globe (D-10, D-14).
- **AC11:** reads the saved expedition back through `ShopApi`, the cross-check.
- **Shared data:** the suite saves to `apps/mock-api/data/expeditions.json`, and AC3 and AC4 count
  what it holds. It takes the shared expeditions lock in `lib/support/expeditions-lock.ts`, which uses
  the same path as the other three specs' copies. It runs in one worker and puts the file back after
  each test. It uses `mode: 'default'`, not `'serial'`, so one failure does not skip the rest.

## Faults

Each fault was applied alone, the suite was run in `saturday-runner:debian`, and the fault was
reverted. Faults M1 to M3 are in the mock API. M4 to M10 are in `ReagentMapView.vue`, because what
they break exists only in the view.

| Fault | Scenarios that failed |
|---|---|
| M1 Sites published about 5 km from their place | AC1, AC2 (both examples) |
| M2 A seeded expedition served about 5 km from its sites | AC3, AC4 |
| M3 A saved line stored with one end about 5 km away | AC11 |
| M4 The line tool ignores the second point | AC6, AC8, AC11, AC12, AC13, AC15 |
| M5 Flying to a site does nothing | AC5 |
| M6 One point can be saved | AC7 |
| M7 A click that misses the globe places a point at 0, 0 | AC9, AC10 |
| M8 Saving skips the sign-in check | AC12 |
| M9 Starting a new line keeps the old points | AC14 |
| M10 An unsaved line is kept in memory and redrawn when the map reopens | AC15 |

Every scenario failed on at least one fault, and no fault survived. The order names four faults: an
endpoint moved 5 km (M3), a seeded expedition at the wrong place (M2), a site at the wrong position
(M1), and the line tool ignoring the second point (M4). M5 to M10 cover the scenarios those four
cannot reach.

M10 needs the reagent map to be left through the app's own navigation. A full page load clears
memory, so a test that reloads would miss it. AC15 leaves through the Shop link and comes back
through the Reagent Map link.

The faults were applied by a throwaway script, not committed.

## Suite runs

| Run | debian | ubi9 |
|---|---|---|
| Full suite | 99 passed | 99 passed |
| Every reagent map spec, 3 workers, `--repeat-each=3` | 201 passed | 201 passed |
