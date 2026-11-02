---
feature: "golden-workspace-packages"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: golden-workspace-packages

## ye-olde-magic-shop's package.json uses workspace dependencies for all three Saturday packages

Verifies: apps/ye-olde-magic-shop/package.json MUST declare @orieken/saturday-core, @orieken/saturday-cucumber and @orieken/saturday-playwright-heatmap as workspace:* dependencies, with no dependency in the file pointing at a path inside the app.

```gherkin
Scenario: ye-olde-magic-shop's package.json uses workspace dependencies for all three Saturday packages
  Given apps/ye-olde-magic-shop/package.json currently declares file: paths for saturday-core, saturday-cucumber, and saturday-playwright-heatmap
  When the workspace package dependency migration is applied
  Then apps/ye-olde-magic-shop/package.json declares @orieken/saturday-core, @orieken/saturday-cucumber, and @orieken/saturday-playwright-heatmap as workspace:* dependencies
  And no dependency entry in apps/ye-olde-magic-shop/package.json points at a path inside the app
```

## Vendored framework copies are removed and any uniquely covering tests move into their packages

Verifies: The vendored directories apps/ye-olde-magic-shop/saturday-core/, saturday-cucumber/ and saturday-playwright-heatmap/ MUST be gone. Any test inside them that covers framework code and has no content-equivalent already under packages/ MUST move there instead of being deleted.

```gherkin
Scenario: Vendored framework copies are removed and any uniquely covering tests move into their packages
  Given apps/ye-olde-magic-shop/saturday-core/, saturday-cucumber/, and saturday-playwright-heatmap/ exist with vendored framework code and tests
  When the vendored directories are removed in favor of the workspace packages
  Then apps/ye-olde-magic-shop/saturday-core/, saturday-cucumber/, and saturday-playwright-heatmap/ no longer exist
  And every vendored test that covered framework behavior not already covered under packages/*/tests/** now lives under packages/*/tests/** instead of being deleted
```

## Every Cucumber scenario and Playwright spec that passed before the migration still passes after it

Verifies: Every Cucumber scenario and Playwright spec that passed before the change MUST still pass after it, with no .feature file or spec assertion changed. The before and after results MUST be recorded in the delivery record.

```gherkin
Scenario: Every Cucumber scenario and Playwright spec that passed before the migration still passes after it
  Given the full set of Cucumber scenarios and Playwright specs in apps/ye-olde-magic-shop pass on the current codebase before the migration
  When the workspace package migration and vendored directory removal are applied
  Then the same set of Cucumber scenarios and Playwright specs pass after the migration with no .feature file or spec assertion changed
  And the before and after results are recorded in the delivery record
```

## Dockerfiles that copied vendored packages are updated to build from the workspace or removed

Verifies: apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner, which both COPY the vendored directories, MUST each build from the workspace instead, or be removed with the reason recorded.

```gherkin
Scenario: Dockerfiles that copied vendored packages are updated to build from the workspace or removed
  Given apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner each COPY the vendored saturday-core, saturday-cucumber, or saturday-playwright-heatmap directories
  When the vendored directories are removed from the app
  Then apps/ye-olde-magic-shop/Dockerfile and Dockerfile.runner each either build from the workspace packages or are removed
  And the reason for removal, if that path is chosen, is recorded in the delivery record
```

## The Chrome UI suite script runs clean on a fresh clone after the migration

Verifies: Running scripts/ui-suite-in-chrome.sh from the repository root on a fresh clone MUST exit 0, with its output recorded in the delivery record.

```gherkin
Scenario: The Chrome UI suite script runs clean on a fresh clone after the migration
  Given a fresh clone of the repository with the workspace package migration applied and pnpm-lock.yaml regenerated
  When scripts/ui-suite-in-chrome.sh is run from the repository root
  Then the script exits 0
  And its output is recorded in the delivery record
```

## Changed code stays within complexity, coverage, and file-scope limits

Verifies: Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except this run's own delivery record under docs/features/golden-workspace-packages/.

```gherkin
Scenario: Changed code stays within complexity, coverage, and file-scope limits
  Given the set of functions and lines changed or added by this migration
  When the change is reviewed for complexity, test coverage, and file scope
  Then every new or changed function has cyclomatic complexity below 7
  And every changed line is covered by a test that can fail
  And no file outside scope.paths changed except this run's own delivery record under docs/features/golden-workspace-packages/
```
