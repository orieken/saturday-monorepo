---
feature: "ui-suite-real-chrome"
bounded_context: "Test Automation"
domain_terms: ["Runner image"]
files_touched: ["apps/ye-olde-magic-shop/playwright.config.ts", "apps/ye-olde-magic-shop/package.json", "apps/ye-olde-magic-shop/tests/webgl.spec.ts", "scripts/ui-suite-in-chrome.sh", ".dockerignore", "docs/features/ui-suite-real-chrome/"]
issue_refs: ["ui-suite-real-chrome"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: ui-suite-real-chrome

## Summary

Wire the demo app's (ye-olde-magic-shop) Playwright suite to the already-built, already-verified Saturday runner image (containers/saturday-runner/, debian target) so the map/Cesium tests run headed in real Chrome under Xvfb with software WebGL (ANGLE SwiftShader), instead of Playwright's bundled Chromium which cannot provide a WebGL context the way this suite needs. A first attempt was killed at a stage timeout mid-build and left working drafts on commit 034adbd (branch salvage/ui-suite-real-chrome): apps/ye-olde-magic-shop/playwright.config.ts (new 'chrome' project), apps/ye-olde-magic-shop/tests/webgl.spec.ts (webgl2 context assertion), scripts/ui-suite-in-chrome.sh (build/run/health-check/runner-check/suite/report orchestration), a package.json test:chrome script, and a .dockerignore fix so the app's own lib/ source directory isn't excluded from the Docker build context. This order verifies and finishes those drafts against the acceptance criteria rather than rewriting them from scratch; it does not touch the runner image itself.

### Acceptance Criteria

- playwright.config.ts MUST define a 'chrome' project with channel 'chrome', headless false, and launch args --use-gl=angle --use-angle=swiftshader, and the container MUST run that project. The bundled-Chromium 'chromium' project MUST remain for local runs. The chrome project's slowMo and video MUST come from environment variables and default to off.
  - Example: RUNNER_SLOW_MO unset -> slowMo resolves to 0 (off)
  - Example: RUNNER_VIDEO unset -> video resolves to 'off'
  - Example: RUNNER_VIDEO=1 -> video resolves to 'on'
- scripts/ui-suite-in-chrome.sh, run from the repo root on a machine with Docker, MUST build the image from containers/saturday-runner/Containerfile with --platform linux/amd64 --target debian without changing that Containerfile, mount the checkout at /work with a named pnpm store volume, install dependencies with --frozen-lockfile, start the app and the mock API, run runner-check and then the chrome project in that order, copy the Playwright report to artifacts/ui-suite/, and exit with the suite's exit code. The app and mock API MUST make no request off the container after dependencies are installed.
- tests/webgl.spec.ts MUST assert a webgl2 context is available on the running app's page, and MUST fail when the browser is launched with --disable-gpu --disable-software-rasterizer added to the chrome project's launch args. That launch MUST be used once as a test-of-the-test and the result recorded in the delivery record.
- If the image fails to build, or the app is not healthy within 180 seconds, the script MUST exit non-zero and state which of the two happened, using an exit code and message distinct from a Playwright test failure's exit code.
- .dockerignore's bare 'lib' pattern MUST no longer exclude apps/ye-olde-magic-shop/lib/ from the Docker build context, without removing that pattern's effect on other lib/ directories (e.g. package build output).
- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change, except this run's own delivery record under docs/features/ui-suite-real-chrome/.

### Non-Functional Requirements

- **performance**: The app and mock API must report healthy before the suite starts running against them; amd64 emulation on Apple Silicon is slow, so the budget is generous. (app/mock-api health check timeout 180s)
- **security**: With RUNNER_VIEWER unset, the container must listen on no port; scripts/ui-suite-in-chrome.sh must never set RUNNER_VIEWER and must never print the VNC viewer's password, even though the script never touches the viewer today.
- **other**: The app and mock API processes must make no network request off the container after pnpm install completes (no telemetry, no package re-fetch mid-run).

## Proposed Fitness Functions

- **The 'chrome' Playwright project resolves with channel='chrome', headless=false, the two ANGLE/SwiftShader launch args, and env-driven slowMo/video defaulting to off.**
  - **Verification**: QA test asserting the resolved project config (e.g. `npx playwright test --project=chrome --list` or a small Node script reading the exported config) matches these fields.
  - **Owner**: qa-engineer
- **scripts/ui-suite-in-chrome.sh produces a distinct exit code and message for 'no docker', 'image build failed', and 'app unhealthy within 180s', separate from a Playwright test-failure exit code.**
  - **Verification**: Shell-level test (e.g. stubbing docker/build/health-check calls) exercising each failure branch and asserting its exit code and stderr message.
  - **Owner**: qa-engineer
- **tests/webgl.spec.ts passes under the chrome project's normal launch args and fails when --disable-gpu --disable-software-rasterizer is added to those args.**
  - **Verification**: Manual run of both variants inside the runner container; both outputs recorded in docs/features/ui-suite-real-chrome/. Not continuously run in CI, since it requires mutating launch args by hand.
  - **Owner**: qa-engineer
- **.dockerignore no longer excludes apps/ye-olde-magic-shop/lib/ from the build context, and no earlier broader pattern re-excludes a parent directory of that path.**
  - **Verification**: Build the image and confirm apps/ye-olde-magic-shop/lib/site/*.ts exist inside the build context/container (or a .dockerignore-semantics unit test).
  - **Owner**: qa-engineer
- **Running `npm run test:e2e` / `npm run test:pw` locally (outside the container, no --project flag) continues to pass without attempting to launch the headed chrome project.**
  - **Verification**: Local run of both npm scripts on a developer machine without Xvfb; recorded as a regression check.
  - **Owner**: qa-engineer
- **Cyclomatic complexity below 7 for every new/changed function; changed lines covered by a test that can fail.**
  - **Verification**: analyze-complexity skill + ADR-008 coverage gate (>=85% of changed lines).
  - **Owner**: qa-engineer

## Out of Scope

- Any change to containers/saturday-runner/ (the image, its locks, keys, entrypoint, or runner-check) — a separate order if the suite needs one.
- The polyglot target (JDK, Python, Gradle, Go) — order runner-image-polyglot.
- A UBI9 variant — order runner-image-ubi9.
- testing-with-certs adopting this image as its base — a separate order in that repo.
- Changing any existing test's behaviour (assertions in map.spec.ts, cart.spec.ts, magic-shop.spec.ts, heatmap.spec.ts, example.spec.ts stay as-is).
- Wiring scripts/ui-suite-in-chrome.sh into .claude/delivery-policy.yaml's preShipChecks — spec states a human does this after merge; not part of this pipeline run.
- Updating DOMAIN_DICTIONARY.md for the spec's 'Runner image' term — that term is pre-existing infrastructure vocabulary already authoritatively defined in docs/runner-image.md (itself out of scope), not a new business concept this spec introduces, and DOMAIN_DICTIONARY.md is outside scope.paths.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: containers/saturday-runner (the runner image) — read-only dependency; this order consumes runner-check and the image's launch-arg contract without modifying the image.

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/playwright.config.ts` — Add the 'chrome' project (channel, headless, ANGLE/SwiftShader launch args, env-driven slowMo/video) alongside the existing 'chromium' project; verify default (no --project) runs don't pick up the headed chrome project locally.
- `apps/ye-olde-magic-shop/package.json` — Add a test:chrome script (`playwright test --project=chrome`) consistent with existing test:e2e/test:pw/test:bdd scripts.
- `apps/ye-olde-magic-shop/tests/webgl.spec.ts` — New spec asserting a webgl2 context exists on the running app's page; must be provably able to fail (AC3).
- `scripts/ui-suite-in-chrome.sh` — New script: builds the image, mounts the repo, installs deps, starts app+mock-api, waits for health, runs runner-check then the chrome project, copies the report, and propagates the right exit code.
- `.dockerignore` — Un-exclude apps/ye-olde-magic-shop/lib/ from the bare 'lib' pattern so the Site-Centric pattern classes under it are present in the Docker build context.
- `docs/features/ui-suite-real-chrome/` — Permanent delivery record: AC3's test-of-the-test evidence, AC4's negative-path evidence, and the end-to-end script run's output.

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Verify/finish apps/ye-olde-magic-shop/playwright.config.ts against commit 034adbd's draft: 'chrome' project with channel 'chrome', headless false, args ['--use-gl=angle', '--use-angle=swiftshader'], slowMo from Number(process.env.RUNNER_SLOW_MO ?? 0) and video from process.env.RUNNER_VIDEO === '1' ? 'on' : 'off'; guard against RUNNER_SLOW_MO resolving to NaN on malformed input.
2. Ensure the default Playwright invocation (no --project flag), as used by the existing test:e2e/test:pw npm scripts, does not attempt to launch the headed 'chrome' project — scope project selection so those scripts keep running only 'chromium' unless --project=chrome (or scripts/ui-suite-in-chrome.sh) is used explicitly.
3. Add a test:chrome script to apps/ye-olde-magic-shop/package.json: "playwright test --project=chrome".
4. Create apps/ye-olde-magic-shop/tests/webgl.spec.ts (per the 034adbd draft) asserting canvas.getContext('webgl2') is non-null on the app's running page.
5. Create scripts/ui-suite-in-chrome.sh (per the 034adbd draft) with: docker-presence check; `docker build --platform linux/amd64 --target debian -f containers/saturday-runner/Containerfile ...` without editing that Containerfile; a named pnpm store volume; `pnpm install --frozen-lockfile`; starting the app+mock-api (npm run dev, which already runs both via concurrently); polling health on localhost:8000/ and localhost:8001/api/health up to 180s; running runner-check then `playwright test --project=chrome`; copying the report to artifacts/ui-suite/; exiting with the suite's own exit code. Give the three failure paths (no docker, build failed, app unhealthy) distinct exit codes and messages, each clearly different from a test-failure exit code.
6. Fix .dockerignore: add a negation for apps/ye-olde-magic-shop/lib/ immediately after the bare 'lib' line, and confirm no earlier, broader pattern in the file re-excludes a parent directory of that path (none currently does).

### QA Tasks

1. Run the chrome project (apps/ye-olde-magic-shop/tests/webgl.spec.ts and the existing suite) inside the actual runner container via scripts/ui-suite-in-chrome.sh and confirm it passes.
2. Prove webgl.spec.ts can fail (AC3): temporarily add --disable-gpu --disable-software-rasterizer to the chrome project's launch args, re-run, confirm the test fails, then revert; record both outputs in docs/features/ui-suite-real-chrome/.
3. Run scripts/ui-suite-in-chrome.sh end-to-end on a machine with Docker; record exit code, timing, and the 180s health-check headroom actually used, in the delivery record.
4. Force each of the three failure paths (docker absent, image build failure, app never becomes healthy) and confirm each produces a distinct, correctly-labeled non-zero exit, separate from a Playwright test-failure exit code.
5. Confirm `npm run test:e2e` and `npm run test:pw`, run locally outside the container, still pass without attempting to launch the headed chrome project — regression check against the existing suite's behavior.
6. Verify cyclomatic complexity below 7 for every new/changed function (script functions, webgl.spec.ts) and that every changed line is covered by a test able to fail (ADR-008, >=85% of changed lines).
7. Confirm no file outside scope.paths (apps/ye-olde-magic-shop/playwright.config.ts, apps/ye-olde-magic-shop/package.json, apps/ye-olde-magic-shop/tests/**, .dockerignore, scripts/**, docs/features/ui-suite-real-chrome/**) changed.

### Tech Writer Tasks

1. Within docs/features/ui-suite-real-chrome/ only (docs/runner-image.md is out of scope for this order), write up how to run scripts/ui-suite-in-chrome.sh as the factory's pre-ship check, and record AC3's test-of-the-test evidence and AC4's negative-path evidence captured by qa-engineer.

### DevOps Tasks

None

## Edge Cases and Risks

- Default `playwright test` invocation (no --project flag) would, without guarding, resolve both the existing 'chromium' project and the new headed 'chrome' project.: Scope project selection so npm run test:e2e / test:pw keep exercising only 'chromium' by default; 'chrome' runs only via --project=chrome or scripts/ui-suite-in-chrome.sh. Otherwise local dev machines without a display, or CI machines without the Chrome channel installed, would start failing a previously-working script.
- RUNNER_SLOW_MO is set to a non-numeric value.: Number(process.env.RUNNER_SLOW_MO ?? 0) would resolve to NaN and get passed to Playwright's slowMo launch option; validate or default safely instead of passing NaN through.
- runner-check fails inside the container (WebGL genuinely unavailable in the image).: The 034adbd draft exits 1 here, the same value a generic Playwright test failure is likely to use. AC4 only requires distinctness for 'build failed' and 'app unhealthy', so this isn't a blocking gap, but developer should consider a dedicated exit code for clearer diagnosis.
- Apple Silicon amd64 emulation slows image build and app startup.: Keep the full 180s health-check budget (already generous per docs/runner-image.md's own build record); do not shorten it to speed up local iteration.
- The app's existing map view fetches real map tiles from tile.openstreetmap.org client-side at /locations.: Pre-existing behavior, unrelated to this order (out of scope: changing any existing test's behaviour) and not what AC2's 'app and mock API make no request off the container' governs — that clause is about the app/mock-api server processes, not browser-initiated page requests. Note as a known source of flakiness if the chrome project runs offline, but do not change it here.
- The app process crashes shortly after starting, rather than just being slow to become healthy.: The health-check polling loop in the 034adbd draft can't distinguish 'slow' from 'crashed' — both fall through to the same 180s timeout and the same tail -n 50 of the app log printed to stderr, which is sufficient for AC4's 'says which' requirement (distinct from a build failure).
- The spec's Domain Language section defines 'Runner image', which would normally trigger a DOMAIN_DICTIONARY.md update under the analyst's Ubiquitous Language rule.: Skip: DOMAIN_DICTIONARY.md is outside this order's scope.paths, and 'Runner image' is pre-existing infrastructure vocabulary already canonically defined in docs/runner-image.md (also out of scope), not a new concept this spec introduces.

## Definition of Done

- All six acceptance criteria verified with evidence captured under docs/features/ui-suite-real-chrome/.
- scripts/ui-suite-in-chrome.sh runs end-to-end against a real Docker daemon and exits 0 on a green suite.
- webgl.spec.ts is proven to fail under --disable-gpu --disable-software-rasterizer and to pass under the normal chrome-project launch args; both runs recorded.
- Existing suite's assertions are unchanged; npm run test:e2e / test:pw continue to work locally without attempting the headed chrome project.
- Cyclomatic complexity below 7 and the ADR-008 coverage threshold (>=85% of changed lines) are met for every new/changed function.
- No file outside scope.paths changed, except this run's own delivery record under docs/features/ui-suite-real-chrome/.
