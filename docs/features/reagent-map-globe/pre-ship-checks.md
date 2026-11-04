# Pre-Ship Checks

## ui-suite

- Command: `scripts/ui-suite-in-chrome.sh`
- Duration: 1m58.372997625s
- Status: exit 0

```

[dotenv@17.2.3] injecting env (0) from .env -- tip: ✅ audit secrets and track compliance: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: 🔄 add secrets lifecycle management: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: ⚙️  enable debug logging with { debug: true }
[dotenv@17.2.3] injecting env (0) from .env -- tip: ⚙️  write to custom object with { processEnv: myObject }
[dotenv@17.2.3] injecting env (0) from .env -- tip: 👥 sync secrets across teammates & machines: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: ✅ audit secrets and track compliance: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: 👥 sync secrets across teammates & machines: https://dotenvx.com/ops
  ✓   2 [chrome] › tests/magic-shop.spec.ts:3:5 › magic shop loads and displays wares (5.0s)
  ✓   6 [chrome] › tests/map.spec.ts:9:7 › Ye Olde Magic Shop - Map Locations › map should load and center on the correct region (5.1s)
  ✓   4 [chrome] › tests/map.spec.ts:17:7 › Ye Olde Magic Shop - Map Locations › markers should exist and popups should display correct information (5.6s)
  ✓   3 [chrome] › tests/example.spec.ts:5:5 › homepage smoke test @k6 (5.8s)
  ✓   7 [chrome] › tests/cart.spec.ts:29:7 › Shopping Cart › Increase quantity from the cart page (6.7s)
  ✓   5 [chrome] › tests/cart.spec.ts:12:7 › Shopping Cart › Add the same item twice increases quantity to 2 (6.5s)
  ✓   1 [chrome] › tests/heatmap.spec.ts:4:5 › heatmap verification (7.5s)
  ✓   9 [chrome] › tests/map.spec.ts:48:7 › Ye Olde Magic Shop - Map Locations › should check point in polygon (3.5s)
  ✓   8 [chrome] › tests/map.spec.ts:27:7 › Ye Olde Magic Shop - Map Locations › should be able to draw a line and find markers along it (3.6s)
  ✓  15 [chrome] › tests/playwright-config.spec.ts:33:7 › chrome project configuration (AC1) › resolves channel "chrome" and headed mode (6ms)
  ✓  16 [chrome] › tests/playwright-config.spec.ts:46:7 › chrome project configuration (AC1) › the bundled-Chromium project remains defined for local runs (5ms)
  ✓  18 [chrome] › tests/playwright-config.spec.ts:65:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › slowMo defaults to 0 when RUNNER_SLOW_MO is unset (3ms)
  ✓  17 [chrome] › tests/playwright-config.spec.ts:40:7 › chrome project configuration (AC1) › resolves the ANGLE SwiftShader software-GL launch args (10ms)
  ✓  19 [chrome] › tests/playwright-config.spec.ts:71:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › slowMo resolves to the parsed number when RUNNER_SLOW_MO is numeric (5ms)
  ✓  20 [chrome] › tests/playwright-config.spec.ts:77:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › slowMo falls back to 0 when RUNNER_SLOW_MO is not a number (4ms)
  ✓  22 [chrome] › tests/playwright-config.spec.ts:83:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › video defaults to "off" when RUNNER_VIDEO is unset (3ms)
  ✓  23 [chrome] › tests/playwright-config.spec.ts:89:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › video resolves to "on" when RUNNER_VIDEO is "1" (4ms)
  ✓  24 [chrome] › tests/playwright-config.spec.ts:95:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › video stays "off" for any other RUNNER_VIDEO value (6ms)
  ✓  25 [chrome] › tests/reagent-map-feature-flag.spec.ts:21:7 › isReagentMapEnabled (AC1) › returns true for the literal string "true" (3ms)
  ✓  26 [chrome] › tests/reagent-map-feature-flag.spec.ts:25:7 › isReagentMapEnabled (AC1) › returns false for the literal string "false" (16ms)
  ✓  27 [chrome] › tests/reagent-map-feature-flag.spec.ts:29:7 › isReagentMapEnabled (AC1) › returns false when unset (4ms)
  ✓  28 [chrome] › tests/reagent-map-feature-flag.spec.ts:33:7 › isReagentMapEnabled (AC1) › returns false for the empty string (1ms)
  ✓  29 [chrome] › tests/reagent-map-feature-flag.spec.ts:37:7 › isReagentMapEnabled (AC1) › returns false for near-true case and whitespace variants (AC1 edge case) (10ms)
  ✓  30 [chrome] › tests/reagent-map-feature-flag.spec.ts:47:7 › reagent-map-*.spec.ts scoping (AC6) › the chromium project ignores reagent-map-*.spec.ts (3ms)
  ✓  31 [chrome] › tests/reagent-map-feature-flag.spec.ts:58:7 › reagent-map-*.spec.ts scoping (AC6) › the chrome project defines no such exclusion (4ms)
  ✓  10 [chrome] › tests/map.spec.ts:61:7 › Ye Olde Magic Shop - Map Locations › should verify bounds and visibility (4.1s)
  ✓  13 [chrome] › tests/map.spec.ts:97:7 › Ye Olde Magic Shop - Map Locations › should not have broken tiles (3.5s)
  ✓  12 [chrome] › tests/map.spec.ts:104:7 › Ye Olde Magic Shop - Map Locations › should calculate distance between locations (3.9s)
  ✓  11 [chrome] › tests/map.spec.ts:79:7 › Ye Olde Magic Shop - Map Locations › should use geolocation to find nearest shop (5.3s)
  ✓  14 [chrome] › tests/map.spec.ts:117:7 › Ye Olde Magic Shop - Map Locations › should be able to toggle layer visibility (3.7s)
  ✓  21 [chrome] › tests/reagent-map-a11y-check.spec.ts:3:5 › reagent map a11y check (11.9s)
  ✓  32 [chrome] › tests/reagent-map-globe.spec.ts:53:7 › Reagent map globe (AC2, AC3, AC4) › runs Cesium with no ion token/imagery, pickers off, SCENE3D, and a 6,000km-radius sphere (AC2) (11.4s)
  ✓  33 [chrome] › tests/reagent-map-globe.spec.ts:80:7 › Reagent map globe (AC2, AC3, AC4) › exposes window.__SATURDAY_CESIUM__ and tags every site entity with saturdayMetadata (AC3) (11.4s)
  ✓  34 [chrome] › tests/reagent-map-globe.spec.ts:105:7 › Reagent map globe (AC2, AC3, AC4) › shows all 8 reagent sites exactly once, each within 1km of its published position (AC4) (11.0s)
  ✓  37 [chrome] › tests/reagent-map-globe.spec.ts:157:5 › reagent map view makes no request off the app's own origin or the mock API's origin (NFR) (11.1s)
  ✓  35 [chrome] › tests/reagent-map-globe.spec.ts:129:7 › Reagent map globe (AC2, AC3, AC4) › lists every reagent site by name beside the globe (AC4) (11.8s)
  ✓  39 [chrome] › tests/reagent-map-resilience.spec.ts:46:7 › ReagentMapView onMounted failure paths › rethrows with no entities/sites shown when buildViewer() itself rejects (2.8s)
  ✓  40 [chrome] › tests/webgl.spec.ts:4:7 › WebGL support in the runner browser › provides a webgl2 context on the running app page (3.1s)
  ✓  36 [chrome] › tests/reagent-map-globe.spec.ts:138:7 › Reagent map globe (AC2, AC3, AC4) › Frostcap Hollow (near the pole) and Tidewrack Shoals (on the antimeridian) stay single entities after flying to each (AC4 edge case) (13.8s)
  ✓  38 [chrome] › tests/reagent-map-resilience.spec.ts:31:7 › ReagentMapView onMounted failure paths › destroys the viewer and leaves no canvas when fetchReagentSites() rejects after the viewer resolves, across repeated mounts (5.3s)

  40 passed (31.1s)
```

