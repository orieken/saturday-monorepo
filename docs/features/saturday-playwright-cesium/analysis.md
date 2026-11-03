---
feature: "saturday-playwright-cesium"
bounded_context: "Test Automation"
domain_terms: ["Aurelis", "Reagent Site", "Seam contract", "Settled"]
files_touched: ["packages/saturday-playwright-cesium/README.md", "packages/saturday-playwright-cesium/package.json", "packages/saturday-playwright-cesium/tsconfig.json", "packages/saturday-playwright-cesium/vitest.config.ts", "packages/saturday-playwright-cesium/playwright.config.ts", "packages/saturday-playwright-cesium/src/index.ts", "packages/saturday-playwright-cesium/src/fixture.ts", "packages/saturday-playwright-cesium/src/geometry.ts", "packages/saturday-playwright-cesium/fixtures/index.html", "packages/saturday-playwright-cesium/fixtures/serve.mjs", "packages/saturday-playwright-cesium/fixtures/assets/fixture-globe.js", "packages/saturday-playwright-cesium/fixtures/assets/aurelis.jpg", "packages/saturday-playwright-cesium/tests/fixture.spec.ts", "packages/saturday-playwright-cesium/tests/geometry.spec.ts", "pnpm-lock.yaml", "docs/features/saturday-playwright-cesium/"]
issue_refs: ["saturday-playwright-cesium"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: saturday-playwright-cesium

## Summary

@orieken/saturday-playwright-cesium gives Playwright tests a seam (window.__SATURDAY_CESIUM__ = { viewer, Cesium }) to act on a Cesium globe like a user and read back entities, screen projections and camera state without depending on pixels, following the same seam-contract pattern as @orieken/saturday-playwright-leaflet. The first developer-stage run was killed by the 30-minute stage timeout and finished by hand on branch salvage/saturday-playwright-cesium (commit 926ad56, not yet applied to this worktree); its own fixture suite (7 Vitest + 11 Playwright tests) covers all seven helpers (waitForSettled, getEntities, toScreen, clickAt, flyTo, metersPerPixelAt, globeRadiusMeters) and three radius-parameterized geometry functions (distance, lineLength, pointInCircle) against a static fixture globe named Aurelis (radius 6,000 km, baseLayerPicker/geocoder/ion off, local equirectangular image, no CDN). This order's job is to apply that salvage diff verbatim, verify every acceptance criterion against the real code rather than rewriting it, and resolve the one undecided structural question the salvage code already assumed an answer to: whether the geometry helpers stay package-local or move to a new packages/saturday-geo.

### Acceptance Criteria

- An app under test exposes window.__SATURDAY_CESIUM__ = { viewer, Cesium } and attaches saturdayMetadata (id, name, kind) to each entity it wants tests to find, as documented in the package's own README seam contract.
  - Example: README.md section 1 shows (window as any).__SATURDAY_CESIUM__ = { viewer, Cesium } and entity.saturdayMetadata = { id, name, kind }
  - Example: fixtures/assets/fixture-globe.js builds its own viewer and sets window.__SATURDAY_CESIUM__ after adding entities, proving the seam can attach asynchronously
- packages/saturday-playwright-cesium/fixtures/ serves a static page with a Cesium Viewer on a spherical ellipsoid of a declared radius, one equirectangular image covering the whole globe from the package's own files, baseLayerPicker/geocoder/ion imagery off, known point/line/circle entities carrying metadata, and makes no request off its own origin.
  - Example: fixtures/assets/fixture-globe.js: new Cesium.Ellipsoid(6_000_000, 6_000_000, 6_000_000); baseLayerPicker: false, geocoder: false; SingleTileImageryProvider.fromUrl('/assets/aurelis.jpg', ...)
  - Example: Four saturdayMetadata-tagged entities (reagent-site-1 point, supply-route-1 line, exclusion-zone-1 circle, mystery-site-1 point with a deliberately unrecognized kind) plus one untagged point to prove exclusion
  - Example: tests/fixture.spec.ts's 'fixture page makes no request off its own origin' test already asserts every http(s) request's hostname is 127.0.0.1
- cesium.waitForSettled() resolves once the globe's tiles are loaded, the camera has stopped, and no render is pending, and rejects with a message naming the stuck condition (tiles/camera/render/seam) after a timeout.
  - Example: fixture.ts's isSettled()/describeUnsettledConditions() check globe.tilesLoaded, a cameraMoving flag and a renderPending flag tracked via moveStart/moveEnd/postRender listeners
  - Example: fixture.spec.ts already has one test asserting resolution and one forcing cameraMoving=true to assert the rejection message contains 'camera'
- cesium.getEntities() returns each saturdayMetadata-tagged entity's metadata and geometry in degrees and metres: point position, line endpoints, circle centre and radius.
  - Example: fixture.ts's extractEntities() dispatches on entity.ellipse/polyline/position, never on saturdayMetadata.kind, and excludes untagged entities
  - Example: fixture.spec.ts asserts all four tagged entities' geometry, including the mystery-site-1 entity whose kind is 'unexpected-kind-value'
- cesium.toScreen(lat, lng) returns page coordinates, or null when the point is on the far side of the globe or off-screen.
  - Example: fixture.ts's projectToScreen() uses an EllipsoidalOccluder for far-side occlusion and a canvas-bounds check for off-screen, returning page (not canvas) coordinates via getBoundingClientRect()
  - Example: fixture.spec.ts has one passing case, one far-side-null case, and one off-canvas-but-visible-null case
- cesium.clickAt(lat, lng) clicks through page.mouse, and a fixture test shows the globe's pick returned that location within the tolerance metersPerPixelAt gives.
  - Example: fixture.ts's clickAt() calls toScreen() then page.mouse.click(); throws if the point isn't visible
  - Example: fixture.spec.ts installs a LEFT_CLICK pick handler and asserts the picked lat/lng is within metersPerPixelAt(lat,lng) * 2 metres of the requested point, for both a near-centre and a near-horizon case
- cesium.flyTo(lat, lng, heightMeters) moves the camera and waits until settled.
  - Example: fixture.ts's flyTo() calls camera.flyTo({ destination, duration: 0 }) then waitForSettled()
  - Example: fixture.spec.ts asserts the camera's resulting lat/lng/height match the requested destination within tolerance
- distance, lineLength and point-in-circle checks take the globe's radius from the seam rather than assuming Earth's, and are unit-tested in Node without a browser; whether they live in this package or a shared saturday-geo is decided at the design gate.
  - Example: src/geometry.ts's distance/lineLength/pointInCircle all take an explicit radiusMeters/globeRadiusMeters parameter, zero Cesium or Playwright imports
  - Example: tests/geometry.spec.ts runs under vitest (environment: node) and explicitly proves Earth's radius (6,371,000m) gives a different, wrong result than Aurelis's (6,000,000m)
  - Example: geometry.ts's own header comment already asserts 'see architecture-notes.md's structural decision to keep this package-local' — but no architecture-notes.md exists for this feature yet, so the architect must ratify or override this, not treat it as already decided
- The package's test script runs the fixture suite in real Chrome (channel: 'chrome', headed, --use-gl=angle --use-angle=swiftshader) and passes with the literal AC9 docker command after dependencies are installed and the package is built; record its output in the delivery record.
  - Example: package.json's test script: 'vitest run && playwright test --project=chrome'
  - Example: playwright.config.ts's 'chrome' project matches RUNNER_IMAGE_SOFTWARE_GL_ARGS = ['--use-gl=angle', '--use-angle=swiftshader'] literally, per docs/runner-image.md and the ui-suite-real-chrome precedent
  - Example: docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-cesium test'
- Each helper's test fails when the fixture's corresponding value changes; for a helper reading nothing from the fixture, a change to the helper's own code that the test catches is the proof instead. Record, per helper, the change made and the failing output in the delivery record.
  - Example: metersPerPixelAt and globeRadiusMeters read nothing fixed from fixture-globe.js's entity list — their fail-proof must come from a deliberate code change (e.g. flipping camera.getPixelSize's args) rather than a fixture edit
  - Example: waitForSettled/getEntities/toScreen/clickAt/flyTo each have a known fixture value (entity coordinates, circle radius, camera destination) that QA can perturb and revert
- The fixture's static server listens on 127.0.0.1 only.
  - Example: fixtures/serve.mjs: createServer(...).listen(port, '127.0.0.1', ...) — its own comment notes saturday-playwright-leaflet's serve.mjs omits the host argument (binding all interfaces) and that call must not be copied here
- Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail, and no file outside scope.paths changes except the run's own delivery record under docs/features/saturday-playwright-cesium/.
  - Example: scope.paths: packages/saturday-playwright-cesium/**, packages/saturday-geo/**, pnpm-lock.yaml, docs/features/saturday-playwright-cesium/**

### Non-Functional Requirements

- **other**: The fixture page and suite MUST make zero requests off their own origin (127.0.0.1).
- **performance**: cesium.waitForSettled() MUST reject with a clear message naming the stuck condition if the globe never settles within its timeout. (waitForSettled default timeout 15s)
- **other**: The package's test script MUST pass fully offline inside the saturday-runner:debian image with --network none, after dependencies are installed and the package is built.
- **other**: Every new or changed function MUST have cyclomatic complexity below the framework-wide budget. (cyclomatic complexity per function 7count)
- **other**: Every changed line MUST be covered by a test that can fail (ADR-008 coverage gate). (changed-line test coverage 85%)

## Proposed Fitness Functions

- **The package's test command passes fully offline, in real Chrome, inside saturday-runner:debian with --network none**
  - **Verification**: the literal AC9 docker run command, exit code checked by devops/CI
  - **Owner**: devops-engineer
- **Zero off-origin network requests during fixture page load and tile load**
  - **Verification**: tests/fixture.spec.ts's existing 'fixture page makes no request off its own origin' test, which filters to http(s) requests and asserts hostname === '127.0.0.1'
  - **Owner**: qa-engineer
- **Geometry helpers never fall back to Earth's radius**
  - **Verification**: tests/geometry.spec.ts's existing tests that assert Aurelis's radius and Earth's radius produce different, independently-verified results for the same inputs
  - **Owner**: qa-engineer
- **Cyclomatic complexity below 7 for every new/changed function in the package**
  - **Verification**: analyze-complexity skill / eslint complexity rule in CI
  - **Owner**: developer
- **Every helper's and geometry function's test demonstrably fails when its fixture value (or, for value-free helpers, the helper's own code) changes**
  - **Verification**: manual fixture/code-break-and-revert per helper, output recorded in the delivery record; no automated mutation tool is wired for this package
  - **Owner**: qa-engineer
- **Fixture server binds 127.0.0.1 only, never all interfaces**
  - **Verification**: code review of fixtures/serve.mjs's .listen(port, '127.0.0.1', ...) call, confirmed indirectly by the --network none run passing
  - **Owner**: code-reviewer

## Out of Scope

- Any change to the demo app — the reagent map view is a separate order, written from the reagent-map spec
- Drawing-tool UX — that is a Flow in each app's tests
- Terrain, 3D tiles and time-dynamic entities
- Adding new helper methods to CesiumFixture beyond the seven already implemented
- Tightening saturday-playwright-leaflet's own serve.mjs to bind 127.0.0.1 — noted as a divergence this package deliberately avoids repeating, but fixing the sibling package is out of this scope

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: containers/saturday-runner (the runner image) — read-only: this feature consumes runner-check's verified ANGLE/SwiftShader launch-arg contract and the debian build target without modifying the image itself, docs/specs/reagent-map-spec — read-only: the fixture planet's name (Aurelis) and radius (6,000 km) are taken from that spec's resolved Q8, not redecided here

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `packages/saturday-playwright-cesium/README.md` — Documents the seam contract, all 7 helpers, the 3 geometry functions, and both test-run commands; apply via the salvage diff, verify against the actual implementation
- `packages/saturday-playwright-cesium/package.json` — Declares cesium 1.146.0 as a dependency and the test script ('vitest run && playwright test --project=chrome'); apply via the salvage diff
- `packages/saturday-playwright-cesium/tsconfig.json` — Scopes the package's own tsc build to src/**/* only, excluding fixtures/ and tests/; apply via the salvage diff
- `packages/saturday-playwright-cesium/vitest.config.ts` — Runs tests/geometry.spec.ts under a Node environment with no browser; its comment cites the not-yet-ratified saturday-geo decision
- `packages/saturday-playwright-cesium/playwright.config.ts` — Defines the 'chrome' project with the literal RUNNER_IMAGE_SOFTWARE_GL_ARGS from docs/runner-image.md, and the fixture webServer on port 4174
- `packages/saturday-playwright-cesium/src/index.ts` — Public entry point re-exporting test, CesiumFixture types, and the three geometry functions
- `packages/saturday-playwright-cesium/src/fixture.ts` — Implements all 7 helpers (waitForSettled, getEntities, toScreen, clickAt, flyTo, metersPerPixelAt, globeRadiusMeters) backing every helper-level acceptance criterion
- `packages/saturday-playwright-cesium/src/geometry.ts` — Implements distance/lineLength/pointInCircle, radius-parameterized per AC8; header comment pre-empts an architect decision that must be ratified or overridden
- `packages/saturday-playwright-cesium/fixtures/index.html` — Static fixture page: fixed 1024x512 canvas offset from the page origin, loads Cesium from CESIUM_BASE_URL='/cesium/' (never a CDN)
- `packages/saturday-playwright-cesium/fixtures/serve.mjs` — Minimal node:http static server binding 127.0.0.1 only, serving Cesium's Build/Cesium assets and the local planet texture
- `packages/saturday-playwright-cesium/fixtures/assets/fixture-globe.js` — Builds the known Aurelis globe (radius 6,000,000m) with baseLayerPicker/geocoder/ion off and the 4 tagged + 1 untagged entities every test asserts against
- `packages/saturday-playwright-cesium/fixtures/assets/aurelis.jpg` — The package's own equirectangular planet texture; must be verified at most 2048x1024px per the spec's resolved Open Question
- `packages/saturday-playwright-cesium/tests/fixture.spec.ts` — 11 Playwright tests covering every helper plus the off-origin-request check; verify, do not rewrite
- `packages/saturday-playwright-cesium/tests/geometry.spec.ts` — 7 Vitest tests covering distance/lineLength/pointInCircle including explicit Earth-vs-Aurelis radius divergence; verify, do not rewrite
- `pnpm-lock.yaml` — Already updated by the salvage commit for cesium and its devDependencies; confirm no drift after git apply --index against this worktree's current lockfile
- `docs/features/saturday-playwright-cesium/` — This order's own delivery record: docker run output, per-helper fail-proof evidence, and the architect's saturday-geo ratification

### Data Model Changes

None

### API Changes

None

### New Dependencies

- cesium@1.146.0 — new to the monorepo; pinned exact version per the salvage package.json, no caret range

## Task List

None

### Developer Tasks

1. Apply the salvage diff verbatim: git diff --binary 68a5901 926ad56 | git apply --index, bringing packages/saturday-playwright-cesium onto disk at the verified commit state; do not rewrite the implementation.
2. Run pnpm install to materialize node_modules/cesium and confirm pnpm-lock.yaml matches the salvage commit's already-updated lockfile with no drift.
3. Apply the architect's ratified decision on geometry.ts's placement: if a new packages/saturday-geo is created, move distance/lineLength/pointInCircle and tests/geometry.spec.ts there and update src/index.ts's re-export and every file citing the old placement; if kept package-local, replace the four 'see architecture-notes.md's structural decision' comments (geometry.ts, fixture-globe.js, serve.mjs, vitest.config.ts) with a concrete reference now that architecture-notes.md exists.
4. Verify fixtures/assets/aurelis.jpg is at most 2048x1024px per the spec's resolved Open Question; replace it if it exceeds that bound.
5. Build the package (pnpm --filter @orieken/saturday-playwright-cesium build) and run pnpm --filter @orieken/saturday-playwright-cesium test locally to confirm the salvage suite (7 Vitest + 11 Playwright tests) passes unmodified.

### QA Tasks

1. Run the AC9 docker command verbatim inside saturday-runner:debian with --network none, after install and build; record the full output in the delivery record.
2. For each of the 7 fixture helpers and the 3 geometry functions, change the fixture's corresponding known value (fixture-globe.js's entity/camera values, or for metersPerPixelAt/globeRadiusMeters, which read no fixed fixture value, the helper's own code), re-run that test, confirm it fails, then revert — record the exact edit and the failing output per helper/function in the delivery record per AC10 and shared/rules/test-repair-contract.md.
3. Verify cyclomatic complexity below 7 for every new/changed function and that changed lines reach the required coverage via tests that can fail.
4. Confirm fixtures/serve.mjs binds 127.0.0.1 only and that the existing 'fixture page makes no request off its own origin' test passes.
5. Confirm no file outside packages/saturday-playwright-cesium/**, packages/saturday-geo/** (if the architect creates it), pnpm-lock.yaml, and docs/features/saturday-playwright-cesium/** changed.

### Tech Writer Tasks

1. Confirm README.md's seam contract, all 7 helper descriptions, the geometry helpers section, and both run commands (local and --network none docker) still match the applied implementation; update any import path the architect's saturday-geo decision changes.
2. If geometry.ts moves to a new packages/saturday-geo, add that package's own README documenting distance/lineLength/pointInCircle independent of saturday-playwright-cesium.

### DevOps Tasks

None

## Edge Cases and Risks

- metersPerPixelAt and globeRadiusMeters read no fixed value from fixture-globe.js's entity list: Their fail-proof evidence per AC10 comes from a deliberate change to the helper's own code (e.g. swapping camera.getPixelSize's width/height arguments), not a fixture edit — document this distinction in the delivery record.
- clickAt's pick tolerance must scale with ground resolution, not a fixed pixel/metre value: Already handled: fixture.spec.ts computes tolerance as metersPerPixelAt(lat,lng) * 2 at both a near-centre and a near-horizon point.
- toScreen must distinguish far-side occlusion from merely off-canvas: Already handled: projectToScreen() checks EllipsoidalOccluder.isPointVisible() before the canvas-bounds check, with a dedicated test for each case.
- waitForSettled must not resolve while the camera is still mid-move or a render is pending after it stops: Already handled: moveStart/moveEnd/postRender listeners drive cameraMoving/renderPending flags; a dedicated test forces cameraMoving=true and asserts the rejection names 'camera'.
- The seam may attach asynchronously (after an imagery provider resolves), not synchronously on page.goto(): Already handled: pollUntilSettled() waits for window.__SATURDAY_CESIUM__ to exist before installing settle tracking, with its own named timeout message.
- Applying the salvage diff via git apply --index may conflict if pnpm-lock.yaml has moved since 68a5901: Developer verifies a clean apply and, on conflict, resolves by re-running pnpm install rather than hand-editing the lockfile.
- The geometry.ts architecture decision is baked into four files before an architect has ratified it: Architect stage must explicitly ratify or override before developer work proceeds; developer then updates all four citing comments to match the real decision.

## Definition of Done

- All 12 acceptance criteria met and verified against the applied salvage implementation
- Architect has explicitly ratified or overridden the geometry.ts package-local-vs-saturday-geo decision, with a fitness function recorded
- pnpm --filter @orieken/saturday-playwright-cesium test passes locally and via the exact AC9 docker command inside saturday-runner:debian; output recorded in the delivery record
- Per-helper and per-geometry-function fail-proof evidence recorded per AC10
- Cyclomatic complexity below 7 for all new/changed functions; changed lines covered by tests that can fail
- fixtures/assets/aurelis.jpg confirmed at most 2048x1024px
- No files outside scope.paths changed, except the run's own delivery record under docs/features/saturday-playwright-cesium/
- code-reviewer has approved the implementation
- README.md accurately documents the applied implementation, including any saturday-geo extraction
