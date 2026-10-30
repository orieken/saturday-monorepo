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
  max-usd: 40
  max-hours: 6
scope:
  paths: [apps/ye-olde-magic-shop/playwright.config.ts, apps/ye-olde-magic-shop/package.json, apps/ye-olde-magic-shop/tests/**, containers/saturday-runner/**, .dockerignore, scripts/**, docs/**]
---

# Run the demo app's UI suite in real Chrome under Xvfb

## Request

> The map tests have to run in a real Chrome container, not Playwright's bundled Chromium. Cesium
> needs WebGL. Combine the best of Saturday's ubi-vnc image and testing-with-certs' runner images
> into the most robust container, with the tools we need.

## Summary

The demo app's Playwright suite runs today in Playwright's bundled Chromium, on the host. This
order builds the Saturday runner image designed in `docs/runner-image.md`, runs the suite in it
headed in real Chrome under Xvfb, proves WebGL works there, and gives the factory one command to
run it as a pre-ship check.

## Acceptance Criteria

- [ ] Given `containers/saturday-runner/Containerfile`, then its `debian` target matches
      `docs/runner-image.md`, "The combined design": the pinned `node:24.19.0-bookworm-slim` digest
      from that page; Chrome and Edge from repositories signed by the keys in
      `containers/saturday-runner/keys/`, checked against `SHA256SUMS` before use; Xvfb, openbox,
      x11vnc, noVNC, `x11-utils`, Mesa GL, fonts and dbus; pnpm pinned; ffmpeg from
      `playwright install ffmpeg`; runs as uid 1000. It contains no `gpgcheck=0`, `--nogpgcheck`,
      `curl | bash`, `COPY . .`, root user or supervisord.
- [ ] Given `containers/saturday-runner/apt.lock`, then it pins every package installed, and
      `containers/saturday-runner/lock-packages.sh` regenerates it from the package list in the same
      directory. The build fails if a pinned version is unavailable.
- [ ] Given the image's entrypoint, when it starts, then it starts Xvfb at 1920x1080x24 with
      `-nolisten tcp`, waits until `xdpyinfo` succeeds (failing after 15 seconds), starts openbox,
      and runs its command with `exec`, so the container's exit code is the command's.
- [ ] Given `RUNNER_VIEWER=1`, then the entrypoint also starts x11vnc with a random password on
      localhost only and noVNC on port 6080, and prints where the password is. Without it, nothing
      listens on any port.
- [ ] Given `runner-check` in the image, when run, then it launches Chrome and Edge headed, creates a
      `webgl2` context in each, writes each renderer string to `/work/evidence/webgl.txt` with the
      browser versions and installed packages, and exits non-zero if either context is null.
- [ ] Given `apps/ye-olde-magic-shop/playwright.config.ts`, then it has a `chrome` project with
      `channel: 'chrome'` and `headless: false`, and that project is what the container runs. The
      bundled-Chromium project may stay for local runs. `slowMo` and `video` are read from
      environment variables and default to off.
- [ ] Given `scripts/ui-suite-in-chrome.sh`, when run from the repository root on a machine with
      Docker, then it builds the image (`--platform linux/amd64`), mounts the checkout at `/work`
      with a named pnpm store volume, installs dependencies with `--frozen-lockfile`, starts the app
      and the mock API, runs `runner-check` and then the `chrome` project, copies the report to
      `artifacts/ui-suite/`, and exits with the suite's exit code. The app and mock API make no
      request off the container after dependencies are installed.
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

- The `polyglot` target (JDK, Python, Gradle, Go): order `runner-image-polyglot`.
- A UBI9 variant.
- testing-with-certs adopting this image as its base; that is a separate order in that repo.
- Changing any existing test's behaviour.

## Domain Language

- **Runner image**: the container the browser suites run in, defined by `docs/runner-image.md`.

## Non-Functional Requirements

- Display readiness within 15 seconds; app health within 180 seconds.
- With `RUNNER_VIEWER` unset, the container listens on no port.

## Trust Boundaries

- Browser packages cross into the image from Google's and Microsoft's repositories: accepted only
  when signed by the vendored keys, whose hashes are recorded.
- The VNC viewer is a remote desktop into the test browser: password-protected, bound to
  localhost, off by default.

## Test Approach

Saturday/Playwright, through `scripts/ui-suite-in-chrome.sh`: `runner-check`, the new WebGL spec, and
the existing suite.

## Open Questions

- None

## Infrastructure / Deploy Notes

After merge, a human adds `{name: ui-suite, command: scripts/ui-suite-in-chrome.sh}` to
`.claude/delivery-policy.yaml` `preShipChecks`.
