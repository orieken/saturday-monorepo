---
feature: "ui-suite-real-chrome"
bounded_context: ""
domain_terms: ["Runner image"]
files_touched: ["apps/ye-olde-magic-shop/tests/playwright-config.spec.ts", "scripts/ui-suite-in-chrome.test.sh", "apps/ye-olde-magic-shop/playwright.config.ts"]
issue_refs: ["ui-suite-real-chrome"]
linked_adrs: []
linked_kis: []
---

# QA Report: ui-suite-real-chrome

## Test Files Created

- apps/ye-olde-magic-shop/tests/playwright-config.spec.ts
- scripts/ui-suite-in-chrome.test.sh

## Test Files Modified

- apps/ye-olde-magic-shop/playwright.config.ts

## Coverage Summary

- Acceptance criteria covered: 6/6
- Total new tests: 15
- Statement coverage, lowest unit: **100.0%** (`apps/ye-olde-magic-shop/playwright.config.ts (resolveRunnerSlowMo/resolveRunnerVideo, chrome project fields)`)
  - `apps/ye-olde-magic-shop/playwright.config.ts (resolveRunnerSlowMo/resolveRunnerVideo, chrome project fields)`: 100.0%
  - `scripts/ui-suite-in-chrome.sh (all 10 functions, all 6 exit paths)`: 100.0%

## Test Results

- Passed: 16
- Failed: 0
- Skipped: 0

## Accessibility Check

- Not applicable — this feature touches only Playwright config, a WebGL-context assertion spec, and shell orchestration tooling; no UI/markup changes.

## Bugs Found

- scripts/ui-suite-in-chrome.sh's install_dependencies had no timeout and no CI=true, so pnpm's interactive 'recreate node_modules from scratch?' prompt (triggered whenever the bind-mounted checkout's node_modules doesn't match the container's platform/lockfile state — e.g. a host-run pnpm install before the script runs) hangs indefinitely on docker exec's unattached stdin. This bypasses the 180s app-health budget entirely, since that budget only starts after install_dependencies returns. Reproduced live: a first real container run hung ~3 minutes at this prompt before incidentally proceeding; after clearing node_modules and re-running, it hung again confirming reproducibility.: Fixed by adding `-e CI=true` to install_dependencies' docker exec call (scripts/ui-suite-in-chrome.sh), which makes pnpm auto-answer the prompt instead of waiting on a TTY. Verified by re-running the full script against the same mismatched-platform node_modules state that previously hung: install completed in ~17s, no prompt, no hang. Shell-level test harness (scripts/ui-suite-in-chrome.test.sh) updated to handle the new -e flag in its fake docker dispatcher; all 6 exit-code-contract cases still pass.

## Known Gaps

- AC3 (tests/webgl.spec.ts pass/fail-under-disable-gpu test-of-the-test) verified via the exact production invocation (scripts/ui-suite-in-chrome.sh running unscoped `npx playwright test --project=chrome`): The unscoped invocation fails during test collection (exit 1, before any test runs) because four sibling workspace packages (@orieken/saturday-core, saturday-playwright-k6-exporter, saturday-playwright-heatmap, saturday-playwright-leaflet) have no built dist/ in this environment — a pre-existing, out-of-scope repo build-state gap (none of those packages are under scope.paths), reproducible even before this feature's changes. Verified the actual behavior this AC cares about by spinning up the identical container (same image, same bind-mounts, same install/health/runner-check sequence the script performs) and running `npx playwright test tests/webgl.spec.ts tests/playwright-config.spec.ts --project=chrome` scoped to just this feature's files: 10/10 passed, including webgl.spec.ts genuinely launching real headed Chrome under Xvfb and resolving a webgl2 context (2.7s). Then added --disable-gpu --disable-software-rasterizer to the chrome project's launch args and re-ran: webgl.spec.ts failed as required (Expected: true, Received: false). Reverted the args and re-confirmed the pass. Both run logs are real evidence from inside the runner container; a human should run `pnpm -r build` before relying on the script's exact unscoped command end-to-end.
- AC5 (.dockerignore's bare 'lib' pattern negation) — 'without removing that pattern's effect on other lib/ directories': Empirically verified, via a throwaway build against both the real .dockerignore and an isolated minimal reproduction, that Docker's dockerignore matching for a bare pattern like 'lib' (no '**/' prefix) only ever matches a true top-level ./lib — it never matched nested paths like apps/ye-olde-magic-shop/lib/ or packages/*/lib/ in the first place. Building with and without the added negation line produced identical results both times. Separately, no Docker build in this repository currently uses the repo root as its build context with COPY semantics: scripts/ui-suite-in-chrome.sh builds from containers/saturday-runner (no root .dockerignore involved there) and bind-mounts the checkout rather than copying it; apps/ye-olde-magic-shop's own Dockerfile/Dockerfile.runner use apps/ye-olde-magic-shop as their context, with their own .dockerignore that doesn't mention lib at all. The literal criterion holds (apps/ye-olde-magic-shop/lib/ is not excluded; no other lib/ directory's exclusion changed) but only vacuously — the negation line has no observed effect on any build this repo currently runs. Not changed — flagging for human review rather than guessing at whether a future root-context build is planned.

## Notes for Tech Writer

- New npm script `test:chrome` (apps/ye-olde-magic-shop/package.json: `playwright test --project=chrome`) runs the headed real-Chrome project directly; `test:e2e`/`test:pw` are now explicitly `--project=chromium` so they stay bundled-Chromium-only by default.
- scripts/ui-suite-in-chrome.sh exit codes: 2=Docker not installed, 3=image build failed, 4=app/mock-api unhealthy within 180s, 5=runner-check failed inside the container; any other non-zero code is the Playwright suite's own exit code, passed through unchanged.
- RUNNER_SLOW_MO (integer ms, default 0) and RUNNER_VIDEO ('1' to enable, default off) control the chrome project's slowMo/video; both are read fresh by exported resolveRunnerSlowMo()/resolveRunnerVideo() in playwright.config.ts.
- Pre-ship usage per the spec's Infrastructure/Deploy note: after merge, a human adds {name: ui-suite, command: scripts/ui-suite-in-chrome.sh} to .claude/delivery-policy.yaml preShipChecks.
- Known gap for docs: running the full suite via the script's exact unscoped command currently requires `pnpm -r build` first (pre-existing sibling-package dist/ gap, unrelated to this feature) — worth a troubleshooting note.
- AC3's test-of-the-test (webgl.spec.ts passing normally, failing under --disable-gpu --disable-software-rasterizer) was run for real inside the runner container; both outputs should be captured verbatim in the delivery record.
