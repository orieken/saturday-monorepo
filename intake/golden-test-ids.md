---
intake: 1
id: golden-test-ids
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
depends-on: [golden-site-file-names]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/lib/**, docs/**]
---

# Give every element the tests touch a data-testid

## Request

> `CartPage` guesses `button:has-text("COMPLETE TRANSACTION")` because the test id was never added.

## Summary

Every element a page object registers is located by `data-testid` or by role, never by text or CSS
structure, so a copy change cannot break a test.

## Acceptance Criteria

- [ ] Given every `registerElement` call in `lib/site/`, then its locator is
      `[data-testid="..."]` or a role-based locator, and the Vue component renders that test id.
- [ ] Given `CartPage`'s complete-transaction button, then it is located by
      `data-testid="complete-transaction-btn"`.
- [ ] Given the change is complete, then no locator in `lib/site/` uses `has-text`, `:text`, `nth-`,
      or a selector built from element nesting.
- [ ] Given the existing suites, then every scenario that passed before passes after, with no
      `.feature` file or assertion changed.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/golden-test-ids/` excepted.

## Out of Scope

- Step definitions; visual or behaviour changes to the app.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Cucumber and Saturday/Playwright, existing suites.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
