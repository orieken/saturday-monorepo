---
intake: 1
id: golden-steps-through-site
type: chore
plan: deliver-feature
status: draft
priority: 2
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [golden-test-ids, golden-flows]
budget:
  max-attempts: 2
  max-usd: 30
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/features/step_definitions/**, apps/ye-olde-magic-shop/features/support/**, apps/ye-olde-magic-shop/lib/**, docs/**]
---

# Make step definitions talk only to the site

## Request

> Step definitions bypass the model: 25 uses of `this.page` and `getByTestId` across five files.

## Summary

Steps call the site, its flows and its pages; locators and page handling live only in the site
model.

## Acceptance Criteria

- [ ] Given `features/step_definitions/`, then no file uses `this.page`, `getByTestId`, `locator(`,
      `$(`, `waitForTimeout` or `setTimeout`.
- [ ] Given a step that needed one of those, then the behaviour moved to a page method or a flow,
      named for what it does (`accountPage.newestOrderStatus()`), and the step calls that.
- [ ] Given login and checkout steps, then they call `loginFlow` and `checkoutFlow`.
- [ ] Given the existing feature files, then every scenario that passed before passes after, with no
      `.feature` file changed and no `Then` asserting less than before.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/golden-steps-through-site/` excepted.

## Out of Scope

- New scenarios.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Cucumber, existing suite.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
