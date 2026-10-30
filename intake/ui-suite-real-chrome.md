---
intake: 1
id: ui-suite-real-chrome
type: chore
plan: deliver-feature
status: draft
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
  paths: [apps/ye-olde-magic-shop/playwright.config.ts, apps/ye-olde-magic-shop/package.json, containers/**, scripts/**, docs/**]
---

# Run the demo app's UI suite in real Chrome under Xvfb

## Request

> The map tests have to run in the real Chrome container, not Playwright's bundled Chromium. Cesium
> needs WebGL, and our containers already launch real Chrome under X11.

## Summary

The demo app's Playwright suite runs today in Playwright's bundled Chromium with no channel set, on
the host. Every later order, Cesium above all, needs it to run in real Chrome under Xvfb, in a
container, with one command the factory can use as a pre-ship check.

## Acceptance Criteria

- [ ] Given `apps/ye-olde-magic-shop/playwright.config.ts`, then it has a `chrome` project with
      `channel: 'chrome'`, and that project is what the container runs. The bundled-Chromium
      project may stay for local runs.
- [ ] Given `slowMo` and `video`, then each is read from an environment variable and defaults to
      off, so the suite does not take a second per action in CI or in the factory.
- [ ] Given one of the existing real-Chrome images under `containers/` (`ubi-vnc` or
      `x11-forwarding-ubi9`), then the change picks it, says why in `docs/`, and makes it build
      from the repository root. `containers/Dockerfile` copies a `supervisord.conf` that is not in
      `containers/`: fix it or remove it, and say which.
- [ ] Given `scripts/ui-suite-in-chrome.sh`, when run from the repository root on a machine with
      Docker, then it builds that image, starts the app and the mock API inside it, runs the
      `chrome` project headed under Xvfb, copies the report to `artifacts/ui-suite/`, and exits with
      the suite's exit code. It needs no network after the image is built.
- [ ] Given the script runs, then a test in the suite asserts WebGL is available: a page creates a
      `webgl2` context, and the context is not null and its renderer is reported.
- [ ] Given the image fails to build or the app does not become healthy within 120 seconds, then
      the script exits non-zero with a message that names which, distinct from a test failure.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/ui-suite-real-chrome/` excepted.

## Out of Scope

- Changing any test's behaviour. Edge, other browsers and the Cucumber runner.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- The script's container startup is bounded: 120 seconds to healthy.

## Trust Boundaries

- None. No auth, no user input, no external calls at test time.

## Test Approach

Saturday/Playwright. The WebGL check and the existing suite, run through the script.

## Open Questions

- Which image: Assumption if unresolved: `x11-forwarding-ubi9`, which has a compose file and docs.

## Infrastructure / Deploy Notes

After merge, a human adds `{name: ui-suite, command: scripts/ui-suite-in-chrome.sh}` to
`.claude/delivery-policy.yaml` `preShipChecks`.
