---
feature: "leaflet-package-tests"
bounded_context: "Test Automation"
domain_terms: []
files_touched: ["packages/saturday-playwright-leaflet/package.json", "packages/saturday-playwright-leaflet/src/fixture.ts", "packages/saturday-playwright-leaflet/fixtures/index.html", "packages/saturday-playwright-leaflet/fixtures/serve.mjs", "packages/saturday-playwright-leaflet/playwright.config.ts", "packages/saturday-playwright-leaflet/tests/fixture.spec.ts", "packages/saturday-playwright-leaflet/README.md", "pnpm-lock.yaml"]
issue_refs: ["leaflet-package-tests"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: leaflet-package-tests

## Summary

The saturday-playwright-leaflet package ships a LeafletFixture Playwright fixture with 11 documented helpers but has no tests of its own; its test script is a literal exit 1 that breaks the monorepo-wide `pnpm -r test`. This order adds a self-contained fixture page under fixtures/ (local Leaflet assets, a local image tile, known markers/polygon/layer control, exposing window.__SATURDAY_LEAFLET_MAP__), a Playwright suite covering every README helper against known fixture values, and the config to run that suite in real Chrome inside the saturday-runner:debian image with --network none — establishing the pattern the Cesium package will later follow. Doing so surfaces a real defect: src/fixture.ts currently loads Turf.js from a public CDN at test time, which cannot work offline, so that load must move to the package's own local @turf/turf dependency as part of this fix.

### Acceptance Criteria

- Given packages/saturday-playwright-leaflet/fixtures/, a static page MUST load Leaflet from the package's own dependencies (no CDN), a local image, and a known set of markers carrying saturdayMetadata, a polygon, and a layer control, and MUST expose window.__SATURDAY_LEAFLET_MAP__.
  - Example: fixtures/index.html references leaflet.js and leaflet.css from the package's own node_modules/leaflet/dist, never a CDN URL
  - Example: At least 3 markers, each with marker.saturdayMetadata = { id, name }
  - Example: One L.polygon with fixed, known coordinates
  - Example: One L.control.layers overlay with a fixed, known name
  - Example: window.__SATURDAY_LEAFLET_MAP__ is set once the map and all layers finish loading
- Given the package's test script, it MUST serve the fixture page and run a Playwright suite with at least one test per helper documented in README.md, and each test MUST assert a known value taken from the fixture.
  - Example: waitForMapLoad, getCenter, getBounds, getMarkers, drawLine, getMarkersAlongLine, distanceBetween, isPointInPolygon, assertPopupContent, isLayerVisible, checkTileFailures — all 11 README helpers get at least one test
- Given each helper's test, it MUST fail when the fixture's corresponding value is changed; the delivery record MUST document, per helper, the fixture change made and the resulting failing output.
- Given the suite, it MUST run in real Chrome (channel: 'chrome') inside the runner image with no network, and MUST pass with the literal docker command below after dependencies are installed and the package is built; the delivery record MUST capture that command's output.
  - Example: docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-leaflet test'
- Given the change is complete, every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except the run's own delivery record under docs/features/leaflet-package-tests/.

### Non-Functional Requirements

None

## Proposed Fitness Functions

- **Suite passes fully offline, in real Chrome, inside saturday-runner:debian with --network none**
  - **Verification**: the literal docker run command from AC4, exit code checked by devops/CI
  - **Owner**: devops-engineer
- **No CDN/network script loads survive in the fixture page or in src/fixture.ts's Turf.js injection**
  - **Verification**: code review for http(s):// URLs in addScriptTag/<script src>/<link href>; confirmed indirectly by the --network none run passing
  - **Owner**: code-reviewer
- **Cyclomatic complexity < 7 for every new/changed function in the package**
  - **Verification**: analyze-complexity skill / eslint complexity rule in CI
  - **Owner**: developer
- **Every helper's test demonstrably fails when its fixture value changes (test-repair-contract 'prove it can fail')**
  - **Verification**: manual fixture-break-and-revert per helper, output recorded in the delivery record; no automated mutation tool is wired for this package
  - **Owner**: qa-engineer

## Out of Scope

- Adding new helper methods to LeafletFixture
- Changing the drawLine helper's behavior or signature

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: None

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `packages/saturday-playwright-leaflet/package.json` — Replace the `exit 1` test script with the real Playwright invocation; add @turf/turf to dependencies
- `packages/saturday-playwright-leaflet/src/fixture.ts` — Replace the CDN-loaded Turf.js script tag with a load from the package's own local @turf/turf dependency so the suite works under --network none
- `packages/saturday-playwright-leaflet/fixtures/index.html` — New static fixture page: local Leaflet assets, local image tile, known markers/polygon/layer control, exposes __SATURDAY_LEAFLET_MAP__
- `packages/saturday-playwright-leaflet/fixtures/serve.mjs` — New minimal static file server (plain node:http, no new dependency) serving the fixture page and its local assets for the webServer config
- `packages/saturday-playwright-leaflet/playwright.config.ts` — New config: chrome project with channel: 'chrome', webServer pointing at the new static server, following apps/ye-olde-magic-shop/playwright.config.ts's pattern
- `packages/saturday-playwright-leaflet/tests/fixture.spec.ts` — New suite: one test per README-documented helper, each asserting a known fixture value
- `packages/saturday-playwright-leaflet/README.md` — Add a short Testing section documenting how to run the new suite
- `pnpm-lock.yaml` — Updated by `pnpm install` after adding @turf/turf as a package dependency

### Data Model Changes

None

### API Changes

None

### New Dependencies

- @turf/turf — added to packages/saturday-playwright-leaflet's own dependencies, replacing the CDN-loaded turf.min.js in src/fixture.ts; required because AC4 runs the suite with --network none and the current page.addScriptTag CDN URL cannot load offline

## Task List

None

### Developer Tasks

1. Create packages/saturday-playwright-leaflet/fixtures/index.html: load Leaflet from the package's own node_modules/leaflet/dist (no CDN script/link tags); render a map whose tile layer is a local image (not a remote tile server); add at least 3 markers each with a fixed marker.saturdayMetadata = { id, name }; add one L.polygon with known coordinates; add one L.control.layers overlay with a known name; set window.__SATURDAY_LEAFLET_MAP__ once ready.
2. In src/fixture.ts, replace the CDN page.addScriptTag call for Turf with a load from the package's own local @turf/turf dependency (serve its browser bundle from the same local static server, or inject a bundle built from the installed package) so every Turf-dependent helper works under --network none.
3. Add @turf/turf to package.json dependencies (pin the same major version, ^6, previously CDN-pinned); run pnpm install to update pnpm-lock.yaml.
4. Write fixtures/serve.mjs: a minimal static file server using node:http/node:fs (no new dependency) serving fixtures/index.html plus the local Leaflet/Turf/image assets it references.
5. Create packages/saturday-playwright-leaflet/playwright.config.ts following apps/ye-olde-magic-shop/playwright.config.ts's pattern: a chrome project using channel: 'chrome' and a webServer block pointing at fixtures/serve.mjs; omit the Cesium-specific RUNNER_IMAGE_SOFTWARE_GL_ARGS launch args since Leaflet needs no WebGL context.
6. Create packages/saturday-playwright-leaflet/tests/fixture.spec.ts with at least one test per each of the 11 README-documented helpers, each test asserting a literal known value from the fixture (marker names/coordinates, polygon bounds, layer name, center, computed distance).
7. Update package.json's test script to invoke the new Playwright suite (e.g. playwright test --project=chrome), removing the exit 1 placeholder.
8. Verify (no change expected) that tsconfig.json's rootDir/include scoping still excludes fixtures/ and tests/ from the package's own tsc build.

### QA Tasks

1. Run the new suite locally and via the exact AC4 docker command inside saturday-runner:debian; record the full output in the delivery record.
2. For each of the 11 README helpers, change the fixture's corresponding known value, re-run that helper's test, confirm it fails, then revert the change — record the fixture edit and the failing output per helper in the delivery record, per shared/rules/test-repair-contract.md's 'prove it can fail' requirement.
3. Verify cyclomatic complexity below 7 for every new/changed function and that changed lines reach the required coverage via tests that can fail.
4. Confirm no file outside packages/saturday-playwright-leaflet/**, pnpm-lock.yaml, and docs/features/leaflet-package-tests/** changed.

### Tech Writer Tasks

1. Add a short 'Testing' section to packages/saturday-playwright-leaflet/README.md describing the fixture page, how to run `pnpm test`, and how it runs inside the runner image.

### DevOps Tasks

None

## Edge Cases and Risks

- Turf.js script fails to load under --network none because fixture.ts still CDN-loads it: Load Turf from the package's own local @turf/turf dependency instead; verified by the AC4 docker run passing
- checkTileFailures()'s happy-path test needs zero tile errors with no network tile server available: Fixture's tile layer uses a local image via the fixture's own static server, not a remote tile host
- checkTileFailures()'s failing-path test (AC3) needs a tile load to fail without any external network access: Point the broken-fixture variant at a URL that 404s against the fixture's own local static server
- Package is tested before `pnpm build` has run, so dist/ is stale or missing: AC4's docker command explicitly runs after install and build; developer confirms the suite exercises the package consistently with how apps/ye-olde-magic-shop consumes it

## Definition of Done

- All 5 acceptance criteria met and verified
- pnpm --filter @orieken/saturday-playwright-leaflet test passes locally and via the exact AC4 docker command inside saturday-runner:debian; output recorded in the delivery record
- Per-helper fixture-break-and-revert evidence recorded per AC3
- Cyclomatic complexity below 7 for all new/changed functions; changed lines covered by tests that can fail
- No files outside scope.paths changed, except the run's own delivery record under docs/features/leaflet-package-tests/
- code-reviewer has approved the implementation
- README.md documents how to run the new suite
