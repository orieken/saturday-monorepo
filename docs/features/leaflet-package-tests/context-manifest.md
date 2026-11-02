# Context Manifest: leaflet-package-tests

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `packages/saturday-playwright-leaflet/README.md` — Lists the 11 documented helpers (waitForMapLoad, getCenter, getBounds, getMarkers, drawLine, getMarkersAlongLine, distanceBetween, isPointInPolygon, assertPopupContent, isLayerVisible, checkTileFailures) — AC2 requires at least one test per README helper, so this is the literal checklist for the suite.
- `packages/saturday-playwright-leaflet/src/fixture.ts` — The LeafletFixture implementation under test. Note it actually exposes 13 methods — dragMap and zoomMap exist here but are undocumented in the README; checkTileFailures reads a DOM class (.leaflet-tile-error) rather than intercepting network, which matters for a --network none run.
- `packages/saturday-playwright-leaflet/src/index.ts` — Package's public entry point (re-exports `test` and `LeafletFixture`) — confirms what the fixture page's test files will import.
- `packages/saturday-playwright-leaflet/package.json` — Current `test` script is the literal `exit 1` this order exists to fix (AC in the spec); also the dependency list the fixture page/test runner additions must fit into (currently only leaflet, @types/leaflet, @playwright/test, typescript — no static-file server dependency exists yet).
- `packages/saturday-playwright-leaflet/tsconfig.json` — rootDir/include are scoped to src/**/* only — confirms new fixtures/ and test files won't be pulled into the package's own `tsc` build and need their own compilation path (Playwright's test runner, not this tsconfig).
- `apps/ye-olde-magic-shop/tests/map.spec.ts` — Existing consumer-app test that already exercises every one of the 11 README helpers (plus dragMap via zoomMap) against a live Leaflet map — the closest thing to a worked answer for what each new fixture-page test should assert (marker names, polygon bounds, layer-control label, distance tolerances).
- `apps/ye-olde-magic-shop/playwright.config.ts` — Only existing example in the repo of a `channel: 'chrome'` project plus the runner-image launch-arg pattern (resolveRunnerSlowMo/resolveRunnerVideo) and a webServer block — AC4 requires the new package's own suite to run the same way, headed, in real Chrome, inside the runner image.
- `docs/runner-image.md` — Read lines 1-45 and 160-180 only (file is 394 lines; the WebGL/ANGLE build-record and UBI9 sections past that are not relevant — Leaflet needs no GL context). Establishes the docker run contract (repo mounted at /work, non-root, --network none expectation, exec-based exit code) that AC4's literal docker command depends on.

## Token Budget

- **Status**: OK — ~11474 tokens, 9.6% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `packages/saturday-playwright-leaflet/README.md` | 3202 | 800 |
| `packages/saturday-playwright-leaflet/src/fixture.ts` | 10555 | 2638 |
| `packages/saturday-playwright-leaflet/src/index.ts` | 83 | 20 |
| `packages/saturday-playwright-leaflet/package.json` | 555 | 138 |
| `packages/saturday-playwright-leaflet/tsconfig.json` | 331 | 82 |
| `apps/ye-olde-magic-shop/tests/map.spec.ts` | 5064 | 1266 |
| `apps/ye-olde-magic-shop/playwright.config.ts` | 2378 | 594 |
| `docs/runner-image.md` | 23746 | 5936 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `No KI/ADR corpus exists in this repo (no shared/knowledge/, .claude/knowledge/, or docs/adrs(/adr)/ directory found) — search-ki lookup was skipped for lack of a corpus, not skipped by choice.`
- `Blast radius was not graph-derived: no mcp__loom__* tools are registered in this session. Dependents were instead found by direct Grep: apps/ye-olde-magic-shop/tests/map.spec.ts and src/views/LocationsView.vue are the only consumers of @orieken/saturday-playwright-leaflet in the repo.`
- `Prior deliveries in this bounded context (Test Automation, per 'Owning Context' in docs/features/ui-suite-real-chrome/analysis.md and docs/features/golden-workspace-packages/analysis.md): both found, neither has a retrospective.md, so no retrospective lessons to surface. ui-suite-real-chrome is also this order's declared depends-on and is the source of the channel:'chrome'/runner-image pattern pinned above.`
- `apps/ye-olde-magic-shop/src/views/LocationsView.vue — excluded as a different bounded context (the demo app, not Test Automation); the spec documents no Context Crossing for this order. map.spec.ts's inline comments already carry the coordinate/metadata values needed without reading the app's own source.`
- `scripts/ui-suite-in-chrome.sh — excluded: the spec gives the exact docker run command AC4 must pass verbatim, so the wrapper script's own argument parsing isn't load-bearing for this order.`
- `packages/saturday-playwright-k6-exporter/tests/fixture.spec.ts — excluded: that package's test pattern mocks `page` and unit-tests fixture logic in Jest, the opposite of AC4's requirement for a real end-to-end Playwright/Chrome run against a served page, so it would mislead rather than help as an exemplar.`
- `containers/saturday-runner/** (Containerfile, entrypoint.sh, apt.lock, etc.) — excluded: this order only consumes the already-built saturday-runner:debian image per AC4's docker command; docs/runner-image.md stands in as the reference surface, same exclusion ui-suite-real-chrome's own manifest made for the same reason.`
