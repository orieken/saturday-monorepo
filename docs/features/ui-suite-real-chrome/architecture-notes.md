---
feature: "ui-suite-real-chrome"
bounded_context: "Test Automation"
domain_terms: ["Runner image"]
files_touched: ["apps/ye-olde-magic-shop/playwright.config.ts", "apps/ye-olde-magic-shop/tests/webgl.spec.ts", "scripts/ui-suite-in-chrome.sh", ".dockerignore"]
issue_refs: ["ui-suite-real-chrome"]
linked_adrs: []
linked_kis: []
---

# Architecture Notes: ui-suite-real-chrome

## Structural Decisions

### Add the 'chrome' project inside the existing projects array in apps/ye-olde-magic-shop/playwright.config.ts, alongside 'chromium'. Do not split into a second config file or a new package.
- **Rationale**: The config already has exactly the shape needed — a `projects` array with one entry. Adding a sibling entry is the established Playwright multi-project pattern already in use here; a second file would duplicate webServer/use/reporter settings for no benefit.
- **Trade-offs**:
  - Both projects share one config file's global `use` block, so any future global option change affects both projects and must be checked against each.
- **Fitness Function**: The 'chrome' Playwright project resolves with channel='chrome', headless=false, args ['--use-gl=angle','--use-angle=swiftshader'], and env-driven slowMo/video defaulting to off; a default (no --project) invocation resolves only 'chromium'. — verified by `npx playwright test --list --project=chrome` or an equivalent Node script reading the exported config, asserted by qa-engineer.
### The chrome project's launch args MUST be exactly `['--use-gl=angle', '--use-angle=swiftshader']` — the same literal array as `SOFTWARE_GL_ARGS` in containers/saturday-runner/runner-check/index.js, not a re-derived equivalent.
- **Rationale**: runner-check is the verified, baked-in contract for what makes webgl2 work in this image (containers/saturday-runner/runner-check/index.js:19, confirmed 2026-07-07 per docs/runner-image.md). The app suite's chrome project is a second consumer of that same contract; if the two drift, a later ANGLE/SwiftShader fix applied to one and not the other reintroduces the null-context bug silently.
- **Trade-offs**:
  - Couples the app's test config to a value defined in a container-image package this feature is scoped not to touch — if runner-check's args ever change, this project must change too, by inspection rather than by import (Playwright config and a container-internal Node script don't share a module boundary).
- **Fitness Function**: apps/ye-olde-magic-shop/playwright.config.ts's chrome project args equal containers/saturday-runner/runner-check/index.js's SOFTWARE_GL_ARGS. — verified by AC3's webgl.spec.ts pass/fail-under-disable-gpu test already exercises this at the behavior level; record the literal-array match as a one-line note in the delivery record for a human to catch future drift (judgment-only — no automated cross-file string-equality check is being added for two values in different languages/runtimes).
### Place the new orchestration script at scripts/ui-suite-in-chrome.sh, following scripts/update-framework.sh's conventions exactly: `set -euo pipefail`, `trap '...' EXIT` for any temp-resource cleanup, small single-purpose functions (build, start, wait_for_health, run_suite, copy_report), and a `main` function at the bottom calling them in order.
- **Rationale**: This repo already has one hand-rolled root-level orchestration script with an established shape. Matching it keeps scripts/ consistent rather than introducing a second convention for the same kind of file.
- **Trade-offs**:
  - Bash, not a typed language — correctness relies on shellcheck/review discipline, not a compiler.
- **Fitness Function**: scripts/ui-suite-in-chrome.sh produces a distinct exit code and stderr message for each of: no docker, image build failed, app unhealthy within 180s — each different from the Playwright suite's own exit code, which the script passes through unchanged on the normal path. — verified by Shell-level test stubbing docker/build/health-check calls and asserting exit code + message per branch (qa-engineer, per analyst's proposed fitness function).
### Exit-code namespace: reserve 2 = Docker not present, 3 = image build failed, 4 = app/mock-api unhealthy after 180s. Any other non-zero exit is the Playwright suite's own code, passed through verbatim (no re-mapping, no catch-all exit 1 for orchestration failures).
- **Rationale**: AC4 requires the three failure cases be distinguishable from "a Playwright test failure" (which exits 1). Reserving 2-4 for orchestration failures and never touching the suite's own exit code is the simplest rule that satisfies this without a lookup table or an exit-code registry.
- **Trade-offs**:
  - If a future failure mode needs its own code, the next free number is 5, decided ad hoc rather than from a central enum — acceptable at this script's current size (three failure modes).
- **Fitness Function**: Exit codes 2, 3, 4 are used only for 'no docker', 'build failed', 'unhealthy' respectively; the suite's own exit code is never remapped. — verified by Same shell-level test as the decision above; add one assertion that a forced Playwright test failure exits 1 (or whatever Playwright's own code is), unchanged by the script.
### The app/mock-API health-check wait in ui-suite-in-chrome.sh uses a bounded `until curl ...; do sleep N; done` loop with an explicit counter against the 180s budget — the same shape already used by containers/saturday-runner/entrypoint.sh's `wait_for_display` (XDPYINFO_TIMEOUT_SECONDS=15, entrypoint.sh:14-24). This is classified as a local-process readiness gate, not the request-retry logic architecture-guardrails.md #5 targets, and is not required to use a CircuitBreaker or ExponentialBackoffStrategy.
- **Rationale**: Guardrail #5 bans custom retry loops for resilience against a remote/unreliable dependency in application business logic, and requires an explicit timeout on every network call — this loop already has one (180s, NFR-stated). Distinguishing factor: this is a one-shot CI/local bootstrap script polling a process it itself just started, not a production code path calling out repeatedly to a service that might fail. The identical shape exists unchallenged in entrypoint.sh for the same reason (waiting for Xvfb to come up). Introducing an ExponentialBackoffStrategy abstraction for a fixed-interval local-process poll would be the premature generalization this framework's own design-principles.md warns against — the guardrail's intent is served by the explicit bound, not by borrowing a resilience pattern meant for unreliable remote calls.
- **Trade-offs**:
  - This is a judgment call, not a mechanically-checked boundary — a future reviewer skimming guardrail #5 without this note could flag the loop as a violation. The rationale above is the mitigation; restate it as a one-line comment at the loop site in the script itself.
- **Fitness Function**: The health-check loop has an explicit, named upper bound (180s) that the script enforces and reports on timeout; it does not retry indefinitely. — verified by The shell-level test's 'unhealthy within 180s' branch (same test as AC4) proves the bound is enforced, not just configured.
### .dockerignore gets a negation line `!apps/ye-olde-magic-shop/lib` immediately after the bare `lib` line (line 13), not a rewrite of the `lib` pattern itself and not a move to a narrower pattern elsewhere in the file.
- **Rationale**: Docker's ignore-pattern evaluation is order-sensitive and negations only work when a later, more specific rule un-excludes a path a broader earlier rule excluded. Placing the negation right after the pattern it corrects keeps the fix legible at the point of the problem, and the analyst already confirmed no earlier, broader pattern in this file re-excludes a parent of this path.
- **Trade-offs**:
  - Any other app under apps/*/lib/ that needs its own source directory preserved will need its own negation line added the same way — this is not a general-purpose fix for 'lib under apps/'.
- **Fitness Function**: .dockerignore no longer excludes apps/ye-olde-magic-shop/lib/ from the Docker build context, while still excluding other lib/ directories (e.g. package build output) from it. — verified by Build the image and confirm apps/ye-olde-magic-shop/lib/site/*.ts are present in the container; a .dockerignore-semantics unit test if one exists.

## Component Placement

- `'chrome' Playwright project` → framework layer, `apps/ye-olde-magic-shop/playwright.config.ts` — Test-runner configuration, sibling to the existing 'chromium' project entry.
- `webgl.spec.ts` → test layer, `apps/ye-olde-magic-shop/tests/webgl.spec.ts` — New Playwright spec, same directory as the suite's other specs.
- `ui-suite-in-chrome.sh` → tooling layer, `scripts/ui-suite-in-chrome.sh` — Root-level orchestration script, same location and shape as scripts/update-framework.sh.
- `.dockerignore negation line` → config layer, `.dockerignore` — One-line correction to an existing build-context exclusion rule.

## Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: containers/saturday-runner (the runner image) — read-only: this feature consumes runner-check's verified ANGLE/SwiftShader launch-arg contract and the image's debian build target without modifying the image, its Containerfile, entrypoint, or runner-check itself.

## Stability Design

- **Fail Fast with an explicit, named timeout** on scripts/ui-suite-in-chrome.sh's app/mock-API health-check wait — 180s bound (NFR-stated, amd64-emulation-aware), distinct exit code 4 on timeout, no indefinite retry — see the structural decision above on why this is not a CircuitBreaker/ExponentialBackoffStrategy case.

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] Clean Architecture dependency direction (Entities/UseCases/Adapters/Frameworks, architecture-guardrails.md #1) — None of the five changed paths (a Playwright config, a spec file, a shell script, a .dockerignore line, and a docs delivery record) sit inside a layered application codebase — they are test/tooling/config surfaces. mcp__loom__dependencies_of was not invoked because no directory under domain/application/adapter boundaries is touched; confirmed by Glob/Read of the affected paths instead.
- [NOT_APPLICABLE] Observability boundary — no OTel instrumentation inside domain/page logic (architecture-guardrails.md #8) — This feature adds no instrumentation; the existing @orieken/saturday-playwright-otel-reporter wiring in playwright.config.ts (ENABLE_OTEL branch) is untouched.

## Anti-Pattern Check

- Premature generalization: not found — Considered and rejected introducing a generic 'wait-for-ready' helper or a backoff-strategy abstraction for the health check; the fixed-bound loop matching entrypoint.sh's existing pattern is simpler and sufficient.
- Shotgun surgery: not found — Five files change, each for one reason; no change ripples into unrelated files.
- Leaky abstraction: not found — No new abstraction is introduced; this is additive config/script/test work inside existing shapes.

## Fitness Functions

- **The 'chrome' Playwright project resolves with channel='chrome', headless=false, args ['--use-gl=angle','--use-angle=swiftshader'], and env-driven slowMo/video defaulting to off; a default (no --project) invocation resolves only 'chromium'.**
  - **Verification**: `npx playwright test --list --project=chrome` or an equivalent Node script reading the exported config.
  - **Owner**: qa-engineer
- **scripts/ui-suite-in-chrome.sh produces exit codes 2/3/4 for 'no docker'/'build failed'/'unhealthy', and passes through the Playwright suite's own exit code unchanged on any other failure.**
  - **Verification**: Shell-level test stubbing docker/build/health-check calls, asserting exit code + stderr message per branch, plus one assertion that a forced test failure keeps Playwright's own exit code.
  - **Owner**: qa-engineer
- **tests/webgl.spec.ts passes under the chrome project's normal launch args and fails when --disable-gpu --disable-software-rasterizer is added to those args.**
  - **Verification**: Manual run of both variants inside the runner container; both outputs recorded in docs/features/ui-suite-real-chrome/.
  - **Owner**: qa-engineer
- **.dockerignore no longer excludes apps/ye-olde-magic-shop/lib/ from the build context, while still excluding other lib/ directories from it.**
  - **Verification**: Build the image and confirm apps/ye-olde-magic-shop/lib/site/*.ts exist inside the container.
  - **Owner**: qa-engineer
- **Cyclomatic complexity below 7 for every new/changed function; changed lines covered by a test that can fail.**
  - **Verification**: analyze-complexity skill + ADR-008 coverage gate (>=85% of changed lines).
  - **Owner**: qa-engineer

## Refactoring Opportunities

None

## Developer Handoff Notes

- Start from commit 034adbd's drafts (git show 034adbd:<path> for each of the five files) and check each against the acceptance criteria and the structural decisions above — do not rewrite from scratch.
- Match containers/saturday-runner/runner-check/index.js:19's SOFTWARE_GL_ARGS literally: ['--use-gl=angle', '--use-angle=swiftshader']. Leave a one-line comment in playwright.config.ts at the chrome project noting this value must track runner-check's, since the two can't share a module across the container/app boundary.
- Guard RUNNER_SLOW_MO parsing against NaN: `Number(process.env.RUNNER_SLOW_MO ?? 0)` can produce NaN on malformed input (e.g. RUNNER_SLOW_MO='abc'); fall back to 0 explicitly rather than passing NaN to slowMo.
- Keep project selection scoped so the existing test:e2e/test:pw npm scripts (no --project flag) only ever resolve 'chromium' — don't let Playwright's default project list pick up 'chrome' implicitly.
- In ui-suite-in-chrome.sh, model the health-check wait on entrypoint.sh's wait_for_display function (same repo, containers/saturday-runner/entrypoint.sh:14-24): a named timeout constant, an until-loop with sleep, and a clear exit message on timeout — reuse that shape rather than inventing a new one.
- Never set RUNNER_VIEWER in the script and never log anything from a code path that could echo the VNC password — entrypoint.sh already only starts the viewer when RUNNER_VIEWER=1 (entrypoint.sh:51-53), so the script's job is simply to never set that variable.
- Exit codes: reserve 2 (no docker), 3 (build failed), 4 (unhealthy); pass through the Playwright suite's own exit code on every other path, don't remap it.

## Open Architectural Questions

- Should a runner-check failure inside the container get its own distinct exit code (5), separate from both the orchestration codes (2-4) and the Playwright suite's own code? AC4 names only build-failure and app-unhealthy; architecture leaves this as the developer's call, with a light lean toward giving it its own code for the same reason 2/3/4 are separated — a human reading a bare '1' can't tell runner-check failed from an actual test assertion failing.
