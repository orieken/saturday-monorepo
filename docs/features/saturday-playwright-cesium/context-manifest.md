# Context Manifest: saturday-playwright-cesium

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `packages/saturday-playwright-cesium/README.md` — Not yet on disk — created by applying commit 926ad56 (`git diff --binary 68a5901 926ad56 | git apply --index`, per the spec's Starting Point). Documents the seam contract (`window.__SATURDAY_CESIUM__ = { viewer, Cesium }`, `entity.saturdayMetadata = { id, name, kind }`) and all 7 helpers named in the acceptance criteria: waitForSettled, getEntities, toScreen, clickAt, flyTo, metersPerPixelAt, globeRadiusMeters. Confirms the fixture planet is named 'Aurelis', radius 6,000 km, and states the test command and the --network none docker contract verbatim.
- `packages/saturday-playwright-cesium/src/geometry.ts` — Not yet on disk (same apply step as README.md). Implements distance/lineLength/pointInCircle, radius-parameterized per AC8. Its own header comment already asserts a structural decision — 'see architecture-notes.md's structural decision to keep this package-local rather than a new packages/saturday-geo' — pre-empting the spec's open question about saturday-geo. No architecture-notes.md exists yet for this feature, so this is a decision baked into the salvage code that the architect stage must ratify or override, not re-derive blind.
- `packages/saturday-playwright-cesium/src/fixture.ts` — Not yet on disk (same apply step). The CesiumFixture implementation backing every helper-level acceptance criterion (waitForSettled, getEntities, toScreen, clickAt, flyTo) — needed to confirm each AC is actually implemented as described, since this order verifies existing code rather than writing it.
- `packages/saturday-playwright-cesium/tests/fixture.spec.ts` — Not yet on disk (same apply step). The existing Playwright suite (11 tests per the spec). The Starting Point says 'verify each criterion... do not rewrite it' — the analyst needs to see what's already tested to plan the QA fixture-break-and-revert verification task per helper, not assume tests are missing.
- `packages/saturday-playwright-cesium/tests/geometry.spec.ts` — Not yet on disk (same apply step). The existing Vitest suite (7 tests per the spec) for the Node-only geometry helpers — same verification-not-authoring concern as fixture.spec.ts.
- `packages/saturday-playwright-leaflet/README.md` — Exists today. The cesium package's own (not-yet-applied) README explicitly says it follows 'the same seam-contract pattern as @orieken/saturday-playwright-leaflet' — this is the sibling pattern being mirrored, useful to confirm consistency.
- `docs/features/leaflet-package-tests/analysis.md` — Nearly identical prior feature in the same bounded context (Owning Context: Test Automation, confirmed by grep across docs/features/*/analysis.md) — the proposed fitness functions (offline/real-Chrome run, complexity<7, per-helper fail-proof) and task-breakdown shape here are the direct template for this order's analysis.md.
- `docs/features/ui-suite-real-chrome/architecture-notes.md` — Same bounded context; this is the only prior feature that put a `channel: 'chrome'` + ANGLE/SwiftShader WebGL project into CI, which Cesium's AC9 (real Chrome, --use-gl=angle --use-angle=swiftshader) depends on directly — it is also this order's declared depends-on.
- `docs/runner-image.md` — 394 lines, under the 500-line range-constraint threshold, but focus on lines 1-100 (the ANGLE/SwiftShader software-WebGL design, directly load-bearing for this package unlike for Leaflet) and the Debian build's verify record (~225-245, showing a real webgl2 context renders in the image) — AC9's docker command and WebGL assumption both rest on this.
- `docs/specs/reagent-map-spec/discovery.md` — Line 65 (Q8) is the canonical source for the fixture planet: 'A sphere of radius 6,000 km' — resolves the spec's own Open Question naming this file as the authority for Aurelis's radius, which the fixture, README, and geometry tests all depend on matching.

## Token Budget

- **Status**: INCOMPLETE — ~15117 tokens, 12.6% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `packages/saturday-playwright-cesium/README.md` | — | not measured: cannot read packages/saturday-playwright-cesium/README.md |
| `packages/saturday-playwright-cesium/src/geometry.ts` | — | not measured: cannot read packages/saturday-playwright-cesium/src/geometry.ts |
| `packages/saturday-playwright-cesium/src/fixture.ts` | — | not measured: cannot read packages/saturday-playwright-cesium/src/fixture.ts |
| `packages/saturday-playwright-cesium/tests/fixture.spec.ts` | — | not measured: cannot read packages/saturday-playwright-cesium/tests/fixture.spec.ts |
| `packages/saturday-playwright-cesium/tests/geometry.spec.ts` | — | not measured: cannot read packages/saturday-playwright-cesium/tests/geometry.spec.ts |
| `packages/saturday-playwright-leaflet/README.md` | 4600 | 1150 |
| `docs/features/leaflet-package-tests/analysis.md` | 11408 | 2852 |
| `docs/features/ui-suite-real-chrome/architecture-notes.md` | 15105 | 3776 |
| `docs/runner-image.md` | 23746 | 5936 |
| `docs/specs/reagent-map-spec/discovery.md` | 5614 | 1403 |

> 5 pinned file(s) could not be measured, so the total above is a floor, not a figure.

## Knowledge Items

None

## ADRs

None

## Pruned

- `No KI/ADR corpus exists in this repo (no shared/knowledge/, .claude/knowledge/, or docs/adrs(/adr)/ directory found) — search-ki lookup was skipped for lack of a corpus, not skipped by choice.`
- `Blast radius was not graph-derived: no mcp__loom__* tools are registered in this session. A Grep-based dependents check is also moot here: packages/saturday-playwright-cesium does not exist in this worktree yet, so nothing in the repo can import it today — there is no blast radius to compute until the package exists.`
- `Prior deliveries in this bounded context (Test Automation): leaflet-package-tests, ui-suite-real-chrome, and golden-workspace-packages all found via grep for 'Owning Context'/'Test Automation' across docs/features/*/analysis.md. None of the three (nor any feature in the archive) has a retrospective.md, so there are no retrospective lessons to surface.`
- `packages/saturday-playwright-cesium/package.json and playwright.config.ts (salvage commit 926ad56, not yet applied) — excluded from the pin list to stay at the 10-file cap; their load-bearing facts (cesium 1.146.0, test script runs vitest then playwright, the chrome project's GL launch args) are already fully covered inline in the README.md pin and docs/runner-image.md.`
- `packages/saturday-playwright-leaflet/src/fixture.ts and docs/features/leaflet-package-tests/context-manifest.md — excluded: the leaflet README pin already captures the seam-contract pattern being mirrored, and this order's own fitness-function/task-breakdown template is better served by leaflet-package-tests/analysis.md directly rather than its context manifest.`
- `pnpm-workspace.yaml and root package.json — checked (workspace globs are apps/* and packages/*, no packageManager-pin conflict); excluded from pins since they add no information beyond confirming the new package slots into the existing packages/* glob with no special-case needed.`
- `DOMAIN_DICTIONARY.md — checked for 'seam', 'settled', 'cesium', 'geo'; no entries exist, so no ubiquitous-language constraint applies and nothing to pin from it.`
