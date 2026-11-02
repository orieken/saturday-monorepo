---
feature: "golden-workspace-packages"
bounded_context: "Test Automation"
domain_terms: ["Test Automation"]
files_touched: ["apps/ye-olde-magic-shop/package.json", "apps/ye-olde-magic-shop/saturday-core/**", "apps/ye-olde-magic-shop/saturday-cucumber/**", "apps/ye-olde-magic-shop/saturday-playwright-heatmap/**", "apps/ye-olde-magic-shop/Dockerfile", "apps/ye-olde-magic-shop/Dockerfile.runner", "packages/saturday-cucumber/src/config/index.ts", "packages/saturday-cucumber/tests/config.spec.ts", "pnpm-lock.yaml", "docs/features/golden-workspace-packages/"]
issue_refs: ["intake/golden-workspace-packages.md", "#7 (fix/cart-spec-magicshop-import)", "salvage/golden-workspace-packages @ f372a30"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: golden-workspace-packages

## Summary

The demo app apps/ye-olde-magic-shop currently depends on file:./saturday-core, file:./saturday-cucumber and file:./saturday-playwright-heatmap — vendored copies of the framework that have drifted from packages/ (confirmed: the vendored saturday-cucumber config reads process.env.HEADLESS, the workspace package does not). This feature switches those three devDependencies to workspace:*, deletes the vendored directories and both app Dockerfiles that COPY them, moves the one confirmed behavior gap (HEADLESS default) into packages/saturday-cucumber with a new test, and regenerates pnpm-lock.yaml — so the app's golden Cucumber/Playwright suite exercises the real framework packages instead of a stale copy. A prior attempt already did this work on commit f372a30 (branch salvage/golden-workspace-packages) and parked only on a since-fixed cart.spec.ts bug (#7); the task is to re-apply that diff on top of current main, verify every test file it deletes already has a content-equivalent in packages/*/tests/**, resolve the two Dockerfiles (fix to build from the workspace, or remove with a recorded reason), and prove scripts/ui-suite-in-chrome.sh now exits 0 where today it fails at collection for lack of a built dist/.

### Acceptance Criteria

- apps/ye-olde-magic-shop/package.json MUST declare @orieken/saturday-core, @orieken/saturday-cucumber and @orieken/saturday-playwright-heatmap as workspace:* dependencies, with no dependency in the file pointing at a path inside the app.
  - Example: Before: "@orieken/saturday-core": "file:./saturday-core" — After: "@orieken/saturday-core": "workspace:*" (matches the existing pattern already used for @orieken/saturday-playwright-k6-exporter and @orieken/saturday-playwright-leaflet in the same file)
- The vendored directories apps/ye-olde-magic-shop/saturday-core/, saturday-cucumber/ and saturday-playwright-heatmap/ MUST be gone. Any test inside them that covers framework code and has no content-equivalent already under packages/ MUST move there instead of being deleted.
  - Example: By filename, all 13 vendored test files (e.g. tests/ml/coverage.spec.ts, saturday-cucumber/tests/world.spec.ts, saturday-playwright-heatmap/tests/tracker.spec.ts) already have a same-named file under packages/*/tests/** — content parity still needs byte-level confirmation before trusting the deletion, per the Open Questions below.
- Every Cucumber scenario and Playwright spec that passed before the change MUST still pass after it, with no .feature file or spec assertion changed. The before and after results MUST be recorded in the delivery record.
- apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner, which both COPY the vendored directories, MUST each build from the workspace instead, or be removed with the reason recorded.
  - Example: Neither file is referenced by any GitHub Actions workflow, docker-compose file, or script in this repo — the only other Dockerfile.runner reference found (prototype-runner/local-cluster/run-local-k8s.sh) uses a generic $RUNNER_DIR unrelated to this app. Removal with that reason recorded is a valid reading of the criterion's 'or', matching what the salvaged commit already did.
- Running scripts/ui-suite-in-chrome.sh from the repository root on a fresh clone MUST exit 0, with its output recorded in the delivery record.
  - Example: The script's install_dependencies step already runs `pnpm --filter "@orieken/ye-olde-magic-shop^..." -r build` before the suite — once the three deps are workspace:* and pnpm-lock.yaml is regenerated, this step builds their dist/ and the collection failure described in the spec should clear without any script change.
- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except this run's own delivery record under docs/features/golden-workspace-packages/.

### Non-Functional Requirements

None

## Proposed Fitness Functions

- **apps/ye-olde-magic-shop/package.json never re-introduces a file: dependency on a vendored framework copy**
  - **Verification**: judgment-only — code review checks devDependencies at review time; the spec's Infrastructure/Deploy Notes say None, so no new CI check is being wired this run (approval-gates.md #7 would gate that decision separately)
  - **Owner**: code-reviewer

## Out of Scope

- Renaming files, adding flows, or changing step definitions — explicitly deferred to later orders per the spec.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: apps/ye-olde-magic-shop (consumer app) depends on packages/saturday-core, packages/saturday-cucumber, packages/saturday-playwright-heatmap (framework packages) — this feature makes that dependency a real workspace link instead of a drifted vendored copy.

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/package.json` — Switch the three vendored-copy devDependencies to workspace:*
- `apps/ye-olde-magic-shop/saturday-core/**` — Vendored copy to remove; confirm each test file's content matches its packages/saturday-core/tests/** counterpart before deleting
- `apps/ye-olde-magic-shop/saturday-cucumber/**` — Vendored copy to remove; its config/index.ts HEADLESS behavior is confirmed drifted and must land in the package first
- `apps/ye-olde-magic-shop/saturday-playwright-heatmap/**` — Vendored copy to remove; confirm each test file's content matches its packages/saturday-playwright-heatmap/tests/** counterpart
- `apps/ye-olde-magic-shop/Dockerfile` — COPYs the vendored directories; fix to build from workspace or remove with reason recorded
- `apps/ye-olde-magic-shop/Dockerfile.runner` — Same vendored-copy coupling as Dockerfile, used only by scripts/ui-suite-in-chrome.sh's image build step indirectly (not referenced directly — confirm before removing)
- `packages/saturday-cucumber/src/config/index.ts` — Move the confirmed HEADLESS env-var default behavior from the vendored copy into the workspace package
- `packages/saturday-cucumber/tests/config.spec.ts` — New test proving the HEADLESS default — must demonstrate it can fail, per the Test Repair Contract's rule for new tests
- `pnpm-lock.yaml` — Regenerate against current main after the dependency and directory changes
- `docs/features/golden-workspace-packages/` — Delivery record: before/after suite results, the Dockerfile decision and its reason, and resolution of the Open Questions below — the one path allowed to change outside scope.paths

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Re-apply the salvaged diff on top of current main: git diff --binary e9e4ae6 f372a30 | git apply --index, from the repository root, inside this worktree.
2. Resolve any apply conflict against the #7 cart.spec.ts fix that landed after e9e4ae6 — reconcile manually rather than discarding either change.
3. Run pnpm install from the repository root to regenerate pnpm-lock.yaml against current main, rather than trusting the salvaged lockfile diff byte-for-byte.
4. For each of the 13 vendored test files the salvaged diff deletes, diff its content against the same-named file under packages/saturday-core/tests/**, packages/saturday-cucumber/tests/**, or packages/saturday-playwright-heatmap/tests/** — confirm genuine content parity, not just a filename match, before trusting the deletion.
5. If any vendored test or source file has behavior or coverage the package version lacks beyond the already-confirmed HEADLESS case, move that behavior into the package with a test, in the same order as the HEADLESS fix, and name it in the delivery record.
6. Decide apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner: confirm neither is invoked by any workflow, script, or compose file in this repo (none found as of this analysis), then either remove both with that reason recorded in the delivery record, or rewrite both to build from the workspace instead of COPYing the vendored directories.
7. Verify packages/saturday-cucumber/src/config/index.ts's headless default change (process.env.HEADLESS === 'true') and its new config.spec.ts match the project's existing SaturdayConfig shape and do not alter any other default.

### QA Tasks

1. Before applying any change, run scripts/ui-suite-in-chrome.sh from the repository root on current main and record its exact failing output (collection failure) for the delivery record's before/after comparison.
2. After the change, run scripts/ui-suite-in-chrome.sh again and confirm it exits 0; record its output.
3. Run the app's existing test:bdd and test:chrome suites before and after the change; confirm every scenario that passed before still passes, with zero .feature or spec assertion diffs.
4. Run packages/saturday-cucumber's vitest suite and confirm config.spec.ts passes and exercises both HEADLESS=true and HEADLESS unset/false branches as a test that can fail.
5. Confirm every changed line in packages/saturday-cucumber/src/config/index.ts is covered at or above 85%, per ADR-008.
6. If any spec fails through the script for a reason unrelated to the vendored packages (for example map.spec.ts needing network access to OpenStreetMap tiles), do not change the spec or the script — stop at the ship gate and name the spec and its cause in the delivery record.

### Tech Writer Tasks

1. Write the delivery record under docs/features/golden-workspace-packages/ capturing: before/after scenario-pass counts from scripts/ui-suite-in-chrome.sh and the Cucumber/Playwright suites, the Dockerfile decision and its recorded reason, and the resolution (or carried-forward assumption) for each Open Question below.
2. No README, API doc, or ADR changes are expected from this change — state that explicitly in the delivery record rather than leaving it unaddressed.

### DevOps Tasks

1. None required by this spec — Infrastructure / Deploy Notes state None, and no new CI wiring is requested. Confirm the regenerated pnpm-lock.yaml does not break any existing CI install step that pins to the prior lockfile hash.

## Edge Cases and Risks

- git apply of the salvaged diff conflicts because main advanced past e9e4ae6 (notably the cart.spec.ts fix in #7): Reconcile manually, preserving both the cart.spec.ts import fix and the salvaged workspace/vendored-removal changes — do not discard either.
- A vendored test file's content differs from its same-named packages/ counterpart in a way that adds real coverage the package lacks: Per the spec's Open Questions, move that behavior into the package with a test, in this order, and name it in the delivery record — do not delete the vendored file first.
- A golden spec fails through scripts/ui-suite-in-chrome.sh for a reason unrelated to the vendored packages (e.g. the map page's OpenStreetMap tiles needing network access): Do not change the spec or the script; stop at the ship gate and name the spec and its cause in the delivery record, per the spec's own Open Questions resolution.
- apps/ye-olde-magic-shop/Dockerfile or Dockerfile.runner turn out to have an external CI consumer not visible in this repo: If found, rewrite both to build from the workspace instead of removing them; if none is found by delivery time, default to removal with the reason recorded, matching the salvaged commit's choice.

## Definition of Done

- All six acceptance criteria above are verified true, with before/after evidence recorded in docs/features/golden-workspace-packages/.
- apps/ye-olde-magic-shop/package.json has no remaining file: dependency on a vendored framework copy.
- apps/ye-olde-magic-shop/saturday-core/, saturday-cucumber/, and saturday-playwright-heatmap/ no longer exist, and every test with real coverage they held now lives under packages/*/tests/**.
- scripts/ui-suite-in-chrome.sh exits 0 on a fresh clone, with its output recorded.
- apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner either build from the workspace or are removed with the reason recorded.
- Every new or changed function has cyclomatic complexity below 7 and every changed line is covered by a test that can fail (ADR-008 mutation expectation).
- No file outside scope.paths changed, except this run's own delivery record under docs/features/golden-workspace-packages/.
