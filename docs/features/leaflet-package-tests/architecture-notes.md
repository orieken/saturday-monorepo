---
feature: "leaflet-package-tests"
bounded_context: "Test Automation"
domain_terms: ["LeafletFixture", "saturdayMetadata", "Site-Centric pattern"]
files_touched: ["packages/saturday-playwright-leaflet/fixtures/", "packages/saturday-playwright-leaflet/fixtures/", "packages/saturday-playwright-leaflet/ (root)", "packages/saturday-playwright-leaflet/tests/", "packages/saturday-playwright-leaflet/src/"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: leaflet-package-tests

## Structural Decisions

### Replace the CDN-loaded Turf.js script tag in src/fixture.ts with a load from the package's own local @turf/turf dependency, served statically alongside the fixture page.
- **Rationale**: AC4 requires the suite to pass under `--network none`. A page.addScriptTag pointed at https://cdn.jsdelivr.net cannot resolve offline, which would fail every helper that depends on window.turf (getMarkersAlongLine, isPointInPolygon, distanceBetween). fixture.ts is this package's adapter boundary (a Playwright fixture wrapping Leaflet's window global); sourcing its runtime library locally instead of from a CDN is a dependency-sourcing fix, not a layer-boundary change.
- **Trade-offs**:
  - Requires shipping a browser-consumable Turf bundle (UMD/IIFE) from node_modules or building one at pretest time, instead of one CDN URL.
  - @turf/turf becomes a real runtime dependency of the published package, increasing its install footprint for every consumer even if they never call the Turf-backed helpers.
- **Fitness Function**: No CDN/network script or link tag reaches the browser from the fixture page or fixture.ts's library loading — verified by code-reviewer greps for http(s):// inside addScriptTag/<script src>/<link href>; confirmed end-to-end by the AC4 docker run with --network none exiting 0
### Move @turf/turf into package.json's `dependencies` block (not devDependencies), pinned to ^6 to match the previously CDN-pinned major version.
- **Rationale**: The code that loads Turf executes inside the browser context under test at runtime, not just at build/test-authoring time on the Node side — that makes it a runtime dependency of the published package, the same category leaflet is already in.
- **Trade-offs**:
  - Every consumer of @orieken/saturday-playwright-leaflet now pulls @turf/turf transitively, whether or not they exercise the three Turf-backed helpers.
- **Fitness Function**: package.json lists @turf/turf under dependencies, pinned ^6, and pnpm-lock.yaml reflects it under this package's own dependency block — verified by code-reviewer greps package.json; pnpm install's lockfile diff is reviewed in the same PR
### Keep all new test infrastructure (fixtures/index.html, fixtures/serve.mjs, playwright.config.ts, tests/fixture.spec.ts) inside packages/saturday-playwright-leaflet/ itself, following apps/ye-olde-magic-shop/playwright.config.ts's existing chrome/channel/webServer pattern verbatim rather than extracting a shared test-harness package.
- **Rationale**: This package is the sole owner and sole consumer of the LeafletFixture it's testing; the fixture page has no reuse value outside verifying this package's own public API. CLAUDE.md's architect-skip guidance applies directly: purely additive test infrastructure within an established pattern, no new package or base class.
- **Fitness Function**: judgment-only — This package is the sole owner and sole consumer of the LeafletFixture it's testing; the fixture page has no reuse value outside verifying this package's own public API. CLAUDE.md's architect-skip guidance applies directly: purely additive test infrastructure within an established pattern, no new package or base class.
### No change to tsconfig.json's existing rootDir ("./src") / include (["src/**/*"]) scoping — fixtures/ and tests/ stay outside the package's own tsc build by construction.
- **Rationale**: Verified directly by reading tsconfig.json: rootDir and include already scope compilation to src/ only, so the new fixtures/ and tests/ directories are structurally excluded from dist/ without any edit.
- **Fitness Function**: fixtures/ and tests/ never appear in the compiled dist/ output — verified by pnpm --filter @orieken/saturday-playwright-leaflet build; inspect dist/ for absence of fixtures/tests files

## Component Placement

- `fixtures/index.html` → Frameworks/Drivers (test harness, not shipped) layer, `packages/saturday-playwright-leaflet/fixtures/` — Static page used only to drive the Playwright suite against a known Leaflet map; excluded from tsc build by existing tsconfig scoping.
- `fixtures/serve.mjs` → Frameworks/Drivers (test infrastructure) layer, `packages/saturday-playwright-leaflet/fixtures/` — Minimal node:http static server for the Playwright webServer block; no new dependency, no production code path.
- `playwright.config.ts` → Frameworks/Drivers (test infrastructure) layer, `packages/saturday-playwright-leaflet/ (root)` — Follows apps/ye-olde-magic-shop/playwright.config.ts's existing chrome-project/webServer pattern exactly; no new structural shape introduced.
- `tests/fixture.spec.ts` → Adapters (E2E/UI test level per testing-conventions.md) layer, `packages/saturday-playwright-leaflet/tests/` — Exercises the package's public LeafletFixture API only, via the exported `test` fixture — never reaches into src/ internals directly.
- `src/fixture.ts (modified: Turf.js sourcing only)` → Adapters layer, `packages/saturday-playwright-leaflet/src/` — Already the adapter boundary between Playwright's page API and Leaflet's window-exposed map; the Turf-sourcing change stays within that same boundary.

## Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: None

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [PASS] fixture.ts (adapter layer) may depend on leaflet and @turf/turf directly; no domain layer exists in this package to violate — Confirmed by reading src/fixture.ts directly (mcp__loom__* graph tools were not invoked — no indication `loom index` has run for this worktree, so this used the documented fallback of reading imports/config by hand).
- [PASS] Test-only code (fixtures/, tests/) must not leak into the published package's dist/ output — tsconfig.json already sets rootDir: './src' and include: ['src/**/*'], structurally excluding fixtures/ and tests/ without any change.

## Anti-Pattern Check

- Premature generalization: not found — No new shared package, base class, or abstraction introduced; everything lives inside the one package it tests, matching CLAUDE.md's architect-skip criterion for additive work in an established pattern.
- Leaky abstraction: not found — LeafletFixture's 11-helper public API is unchanged; only Turf's internal sourcing (CDN to local) changes, invisible to callers.
- God object: not found
- Shotgun surgery: not found — Single package touched; the pnpm-lock.yaml update is a mechanical side-effect of `pnpm install` adding one dependency, not a structural ripple across files.

## Fitness Functions

- **Suite passes fully offline, in real Chrome, inside saturday-runner:debian with --network none**
  - **Verification**: the literal docker run command from AC4; exit code checked by devops-engineer in CI
  - **Owner**: devops-engineer
- **No CDN/network script or link tag reaches the browser from the fixture page or fixture.ts's library loading**
  - **Verification**: code-reviewer greps for http(s):// inside addScriptTag/<script src>/<link href>; confirmed end-to-end by the --network none run passing
  - **Owner**: code-reviewer
- **@turf/turf listed under dependencies (not devDependencies), pinned ^6**
  - **Verification**: code-reviewer greps package.json; pnpm-lock.yaml records it under this package's own dependencies block
  - **Owner**: code-reviewer
- **fixtures/ and tests/ excluded from tsc build output (dist/)**
  - **Verification**: pnpm --filter @orieken/saturday-playwright-leaflet build; ls dist/ contains no fixtures/tests files
  - **Owner**: devops-engineer
- **Cyclomatic complexity < 7 for every new/changed function in the package**
  - **Verification**: analyze-complexity skill / eslint complexity rule in CI
  - **Owner**: developer
- **Every helper's test demonstrably fails when its fixture value changes**
  - **Verification**: manual fixture-break-and-revert per helper, output recorded in the delivery record; no automated mutation tool is wired for this package
  - **Owner**: qa-engineer

## Refactoring Opportunities

- Replace raw `any` with typed narrowing (architecture-guardrails.md #4) on `src/fixture.ts: LeafletFixture's `any`-typed params/returns (metadata?: any, options?: any, (window as any), (layer: any))` — Pre-existing Boy Scout Rule debt, not introduced by this order. Deferred: a proper fix means defining typed LeafletMetadata/LeafletLayerOptions contracts, which changes this package's public API surface and needs its own RFC — doing it as a drive-by edit alongside the narrowly-scoped Turf CDN fix risks destabilizing the 11 helper tests this order is adding.

## Developer Handoff Notes

- Touch only the Turf.js script-tag line and its immediate surrounding code in src/fixture.ts for this order — do not refactor the file's existing `any` usage (see refactoringOpportunities); that's a separate, public-API-affecting change.
- Verify the installed @turf/turf ^6 version's node_modules layout actually exposes a browser-ready UMD/IIFE bundle (e.g. node_modules/@turf/turf/turf.min.js) before wiring the local <script> tag in fixtures/serve.mjs; if it doesn't, build one with esbuild at pretest time rather than committing a prebuilt copy.
- tsconfig.json needs no change — rootDir: './src' and include: ['src/**/*'] already exclude fixtures/ and tests/ from the compiled package; confirmed by direct read, not by the context graph.
- @turf/turf goes into package.json's dependencies block, not devDependencies — it runs in the browser under test, a runtime consumption path.
- mcp__loom__* graph tools were not invoked for this review (no sign `loom index` has run against this worktree); all layer-boundary findings above came from direct inspection of tsconfig.json and src/fixture.ts.

## Open Architectural Questions

None
