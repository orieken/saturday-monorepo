---
intake: 1
id: ui-suite-real-chrome
type: chore
plan: deliver-feature
status: ready
priority: 1
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: []
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/playwright.config.ts, apps/ye-olde-magic-shop/package.json, apps/ye-olde-magic-shop/tests/**, .dockerignore, scripts/**, docs/features/ui-suite-real-chrome/**]
---

# Run the demo app's UI suite in real Chrome under Xvfb

## Request

> The map tests have to run in a real Chrome container, not Playwright's bundled Chromium. Cesium
> needs WebGL. Combine the best of Saturday's ubi-vnc image and testing-with-certs' runner images
> into the most robust container, with the tools we need.

## Summary

The Saturday runner image already exists: `containers/saturday-runner/` (debian target), built and
verified by hand per `docs/runner-image.md`, "Debian variant: first build record". This order does
not change the image. It wires the demo app's Playwright suite to it: a `chrome` project that runs
headed in real Chrome under Xvfb, a spec that proves WebGL works there, and one command the factory
can run as a pre-ship check.

## Starting point

The first attempt (killed at the stage timeout while building the image) left working drafts of
every file this order needs. They are on commit `034adbd` (branch `salvage/ui-suite-real-chrome`),
not on `main`. Start from them and check each against the criteria below; do not rewrite them from
scratch:

```bash
git show 034adbd:scripts/ui-suite-in-chrome.sh
git show 034adbd:apps/ye-olde-magic-shop/playwright.config.ts
git show 034adbd:apps/ye-olde-magic-shop/tests/webgl.spec.ts
git show 034adbd:apps/ye-olde-magic-shop/package.json
git show 034adbd:.dockerignore
```

## Acceptance Criteria

- [ ] Given `apps/ye-olde-magic-shop/playwright.config.ts`, then it has a `chrome` project with
      `channel: 'chrome'`, `headless: false`, and launch args `--use-gl=angle
      --use-angle=swiftshader` (the same args `runner-check` uses; see `docs/runner-image.md`), and
      that project is what the container runs. The bundled-Chromium project may stay for local runs.
      `slowMo` and `video` are read from environment variables and default to off.
- [ ] Given `scripts/ui-suite-in-chrome.sh`, when run from the repository root on a machine with
      Docker, then it builds the image from `containers/saturday-runner/Containerfile`
      (`--platform linux/amd64 --target debian`) without changing it, mounts the checkout at
      `/work` with a named pnpm store volume, installs dependencies with `--frozen-lockfile`,
      starts the app and the mock API, runs `runner-check` and then the `chrome` project, copies the
      report to `artifacts/ui-suite/`, and exits with the suite's exit code. The app and mock API
      make no request off the container after dependencies are installed.
- [ ] Given the suite, then `tests/webgl.spec.ts` asserts a `webgl2` context is available in the
      running app's page, and fails when the browser is launched with `--disable-gpu
      --disable-software-rasterizer`. Use that launch as the test of the test, recorded in the
      delivery record.
- [ ] Given the image fails to build, or the app is not healthy within 180 seconds (amd64 emulation
      on Apple Silicon is slow), then the script exits non-zero and says which, distinct from a test
      failure.
- [ ] Given `.dockerignore`, then the bare `lib` pattern no longer excludes
      `apps/ye-olde-magic-shop/lib/`.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under
      `docs/features/ui-suite-real-chrome/` excepted.

## Out of Scope

- Any change to `containers/saturday-runner/` (the image, its locks, keys, entrypoint or
  `runner-check`). If the suite needs one, stop and say so; it is a separate order.
- The `polyglot` target (JDK, Python, Gradle, Go): order `runner-image-polyglot`.
- A UBI9 variant: order `runner-image-ubi9`.
- testing-with-certs adopting this image as its base; that is a separate order in that repo.
- Changing any existing test's behaviour.

## Domain Language

- **Runner image**: the container the browser suites run in, defined by `docs/runner-image.md`.

## Non-Functional Requirements

- App health within 180 seconds.
- With `RUNNER_VIEWER` unset, the container listens on no port; the script never sets it.

## Trust Boundaries

- The VNC viewer is a remote desktop into the test browser: password-protected, bound to
  localhost, off by default. The script must not enable it or print its password.

## Test Approach

Saturday/Playwright, through `scripts/ui-suite-in-chrome.sh`: `runner-check`, the new WebGL spec, and
the existing suite.

## Open Questions

- None

## Infrastructure / Deploy Notes

After merge, a human adds `{name: ui-suite, command: scripts/ui-suite-in-chrome.sh}` to
`.claude/delivery-policy.yaml` `preShipChecks`.
