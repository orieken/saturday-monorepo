---
intake: 1
id: ui-suite-consolidated-report
type: feature
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
depends-on: [ui-suite-run-evidence]
budget:
  max-attempts: 2
  max-usd: 30
  max-hours: 3
scope:
  paths: [packages/saturday-playwright-report/**, apps/ye-olde-magic-shop/package.json, pnpm-lock.yaml, scripts/ui-suite-in-chrome.sh, scripts/ui-suite-in-chrome.test.sh, docs/runner-image.md, docs/features/ui-suite-consolidated-report/**]
---

# Add a consolidated UI test report with videos

## Request

> testing-with-certs has a consolidated test report with embedded videos. Add that reporting to
> Saturday. The UI suite runs in the debian and ubi9 runner images.

## Summary

Generate one page, `artifacts/ui-suite/index.html`, that shows the debian and ubi9 results together.
Each test has its result and an inline video player. The page opens from disk with no server.

## Starting point

Use the testing-with-certs report as the reference. Do not copy its Go operator: Saturday has no
operator, and its stack for this work is TypeScript. Read these files in that repository:

- `operator/matrix_videos.go`: finds videos and writes the page.
- `operator/matrix_gallery.html`: the page, with cards, filters and search.
- `operator/matrix_video_metadata.go`: links each video to its test and result.
- `operator/matrix_merge.go`: merges several runs into one report.
- `docs/completion/48-matrix-video-gallery.md` and `49-matrix-report-style-spec-results.md`: the
  decisions and their reasons.

Keep these behaviours from that report:

- Native `<video>` players with relative links, `preload="none"` and no autoplay. Videos are files
  beside the page. They are not base64 in the page.
- Empty recordings are excluded, symlinks are rejected, and all labels are HTML-escaped.
- The page uses the real result of each test. A passed suite does not turn a skipped test into a pass.
- A video that cannot be linked to a test is labelled as unlinked. The generator never guesses from
  file order.
- Generated files are written to a temporary path and then renamed.

Change these behaviours:

- Read test metadata from `results.json`, the documented Playwright format. Do not decode the zip
  inside the HTML report, which is an internal format.
- Load no web font. The testing-with-certs page loads Google Fonts, but the UI suite makes no request
  off the container and this page MUST work offline.
- Keep one copy of each video on the page. Link to `test-results/`, so video counts equal recordings.

## Acceptance Criteria

- [ ] Given `packages/saturday-playwright-report`, then it provides a CLI that takes a root directory,
      reads every `<variant>/run.json` and `<variant>/results.json` below it, and writes
      `index.html` in the root.
- [ ] Given `scripts/ui-suite-in-chrome.sh`, when a run finishes, then it runs the CLI in the container
      and rebuilds the page from all evidence directories. The suite exit code wins. If the suite
      passes and the CLI fails, the script exits with a new, distinct code and says why.
- [ ] Given a failed suite, then the script still writes the page, and the page shows the failures and
      their videos.
- [ ] Given the page, then its summary shows, for each variant, the commit SHA, dirty state, run time,
      and passed, failed, flaky and skipped counts.
- [ ] Given the two variants ran on different commits, or one ran from a dirty checkout, then the
      summary states this as a warning. Given one variant has no evidence, then the summary states
      that it did not run.
- [ ] Given a test, then its card shows the spec file, test title, project, variant, result, duration,
      video player, download link and a link to that variant's Playwright HTML report.
- [ ] Given a test with several attempts, then the card shows each attempt with its own result and
      video. Given a test with no recording, then the card says so and still shows its result.
- [ ] Given the filters, then the reader can filter by variant, project and result, and can search by
      text. A filter shows the reader when it matches nothing.
- [ ] Given a test whose result differs between the variants, then a filter shows only those tests.
- [ ] Given attachment paths in `results.json`, then they are absolute paths from the container, so the
      generator rebases them onto `<variant>/test-results/`. It MUST reject a path that escapes the
      evidence directory.
- [ ] Given a `results.json` over 20 MB, or one that is not valid, then the CLI exits non-zero and names
      the file.
- [ ] Given error text and steps from a test, then the page does not show them. The reader follows the
      link to the Playwright HTML report for those details.
- [ ] Given the page opened from `file://` in Chrome, then it makes no network request and no `fetch`,
      and one video from each variant reaches ready state 4 with `currentTime` advancing. Record the
      check in the delivery record.
- [ ] Given a test injected to fail in one variant, then the page shows that test as failed in that
      variant, and the differing-results filter finds it. Record the output, and revert the change.
- [ ] Given `docs/runner-image.md`, then it says where the report is, how to rebuild it and that the
      report stays local. It follows `.claude/rules/writing-style.md`.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under
      `docs/features/ui-suite-consolidated-report/` excepted.

## Out of Scope

- Publishing the report, for example to GitHub Pages, and any CI workflow. That needs its own order,
  and it needs the provenance and privacy checks in testing-with-certs' `docs/pages-report.md`.
- Cucumber recordings. The UI suite script runs only the Playwright `chrome` project.
- An Edge project.
- Authentication filters and the operator's other matrix dimensions (OS, runner, user, suite).
- Any change to `containers/saturday-runner/` or to testing-with-certs.

## Domain Language

- **Consolidated report**: `artifacts/ui-suite/index.html`, one page over all evidence directories.
- **Variant** and **Evidence directory**: as order `ui-suite-run-evidence` defines them.

## Non-Functional Requirements

- The page MUST work when the `artifacts/ui-suite/` directory moves as a whole, so all links are relative.
- Building the page from a full run takes less than 10 seconds in the runner image.

## Trust Boundaries

- Test titles and file paths are untrusted input to an HTML page, so the generator escapes them and
  encodes URLs. Videos show synthetic demo data and stay on the local machine.

## Test Approach

Vitest in the package, with fixture evidence that includes both variants, a retry, a skip, a missing
video and a path that escapes. `scripts/ui-suite-in-chrome.test.sh` covers the script's new behaviour.
Real runs of the script in both images prove the page.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
