# Implementation notes: reagent-map-visual-baselines

## How this order finished

The factory's developer stage ran into the 30-minute stage limit twice, and the order parked. By then the
stage had written the tests, the config, the script flag and both images' baselines. The work was saved to
`salvage/reagent-map-visual-baselines` and finished by hand. The verification runs alone take longer than
30 minutes: one suite run takes about two minutes, and the order asks for ten.

## What changed

- `tests/reagent-map-visual.spec.ts`: four `@visual` tests. They cover the page header with its nav link,
  the site list, the error state when `/api/reagent-sites` fails, and the globe at its default view. The
  globe test allows a 1% pixel difference (`maxDiffPixelRatio: 0.01`).
- `playwright.config.ts`: `snapshotPathTemplate` puts baselines under `.playwright-baselines/<image>/`.
  `resolveRunnerVariant()` reads `RUNNER_VARIANT`, and throws on a value that is not `debian` or `ubi9`.
- `playwright.config.ts`: `updateSnapshots: 'none'`. Playwright's default, `missing`, writes a missing
  baseline. The stage did not set this, so the first hand fix was to add it.
- `scripts/ui-suite-in-chrome.sh`: sets `RUNNER_VARIANT`, and takes `--update-baselines`, which passes
  `--update-snapshots=all`. The flag works alone or with `--ubi9`, in either order.
- `src/views/ReagentMapView.vue`: an alert when the sites fail to load. See "Outside scope" below.
- `docs/testing-cesium.md` and `docs/runner-image.md`: how and when to regenerate baselines.

## Outside scope

The order asks for a screenshot of the view's error state, but the view had none. It only threw. The stage
added a `role="alert"` element with `data-testid="reagent-map-error"`. That file is not in `scope.paths`, and
"changing the view" is out of scope. Oscar approved the design that named this element at confirm-design.

## Evidence

### Five runs in each image

| Image | Run | Exit | Result |
|---|---|---|---|
| debian | 1–5 | 0 | 64 passed each time |
| ubi9 | 1 | 0 | 64 passed |
| ubi9 | 2 | 1 | 62 passed, 2 failed: `map.spec.ts` geolocation (page.goto timeout), `reagent-map-pan-zoom.spec.ts` drag |
| ubi9 | 3, 4 | 0 | 64 passed |
| ubi9 | 5 | 1 | 63 passed, 1 failed: `reagent-map-pan-zoom.spec.ts` drag |

The four `@visual` tests passed in all ten runs. The criterion "the suite passes five runs in a row" is
not met in ubi9, because of two tests this order did not touch. The drag test missed by 312 px, as if
Cesium's inertia ran after the release. A later run in debian also failed a different pan-zoom test. These
flakes need their own fix.

### Test of the tests

One run per image, with two one-line changes to `ReagentMapView.vue`:

```diff
-        class="w-64 shrink-0 text-dnd-parchment"
+        class="w-64 shrink-0 text-dnd-parchment italic"
-const DEFAULT_VIEW = { lat: 10, lng: 20, heightMeters: 12_000_000 };
+const DEFAULT_VIEW = { lat: 10, lng: 25, heightMeters: 12_000_000 };
```

| Image | Site list | Globe default view |
|---|---|---|
| debian | failed: 2,792 pixels different (ratio 0.02) | failed: 13,637 pixels different (ratio 0.03) |
| ubi9 | failed: 2,792 pixels different (ratio 0.02) | failed: 13,637 pixels different (ratio 0.03) |

The header and error-state tests still passed. Both changes were then reverted.

### A missing baseline fails

With `debian/reagent-map-header.png` moved away, the header test failed with "A snapshot doesn't exist",
and the script did not write the file.

### Unit tests

- `scripts/ui-suite-in-chrome.test.sh`: all cases pass, including the new ones for `RUNNER_VARIANT` and
  `--update-baselines`.
- `tests/playwright-config.spec.ts`: `resolveRunnerVariant()` for unset, `debian`, `ubi9` and an unknown value,
  plus `updateSnapshots` and `snapshotPathTemplate`. These run in the suite above.
