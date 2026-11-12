---
feature: "reagent-map-visual-baselines"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: reagent-map-visual-baselines

## Site list visual baseline is captured and matched

Verifies: `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` MUST hold `@visual`-tagged `toHaveScreenshot()` tests for the site list, the page header with its nav link, and the view's error state when `/api/reagent-sites` fails through route interception.

```gherkin
Scenario: Site list visual baseline is captured and matched
  Given The reagent map page has loaded with the site list populated
  And RUNNER_VARIANT is set to a supported runner image
  When The site list region is captured as a screenshot in a @visual-tagged test
  Then The captured screenshot matches the stored baseline for the current runner image within the stated tolerance
```

## Page header and nav link visual baseline is captured and matched

Verifies: `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` MUST hold `@visual`-tagged `toHaveScreenshot()` tests for the site list, the page header with its nav link, and the view's error state when `/api/reagent-sites` fails through route interception.

```gherkin
Scenario: Page header and nav link visual baseline is captured and matched
  Given The reagent map page has loaded
  And RUNNER_VARIANT is set to a supported runner image
  When The page header including its nav link is captured as a screenshot in a @visual-tagged test
  Then The captured screenshot matches the stored baseline for the current runner image within the stated tolerance
```

## Error state visual baseline is captured when the reagent sites request fails

Verifies: `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` MUST hold `@visual`-tagged `toHaveScreenshot()` tests for the site list, the page header with its nav link, and the view's error state when `/api/reagent-sites` fails through route interception.

```gherkin
Scenario: Error state visual baseline is captured when the reagent sites request fails
  Given The reagent map page is loading
  And The `/api/reagent-sites` request is intercepted and aborted, mirroring reagent-map-resilience.spec.ts's REAGENT_SITES_API_PATTERN
  And RUNNER_VARIANT is set to a supported runner image
  When The rendered error state is captured as a screenshot in a @visual-tagged test
  Then A visible element with data-testid="reagent-map-error" and role="alert" is present
  And The captured screenshot matches the stored baseline for the current runner image within the stated tolerance
```

## Globe default view visual baseline is captured within tolerance

Verifies: One `@visual` test MUST cover `#cesiumContainer` at the fixed default view, after `cesium.waitForSettled()`, with a stated `maxDiffPixelRatio` that never exceeds 1%.

```gherkin
Scenario: Globe default view visual baseline is captured within tolerance
  Given The globe viewer has finished building and has settled via cesium.waitForSettled()
  And The camera is at the fixed default view
  And RUNNER_VARIANT is set to a supported runner image
  When #cesiumContainer is captured as a screenshot in a @visual-tagged test with maxDiffPixelRatio set to at most 0.01
  Then The captured screenshot matches the stored baseline within a tolerance that never exceeds 0.01
```

## Globe test flakiness is resolved without loosening tolerance

Verifies: If the globe `@visual` test is flaky under the five-run repeat criterion, the developer MUST remove it, mask the canvas in the other `@visual` tests, and record why in the delivery record, instead of raising the tolerance past 1%.

```gherkin
Scenario: Globe test flakiness is resolved without loosening tolerance
  Given The globe default-view visual test has been run five consecutive times in each runner image
  When The five-run results are evaluated for flakiness
  Then If any run differs, the globe visual test is removed, the canvas region is masked in the remaining @visual tests, and the reason is recorded in the delivery record
  And No @visual test's maxDiffPixelRatio is raised above 0.01 as a response to the flakiness
```

## Baselines resolve to the correct per-image directory

Verifies: `playwright.config.ts` MUST store baselines per runner image under `.playwright-baselines/{debian,ubi9}/` using `snapshotPathTemplate`, keyed by a `RUNNER_VARIANT` environment variable that `scripts/ui-suite-in-chrome.sh` sets.

```gherkin
Scenario: Baselines resolve to the correct per-image directory
  Given playwright.config.ts defines snapshotPathTemplate keyed by RUNNER_VARIANT
  When A @visual test runs with RUNNER_VARIANT set to 'debian' or to 'ubi9'
  Then The test reads and writes its baseline under .playwright-baselines/debian/ or .playwright-baselines/ubi9/ matching the value of RUNNER_VARIANT
```

## An unrecognized RUNNER_VARIANT fails loudly instead of falling back

Verifies: `playwright.config.ts` MUST store baselines per runner image under `.playwright-baselines/{debian,ubi9}/` using `snapshotPathTemplate`, keyed by a `RUNNER_VARIANT` environment variable that `scripts/ui-suite-in-chrome.sh` sets.

```gherkin
Scenario: An unrecognized RUNNER_VARIANT fails loudly instead of falling back
  Given RUNNER_VARIANT is set to a value other than 'debian' or 'ubi9'
  When A @visual test attempts to resolve its baseline path
  Then The test run fails with an error naming the unrecognized RUNNER_VARIANT value
  And No baseline comparison silently proceeds against another image's baseline directory
```

## Visual tests skip with a stated reason when RUNNER_VARIANT is unset

Verifies: When `RUNNER_VARIANT` is unset, as on a developer's Mac, every `@visual` test MUST skip with a message stating why. When it is set, the tests MUST never skip, and a missing baseline MUST fail the test.

```gherkin
Scenario: Visual tests skip with a stated reason when RUNNER_VARIANT is unset
  Given RUNNER_VARIANT is unset in the test environment
  When Any @visual-tagged test runs
  Then The test is skipped
  And The skip message states that RUNNER_VARIANT is unset
```

## A missing baseline fails the test when RUNNER_VARIANT is set

Verifies: When `RUNNER_VARIANT` is unset, as on a developer's Mac, every `@visual` test MUST skip with a message stating why. When it is set, the tests MUST never skip, and a missing baseline MUST fail the test.

```gherkin
Scenario: A missing baseline fails the test when RUNNER_VARIANT is set
  Given RUNNER_VARIANT is set to a supported runner image
  And The expected baseline file for a @visual test does not exist on disk
  When The @visual test runs
  Then The test fails
  And The test does not skip
```

## Baselines are only regenerated when --update-baselines is passed

Verifies: `scripts/ui-suite-in-chrome.sh --update-baselines`, alone or combined with `--ubi9`, MUST regenerate that image's baselines by passing `--update-snapshots` to Playwright. Without the flag, the script MUST never write a baseline. `ui-suite-in-chrome.test.sh` MUST cover the flag.

```gherkin
Scenario: Baselines are only regenerated when --update-baselines is passed
  Given scripts/ui-suite-in-chrome.sh is invoked without --update-baselines
  When The script runs the @visual suite
  Then No baseline file under .playwright-baselines/ is created or overwritten
```

## --update-baselines regenerates baselines for the targeted image

Verifies: `scripts/ui-suite-in-chrome.sh --update-baselines`, alone or combined with `--ubi9`, MUST regenerate that image's baselines by passing `--update-snapshots` to Playwright. Without the flag, the script MUST never write a baseline. `ui-suite-in-chrome.test.sh` MUST cover the flag.

```gherkin
Scenario: --update-baselines regenerates baselines for the targeted image
  Given scripts/ui-suite-in-chrome.sh is invoked with --update-baselines, optionally combined with --ubi9
  When The script runs
  Then Playwright is invoked with --update-snapshots
  And Only the targeted runner image's baselines under .playwright-baselines/ are written
  And ui-suite-in-chrome.test.sh's assertions for this flag pass
```

## The visual suite passes five consecutive runs in both runner images

Verifies: The suite MUST pass five consecutive runs in each image (debian and ubi9), and the run output MUST be recorded in the delivery record.

```gherkin
Scenario: The visual suite passes five consecutive runs in both runner images
  Given The @visual suite is runnable in the debian runner image
  And The @visual suite is runnable in the ubi9 runner image
  When Each image's suite is executed five consecutive times
  Then All five runs pass in the debian image
  And All five runs pass in the ubi9 image
  And The run output is recorded in the delivery record
```

## A one-line site list style change fails the matching visual test in both images

Verifies: A one-line style change to the site list MUST fail the matching `@visual` test in both images. A moved default camera on the globe test, if that test stays, MUST do the same. The failing output MUST be recorded, then the change reverted.

```gherkin
Scenario: A one-line site list style change fails the matching visual test in both images
  Given The stored visual baselines are unchanged from the last approved run
  And A one-line style change is applied to the site list
  When The @visual suite runs in the debian and ubi9 images
  Then The matching @visual test fails in both images
  And The failing output is recorded in the delivery record before the change is reverted
```

## A moved default camera fails the globe visual test in both images, if that test is kept

Verifies: A one-line style change to the site list MUST fail the matching `@visual` test in both images. A moved default camera on the globe test, if that test stays, MUST do the same. The failing output MUST be recorded, then the change reverted.

```gherkin
Scenario: A moved default camera fails the globe visual test in both images, if that test is kept
  Given The globe default-view visual test was kept rather than removed for flakiness
  And The default camera position is moved
  When The @visual suite runs in the debian and ubi9 images
  Then The globe visual test fails in both images
  And The failing output is recorded in the delivery record before the change is reverted
```

## Baseline regeneration is documented and the Known Gaps entry is closed

Verifies: `docs/testing-cesium.md` and `docs/runner-image.md` MUST document how and when to regenerate baselines, including after a browser lock bump. `testing-cesium.md`'s 'No visual baselines' Known Gaps entry MUST be closed. Both documents MUST follow `.claude/rules/writing-style.md`.

```gherkin
Scenario: Baseline regeneration is documented and the Known Gaps entry is closed
  Given docs/testing-cesium.md currently lists 'No visual baselines' under Known Gaps
  And docs/runner-image.md exists
  When The visual baseline feature ships
  Then docs/testing-cesium.md's 'No visual baselines' Known Gaps entry is closed
  And docs/runner-image.md documents how to regenerate baselines, including after a browser lock bump
  And Both documents conform to .claude/rules/writing-style.md
```

## Changed code stays within complexity, coverage, and scope limits

Verifies: Every new or changed function MUST have cyclomatic complexity below 7. Every changed line MUST be covered by a test that can fail. No file outside `scope.paths` MUST change, except this run's own delivery record under `docs/features/reagent-map-visual-baselines/`.

```gherkin
Scenario: Changed code stays within complexity, coverage, and scope limits
  Given The feature's implementation is complete
  When Every new or changed function is analyzed for cyclomatic complexity and every changed line is checked for failing-capable test coverage
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test that can fail
  And No file outside scope.paths has changed except this run's own delivery record under docs/features/reagent-map-visual-baselines/
```

## Error-state element renders when fetching reagent sites fails

Verifies: `ReagentMapView.vue` MUST render a visible, stable error-state element when `fetchReagentSites()` rejects, since none exists today and the error-state screenshot in the first acceptance criterion needs something stable to capture.

```gherkin
Scenario: Error-state element renders when fetching reagent sites fails
  Given buildViewer() resolves successfully
  And fetchReagentSites() rejects
  When ReagentMapView.vue's onMounted logic handles the rejected sitesResult
  Then A visible element with data-testid="reagent-map-error" and role="alert" is rendered just before the existing throw of sitesResult.reason
  And The existing pageerror event, newViewer.destroy() call, and the zero-canvas/zero-list-item assertions in reagent-map-resilience.spec.ts still pass
```

## Error-state element does not render when the globe viewer itself fails to build

Verifies: `ReagentMapView.vue` MUST render a visible, stable error-state element when `fetchReagentSites()` rejects, since none exists today and the error-state screenshot in the first acceptance criterion needs something stable to capture.

```gherkin
Scenario: Error-state element does not render when the globe viewer itself fails to build
  Given buildViewer() rejects
  When ReagentMapView.vue's onMounted logic runs
  Then The viewer rejection is rethrown before the sitesResult branch is reached
  And The data-testid="reagent-map-error" element is not rendered for this failure path
```
