---
feature: "golden-workspace-packages"
bounded_context: "Test Automation"
domain_terms: ["SaturdayConfig", "BaseSite", "BasePage", "BaseElement", "BaseFlow"]
files_touched: ["packages/saturday-cucumber/src/config/index.ts", "packages/saturday-cucumber/tests/config.spec.ts"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: golden-workspace-packages

## Structural Decisions

### Switch @orieken/saturday-core, @orieken/saturday-cucumber and @orieken/saturday-playwright-heatmap in apps/ye-olde-magic-shop/package.json from file:./<dir> to workspace:*, matching the existing workspace:* entries already present for @orieken/saturday-playwright-k6-exporter and @orieken/saturday-playwright-leaflet in the same devDependencies block.
- **Rationale**: Confirmed by direct read of apps/ye-olde-magic-shop/package.json: two sibling framework deps in the same file already use workspace:*. This is the established pattern for every other Saturday package this app consumes — the three file: deps are the outliers, not a new choice. No new abstraction, package, or layer boundary is introduced; the dependency direction (app -> framework package) is unchanged and already correct.
- **Trade-offs**:
  - The app loses its ability to diverge from the framework package without that divergence being visible as a diffable change inside packages/ — this is the explicit goal of the feature, not a cost to mitigate.
  - Any future intentional per-app override of Saturday-core/cucumber/heatmap behavior must go through the package's own configure() API or a documented extension point, not a vendored fork.
- **Fitness Function**: judgment-only — Confirmed by direct read of apps/ye-olde-magic-shop/package.json: two sibling framework deps in the same file already use workspace:*. This is the established pattern for every other Saturday package this app consumes — the three file: deps are the outliers, not a new choice. No new abstraction, package, or layer boundary is introduced; the dependency direction (app -> framework package) is unchanged and already correct.
### Fix the confirmed HEADLESS default drift inside packages/saturday-cucumber/src/config/index.ts (defaultConfig.browser.headless), changing the hardcoded false to process.env.HEADLESS === 'true', matching the vendored copy's already-correct behavior. Do not change any other field in SaturdayConfig or defaultConfig.
- **Rationale**: Read both files directly: apps/ye-olde-magic-shop/saturday-cucumber/src/config/index.ts:26 reads process.env.HEADLESS; packages/saturday-cucumber/src/config/index.ts:26 hardcodes false. Every other line in both files is byte-identical. This is the one behavior the package genuinely lacks, confirmed, not assumed.
- **Trade-offs**:
  - Changes the workspace package's default behavior for every consumer, not just this app — any other app currently relying on the package's unconditional headless:false default will start running headed unless it sets HEADLESS=true or overrides via configure(). No other in-repo consumer of packages/saturday-cucumber was found as part of this check; the developer should grep for other consumers before landing this.
- **Fitness Function**: packages/saturday-cucumber's default browser.headless value follows process.env.HEADLESS === 'true' (false when unset or any value other than the literal string 'true'), with no regression to any other defaultConfig field. — verified by packages/saturday-cucumber/tests/config.spec.ts (new) — vitest asserting getConfig().browser.headless reflects HEADLESS unset/'true'/'false' cases; must be shown to fail against the current hardcoded false before the fix lands, per the Test Repair Contract's rule for new tests.

## Component Placement

- `HEADLESS env-var default for browser.headless` → Frameworks & Drivers (test-runner configuration) layer, `packages/saturday-cucumber/src/config/index.ts` — This is the canonical location already — the vendored copy being deleted is a drifted fork of this exact file. No new layer or package is created; the fix lands in the file that already owns SaturdayConfig.
- `config.spec.ts (HEADLESS default test)` → Frameworks & Drivers (test) layer, `packages/saturday-cucumber/tests/config.spec.ts` — Mirrors the existing tests/ layout in the same package (coverage.spec.ts, world.spec.ts already live there); no new test directory structure needed.

## Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: apps/ye-olde-magic-shop (consumer app, outer/application layer for this feature's purposes) -> packages/saturday-core, packages/saturday-cucumber, packages/saturday-playwright-heatmap (framework packages) — the crossing already exists today via the vendored copies; this feature only changes the link from a drifted file: copy to a live workspace:* link. Direction of dependency (app depends on framework, never the reverse) is unchanged.

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [PASS] Consumer app (apps/ye-olde-magic-shop) may depend on framework packages (packages/saturday-core, packages/saturday-cucumber, packages/saturday-playwright-heatmap); framework packages must never depend back on the app. — Checked by reading apps/ye-olde-magic-shop/package.json directly (mcp__loom__dependencies_of was not used — graph tool availability not confirmed in this session, so imports/package.json were read directly and that fallback is being stated here per process step 3). The three deps being switched are devDependencies the app already requires; nothing in packages/saturday-core, packages/saturday-cucumber, or packages/saturday-playwright-heatmap references apps/ye-olde-magic-shop. Direction is correct both before and after this change.
- [PASS] packages/saturday-cucumber's config module (use-case/framework boundary — it is pure configuration state, not a domain entity) may read process.env directly. — The vendored copy already does this at the same file path with no domain-entity import violation; the package version gains parity, not a new boundary crossing.

## Anti-Pattern Check

- Premature generalization: not found — No new abstraction is introduced — the fix makes an existing field's behavior consistent with its already-published TypeScript interface (SaturdayConfig.browser.headless was always meant to be configurable; the package's own configure() API already allows overriding it, the hardcoded default was simply wrong).
- Distributed monolith / shotgun surgery: not found — Three dependency entries change in one file, plus a one-line default and a new test file. No ripple into unrelated files.

## Fitness Functions

- **packages/saturday-cucumber's default browser.headless value follows process.env.HEADLESS === 'true' (false when unset or any value other than the literal string 'true'), with no regression to any other defaultConfig field.**
  - **Verification**: packages/saturday-cucumber/tests/config.spec.ts (new) — vitest asserting getConfig().browser.headless reflects HEADLESS unset/'true'/'false' cases; must be shown to fail against the current hardcoded false before the fix lands, per the Test Repair Contract's rule for new tests.
  - **Owner**: developer
- **apps/ye-olde-magic-shop/package.json never re-introduces a file: dependency on a vendored framework copy.**
  - **Verification**: judgment-only — code review checks devDependencies at review time (no new CI check being wired this run; carried forward unchanged from the analyst's proposed fitness function).
  - **Owner**: code-reviewer

## Refactoring Opportunities

None

## Developer Handoff Notes

- No RFC is being produced for this feature: the dependency switch matches an already-established in-file pattern (workspace:* for k6-exporter and leaflet), and the HEADLESS fix is a one-line default correction scoped to a single file plus one new test. Neither crosses a layer boundary that doesn't already exist, nor introduces a new abstraction. Reversibility: Cheap for both changes (a git revert fully undoes either).
- Before changing packages/saturday-cucumber/src/config/index.ts, grep the monorepo for other consumers of @orieken/saturday-cucumber's defaultConfig/getConfig to confirm no other app relies on the current unconditional headless:false default.
- Place the new test at packages/saturday-cucumber/tests/config.spec.ts, following the existing flat-file convention in that tests/ directory (coverage.spec.ts, world.spec.ts) — no new subdirectory needed.
- This architecture review did not use the mcp__loom__dependencies_of/dependents_of/describe graph tools — their availability was not established in this session, so all boundary checks above were done by direct file reads, as the process fallback requires. Confirm independently if the graph is available to you.
- All other developer tasks (re-applying the salvaged diff, resolving the Dockerfile question, regenerating pnpm-lock.yaml, confirming the 13 vendored test files' content parity) are execution detail already fully specified in analysis.md — no structural decision from this stage changes any of them.

## Open Architectural Questions

None
