---
feature: "ui-suite-real-chrome"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: ui-suite-real-chrome

## The chrome project runs against the real Chrome browser, not the bundled Chromium

Verifies: playwright.config.ts MUST define a 'chrome' project with channel 'chrome', headless false, and launch args --use-gl=angle --use-angle=swiftshader, and the container MUST run that project. The bundled-Chromium 'chromium' project MUST remain for local runs.

```gherkin
Scenario: The chrome project runs against the real Chrome browser, not the bundled Chromium
  Given The playwright config defines a 'chrome' project alongside the existing 'chromium' project
  When The suite runs the 'chrome' project
  Then The browser launched is the system Chrome identified by channel 'chrome'
  And The browser window runs headed, not headless
```

## The chrome project launches with the ANGLE SwiftShader GL backend

Verifies: playwright.config.ts MUST define a 'chrome' project with channel 'chrome', headless false, and launch args --use-gl=angle --use-angle=swiftshader, and the container MUST run that project. The bundled-Chromium 'chromium' project MUST remain for local runs.

```gherkin
Scenario: The chrome project launches with the ANGLE SwiftShader GL backend
  When The suite runs the 'chrome' project
  Then The browser launch arguments include --use-gl=angle and --use-angle=swiftshader
```

## Slow motion defaults to off when RUNNER_SLOW_MO is not set

Verifies: The chrome project's slowMo and video MUST come from environment variables and default to off.

```gherkin
Scenario: Slow motion defaults to off when RUNNER_SLOW_MO is not set
  Given The RUNNER_SLOW_MO environment variable is unset
  When The chrome project resolves its launch options
  Then slowMo resolves to 0 (off)
```

## Video recording defaults to off when RUNNER_VIDEO is not set

Verifies: The chrome project's slowMo and video MUST come from environment variables and default to off.

```gherkin
Scenario: Video recording defaults to off when RUNNER_VIDEO is not set
  Given The RUNNER_VIDEO environment variable is unset
  When The chrome project resolves its launch options
  Then video resolves to 'off'
```

## Video recording turns on when RUNNER_VIDEO is set to 1

Verifies: The chrome project's slowMo and video MUST come from environment variables and default to off.

```gherkin
Scenario: Video recording turns on when RUNNER_VIDEO is set to 1
  Given The RUNNER_VIDEO environment variable is set to '1'
  When The chrome project resolves its launch options
  Then video resolves to 'on'
```

## The bundled-Chromium project still runs for local development

Verifies: playwright.config.ts MUST define a 'chrome' project with channel 'chrome', headless false, and launch args --use-gl=angle --use-angle=swiftshader, and the container MUST run that project. The bundled-Chromium 'chromium' project MUST remain for local runs.

```gherkin
Scenario: The bundled-Chromium project still runs for local development
  Given A developer runs the suite locally without Docker
  When The developer runs the default Playwright test command (test:e2e or test:pw) with no --project flag
  Then Only the 'chromium' project executes
  And The 'chrome' project does not launch
```

## The containerized UI suite script builds the merged runner image and runs the suite inside it

Verifies: scripts/ui-suite-in-chrome.sh, run from the repo root on a machine with Docker, MUST build the image from containers/saturday-runner/Containerfile with --platform linux/amd64 --target debian without changing that Containerfile, mount the checkout at /work with a named pnpm store volume, install dependencies with --frozen-lockfile, start the app and the mock API, run runner-check and then the chrome project in that order, copy the Playwright report to artifacts/ui-suite/, and exit with the suite's exit code.

```gherkin
Scenario: The containerized UI suite script builds the merged runner image and runs the suite inside it
  Given The repo is checked out on a machine with Docker
  When scripts/ui-suite-in-chrome.sh is run from the repo root
  Then The image builds from containers/saturday-runner/Containerfile with --platform linux/amd64 --target debian
  And Dependencies install inside the container with --frozen-lockfile
  And The app and the mock API start inside the container
  And runner-check runs before the chrome project
  And The chrome project runs only after runner-check succeeds
```

## The containerized run produces a Playwright report and exits with the suite's result

Verifies: scripts/ui-suite-in-chrome.sh, run from the repo root on a machine with Docker, MUST ... copy the Playwright report to artifacts/ui-suite/, and exit with the suite's exit code.

```gherkin
Scenario: The containerized run produces a Playwright report and exits with the suite's result
  Given scripts/ui-suite-in-chrome.sh has finished running the chrome project
  When The script completes
  Then The Playwright report is copied to artifacts/ui-suite/
  And The script's exit code matches the chrome project's own exit code
```

## The app and mock API make no requests outside the container once dependencies are installed

Verifies: The app and mock API MUST make no request off the container after dependencies are installed.

```gherkin
Scenario: The app and mock API make no requests outside the container once dependencies are installed
  Given Dependencies have finished installing inside the container
  When The app server and mock API server start and serve the chrome project's requests
  Then Neither the app server nor the mock API server makes a network request to a host outside the container
```

## The running app exposes a WebGL2 rendering context

Verifies: tests/webgl.spec.ts MUST assert a webgl2 context is available on the running app's page.

```gherkin
Scenario: The running app exposes a WebGL2 rendering context
  Given The chrome project is running against the app inside the container
  When The app's page loads
  Then A webgl2 context is available on the page
```

## Disabling the GPU and software rasterizer makes the WebGL test fail (test-of-the-test)

Verifies: tests/webgl.spec.ts MUST fail when the browser is launched with --disable-gpu --disable-software-rasterizer added to the chrome project's launch args. That launch MUST be used once as a test-of-the-test and the result recorded in the delivery record.

```gherkin
Scenario: Disabling the GPU and software rasterizer makes the WebGL test fail (test-of-the-test)
  Given The chrome project's launch args temporarily include --disable-gpu and --disable-software-rasterizer
  When tests/webgl.spec.ts runs against that launch
  Then The test fails, proving it can detect a missing WebGL2 context
  And The failing run's output is captured in the delivery record alongside the normal passing run
```

## A failed image build is reported distinctly from a test failure

Verifies: If the image fails to build ... the script MUST exit non-zero and state which of the two happened, using an exit code and message distinct from a Playwright test failure's exit code.

```gherkin
Scenario: A failed image build is reported distinctly from a test failure
  Given containers/saturday-runner/Containerfile fails to build
  When scripts/ui-suite-in-chrome.sh runs
  Then The script exits non-zero
  And The script states that the image build failed
  And That exit code differs from the exit code a Playwright test failure produces
```

## An app that never becomes healthy is reported distinctly from a test failure

Verifies: If ... the app is not healthy within 180 seconds, the script MUST exit non-zero and state which of the two happened, using an exit code and message distinct from a Playwright test failure's exit code.

```gherkin
Scenario: An app that never becomes healthy is reported distinctly from a test failure
  Given The app does not respond healthy within 180 seconds of starting
  When scripts/ui-suite-in-chrome.sh runs
  Then The script exits non-zero
  And The script states that the app health check timed out
  And That exit code differs from the exit code a Playwright test failure produces
```

## The shop app's lib directory is included in the Docker build context

Verifies: .dockerignore's bare 'lib' pattern MUST no longer exclude apps/ye-olde-magic-shop/lib/ from the Docker build context, without removing that pattern's effect on other lib/ directories.

```gherkin
Scenario: The shop app's lib directory is included in the Docker build context
  Given .dockerignore previously excluded apps/ye-olde-magic-shop/lib/ via a bare 'lib' pattern
  When The runner image builds
  Then apps/ye-olde-magic-shop/lib/ is present in the build context
```

## Other lib/ build-output directories remain excluded from the Docker build context

Verifies: .dockerignore's bare 'lib' pattern MUST no longer exclude apps/ye-olde-magic-shop/lib/ from the Docker build context, without removing that pattern's effect on other lib/ directories.

```gherkin
Scenario: Other lib/ build-output directories remain excluded from the Docker build context
  Given A package's build output lives in a lib/ directory elsewhere in the repo
  When The runner image builds
  Then That other lib/ directory is still excluded from the build context
```

## No file outside the order's declared scope changes

Verifies: Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change, except this run's own delivery record under docs/features/ui-suite-real-chrome/.

```gherkin
Scenario: No file outside the order's declared scope changes
  Given The feature's scope.paths lists playwright.config.ts, package.json, tests/**, .dockerignore, scripts/**, and docs/features/ui-suite-real-chrome/**
  When The feature's changes are complete
  Then Every changed file's path matches one of the declared scope paths
  And Any file outside that list that changed lives only under docs/features/ui-suite-real-chrome/
```
