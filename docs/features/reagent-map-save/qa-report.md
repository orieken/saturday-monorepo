# QA report: reagent-map-save

## How this order finished

The run passed security review. Its QA stage then hit the 30-minute stage limit and was stopped,
so it reported no spend, and the order parked on an unmeasured budget. By then the stage had written
`reagent-map-save.spec.ts` and the lock described below. The work was saved to
`salvage/reagent-map-save`, replayed onto `main`, and finished by hand.

## Change made during QA

The double-click test clicked the save button twice in parallel. The first save removes the button,
so the second `click()` waited for it until the test timed out. That happened once in each image in
three-times-repeated runs. The test now uses one real `dblclick()`, which sends both clicks before the
first save returns.

## Keeping the writers away from the readers

Saving writes `apps/mock-api/data/expeditions.json`. `reagent-map-api.spec.ts` and
`reagent-map-save.spec.ts` write it, and `reagent-map-expeditions.spec.ts` counts what it holds. Each
of the three runs its tests serially and holds one file lock, `reagent-map-expeditions.lock` in the
OS temporary directory, from `beforeAll` to `afterAll`. A lock older than 180 s counts as stale. The
two writers put the seed file back after each test.

Each file keeps its own copy of the lock functions. Importing them from another spec file would also
run that file's `test.describe`, which Playwright rejects.

## Criteria

| Criterion | Test in `reagent-map-save.spec.ts` |
|---|---|
| Elara saves, reopens the map, sees it; the API holds it within 1 km (AC11) | "Elara, signed in, sees her saved expedition again…" |
| Signed out: asked to sign in, nothing stored, no request sent (AC12) | "saving without being signed in…" |
| Zero-length line: the shop's message, nothing stored (AC13) | "saving a zero-length line…" |
| Unsaved line gone after leaving and reopening (AC15) | "an unsaved line is gone…" |
| A double-click saves once | "rapid double-clicking save…" |
| `.env` turns the map on | `.env` sets `VITE_FEATURE_REAGENT_MAP=true`; `reagent-map-feature-flag.spec.ts` still covers the flag turning it off |
| No baseline changed | The four `@visual` tests pass unchanged |

## Test of the tests

Each break below was made in `ReagentMapView.vue`, run against the save spec in debian, and then reverted.

| Break | Result |
|---|---|
| No sign-in check | AC12 failed |
| Saves `from, from` instead of `from, to` | AC11 failed |
| Shows a fixed message instead of the server's | AC13 failed |
| No "already saving" check in `saveLine` | Nothing failed: the button is also disabled while saving |
| No "already saving" check and no `disabled` while saving | Double-click failed: 2 saves, expected 1 |

## Suite runs

| Run | debian | ubi9 |
|---|---|---|
| Full suite | 82 passed | 82 passed |
| Reagent map specs, 3 workers, `--repeat-each=3`, with the fixed test | 153 passed | 153 passed |

The isolation run before the test fix also had four load timeouts in debian, all in existing specs
at the start of the run: tiles not loaded within 15 s, and a header screenshot not stable within 5 s.
The re-run had none.

## Known gaps

- After a save, the button disappears and keyboard focus falls back to the page (accessibility review).
- Nothing on the server limits saves from outside the UI (security review, existing).
