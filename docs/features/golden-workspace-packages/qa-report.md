---
feature: "golden-workspace-packages"
bounded_context: ""
domain_terms: ["SaturdayConfig"]
files_touched: []
issue_refs: ["intake/golden-workspace-packages.md", "#7 (fix/cart-spec-magicshop-import)", "salvage/golden-workspace-packages @ f372a30"]
linked_adrs: []
linked_kis: []
---

# QA Report: golden-workspace-packages

## Test Files Created

None

## Test Files Modified

None

## Coverage Summary

- Acceptance criteria covered: 5/6
- Total new tests: 4
- Statement coverage, lowest unit: **89.2%** (`packages/saturday-core`)
  - `packages/saturday-cucumber`: 97.0%
  - `packages/saturday-core`: 89.2%
  - `packages/saturday-playwright-heatmap`: 100.0%

## Test Results

- Passed: 234
- Failed: 0
- Skipped: 0

## Accessibility Check

None

## Bugs Found

None

## Known Gaps

- AC4 — apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner each build from the workspace or are removed with the reason recorded.: A real external consumer of the removed Dockerfile exists and was missed: prototype-runner/docker-compose.yml's web-app service builds with `context: ../apps/ye-olde-magic-shop` and no explicit `dockerfile:` key, so Docker Compose looks for the implicit default `Dockerfile` in that directory — now deleted. Reproduced directly: `docker compose build web-app` in prototype-runner/ fails with 'failed to read dockerfile: open Dockerfile: no such file or directory'. The developer's verification grepped literally for the strings 'Dockerfile'/'Dockerfile.runner' across workflows/scripts/compose files, which does not catch an implicit default-filename reference with no literal mention. Per analysis.md's own edge case ('an external CI consumer not visible in this repo' → rewrite rather than remove), the correct resolution is to make apps/ye-olde-magic-shop/Dockerfile build from the workspace — but doing that requires changing prototype-runner/docker-compose.yml's build context to the repo root (plus a pnpm-workspace-aware multi-stage rewrite), and prototype-runner/** is outside this run's scope.paths. Fixing it would itself violate AC6's 'no file outside scope.paths changes' constraint, so I did not make that change. This needs a human/architect call: either approve a scope.paths extension to include prototype-runner/docker-compose.yml so the Dockerfile can be rewritten properly, or confirm prototype-runner is already stale (its own WORKFLOW.md describes a `final/` directory layout that no longer exists after the 'extract starter kits to standalone repos' refactor, commit 3f650ab) and accept this as pre-existing drift the ship gate should note but not block on. Dockerfile.runner has no other consumer found anywhere in the repo (workflows, scripts, compose files) — its removal is unaffected by this finding.

## Notes for Tech Writer

- AC1 verified: apps/ye-olde-magic-shop/package.json declares @orieken/saturday-core, @orieken/saturday-cucumber and @orieken/saturday-playwright-heatmap as workspace:* devDependencies; no dependency points at a path inside the app.
- AC2 verified: all three vendored directories are gone. Independently diffed the two vendored test files the developer flagged as non-identical (saturday-cucumber's coverage.spec.ts and world.spec.ts) against their packages/ counterparts — both package versions are strict supersets (same assertions plus additional tabManager/siteManager/cleanup coverage). No vendored test coverage was lost.
- AC3 verified after the migration: apps/ye-olde-magic-shop's Cucumber suite (test:bdd) — 16/16 scenarios, 85/85 steps passed, with the app dev server and mock-api running locally. The dockerized Chrome Playwright suite (via scripts/ui-suite-in-chrome.sh) — 24/24 specs passed, covering magic-shop, map, cart, heatmap, webgl and playwright-config specs. No .feature file or spec assertion was touched by this change (confirmed via git status/diff). 'Before' state was not independently re-run (would require reverting a 113-file diff and reinstalling); relying on the spec's own documented starting condition ('fails at collection, vendored copies have no dist/ after fresh install') corroborated by the developer's independent confirmation — the full green run after migration is the regression check that matters here.
- AC4 NOT fully satisfied — see knownGaps. Dockerfile.runner's removal is clean (no other consumer found); Dockerfile's removal broke prototype-runner/docker-compose.yml's web-app service, a consumer the developer's grep missed because it relies on Docker's implicit default filename rather than a literal path string.
- AC5 verified: scripts/ui-suite-in-chrome.sh exits 0 from a fresh install, 24/24 Chrome Playwright tests pass. First attempt failed with an unrelated 'Cannot find module @rollup/rollup-linux-x64-gnu' error inside the Docker build step — traced to a stale `ui-suite-pnpm-store` Docker volume (created 2026-10-30, predating this change) rather than anything in this diff. Removed that volume and re-ran clean, which passed. Recording both runs for the delivery record: first run (stale volume) failed at `packages/saturday-core build` with the rollup native-binding error, exit 1; second run (clean volume) completed all 16/16 build steps and 24/24 Playwright tests, exit 0.
- AC6 verified: the only changed production line (packages/saturday-cucumber/src/config/index.ts:26) has cyclomatic complexity 1. Independently re-verified the mutation-can-fail claim by temporarily reverting the line to the old hardcoded `headless: false` and re-running the suite — config.spec.ts's 'is true when HEADLESS is the literal string true' case failed as expected (1 failed/22), then passed again (22/22) once restored. No file outside scope.paths changed in the final diff — two test-run side effects (apps/mock-api/data/orders.json and apps/ye-olde-magic-shop/heatmap-data/*.json, mutated by running the Cucumber/Playwright suites locally) were reverted with `git checkout` before finishing, since mock-api/data is outside scope.paths entirely and heatmap-data is run-artifact noise, not a feature change.
- packages/saturday-cucumber/src/config/index.ts's changed line is excluded from the package's own vitest coverage % report by a pre-existing, project-wide `exclude: ['**/index.ts', ...]` rule in vitest.config.ts (confirmed by reading the file directly) — not introduced by this change. The 97.05% figure above is the package's overall statement coverage with that exclusion in place; the changed line's own coverage is proven by the fail/pass mutation check above instead.
- packages/saturday-core and packages/saturday-playwright-heatmap had no source lines changed by this feature (only consumed via workspace symlink instead of a vendored copy) — their suites (166/166 and 6/6 respectively) are reported to confirm they still pass green after the dependency-wiring change, not because new lines there need coverage.
