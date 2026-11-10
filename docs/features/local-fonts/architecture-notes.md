---
feature: "local-fonts"
bounded_context: "ye-olde-magic-shop frontend — typography and font-serving"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop (dependency + import in src/main.ts)"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: local-fonts

## Structural Decisions

### Self-host Cinzel and Inter via @fontsource/cinzel and @fontsource/inter, imported directly in apps/ye-olde-magic-shop/src/main.ts, rather than introducing a fonts-loading abstraction or a new shared package.
- **Rationale**: main.ts already imports a third-party package's CSS directly before mount (leaflet/dist/leaflet.css at main.ts:6) for exactly this reason — side-effect CSS that must be present before the app renders. @fontsource ships per-weight CSS files that register @font-face rules against bundler-resolved local font files, so importing them the same way extends an existing, working pattern instead of adding a new one. No new layer, base class, or package is justified by two additional CSS imports.
- **Trade-offs**:
  - Four separate imports (cinzel/400, cinzel/700, inter/400, inter/600) in main.ts instead of one line — accepted because @fontsource's per-weight file structure is what makes exact weight verification in the new test possible; a single barrel import would pull in weights the design doesn't use.
- **Fitness Function**: src/styles.css contains no @import pointing at any origin other than the app's own build output — verified by Playwright: tests/reagent-map-globe.spec.ts's no-request test, with PRE_EXISTING_OFF_ORIGIN_HOSTS and its skip branch removed, run via scripts/ui-suite-in-chrome.sh (and --ubi9) in CI
### Pin @fontsource/cinzel and @fontsource/inter to exact versions (5.3.0, no ^ or ~) in apps/ye-olde-magic-shop/package.json.
- **Rationale**: cesium is already exact-pinned in this same package.json ("cesium": "1.146.0", no range operator) while every other dependency uses ^. That establishes the local convention: packages whose CSS/asset output is asserted against in tests get exact pins so a minor bump can't silently change what's served without the test catching it first. Font files are exactly this kind of asset-shape-sensitive dependency — the new document.fonts.load test asserts on the exact weights @fontsource ships.
- **Trade-offs**:
  - Exact pins mean Dependabot/Renovate-style patch bumps won't land automatically; accepted because the AC explicitly requires it and the existing cesium precedent already accepts this cost for the same reason.
- **Fitness Function**: @fontsource/cinzel and @fontsource/inter versions in package.json carry no ^ or ~ prefix — verified by code-reviewer inspects apps/ye-olde-magic-shop/package.json and the pnpm-lock.yaml diff for exact version strings
### Do not introduce a fonts-loading utility, composable, or abstraction layer for this change.
- **Rationale**: The entire change is: remove one @import line, add four @fontsource CSS imports at the existing app-entry side-effect-import site, and tighten one test. There is no second call site, no conditional loading logic, and no future variability this chore needs to anticipate. Introducing an abstraction here would be premature generalization for a use case that doesn't exist yet (design-principles.md #7).
- **Fitness Function**: judgment-only — The entire change is: remove one @import line, add four @fontsource CSS imports at the existing app-entry side-effect-import site, and tighten one test. There is no second call site, no conditional loading logic, and no future variability this chore needs to anticipate. Introducing an abstraction here would be premature generalization for a use case that doesn't exist yet (design-principles.md #7).

## Component Placement

- `@fontsource/cinzel, @fontsource/inter (400/700 and 400/600 weight CSS)` → framework layer, `apps/ye-olde-magic-shop (dependency + import in src/main.ts)` — Font assets are a framework/presentation-layer concern, imported at the same side-effect site as leaflet's CSS; no domain or use-case code is affected.

## Bounded Context

- **Owning Context**: ye-olde-magic-shop frontend — typography and font-serving
- **Context Crossings**: reagent-map-globe's E2E no-request test (tests/reagent-map-globe.spec.ts) loses its PRE_EXISTING_OFF_ORIGIN_HOSTS exception as a direct consequence of this change

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] No OTel/instrumentation or domain logic introduced by this change — scope is CSS imports, a dependency pin, and a test file — Checked via Read of src/main.ts and src/styles.css directly (loom graph not queried — change has no new module boundaries to check).

## Anti-Pattern Check

- Premature generalization: not found — Considered and rejected introducing a fonts-loading helper; direct @fontsource imports at the existing entry-point side-effect-import site are sufficient and consistent with how leaflet's CSS is already loaded.
- Leaky abstraction: not found — No abstraction is being introduced, so none to leak.

## Fitness Functions

- **No page makes a request off the app's own origin (fonts included — no host exception remains in tests/reagent-map-globe.spec.ts)**
  - **Verification**: Playwright no-request test via scripts/ui-suite-in-chrome.sh, run in both the debian and --ubi9 images
  - **Owner**: devops-engineer
- **Headings render in self-hosted Cinzel and body text in self-hosted Inter at the exact pinned weights (400/700 Cinzel, 400/600 Inter)**
  - **Verification**: Playwright tests/local-fonts.spec.ts using document.fonts.load per family/weight pair and comparing matched FontFace.weight to the requested weight — never document.fonts.check, which the analyst confirmed returns a false positive even with the imports removed
  - **Owner**: qa-engineer
- **@fontsource/cinzel and @fontsource/inter stay pinned to exact versions (no ^ or ~ range), consistent with the existing cesium exact-pin precedent in the same package.json**
  - **Verification**: code-reviewer inspects apps/ye-olde-magic-shop/package.json and the pnpm-lock.yaml diff for exact version strings
  - **Owner**: code-reviewer

## Refactoring Opportunities

None

## Developer Handoff Notes

- Reuse the verified salvage commit rather than rewriting: git diff --binary 83beb6e 3629f44 -- apps/ye-olde-magic-shop/package.json apps/ye-olde-magic-shop/src/main.ts apps/ye-olde-magic-shop/src/styles.css apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts apps/ye-olde-magic-shop/tests/local-fonts.spec.ts | git apply --index
- main.ts already imports leaflet/dist/leaflet.css as a direct side-effect import before app.mount() (main.ts:6) — add the four @fontsource imports the same way, in the same file, before the mount call.
- package.json already exact-pins cesium ("cesium": "1.146.0", no ^) alongside otherwise ^-ranged dependencies — add @fontsource/cinzel and @fontsource/inter at "5.3.0" the same way, not with ^5.3.0.
- docs/testing-cesium.md must be hand-ported, not copied from the salvage commit verbatim: main added a Viewer-ellipsoid table row and a wheel-at-centre-of-view gap since the salvage commit (83beb6e) — only remove the Google Fonts 'Still open' paragraph and the matching 'Known gaps' bullet, leave the rest untouched.
- tests/local-fonts.spec.ts must use document.fonts.load, never document.fonts.check — the analyst's salvage run showed document.fonts.check returning true even with all @fontsource imports removed, which would make the test unable to fail.

## Open Architectural Questions

None
