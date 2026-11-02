# Pre-Ship Checks

## ui-suite

- Command: `scripts/ui-suite-in-chrome.sh`
- Duration: 1m32.513275375s
- Status: exit 0

```
packages/saturday-cucumber build: [32mDTS[39m [1mdist/index.d.cts [22m[32m2.42 KB[39m
packages/saturday-cucumber build: [32mDTS[39m [1mdist/index.d.ts  [22m[32m2.42 KB[39m
packages/saturday-cucumber build: Done
packages/saturday-playwright-k6-exporter build: [32mDTS[39m ⚡️ Build success in 2961ms
packages/saturday-playwright-k6-exporter build: [32mDTS[39m [1mdist/cli.d.cts     [22m[32m20.00 B[39m
packages/saturday-playwright-k6-exporter build: [32mDTS[39m [1mdist/fixture.d.cts [22m[32m550.00 B[39m
packages/saturday-playwright-k6-exporter build: [32mDTS[39m [1mdist/index.d.cts   [22m[32m1.64 KB[39m
packages/saturday-playwright-k6-exporter build: [32mDTS[39m [1mdist/cli.d.ts      [22m[32m20.00 B[39m
packages/saturday-playwright-k6-exporter build: [32mDTS[39m [1mdist/fixture.d.ts  [22m[32m549.00 B[39m
packages/saturday-playwright-k6-exporter build: [32mDTS[39m [1mdist/index.d.ts    [22m[32m1.64 KB[39m
packages/saturday-playwright-k6-exporter build: Done
runner-check: Chrome 156.0.8078.79 — ANGLE (Google, Vulkan 1.3.0 (SwiftShader Device (Subzero) (0x0000C0DE)), SwiftShader driver)
runner-check: Edge 156.0.4314.56 — ANGLE (Google, Vulkan 1.3.0 (SwiftShader Device (Subzero) (0x0000C0DE)), SwiftShader driver)
[dotenv@17.2.3] injecting env (5) from .env -- tip: 🔑 add access controls to secrets: https://dotenvx.com/ops

Running 24 tests using 7 workers

[dotenv@17.2.3] injecting env (0) from .env -- tip: 🗂️ backup and recover secrets: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: ⚙️  enable debug logging with { debug: true }
[dotenv@17.2.3] injecting env (0) from .env -- tip: ⚙️  load multiple .env files with { path: ['.env.local', '.env'] }
[dotenv@17.2.3] injecting env (0) from .env -- tip: 🗂️ backup and recover secrets: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: 🔄 add secrets lifecycle management: https://dotenvx.com/ops
[dotenv@17.2.3] injecting env (0) from .env -- tip: 🔐 encrypt with Dotenvx: https://dotenvx.com
[dotenv@17.2.3] injecting env (0) from .env -- tip: 🔐 encrypt with Dotenvx: https://dotenvx.com
  ✓   7 [chrome] › tests/map.spec.ts:9:7 › Ye Olde Magic Shop - Map Locations › map should load and center on the correct region (4.9s)
  ✓   1 [chrome] › tests/example.spec.ts:5:5 › homepage smoke test @k6 (5.2s)
  ✓   3 [chrome] › tests/map.spec.ts:17:7 › Ye Olde Magic Shop - Map Locations › markers should exist and popups should display correct information (5.2s)
  ✓   4 [chrome] › tests/magic-shop.spec.ts:3:5 › magic shop loads and displays wares (5.4s)
  ✓   2 [chrome] › tests/cart.spec.ts:29:7 › Shopping Cart › Increase quantity from the cart page (6.4s)
  ✓   5 [chrome] › tests/cart.spec.ts:12:7 › Shopping Cart › Add the same item twice increases quantity to 2 (6.9s)
  ✓   6 [chrome] › tests/heatmap.spec.ts:4:5 › heatmap verification (7.1s)
  ✓   9 [chrome] › tests/map.spec.ts:61:7 › Ye Olde Magic Shop - Map Locations › should verify bounds and visibility (3.9s)
  ✓  15 [chrome] › tests/playwright-config.spec.ts:33:7 › chrome project configuration (AC1) › resolves channel "chrome" and headed mode (3ms)
  ✓  16 [chrome] › tests/playwright-config.spec.ts:40:7 › chrome project configuration (AC1) › resolves the ANGLE SwiftShader software-GL launch args (4ms)
  ✓  17 [chrome] › tests/playwright-config.spec.ts:46:7 › chrome project configuration (AC1) › the bundled-Chromium project remains defined for local runs (5ms)
  ✓  18 [chrome] › tests/playwright-config.spec.ts:65:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › slowMo defaults to 0 when RUNNER_SLOW_MO is unset (1ms)
  ✓  10 [chrome] › tests/map.spec.ts:48:7 › Ye Olde Magic Shop - Map Locations › should check point in polygon (4.0s)
  ✓  19 [chrome] › tests/playwright-config.spec.ts:71:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › slowMo resolves to the parsed number when RUNNER_SLOW_MO is numeric (1ms)
  ✓  20 [chrome] › tests/playwright-config.spec.ts:77:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › slowMo falls back to 0 when RUNNER_SLOW_MO is not a number (7ms)
  ✓  22 [chrome] › tests/playwright-config.spec.ts:83:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › video defaults to "off" when RUNNER_VIDEO is unset (1ms)
  ✓  23 [chrome] › tests/playwright-config.spec.ts:89:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › video resolves to "on" when RUNNER_VIDEO is "1" (6ms)
  ✓   8 [chrome] › tests/map.spec.ts:27:7 › Ye Olde Magic Shop - Map Locations › should be able to draw a line and find markers along it (4.3s)
  ✓  24 [chrome] › tests/playwright-config.spec.ts:95:7 › RUNNER_SLOW_MO / RUNNER_VIDEO resolution (AC1 edge cases) › video stays "off" for any other RUNNER_VIDEO value (3ms)
  ✓  12 [chrome] › tests/map.spec.ts:97:7 › Ye Olde Magic Shop - Map Locations › should not have broken tiles (3.4s)
  ✓  11 [chrome] › tests/map.spec.ts:79:7 › Ye Olde Magic Shop - Map Locations › should use geolocation to find nearest shop (4.8s)
  ✓  14 [chrome] › tests/map.spec.ts:117:7 › Ye Olde Magic Shop - Map Locations › should be able to toggle layer visibility (3.4s)
  ✓  13 [chrome] › tests/map.spec.ts:104:7 › Ye Olde Magic Shop - Map Locations › should calculate distance between locations (3.5s)
  ✓  21 [chrome] › tests/webgl.spec.ts:4:7 › WebGL support in the runner browser › provides a webgl2 context on the running app page (2.3s)

  24 passed (16.0s)
```

