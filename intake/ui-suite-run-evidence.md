---
intake: 1
id: ui-suite-run-evidence
type: chore
plan: deliver-feature
status: draft
priority: 2
created: 2026-11-04
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: []
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 3
scope:
  paths: [scripts/ui-suite-in-chrome.sh, scripts/ui-suite-in-chrome.test.sh, apps/ye-olde-magic-shop/playwright.config.ts, apps/ye-olde-magic-shop/tests/playwright-config.spec.ts, .gitignore, docs/runner-image.md, docs/features/ui-suite-run-evidence/**]
---

# Keep each runner image's UI suite evidence

## Request

> testing-with-certs has a consolidated test report with embedded videos. Add that reporting to
> Saturday. The UI suite runs in the debian and ubi9 runner images.

## Summary

This order is the first of two. It keeps the evidence that the report needs. Order
`ui-suite-consolidated-report` builds the report from it.

Today the suite cannot feed a consolidated report, for four reasons:

- `copy_report` in `scripts/ui-suite-in-chrome.sh` writes both images to
  `artifacts/ui-suite/playwright-report`, so the ubi9 run deletes the debian evidence.
- `RUNNER_VIDEO` is unset in `preShipChecks`, so the `chrome` project records no video.
- The suite writes no machine-readable results.
- `.gitignore` covers `*.webm` but not `artifacts/`, so a generated file there can be committed.

The runner image needs no change. Both targets already install Playwright's ffmpeg.

## Acceptance Criteria

- [ ] Given `scripts/ui-suite-in-chrome.sh`, when a run starts, then it deletes only
      `artifacts/ui-suite/<variant>/`, where `<variant>` is `debian` or `ubi9`. A run in one image
      MUST NOT change the other image's directory.
- [ ] Given a run that reaches the suite, then `artifacts/ui-suite/<variant>/` holds
      `playwright-report/`, `test-results/` with the videos, `results.json` and `run.json`. The
      script writes these files when the suite fails too.
- [ ] Given `run.json`, then it records the variant, image tag, commit SHA, whether the checkout was
      dirty, start time, finish time and suite exit code.
- [ ] Given `playwright.config.ts`, then it keeps its current reporters and adds the Playwright
      `json` reporter, which writes to `test-results/results.json`.
- [ ] Given the script with `RUNNER_VIDEO` unset, then the suite records video. Given
      `RUNNER_VIDEO=0`, then it records no video. `resolveRunnerVideo` keeps its default of off, so
      a run outside the script does not change.
- [ ] Given both images, then the suite passes three runs in a row with video on, at `--workers=3`.
      Record the duration with video on and off in the delivery record. If video makes the suite
      flaky or more than doubles its duration, stop and say so.
- [ ] Given both images run in turn, then both directories exist, and every recorded video is
      nonempty and identified as WebM. Record the output.
- [ ] Given `.gitignore`, then it ignores `artifacts/`.
- [ ] Given `scripts/ui-suite-in-chrome.test.sh`, then it covers each criterion above that the
      script owns, and it MUST NOT delete evidence from a real run. Its cleanup now removes
      `artifacts/ui-suite`, so point the script at a temporary directory during the test.
- [ ] Given `RUNNER_VARIANT`, then the script sets it for the suite. Order
      `reagent-map-visual-baselines` needs the same variable. Use that name, and reuse the variable
      if that order merges first.
- [ ] Given `docs/runner-image.md`, then it states the evidence layout and the video default, and it
      follows `.claude/rules/writing-style.md`.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under
      `docs/features/ui-suite-run-evidence/` excepted.

## Out of Scope

- The report page: order `ui-suite-consolidated-report`.
- Any change to `containers/saturday-runner/`. If the suite needs one, stop and say so.
- An Edge project, Cucumber recordings and any CI workflow.
- Changing any test's behaviour.

## Domain Language

- **Variant**: `debian` or `ubi9`, as `docs/runner-image.md` defines it.
- **Evidence directory**: `artifacts/ui-suite/<variant>/`, the output of one run in one variant.

## Non-Functional Requirements

- The suite exit code stays the script's exit code. The existing distinct exit codes do not change.
- The suite still runs with `--workers=3`, `--shm-size=1g` and the same 180-second app health limit.

## Trust Boundaries

- None. No auth, no user input and no external calls. The videos show synthetic demo data and stay
  on the local machine.

## Test Approach

`scripts/ui-suite-in-chrome.test.sh` with the stubbed `docker`, `tests/playwright-config.spec.ts`,
and real runs of the script in both images.

## Open Questions

- Video adds time to every pre-ship check: Assumption if unresolved: record video by default, because
  the report is most useful at the ship gate. `RUNNER_VIDEO=0` turns it off.

## Infrastructure / Deploy Notes

- Any tool or person that reads `artifacts/ui-suite/playwright-report` MUST read
  `artifacts/ui-suite/<variant>/playwright-report` after this merges. No other file in the repository
  reads the old path today.
