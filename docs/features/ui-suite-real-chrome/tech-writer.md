---
feature: "ui-suite-real-chrome"
bounded_context: "Test Automation"
domain_terms: ["Runner image"]
files_touched: []
issue_refs: ["ui-suite-real-chrome"]
linked_adrs: []
linked_kis: []
---

# Documentation Report: ui-suite-real-chrome

## Summary

This feature wires the demo app's Playwright suite to the existing Saturday runner image, so the
Cesium map tests run in real Chrome under Xvfb with a working `webgl2` context. It adds one script,
`scripts/ui-suite-in-chrome.sh`, one new Playwright project (`chrome`), and one new spec,
`tests/webgl.spec.ts`, that proves WebGL works. No documentation files needed changes: the code
itself already documents the new behavior at the point of use, and every other candidate doc file
sits outside this order's scope.

## Documentation Updated

None. See "Deferred Documentation" below for why, and "Code-Level Docs" for what already covers
this feature's behavior without a separate doc change.

## Deferred Documentation (Out of Scope)

The feature spec's `scope.paths` lists exactly: `apps/ye-olde-magic-shop/playwright.config.ts`,
`apps/ye-olde-magic-shop/package.json`, `apps/ye-olde-magic-shop/tests/**`, `.dockerignore`,
`scripts/**`, and `docs/features/ui-suite-real-chrome/**`. Acceptance criterion 6 states no file
outside that list may change, with the delivery record as the only exception. Approval Gate 6
("Writing Files out of Boundary") requires explicit human approval for any write outside this set.

Three files would normally get an update from this feature and are not on the allowed list:

- **`apps/ye-olde-magic-shop/README.md`** — could document `test:chrome` and the two environment
  variables as a "run in real Chrome" subsection under "Running Tests".
- **`docs/runner-image.md`** — line 18 says "Order `ui-suite-real-chrome` builds it from this
  page," written before this order existed; its evidence section (around line 227) still says
  `scripts/ui-suite-in-chrome.sh` "was **not** run" and "belongs to the separate ... order," which
  is now stale since that script exists and runs against the image.
- **`docs/features/README.md`** — the delivered-features index; a new row for this feature follows
  the same pattern as the three existing rows once the delivery record lands under
  `docs/features/ui-suite-real-chrome/`.

None of the three are touched here. A human should open a small follow-up change (or extend this
order's scope with explicit approval) to update them once this delivers.

## Documentation Persistence Convention

Per the project's pipeline convention, every artifact in
`.claude/feature-workspace/ui-suite-real-chrome/` — `analysis.md`, `architecture-notes.md`,
`implementation-notes.md`, `security-report.md`, `qa-report.md`, this report, and
`devops-report.md` once written — is copied to `docs/features/ui-suite-real-chrome/` by the
orchestrator after the pipeline finishes. That copy is the permanent record this feature produces;
it needs no separate action here.

## New Capabilities to Document (for the Deferred README Update)

For whoever picks up the deferred `apps/ye-olde-magic-shop/README.md` update, the facts to carry
over:

- **`npm run test:chrome`** runs the suite's `chrome` Playwright project directly, in a headed
  real Chrome, against whatever `channel: 'chrome'` resolves to on the host. `test:e2e` and
  `test:pw` stay pinned to `--project=chromium` and are unaffected.
- **`scripts/ui-suite-in-chrome.sh`**, run from the repository root with Docker available, builds
  the runner image, starts the app and mock API inside it, waits up to 180 seconds for both to
  answer a health check, runs `runner-check`, then the `chrome` project, and copies the Playwright
  report to `artifacts/ui-suite/`.
- **Exit codes** from the script: `2` no Docker, `3` image build failed, `4` app/mock API not
  healthy within 180 seconds, `5` `runner-check` failed inside the container. Any other non-zero
  code is the Playwright suite's own exit code, passed through unchanged.
- **`RUNNER_SLOW_MO`** (integer milliseconds, default `0`) and **`RUNNER_VIDEO`** (`1` enables
  video, default off) control the `chrome` project's `slowMo` and `video` options.

## Code-Level Docs

The implementation already carries the explanatory comments a reader needs, verified against the
actual files rather than assumed:

- `playwright.config.ts` documents, at the `chrome` project definition, that it runs inside the
  Saturday runner image and is only reached via `--project=chrome` or the script — not the default
  `test:e2e`/`test:pw` invocation. The ANGLE/SwiftShader launch-args constant carries a comment
  pointing at `runner-check`'s matching constant as the verified source of truth.
- `scripts/ui-suite-in-chrome.sh` documents `install_dependencies`' `CI=true` fix inline (why a
  bind-mounted `node_modules` can hang `pnpm install` without it) and carries the
  `PRECONDITION:`/`ENFORCED-BY:` comment pair at `wait_for_app_healthy`'s poll loop, per
  `architecture-guardrails.md` #7's judgment-only exception requirement.

No additional module or function doc comments are needed beyond what the developer already wrote.

## Manual Steps Required

- After this change merges, a human adds `{name: ui-suite, command: scripts/ui-suite-in-chrome.sh}`
  to `.claude/delivery-policy.yaml`'s `preShipChecks`, per the spec's Infrastructure/Deploy note.
  This is a config change outside `scope.paths`, so it is not done as part of this delivery.
- Running the script's exact unscoped command (`npx playwright test --project=chrome` with no file
  filter) currently requires `pnpm -r build` first — four sibling workspace packages have no built
  `dist/` in a fresh checkout. This is a pre-existing repository build-state gap, not introduced by
  this feature, and is not fixed here because none of the affected packages are in `scope.paths`.
- The deferred README, runner-image doc, and features-index updates listed above need a human
  decision: fold them into a fast-follow change, or approve widening this order's `scope.paths`.

## Known Gaps Carried Forward from QA

- **AC3** (WebGL test-of-the-test) was verified by running the exact runner container with the
  suite scoped to just this feature's two spec files, not via the script's unscoped invocation,
  because the unscoped command fails at test collection on the pre-existing sibling-package
  build-state gap described above. Both the passing and `--disable-gpu
  --disable-software-rasterizer`-failing runs were captured as real evidence from inside the
  container.
- **AC5** (`.dockerignore`'s `lib` negation) holds literally — `apps/ye-olde-magic-shop/lib/` is no
  longer excluded — but QA found no Docker build in this repository that currently uses the repo
  root as a COPY-based build context, so the fix has no observed effect on any build that runs
  today. Flagged for human review rather than treated as resolved.
